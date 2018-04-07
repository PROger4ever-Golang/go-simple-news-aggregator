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
	)
	v.Check(source.Url,
		revel.Required{},
		revel.MinSize{10},
		revel.MaxSize{200},
	)
	v.Check(source.ArticleUrlsXpath,
		revel.MinSize{1},
		revel.MaxSize{1024},
	)

	v.Check(source.CardXpath,
		revel.MinSize{1},
		revel.MaxSize{1024},
	)
	v.Check(source.TitleXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	)
	v.Check(source.BodyXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	)
	v.Check(source.ImgXpath,
		revel.MinSize{1},
		revel.MaxSize{1024},
	)
	v.Check(source.PublishedXpath,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{1024},
	)

	v.Check(source.PublishedFormat,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{100},
	)
}
