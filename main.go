package main

import (
	"fmt"
	"sync"

	urls "github.com/sendelivery/wikipedia-pagerank/internal"
)

func main() {
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

	count := 0
	for code := range statusCh {
		if code == 200 {
			count++
		}
	}
	fmt.Printf("%d/%d URLs returned status code 200\n", count, len(urls.URLs))
}
