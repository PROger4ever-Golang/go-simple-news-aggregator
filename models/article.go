package models

import "time"

type Article struct {
	Id              int       `orm:"pk;auto"`
	Source          *Source   `orm:"rel(fk)"`
	Url             string    `orm:"size(1024)"`
	Title           string    `orm:"size(1024)"`
	Body            string    `orm:"size(1024)"`
	ImgUrl          string    `orm:"null;size(1024)"`
	PublicationTime time.Time `orm:"type(datetime)"`
	CreatedTime     time.Time `orm:"type(timestamp)"`
	UpdatedTime     time.Time `orm:"type(timestamp)"`
}

// multiple fields unique key
func (u *Article) TableUnique() [][]string {
	return [][]string{
		{"Source", "PublicationTime", "Title"},
	}
}
