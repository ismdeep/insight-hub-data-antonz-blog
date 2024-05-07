package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/insight-hub-data-antonz-blog/internal/schema"
	"github.com/ismdeep/insight-hub-data-antonz-blog/pkg/insight-hub-data/core"
	"golang.org/x/net/html"
	"os"
	"strings"
	"time"
)

// Antonz model
type Antonz struct {
}

// GetBloggerName get blogger name
func (receiver *Antonz) GetBloggerName() string {
	return "antonz.org"
}

// GetAllPageURLs get all page urls
func (receiver *Antonz) GetAllPageURLs() []string {
	return []string{receiver.GetFirstPageURL()}
}

// GetFirstPageURL get first page url
func (receiver *Antonz) GetFirstPageURL() string {
	return "https://antonz.org/all/"
}

// GetLinksFromPage get links from page
func (receiver *Antonz) GetLinksFromPage(pageURL string) ([]string, error) {
	return GetBlogLinks(pageURL,
		`//div[@class="posts"]//div[@class="post-stub"]//a[@class="post-stub__title"]/@href`,
		"https://antonz.org")
}

// GetBlogInfo get blog info
func (receiver *Antonz) GetBlogInfo(blogLink string) (*schema.Blog, error) {
	return GetBlogInfo(blogLink, receiver.GetBloggerName(),
		`//meta[@property="og:title"]/@content`,
		`//meta[@name="author"]/@content`,
		`//article[@class="post"]`,
		func(doc *html.Node) (time.Time, error) {
			s := htmlquery.FindOne(doc, `//footer[@class="post__footer"]/div[@class="row"]//div[@class="post__date"]/time/@datetime`)
			if s == nil {
				return time.Now(), nil
			}
			ss := htmlquery.InnerText(s)
			ss = strings.TrimSpace(ss)
			t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ss)
			if err == nil {
				return t, nil
			}

			t2, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", ss)
			if err == nil {
				return t2, nil
			}

			return time.Now(), errors.New("failed to parse time")
		},
	)
}

func main() {

	f, err := os.OpenFile("data.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}

	s := core.NewStore(f)

	raw, err := os.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}
	if err := s.Load(bytes.NewBuffer(raw)); err != nil {
		panic(err)
	}

	a := Antonz{}
	pageURLs := a.GetAllPageURLs()
	for _, l := range pageURLs {
		links, err := a.GetLinksFromPage(l)
		if err != nil {
			fmt.Printf("failed to get links from page: %v\n", err.Error())
			continue
		}
		for _, link := range links {
			if !core.LinkIsTidy(link) {
				fmt.Printf("link is not tidy: %v\n", link)
				continue
			}

			if s.URLExists(link) {
				continue
			}

			info, err := a.GetBlogInfo(link)
			if err != nil {
				fmt.Printf("failed to get blog info: %v\n", err.Error())
				continue
			}
			if err := s.Save(core.Record{
				Source:      info.Source,
				Link:        info.Link,
				Title:       info.Title,
				Author:      info.Author,
				Content:     info.Content,
				PublishedAt: info.Date,
			}); err != nil {
				fmt.Printf("failed to save record: %v\n", err.Error())
				continue
			}

			fmt.Printf("OK: %v\n", link)
		}
	}
}
