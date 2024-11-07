package main

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/sendelivery/wikipedia-pagerank/internal/article"
	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
	"github.com/sendelivery/wikipedia-pagerank/internal/urls"
)

const NUM_CONCURRENT_FETCHES = 50
const NUM_PAGES = 101

var fetched sync.Map

func main() {
	parentArticlePath := urls.InputWikipediaArticlePath()
	valid := urls.IsValidWikiPath(parentArticlePath)
	if !valid {
		log.Fatalf("Invalid Wikipedia path: %s\nPlease try again.", parentArticlePath)
	}

	// Create an huge buffered channel that we'll use as a queue to hold all the pages we have yet
	// to fetch.
	queue := make(chan string, NUM_PAGES)
	queue <- parentArticlePath

	// `corp` will hold the corpus of wikipedia pages we're building.
	corp := corpus.MakeCorpus(NUM_PAGES)

	// Temporary: close queue after 10 seconds
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	close(queue)
	// }()

	// Todo: make this loop go until the corpus is NUM_PAGES long
	for i := 0; i < NUM_PAGES/NUM_CONCURRENT_FETCHES; i++ {
		// We'll use this to stop our loop until all current 50 fetches and parses have completed.
		var wg sync.WaitGroup
		wg.Add(NUM_CONCURRENT_FETCHES)

		for i := NUM_CONCURRENT_FETCHES; i > 0; i-- {
			var path string

			select {
			case path = <-queue: // acquire a path
			default:
				fmt.Println("queue is empty")
				wg.Done()
				continue
			}

			if _, ok := fetched.Load(path); ok {
				wg.Done()
				continue
			}

			// create a goroutine to fetch, parse, and extract wikipedia links from that path
			go func() {
				defer wg.Done()

				articleHtml, err := urls.GetArticleHTML(path)
				if err != nil {
					fmt.Println(err)
					return
				}

				paths := article.GetWikipediaArticlePaths(articleHtml)

				fmt.Printf("%s has %d wikipedia links\n", path, len(paths))

				corp.Set(path, paths)
				fetched.Store(path, true)

				for _, p := range paths {
					if _, ok := fetched.Load(p); ok {
						continue
					}
					select {
					case queue <- p:
					default:
						goto exit_loop // queue is full
					}
				}
			exit_loop:
			}()
		}

		wg.Wait()
	}

	// Must set cap of urls to the number of pages in the corpus or else the ForEach will overwrite
	// the pointer of the inner slice, and `urls` will be empty.
	urls := make([]string, 0, corp.Size())
	corp.ForEach(func(url string, _ []string) {
		urls = append(urls, url)
	})

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

	// fmt.Println("-- Before consistency check --")
	fmt.Println("Total links in corpus:", corp.GetTotalLinks())
	fmt.Println("Top three articles:")
	fmt.Printf("1. %s with %s as last link\n", urls[0], a[len(a)-1])
	fmt.Printf("2. %s with %s as last link\n", urls[1], b[len(b)-1])
	fmt.Printf("3. %s with %s as last link\n", urls[2], c[len(c)-1])

	corp.EnsureConsistency()
	corp.CheckConsistency()

	fmt.Println("-- Corpus is consistent --")
	fmt.Println("Total links in corpus:", corp.GetTotalLinks())
	fmt.Println("Top three articles:")
	fmt.Printf("1. %s with %s as last link\n", urls[0], a[len(a)-1])
	fmt.Printf("2. %s with %s as last link\n", urls[1], b[len(b)-1])
	fmt.Printf("3. %s with %s as last link\n", urls[2], c[len(c)-1])

}

/*

c := make(chan struct{}, 5)
defer close(c)

statusCh := make(chan int, len(urls.URLs))

var wg sync.WaitGroup
wg.Add(len(urls.URLs))

for i, url := range urls.URLs {
	c <- struct{}{} // acquire a slot

	go func(url string, i int) {
		defer wg.Done()
		defer func() { <-c }() // release a slot

		fmt.Printf("Fetch #%d: %s\n", i, url)
		statusCh <- urls.FetchAndGetStatusCode(url)
	}(url, i)
}

wg.Wait()
close(statusCh)

*/
