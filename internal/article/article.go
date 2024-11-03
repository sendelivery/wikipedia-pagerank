package article

import (
	"errors"
	"fmt"

	"github.com/sendelivery/wikipedia-pagerank/internal/urls"
	"golang.org/x/net/html"
)

func GetWikipediaArticleLinks(article *html.Node) []string {
	links := make([]string, 0)

	/*
		<body>
			<div>
				<a>www.example.com/0</a>
			</div>
			<a>www.example.com/1</a>
			<p>text</p>
		</body>
	*/
	// { "www.example.com/0", "www.example.com/1" }

	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			link, err := getLinkFromATag(node)
			if err != nil {
				fmt.Printf("error when getting link: %e\n", err)
			} else if !urls.IsValidWikiPath(link) {
				fmt.Printf("not a valid wiki article url %s\n", link)
			} else {
				links = append(links, link)
			}
		} else {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				traverseNodes(child)
			}
		}
	}

	traverseNodes(article)
	return links
}

func getLinkFromATag(tag *html.Node) (string, error) {
	fmt.Println(tag.Attr)
	for i := 0; i < len(tag.Attr); i++ {
		attr := tag.Attr[i]
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("node did not contain an href")
}
