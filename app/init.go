package app

import (
	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v2"
	"time"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}
	revel.OnAppStart(func() {
		Dbm := rgorp.Db.Map
		setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
			for col, size := range colSizes {
				t.ColMap(col).MaxSize = size
			}
		}
		setColumnUniques := func(t *gorp.TableMap, colSizes map[string]bool) {
			for col, isUnique := range colSizes {
				t.ColMap(col).Unique = isUnique
			}
		}
		setColumnNotNull := func(t *gorp.TableMap, colSizes map[string]bool) {
			for col, isNotNull := range colSizes {
				t.ColMap(col).SetNotNull(isNotNull)
			}
		}

		t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
		t.ColMap("Password").Transient = true
		setColumnSizes(t, map[string]int{
			"Username": 20,
			"Name":     100,
		})

		t = Dbm.AddTable(models.Source{}).SetKeys(true, "SourceId")
		setColumnUniques(t, map[string]bool{
			"Url": true,
		})
		setColumnSizes(t, map[string]int{
			"Name":             50,
			"Url":              200,
			"ArticleUrlsXpath": 1024,
			"CardXpath":        1024,
			"TitleXpath":       1024,
			"BodyXpath":        1024,
			"ImgXpath":         1024,
			"PublishedXpath":   1024,
		})
		setColumnNotNull(t, map[string]bool{
			"Url":              true,
			"ArticleUrlsXpath": false,
			"CardXpath":        false,
			"TitleXpath":       true,
			"BodyXpath":        true,
			"ImgXpath":         false,
			"PublishedXpath":   true,
		})

		t = Dbm.AddTable(models.Article{}).SetKeys(true, "ArticleId")
		t.ColMap("Source").Transient = true
		t.ColMap("Published").Transient = true
		t.SetUniqueTogether("SourceId", "PublishedStr", "Title")
		setColumnSizes(t, map[string]int{
			"Url":    200,
			"Title":  200,
			"ImgUrl": 200,
		})
		setColumnNotNull(t, map[string]bool{
			"Title":        true,
			"Body":         true,
			"PublishedStr": true,
		})

		rgorp.Db.TraceOn(revel.AppLog)
		Dbm.CreateTables()

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("demo"), bcrypt.DefaultCost)
		demoUser := &models.User{
			Name:           "Demo User",
			Username:       "demo",
			Password:       "demo",
			HashedPassword: bcryptPassword,
		}
		if err := Dbm.Insert(demoUser); err != nil {
			panic(err)
		}

		//TODO: research the standardized articles scheme: http://schema.org/Article
		sources := []*models.Source{
			{
				Name: "Lenta.RU/RSS",
				Url:  "https://lenta.ru/rss",
				//ArticleUrlsXpath: "",
				CardXpath:       `/rss/channel/item`,
				TitleXpath:      `./title`,
				BodyXpath:       `./description`,
				ImgXpath:        `./enclosure/@url`,
				PublishedXpath:  `./pubDate`,
				PublishedFormat: time.RFC1123Z,
			},
			{
				Name:             "ria.ru -> world",
				Url:              "https://ria.ru/world/",
				ArticleUrlsXpath: `//*[contains(@class,"b-list__item")]/a/@href`,
				CardXpath:        `//*[@itemtype="http://schema.org/Article"]`,
				TitleXpath:       `.//*[@itemprop="name"]/span`,
				BodyXpath:        `.//*[@itemprop="articleBody"]`,
				ImgXpath:         `.//img[@itemprop="associatedMedia"]/src`,
				PublishedXpath:   `.//*[@itemprop="dateCreated"]/@datetime`,
				PublishedFormat:  "2006-01-02T15:04",
			},
		}
		for _, source := range sources {
			if err := Dbm.Insert(source); err != nil {
				panic(err)
			}
		}
		now := time.Now()
		articles := []*models.Article{
			{
				SourceId: sources[0].SourceId,

				Url:    "https://lenta.ru/news/2018/04/06/movement/",
				Title:  "It's a title",
				Body:   "It's a body",
				ImgUrl: "//palacesquare.rambler.ru/ocwzjosl/MWF4eWE3LndhaHV1QHsiZGF0YSI6ey/JBY3Rpb24iOiJQcm94eSIsIlJlZmZl/cmVyIjoiaHR0cHM6Ly9sZW50YS5ydS/9uZXdzLzIwMTgvMDQvMDYvbW92ZW1l/bnQvIiwiUHJvdG9jb2wiOiJodHRwcz/oiLCJIb3N0IjoibGVudGEucnUifSwi/bGluayI6Imh0dHBzOi8vaWNkbi5sZW/50YS5ydS9pbWFnZXMvMjAxOC8wNC8w/Ni8xMS8yMDE4MDQwNjExNDQzNjgzNi/9waWNfODY4ZTQ2NTkwYmZkZTJlOTQx/MmViZjI5ZTg3YWVhMWQuanBnIn0%3D/",

				Published: now,
				Source:    sources[0],
			},
		}
		for _, article := range articles {
			if err := Dbm.Insert(article); err != nil {
				panic(err)
			}
		}
	}, 5)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
