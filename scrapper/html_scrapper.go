package scrapper

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"images_urls"`
}

func ExtractPageData(html, pageURL string) PageData {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		log.Printf("extractPageData error while parsing URL: %s ", err)
		return PageData{}
	}

	baseURL := &url.URL{Scheme: parsedURL.Scheme, Host: parsedURL.Host}
	links, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		log.Printf("extractPageData erro while getting page links: %s ", err)
		return PageData{}
	}

	images, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		log.Printf("extractPageData erro while getting page images: %s ", err)
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		Heading:        getHeadingHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  links,
		ImageURLs:      images,
	}
}

func getHeadingHTML(html string) string {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("erro while parsing document at 'getHeadingHTML': %v", err)
		return ""
	}

	h1 := document.Find("h1")
	if h1 != nil && h1.Text() != "" {
		return h1.Text()
	}

	h2 := document.Find("h2")
	if h2 != nil && h2.Text() != "" {
		return h2.Text()
	}

	return ""
}

func getFirstParagraphFromHTML(html string) string {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("erro while parsing document at 'getFirstParagraphFromHTML': %v", err)
		return ""
	}

	mainP := document.Find("main").Find("p").First()
	if mainP != nil && mainP.Text() != "" {
		return mainP.Text()
	}

	firstP := document.Find("p").First()
	if firstP != nil && firstP.Text() != "" {
		return firstP.Text()
	}

	return ""
}

func getURLsFromHTML(html string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	aElements := doc.Find("a")
	if aElements.Size() == 0 {
		return nil, nil
	}

	links := make([]string, 0, aElements.Size())
	aElements.Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if href == "" {
			return
		}

		abs, err := absoluteULR(baseURL, href)
		if err != nil {
			log.Printf("error while parsing href url: %s", err)
			return
		}

		links = append(links, abs)
	})

	return links, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	images := doc.Find("img")
	if images.Size() == 0 {
		return nil, nil
	}

	links := make([]string, 0, images.Size())
	images.Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok || src == "" {
			return
		}

		link, err := absoluteULR(baseURL, src)
		if err != nil {
			log.Printf("getImagesFromHTML FAIL to parse src URL: %s", err)
			return
		}

		links = append(links, link)
	})

	return links, nil
}

func absoluteULR(baseURL *url.URL, raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	url := parsed.String()
	if !parsed.IsAbs() {
		url = baseURL.ResolveReference(parsed).String()
	}

	return url, nil
}
