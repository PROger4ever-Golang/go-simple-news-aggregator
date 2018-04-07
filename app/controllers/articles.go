package controllers

import (
	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
	"github.com/revel/revel"
	"gopkg.in/Masterminds/squirrel.v1"
)

type Articles struct {
	Application
}

func (c Articles) loadArticleSafe(sourceId, id int) *models.Article {
	article := models.Article{}
	builder := c.Db.SqlStatementBuilder.Select("*").
		From("Article").
		Where(squirrel.And{
			squirrel.Eq{"SourceId": sourceId},
			squirrel.Eq{"ArticleId": id},
		})
	if err := c.Txn.SelectOne(&article, builder); err != nil {
		c.Log.Fatal("Unexpected error loading an article", "error", err)
	}
	return &article
}

func (c Articles) Index(sourceId int, size, page uint64) revel.Result {
	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	prevPage := page - 1
	nextPage := page + 1

	var articles []*models.Article
	builder := c.Db.SqlStatementBuilder.Select("*").From("Article").Offset((page - 1) * size).Limit(size)
	if _, err := c.Txn.Select(&articles, builder); err != nil {
		c.Log.Fatal("Unexpected error loading articles", "error", err)
	}
	return c.Render(articles, size, page, prevPage, nextPage)
}

func (c Articles) Show(sourceId, id int) revel.Result {
	article := c.loadArticleSafe(sourceId, id)
	if article == nil {
		return c.NotFound("Article %d/%d does not exist", sourceId, id)
	}
	title := article.Title
	return c.Render(title, article)
}
