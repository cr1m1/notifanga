package main

type Crawler struct {
	users  []*User
	mangas []*Manga
}

func NewCrawler(uarr []*User, marr []*Manga) *Crawler {
	return &Crawler{
		users:  uarr,
		mangas: marr,
	}
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
