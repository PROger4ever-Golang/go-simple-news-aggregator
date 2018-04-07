package jobs

import (
	"fmt"

	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
	"github.com/revel/modules/orm/gorp/app"
)

// Periodically count the articles in the database.
type ArticleCounter struct{}

func (c ArticleCounter) Run() {
	articles, err := gorp.Db.Map.Select(models.Article{},
		`select * from Article`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d articles.\n", len(articles))
}

func init() {
	//TODO: Update articles from sources
	//revel.OnAppStart(func() {
	//	jobs.Schedule("@every 1m", ArticleCounter{})
	//})
}
