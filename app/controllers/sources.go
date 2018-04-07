package controllers

import (
	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
	"github.com/PROger4ever/go-simple-news-aggregator/app/routes"
	"github.com/revel/revel"
	"gopkg.in/Masterminds/squirrel.v1"
)

type Sources struct {
	Application
}

func (c Sources) Index(size, page uint64) revel.Result {
	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	prevPage := page - 1
	nextPage := page + 1

	var sources []*models.Source
	builder := c.Db.SqlStatementBuilder.Select("*").From("Source").Offset((page - 1) * size).Limit(size)
	if _, err := c.Txn.Select(&sources, builder); err != nil {
		c.Log.Fatal("Unexpected error loading sources", "error", err)
	}
	return c.Render(sources, size, page, prevPage, nextPage)
}

func (c Sources) loadSourceById(id int) *models.Source {
	h, err := c.Txn.Get(models.Source{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return h.(*models.Source)
}

func (c Sources) loadSourceByUrl(url string) *models.Source {
	source := models.Source{}
	builder := c.Db.SqlStatementBuilder.Select("*").
		From("Source").
		Where(squirrel.Eq{"Url": url})
	if err := c.Txn.SelectOne(&source, builder); err != nil {
		c.Log.Fatal("Unexpected error loading a source", "error", err)
	}
	return &source
}

func (c Sources) New() revel.Result {
	return c.Render()
}

func (c Sources) Add(source *models.Source) revel.Result {
	existingSource := c.loadSourceByUrl(source.Url) //TODO: check name too?
	if existingSource != nil {
		c.Flash.Error("Source %s does already exist", source.Url)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Sources.New())
	}

	source.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Sources.New())
	}

	err := c.Txn.Insert(source)
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Added.")
	return c.Redirect(routes.Sources.Index(20, 1))
}

func (c Sources) Edit(id int) revel.Result {
	source := c.loadSourceById(id)
	if source == nil {
		return c.NotFound("Source %d does not exist", id)
	}
	title := source.Name
	return c.Render(title, source)
}

func (c Sources) SaveChanges(id int, source *models.Source) revel.Result {
	oldSource := c.loadSourceById(id)
	if oldSource != nil {
		source.SourceId = id
	}

	source.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Sources.Edit(source.SourceId))
	}

	_, err := c.Txn.Update(source)
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Saved.")
	return c.Redirect(routes.Sources.Index(20, 1))
}

//TODO: action for deleting sources
