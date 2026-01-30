package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
)

func ScraperImpl(url, output, logLevel string, parallel, maxDepth int, recursive bool, allowedDomains []string) ([]string, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: stringLevelToSlogLevel(logLevel)}))
	paths := newOutputFiles()
	domains := validateAllowedDomains(allowedDomains)

	if _, err := os.Stat(output); os.IsNotExist(err) {
		err := os.Mkdir(output, 0755)
		if err != nil {
			return nil, err
		}
	}

	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(),
	)

	if len(domains) > 0 {
		c.AllowedDomains = domains
	}

	c.OnRequest(func(r *colly.Request) {
		logger.Debug("Sending new scraping request", "url", r.URL.String())
	})

	if recursive {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			logger.Debug("Found link in the current page", "link", link)
			e.Request.Visit(link)
		})
	}

	c.OnResponse(func(r *colly.Response) {
		bodyStr := string(r.Body)
		flnm := r.FileName()
		logger.Info("Obtained raw file name", "file_name", flnm)
		if flnm == "" || strings.HasPrefix(flnm, ".") {
			now := time.Now().UnixMilli()
			flnm = strconv.Itoa(int(now)) + ".md"
		}
		if !strings.HasSuffix(flnm, ".md") {
			ext := filepath.Ext(flnm)
			flnm = strings.TrimSuffix(flnm, ext) + ".md"
		}
		path := strings.TrimRight(output, "/") + "/" + flnm
		if r.StatusCode >= 200 && r.StatusCode <= 299 {
			markdown, err := htmltomarkdown.ConvertString(bodyStr)
			if err != nil {
				logger.Info("An error occurred while converting HTML content to markdown", "url", r.Request.URL.String(), "error", err.Error())
				return
			}
			err = os.WriteFile(path, []byte(markdown), 0644)
			if err != nil {
				logger.Info("An error occurred while writing the output file", "url", r.Request.URL.String(), "error", err.Error(), "output_file", path)
				return
			}
			paths.append(path)
		} else {
			logger.Info("The page you attempted to scrape returned a non-OK response", "url", r.Request.URL.String(), "status_code", r.StatusCode, "response_body", bodyStr)
		}
	})

	c.Visit(url)
	c.Wait()
	return paths.getFiles(), nil
}
