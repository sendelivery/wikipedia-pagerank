package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sendelivery/wikipedia-pagerank/internal/article"
	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
	"github.com/sendelivery/wikipedia-pagerank/internal/urls"
	"golang.org/x/net/html"
)

const NUM_CONCURRENT_FETCHES = 50
const NUM_PAGES = 10_000

var fetched map[string]bool = make(map[string]bool)

func main() {
	parentArticlePath := urls.InputWikipediaArticlePath()
	valid := urls.IsValidWikiPath(parentArticlePath)
	if !valid {
		log.Fatalf("Invalid Wikipedia path: %s\nPlease try again.", parentArticlePath)
	}

	// Create an unbuffered channel that we'll use as a queue to hold all the pages we have yet to
	// fetch.
	// queue := make(chan string)

	// Create a buffered channel that will hold all of the URLs we are concurrently fetching at any
	// given time while building our corpus.
	fetching := make(chan string, NUM_CONCURRENT_FETCHES)
	fetching <- parentArticlePath

	// `corp` will hold the corpus of wikipedia pages we're building.
	corp := corpus.MakeCorpus(NUM_PAGES)

	// We'll use this to stop the main goroutine until our corpus has been completed.
	// Todo: use context instead?
	// var wg sync.WaitGroup
	// wg.Add(NUM_PAGES)

	for link := range fetching {
		fmt.Printf("Fetching link %s\n", link)
		resp, err := http.Get(urls.BaseURL + link)
		if err != nil {
			log.Fatalf("failed to fetch %s with error %e", link, err)
		}

		fmt.Println("Parsing response body")
		articleHtml, err := html.Parse(resp.Body)
		if err != nil {
			log.Fatalf("failed to parse html for %s with error %e", link, err)
		}

		fmt.Println("Getting wiki links from article")
		links := article.GetWikipediaArticleLinks(articleHtml)
		corp[link] = links

		// Close fetching channel so we can exit the loop
		close(fetching)
	}

	fmt.Println(corp)
	fmt.Println(len(corp["/wiki/Go_(programming_language)"]))
}
