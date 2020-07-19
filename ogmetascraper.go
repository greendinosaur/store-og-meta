package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// OGMetaData stores the key og meta-data from a site
type ogMetaData struct {
	site          string
	ogImage       string
	ogDescription string
	ogWidth       int
	ogHeight      int
}

func findMetaTags(site string) ogMetaData {
	c := colly.NewCollector()
	var siteMetaData ogMetaData
	siteMetaData.site = site
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		metaProperty := strings.ToLower(e.Attr("property"))
		metaContent := e.Attr("content")

		if strings.Contains(metaProperty, "og:image") {
			siteMetaData.ogImage = metaContent
		}
		if strings.Contains(metaProperty, "og:description") {
			siteMetaData.ogDescription = metaContent
		}
		if strings.Contains(metaProperty, "og:width") {
			siteMetaData.ogWidth, _ = strconv.Atoi(metaContent)
		}
		if strings.Contains(metaProperty, "og:height") {
			siteMetaData.ogHeight, _ = strconv.Atoi(metaContent)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(site)

	return siteMetaData
}

func saveMetaTagsToDB(siteOGMetaData ogMetaData) {

	db, err := sql.Open("mysql", os.Getenv("ogmeta_db"))
	errorCheck(err)
	tx, err := db.Begin()
	errorCheck(err)
	defer tx.Rollback()

	stmtStr := "UPDATE Inspiration SET meta_og_image=?, meta_og_description=?, meta_og_height=?, meta_og_width=? WHERE url=?"
	stmt, err := tx.Prepare(stmtStr)
	errorCheck(err)

	defer stmt.Close()

	_, err = stmt.Exec(siteOGMetaData.ogImage, siteOGMetaData.ogDescription, siteOGMetaData.site, siteOGMetaData.ogHeight, siteOGMetaData.ogWidth)
	errorCheck(err)

	err = tx.Commit()
	errorCheck(err)

	defer db.Close()
}

func errorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	sitePtr := flag.String("site", "", "Site whose og meta tags are to be retrieved")
	flag.Parse()
	if len(*sitePtr) > 0 {
		siteOGMetaData := findMetaTags(*sitePtr)
		// now want to save these tags to a database
		fmt.Println(siteOGMetaData)
	} else {
		fmt.Println("no site URL provided")
	}
}
