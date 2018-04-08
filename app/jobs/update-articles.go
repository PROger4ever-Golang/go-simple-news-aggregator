package jobs

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/modules/orm/gorp/app"
	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"gopkg.in/xmlpath.v2"

	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
)

var currentNodePath *xmlpath.Path

func init() {
	currentNodePath = xmlpath.MustCompile(".")
}

type ArticlesUpdater struct{}

func (c *ArticlesUpdater) GetAbsolutePath(base *url.URL, path string) string {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return base.ResolveReference(u).String()
}

func (c *ArticlesUpdater) GetDomRoot(url string) *xmlpath.Node {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var rootNode *xmlpath.Node
	contentTypeHeader := resp.Header.Get("Content-Type")
	if strings.Contains(contentTypeHeader, "html") {
		rootNode, err = xmlpath.ParseHTML(resp.Body)
	} else if strings.Contains(contentTypeHeader, "rss") {
		rootNode, err = xmlpath.Parse(resp.Body)
	} else {
		panic(errors.New("unsupported page content type"))
	}
	if err != nil {
		panic(err)
	}
	return rootNode
}

func (c *ArticlesUpdater) GetArticles(rootNode *xmlpath.Node, source *models.Source, url string) []*models.Article {
	var iter *xmlpath.Iter
	if len(source.CardXpath) == 0 {
		iter = currentNodePath.Iter(rootNode)
	} else {
		iter = xmlpath.MustCompile(source.CardXpath).Iter(rootNode)
	}

	articles := make([]*models.Article, 0, 1)
	for iter.Next() {
		var err error
		cardNode := iter.Node()
		article := models.Article{SourceId: source.SourceId, Source: source}
		articles = append(articles, &article)

		iter := xmlpath.MustCompile(source.TitleXpath).Iter(cardNode)
		if !iter.Next() {
			panic(errors.New("can't find title node in DOM"))
		}
		article.Title = iter.Node().String()

		iter = xmlpath.MustCompile(source.BodyXpath).Iter(cardNode)
		if !iter.Next() {
			panic(errors.New("can't find body node in DOM"))
		}
		article.Body = iter.Node().String()

		iter = xmlpath.MustCompile(source.ImgXpath).Iter(cardNode)
		if iter.Next() {
			imgUrlNode := iter.Node()
			article.ImgUrl = imgUrlNode.String()
		}

		iter = xmlpath.MustCompile(source.PublishedXpath).Iter(cardNode)
		if !iter.Next() {
			panic(errors.New("can't find published node in DOM"))
		}
		article.Published, err = time.Parse(source.PublishedFormat, iter.Node().String())
		if err != nil {
			err = fmt.Errorf("can't parse published date '%s' using the format '%s' because of: %s",
				iter.Node().String(),
				source.PublishedFormat,
				err,
			)
			panic(err)
		}
	}

	if len(articles) == 0 {
		panic(errors.New("can't find article cards in DOM"))
	}
	return articles
}
func (c *ArticlesUpdater) ProcessArticleNode(rootNode *xmlpath.Node, source *models.Source, url string) {
	articles := c.GetArticles(rootNode, source, url)

	for _, article := range articles {
		err := rgorp.Db.Map.Insert(article)
		if err == nil {
			fmt.Printf("The article added: %s\n", article.Title)
			continue
		}
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 2067 {
			fmt.Printf("The article already exists. It's ok: %s\n", article.Title)
		} else {
			panic(err)
		}
	}
}

func (c *ArticlesUpdater) ProcessSource(source *models.Source) {
	rootNode := c.GetDomRoot(source.Url)

	if len(source.ArticleUrlsXpath) == 0 {
		c.ProcessArticleNode(rootNode, source, source.Url)
		return
	}

	path := xmlpath.MustCompile(source.ArticleUrlsXpath)
	articleUrls := path.Iter(rootNode)

	base, err := url.Parse(source.Url)
	if err != nil {
		panic(err)
	}
	for articleUrls.Next() {
		articleUrl := c.GetAbsolutePath(base, articleUrls.Node().String())
		articlePageNode := c.GetDomRoot(articleUrl)
		c.ProcessArticleNode(articlePageNode, source, articleUrl)
	}
}

func (c *ArticlesUpdater) Run() {
	var sources []*models.Source
	builder := gorp.Db.SqlStatementBuilder.Select("*").From("Source")
	if _, err := gorp.Db.Select(&sources, builder); err != nil {
		panic(err)
	}
	for _, source := range sources {
		c.ProcessSource(source)
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 30s", &ArticlesUpdater{})
	})
}
