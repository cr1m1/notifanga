package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func Crawl(m Manga, service *NotifangaService) []*User {
	log.Println("crawler started")
	collector := colly.NewCollector(
		colly.AllowedDomains("mangalib.me"),
	)
	rssCollector := collector.Clone()

	var uarr []*User
	var err error
	isReleased := false

	// marr, err := service.GetAllMangas()
	// if err != nil {
	// 	log.Println("cant get all mangas", err)
	// }

	collector.OnHTML("div.media-sidebar-actions a", func(el *colly.HTMLElement) {
		log.Println("onHTML started")
		rssLink := el.Attr("href")

		fmt.Println(rssLink)

		if strings.Contains(rssLink, "mangalib.me/manga-rss") {
			rssCollector.Visit(rssLink)
		}

		rssCollector.OnXML("rss/channel/item[1]", func(el *colly.XMLElement) {
			log.Println("onXML started")
			newChapter := el.Attr("title")
			newChapterLink := el.Attr("link")
			if newChapter != m.LastChapter {
				log.Println("new chapter")
				isReleased = true
				m.LastChapter = newChapter
				m.LastChapterUrl = newChapterLink
				service.UpdateManga(m)

				uarr, err = service.ListMangaUsers(m)
				if err != nil {
					log.Println("cant use ListMangaUsers", err)
				}
			}
		})

		collector.Visit(m.Url)
	})
	if isReleased {
		return uarr
	}
	return nil
}

func CrawlName(url string) string {
	collector := colly.NewCollector(
		colly.AllowedDomains("mangalib.me"),
	)
	name := ""
	collector.OnHTML("div.media-name__main", func(el *colly.HTMLElement) {
		name = el.Text
	})
	collector.Visit(url)

	return name
}

// func (c *Crawler) Crawl() {
// 	collector := colly.NewCollector(
// 		colly.AllowedDomains("mangalib.me"),
// 	)

// 	rssCollector := collector.Clone()
// }

// func (crawler *Crawler) Crawl() {
// 	c := colly.NewCollector(
// 		colly.AllowedDomains("mangalib.me"),
// 	)

// 	rssCollector := c.Clone()

// 	c.OnHTML("div.media-sidebar-actions a", func(el *colly.HTMLElement) {
// 		rssLink := el.Attr("href")

// 		fmt.Println(rssLink)

// 		if strings.Contains(rssLink, "mangalib.me/manga-rss") {
// 			rssCollector.Visit(rssLink)
// 		}
// 	})

// 	rssCollector.OnXML("rss/channel/item[1]/title", func(el *colly.XMLElement) {
// 		newChapter := el.Text
// 		if newChapter != crawler.LastChapter {
// 			crawler.LastChapter = newChapter
// 			fmt.Println("New chapter released! " + crawler.LastChapter + "\n")
// 		}
// 	})
// }

// func NewCrawler(url string) (*Crawler, error) {
// 	crawler := &Crawler{
// 		MangaUrl:    url,
// 		LastChapter: "",
// 	}

// 	if err := crawler.Validate(); err != nil {
// 		return nil, err
// 	}

// 	return crawler, nil
// }

// func (crawler *Crawler) Validate() error {
// 	if !strings.Contains(crawler.MangaUrl, "mangalib.me/") {
// 		return ErrNotValidUrl
// 	}
// 	chapterCheck, _ := regexp.MatchString(`[a-z]`, crawler.LastChapter)
// 	if chapterCheck {
// 		return ErrNotValidChapter
// 	}
// 	return nil
// }
