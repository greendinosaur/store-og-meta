package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// OGMetaData stores the key og meta-data from a site
type ogMetaData struct {
	site          string
	ogImage       string
	ogDescription string
	ogWidth       int
	ogHeight      int
}

var (
	repo Repository
)

func findMetaTags(site string) ogMetaData {
	c := colly.NewCollector()
	var siteMetaData ogMetaData
	siteMetaData.site = site
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		metaProperty := strings.ToLower(e.Attr("property"))
		metaContent := e.Attr("content")

		/*if strings.Contains(metaProperty, "og:image") {
			fmt.Printf("found image: %s", metaContent)
			siteMetaData.ogImage = metaContent
		}
		if strings.Contains(metaProperty, "og:description") {
			siteMetaData.ogDescription = metaContent
		}
		if strings.Contains(metaProperty, "og:image:width") {
			siteMetaData.ogWidth, _ = strconv.Atoi(metaContent)
		}
		if strings.Contains(metaProperty, "og:image:height") {
			siteMetaData.ogHeight, _ = strconv.Atoi(metaContent)
		}*/
		switch metaProperty {
		case "og:image":
			siteMetaData.ogImage = metaContent
		case "og:description":
			siteMetaData.ogDescription = metaContent
		case "og:image:width":
			siteMetaData.ogWidth, _ = strconv.Atoi(metaContent)
		case "og:image:height":
			siteMetaData.ogHeight, _ = strconv.Atoi(metaContent)

		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(site)

	return siteMetaData
}

func updateDatabase(meta ogMetaData) {
	var err error
	repo, err = NewRepository("mysql", os.Getenv("sql_dns"))
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err.Error())
	}

	defer func() {
		repo.Close()
	}()

	repo.Update(meta)
}

func main() {
	sitePtr := flag.String("site", "", "Site whose og meta tags are to be retrieved")
	flag.Parse()
	if len(*sitePtr) > 0 {
		siteOGMetaData := findMetaTags(*sitePtr)
		fmt.Println(siteOGMetaData)
		if os.Getenv("save") == "true" {
			updateDatabase(siteOGMetaData)
		}
	} else {
		fmt.Println("no site URL provided")
	}
}
