package cmd

import (
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
)

func Scrape(url string) {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		bodyStr := string(r.Body)
		if r.StatusCode >= 200 && r.StatusCode <= 299 {
			markdown, err := htmltomarkdown.ConvertString(bodyStr)
			if err != nil {
				fmt.Printf("An error occurred while converting content to markdown: %s", err.Error())
				return
			}
			fmt.Println(markdown)
		} else {
			fmt.Printf("The page you attempted to scrape returned: %d. Response body: %s\n", r.StatusCode, bodyStr)
		}
	})

	c.Visit(url)
}

func ScraperImpl(url string) (scraped string, err error) {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		bodyStr := string(r.Body)
		if r.StatusCode >= 200 && r.StatusCode <= 299 {
			markdown, errConv := htmltomarkdown.ConvertString(bodyStr)
			if errConv != nil {
				err = errConv
				return
			}
			scraped = markdown
		} else {
			fmt.Printf("The page you attempted to scrape returned: %d. Response body: %s\n", r.StatusCode, bodyStr)
		}
	})

	c.Visit(url)
	return
}
