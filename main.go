package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/alexis-wizeline/boot_crawler/report"
	"github.com/alexis-wizeline/boot_crawler/scrapper"
)

type config struct {
	pages              map[string]scrapper.PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func (c *config) crawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
	}()

	log.Printf("starting crawl for link: %s", rawCurrentURL)
	c.mu.Lock()
	currentPages := len(c.pages)
	c.mu.Unlock()
	if currentPages >= c.maxPages {
		log.Println("max number of pages reached")
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("unable to get parse current url: %s\n", err)
		return
	}

	if c.baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalized, err := scrapper.NormalizeURL(currentURL.String())
	if err != nil {
		log.Printf("Unable to normalize url: %s \n", rawCurrentURL)
		return
	}
	isFirst := c.addPageVisit(normalized)
	if !isFirst {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("Unable to collect html for url: %s\n", rawCurrentURL)
		return
	}

	pageData := scrapper.ExtractPageData(html, rawCurrentURL)
	c.mu.Lock()
	c.pages[normalized] = pageData
	c.mu.Unlock()

	log.Printf("Craw for link: %s, Done\n", rawCurrentURL)
	for _, link := range pageData.OutgoingLinks {
		c.wg.Go(func() {
			c.crawlPage(link)
		})
	}
}

func (c *config) addPageVisit(normalizedURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.pages[normalizedURL]
	if !ok {
		c.pages[normalizedURL] = scrapper.PageData{}
	}

	return !ok
}

func (c *config) writeReport() error {
	return report.WriteJsonReport(c.pages, "report.json")
}

func main() {
	args := os.Args[1:]
	writter := os.Stdout
	defer writter.Close()
	if len(args) < 1 {
		writter.WriteString("no website provided")
		os.Exit(1)
	}

	if len(args) > 3 {
		writter.WriteString("too many arguments provided")
		os.Exit(1)
	}

	maxConcurrency := 25
	maxPages := math.MaxInt

	if len(args) > 1 {
		paramsMaxConcurrency, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(writter, "unable to convert max cocunrrency param in position 1 to int: %s", err)
			os.Exit(1)
		}
		maxConcurrency = paramsMaxConcurrency

		paramsMaxPages, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(writter, "unable to convert max pages param in position 2 to int: %s", err)
			os.Exit(1)
		}
		maxPages = paramsMaxPages
	}

	baseURL := args[0]
	fmt.Fprintf(writter, "starting crawl of: %s\n", baseURL)
	parsed, err := url.Parse(baseURL)
	if err != nil {
		fmt.Fprintf(writter, "unable to parse base url:%s \n", baseURL)
		os.Exit(1)
	}
	pages := map[string]scrapper.PageData{}
	config := config{
		baseURL:            parsed,
		pages:              pages,
		concurrencyControl: make(chan struct{}, maxConcurrency),
		maxPages:           maxPages,
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
	}

	config.wg.Go(func() {
		config.crawlPage(baseURL)
	})

	config.wg.Wait()
	close(config.concurrencyControl)
	if err := config.writeReport(); err != nil {
		log.Fatalf("error writting the report: %s", err)
	}
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Bad Status Code: %v", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("invalid content response: %s", res.Header.Get("content-type"))
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}
