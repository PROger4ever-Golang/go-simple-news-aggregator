package main

import (
	"github.com/PROger4ever/go-simple-news-aggregator/config"
	"github.com/PROger4ever/go-simple-news-aggregator/common"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/PROger4ever/go-simple-news-aggregator/models"
	"fmt"
	"time"
)

func init() {
	conf, err := config.LoadConfig("config.json")
	common.PanicIfError(err, "parsing config")

	// register model
	orm.RegisterModel(new(models.Article))
	orm.RegisterModel(new(models.Source))

	// set default database
	orm.RegisterDataBase("default", "mysql", conf.MySql.Url, 30)
}

func main() {
	o := orm.NewOrm()

	source := models.Source{
		Url:                  "http://xxx.yyy.zzz/aaaa/bbbb.cccc",
		TitleXpath:           "TitleXpath",
		BodyXpath:            "BodyXpath",
		PublicationTimeXpath: "PublicationTimeXpath",
	}
	article := models.Article{
		Title:           "Title",
		Body:            "Body",
		Source:          &source,
		PublicationTime: time.Date(2018, 03, 28, 23, 36, 0, 0, time.UTC),
	}

	// insert
	id, err := o.InsertOrUpdate(&source)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	id, err = o.Insert(&article)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

}
