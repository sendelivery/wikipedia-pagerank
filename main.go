package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
	"sync"
	"time"

	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
	"github.com/sendelivery/wikipedia-pagerank/internal/reporter"
	"github.com/sendelivery/wikipedia-pagerank/internal/scraper"
)

const NUM_CONCURRENT_FETCHES = 5
const NUM_PAGES = 50

// InputWikipediaArticlePath prompts the user to input a Wikipedia article path
// and returns that path.
func InputWikipediaArticlePath() string {
	fmt.Printf("Enter a Wikipedia URL: https://en.wikipedia.org")
	time.Sleep(500)
	path := "/wiki/Go_(programming_language)"
	for _, r := range path {
		fmt.Printf("%c", r)
		d := time.Duration(rand.IntN(20) + 20)
		time.Sleep(d * time.Millisecond)
	}
	fmt.Println()
	return path
}

func main() {
	parentArticlePath := InputWikipediaArticlePath()
	valid := scraper.IsValidWikiPath(parentArticlePath)
	if !valid {
		log.Fatalf("Invalid Wikipedia path: %s\nPlease try again.", parentArticlePath)
	}

	r := reporter.New()

	// Create a buffered channel that we'll use as a queue to hold all the pages we have yet to
	// fetch.
	queue := make(chan string, NUM_PAGES)
	queue <- parentArticlePath

	// `corp` will hold the corpus of wikipedia pages we're building.
	corp := corpus.New(NUM_PAGES)

	r.NewWorkInProgress("Building corpus")

	for corp.Size() < NUM_PAGES {
		if len(queue) == 0 {
			// fmt.Println("Queue is empty... exiting.")
			return
		}

		numStartGoroutines := min(len(queue), NUM_CONCURRENT_FETCHES, NUM_PAGES-corp.Size())

		var wg sync.WaitGroup
		wg.Add(numStartGoroutines)

		for i := numStartGoroutines; i > 0; i-- {
			go func() {
				defer wg.Done()

				path := <-queue
				articles, err := scraper.ScrapeArticlesInWikipediaArticle(path)
				if err != nil {
					// fmt.Println("error when fetching path:", path, "\nwith:", err)
					return
				}

				corp.Set(path, articles)

				for _, articlePath := range articles {
					if _, ok := corp.Get(articlePath); ok {
						continue
					}
					select {
					case queue <- articlePath:
					default:
						// Queue is full.
						goto end_goroutine
					}
				}
			end_goroutine:
			}()
		}

		wg.Wait()
	}

	r.Stop()

	// Must set cap of urls to the number of pages in the corpus or else the ForEach will overwrite
	// the pointer of the inner slice, and `urls` will be empty.
	urls := make([]string, 0, corp.Size())
	corp.ForEach(func(url string, _ []string) {
		urls = append(urls, url)
	})

	fmt.Println("------ Before enforcing consistency ------")
	printTopThreeArticles(&corp, urls)

	r.NewWorkInProgress("Enforcing corpus consistency")
	corp.EnsureConsistency()
	corp.CheckConsistency()
	r.Stop()

	fmt.Println("------ After enforcing consistency ------")
	printTopThreeArticles(&corp, urls)

	fmt.Println("size of corpus:", corp.Size())
}

func printTopThreeArticles(corp *corpus.Corpus, urls []string) {
	slices.SortFunc(urls, func(a, b string) int {
		aArr, _ := corp.Get(a)
		bArr, _ := corp.Get(b)

		return len(bArr) - len(aArr)
	})

	a, ok := corp.Get(urls[0])
	if !ok {
		fmt.Println("a", a, ok, urls[0])
		return
	}
	b, ok := corp.Get(urls[1])
	if !ok {
		fmt.Println("b", b, ok, urls[1])
		return
	}
	c, ok := corp.Get(urls[3])
	if !ok {
		fmt.Println("c", c, ok, urls[2])
		return
	}

	fmt.Printf("Scraped %d articles.\n", corp.Size())

	fmt.Println("Total links in corpus:", corp.GetTotalLinks())
	fmt.Println("Top three articles:")
	fmt.Printf("1. %s with %d links, last link: %s\n", urls[0], len(a), a[len(a)-1])
	fmt.Printf("2. %s with %d links, last link: %s\n", urls[1], len(b), b[len(b)-1])
	fmt.Printf("3. %s with %d links, last link: %s\n", urls[2], len(c), c[len(c)-1])
}
