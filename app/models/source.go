package models

import (
	"fmt"

	"github.com/revel/revel"
)

type Source struct {
	SourceId int

	Name, Url        string
	ArticleUrlsXpath string

	CardXpath  string
	TitleXpath string
	BodyXpath  string
	ImgXpath   string

	PublishedXpath  string
	PublishedFormat string
}

func (source *Source) String() string {
	return fmt.Sprintf("Source(%s,%s)", source.Name, source.Url)
}

func (source *Source) Validate(v *revel.Validation) {
	//TODO: improve Url validation - use regex
	//TODO: add validation error messages
	v.Check(source.Name,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{50},
	).Message("Incorrect Name")
	v.Check(source.Url,
		revel.Required{},
		revel.MinSize{10},
		revel.MaxSize{200},
	).Message("Incorrect Url")
	v.Check(source.ArticleUrlsXpath,
		revel.MinSize{0},
		revel.MaxSize{1024},
	).Message("Incorrect ArticleUrlsXpath")

	v.Check(source.CardXpath,
		revel.MinSize{0},
		revel.MaxSize{1024},
	).Message("Incorrect CardXpath")
	v.Check(source.TitleXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	).Message("Incorrect TitleXpath")
	v.Check(source.BodyXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	).Message("Incorrect BodyXpath")
	v.Check(source.ImgXpath,
		revel.MinSize{1},
		revel.MaxSize{1024},
	).Message("Incorrect ImgXpath")
	v.Check(source.PublishedXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	).Message("Incorrect PublishedXpath")

	v.Check(source.PublishedFormat,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{100},
	).Message("Incorrect PublishedFormat")
}
