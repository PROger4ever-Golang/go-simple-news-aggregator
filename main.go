package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PROger4ever/go-simple-news-aggregator/common"
	"github.com/PROger4ever/go-simple-news-aggregator/config"
	"github.com/PROger4ever/go-simple-news-aggregator/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/xmlpath.v2"
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
	//resp, err := http.Get("https://lenta.ru/parts/news")
	resp, err := http.Get("https://lenta.ru/rss")
	common.PanicIfError(err, "doing a request")
	defer resp.Body.Close()

	var xmlroot *xmlpath.Node
	contentTypeHeader := resp.Header.Get("Content-Type")
	if strings.Contains(contentTypeHeader, "html") {
		xmlroot, err = xmlpath.ParseHTML(resp.Body)
	} else if strings.Contains(contentTypeHeader, "rss") {
		xmlroot, err = xmlpath.Parse(resp.Body)
	} else {
		panic("Unsupported page content type")
	}
	common.PanicIfError(err, "parsing page content")

	xpath := `//title`
	path := xmlpath.MustCompile(xpath)
	iter := path.Iter(xmlroot)
	for iter.Next() {
		fmt.Printf("Found: %s\n", iter.Node().String())
	}
	fmt.Println("done.")
}
