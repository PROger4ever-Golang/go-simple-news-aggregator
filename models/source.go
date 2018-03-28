package models

import "time"

type Source struct {
	Id                   int        `orm:"pk;auto"`
	Url                  string     `orm:"size(1024);unique"`
	ArticlePageUrlsXpath string     `orm:"null;size(1024)"`
	CardXpath            string     `orm:"null;size(1024)"`
	TitleXpath           string     `orm:"size(1024)"`
	BodyXpath            string     `orm:"size(1024)"`
	ImgXpath             string     `orm:"null;size(1024)"`
	PublicationTimeXpath string     `orm:"size(1024)"`
	Articles             []*Article `orm:"reverse(many)"`
	CreatedTime          time.Time  `orm:"type(timestamp)"`
	UpdatedTime          time.Time  `orm:"type(timestamp)"`
}
