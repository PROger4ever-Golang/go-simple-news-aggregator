package models

import (
	"fmt"
	"github.com/revel/revel"
	"gopkg.in/gorp.v2"
	"time"
)

type Article struct {
	ArticleId int
	SourceId  int

	Url    string
	Title  string
	Body   string
	ImgUrl string

	PublishedStr string

	// Transient
	Published time.Time

	Source *Source
}

func (a Article) Validate(v *revel.Validation) {
	v.Required(a.Source)

	v.Check(a.Url,
		revel.Required{},
		revel.MinSize{7},
		revel.MaxSize{200},
	)
	v.Check(a.Title,
		revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{200},
	)
	v.Check(a.Body,
		revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{65535},
	)
	v.Check(a.ImgUrl,
		revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{65535},
	)
	v.Required(a.Published)
}

const SQL_DATE_FORMAT = time.RFC3339

func (a Article) String() string {
	return fmt.Sprintf("Article(%s,%s)", a.Source, a.Title)
}

// These hooks work around two things:
// - Gorp's lack of support for loading relations automatically.
// - Sqlite's lack of support for datetimes.

func (a *Article) PreInsert(_ gorp.SqlExecutor) error {
	a.SourceId = a.Source.SourceId
	a.PublishedStr = a.Published.Format(SQL_DATE_FORMAT)
	return nil
}

func (a *Article) PostGet(exe gorp.SqlExecutor) error {
	var (
		obj interface{}
		err error
	)

	fmt.Printf("Get post %#v\n", a)

	obj, err = exe.Get(Source{}, a.SourceId)
	if err != nil {
		return fmt.Errorf("error loading an article's source (%d): %s", a.SourceId, err)
	}
	a.Source = obj.(*Source)

	if a.Published, err = time.Parse(SQL_DATE_FORMAT, a.PublishedStr); err != nil {
		return fmt.Errorf("error parsing publication in date '%s': %s", a.PublishedStr, err)
	}
	return nil
}
