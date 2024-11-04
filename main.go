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
const NUM_PAGES = 1_001

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

				corp.Store(path, paths)
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

	count := 0
	urls := make([]string, NUM_PAGES/2)
	corp.Range(func(key, value any) bool {
		count++
		url, ok := key.(string)
		if !ok {
			return false
		}
		urls = append(urls, url)
		return true
	})

	slices.SortFunc(urls, func(a, b string) int {
		aVal, _ := corp.Load(a)
		bVal, _ := corp.Load(b)

		aValArr, ok := aVal.([]string)
		if !ok {
			return 0
		}
		bValArr, ok := bVal.([]string)
		if !ok {
			return 0
		}

		return len(bValArr) - len(aValArr)
	})

	a, ok := corp.Load(urls[0])
	if !ok {
		fmt.Println("a")
		return
	}
	b, ok := corp.Load(urls[1])
	if !ok {
		fmt.Println("b")
		return
	}
	c, ok := corp.Load(urls[3])
	if !ok {
		fmt.Println("c")
		return
	}

	aArr, ok := a.([]string)
	if !ok {
		fmt.Println("aArr")
		return
	}
	bArr, ok := b.([]string)
	if !ok {
		fmt.Println("bArr")
		return
	}
	cArr, ok := c.([]string)
	if !ok {
		fmt.Println("cArr")
		return
	}

	fmt.Printf("Scraped %d articles.\n", count)
	fmt.Println("Top three articles:")
	fmt.Printf("1. %s with %d links\n", urls[0], len(aArr))
	fmt.Printf("2. %s with %d links\n", urls[1], len(bArr))
	fmt.Printf("3. %s with %d links\n", urls[2], len(cArr))

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

/*

Scraped 764 articles.
Top three articles:
1. /wiki/Republican_Party_(United_States) with 3399 links
2. /wiki/Zeus with 3380 links
3. /wiki/Russian_invasion_of_Ukraine with 3034 links

*/
