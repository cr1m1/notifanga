package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Crawler struct {
	mangas []*Manga
}

func NewCrawler(marr []*Manga) *Crawler {
	return &Crawler{
		mangas: marr,
	}
}

func (c *Crawler) Crawl(service *NotifangaService) []*User {
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Println(err)
	// }
	// dbconn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Println("cannot get connection with db", err)
	// }
	// defer dbconn.Close()
	// repo, err := NewRepository(dbconn)
	// if err != nil {
	// 	log.Println("cannot create repository", err)
	// }
	// service := NewNotifangaService(repo)

	collector := colly.NewCollector(
		colly.AllowedDomains("mangalib.me"),
	)
	rssCollector := collector.Clone()

	var uarr []*User

	// go func() {
	marr, err := service.GetAllMangas()
	if err != nil {
		log.Println("cant get all mangas", err)
	}

	for i, m := range marr {
		collector.OnHTML("div.media-sidebar-actions a", func(el *colly.HTMLElement) {
			rssLink := el.Attr("href")

			fmt.Println(rssLink)

			if strings.Contains(rssLink, "mangalib.me/manga-rss") {
				rssCollector.Visit(rssLink)
			}

			rssCollector.OnXML("rss/channel/item[1]", func(el *colly.XMLElement) {
				newChapter := el.Attr("title")
				newChapterLink := el.Attr("link")
				if newChapter != m.LastChapter {
					marr[i].LastChapter = newChapter
					marr[i].LastChapterUrl = newChapterLink
					service.UpdateManga(*marr[i])

					uarr, err = service.ListMangaUsers(*marr[i])
					if err != nil {
						log.Println("cant use ListMangaUsers", err)
					}
				}
			})

			collector.Visit(m.Url)
		})
	}
	// }()

	return uarr
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
