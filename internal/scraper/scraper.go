package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

const baseURL = "https://en.wikipedia.org"

const validArticlePathRegex = "^/wiki/[A-Za-z0-9_()%-]+$"

func IsValidWikiPath(url string) bool {
	re := regexp.MustCompile(validArticlePathRegex)
	return re.Match([]byte(url))
}

// Given a Wikipedia article identified by its path, e.g. the `/wiki/Go_(programming_language)`
// part of the `https://en.wikipedia.org/wiki/Go_(programming_language)` URL. Fetch that article's
// HTML and traverse its DOM tree to retrieve the paths of all linked Wikipedia articles as a slice
// of strings.
//
// An error is returned if any part of this process fails.
func ScrapeArticlesInWikipediaArticle(article string) ([]string, error) {
	articleHtml, err := getArticleHTML(article)
	if err != nil {
		// fmt.Println("error when fetching path:", article, "\nwith:", err)
		return nil, err
	}

	paths := extractWikipediaArticleLinks(articleHtml)
	// fmt.Printf("%s has %d wikipedia links\n", article, len(paths))
	return paths, nil
}

func getArticleHTML(path string) (*html.Node, error) {
	// fmt.Printf("Fetching link %s\n", path)
	resp, err := http.Get(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s with error %s", path, err.Error())
	}

	// fmt.Println("Parsing response body")
	articleHtml, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html for %s with error %e", path, err)
	}

	return articleHtml, nil
}

func extractWikipediaArticleLinks(article *html.Node) []string {
	links := make([]string, 0)

	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			link, err := getLinkFromATag(node)
			if err != nil {
				// fmt.Printf("error when getting link: %e\n", err)
			} else if !IsValidWikiPath(link) {
				// fmt.Printf("not a valid wiki article url %s\n", link)
			} else {
				links = append(links, link)
			}
		} else {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				traverseNodes(child)
			}
		}
	}

	// fmt.Println("Scraping wikipedia links from HTML")
	traverseNodes(article)
	return links
}

func getLinkFromATag(tag *html.Node) (string, error) {
	// fmt.Println(tag.Attr)
	for i := 0; i < len(tag.Attr); i++ {
		attr := tag.Attr[i]
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("node did not contain an href")
}