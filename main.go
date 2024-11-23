package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"
	"sync"

	"github.com/sendelivery/wikipedia-pagerank/internal/config"
	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
	"github.com/sendelivery/wikipedia-pagerank/internal/pagerank"
	"github.com/sendelivery/wikipedia-pagerank/internal/scraper"
)

func main() {
	if len(os.Args) != 2 {
		msg := "Usage: %s <Wikipedia article path>\nExample: %s /wiki/Go_(programming_language)\n"
		fmt.Printf(msg, os.Args[0], os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]

	if !scraper.IsValidWikiPath(path) {
		fmt.Printf("Invalid Wikipedia path: %s\nPlease try again.", path)
		os.Exit(1)
	}

	cfg := config.DefaultConfig()
	ctx := config.ContextWithConfig(context.Background(), cfg)

	cfg.Reporter.NewWorkInProgress("Building corpus")
	corp := buildCorpus(ctx, path)
	cfg.Reporter.Stop()

	cfg.Reporter.NewWorkInProgress("Calculating PageRank")
	corp.EnforceConsistency()
	pr := pagerank.CalculatePagerank(corp)
	cfg.Reporter.Stop()

	printResults(corp, pr)
}

func buildCorpus(ctx context.Context, parentArticlePath string) *corpus.Corpus {
	cfg, ok := config.ConfigFromContext(ctx)
	if !ok {
		log.Fatal("Failed to get config from context.")
	}

	// Create a buffered channel that we'll use as a queue to hold all the pages we have yet to
	// fetch.
	queue := make(chan string, cfg.NumPages)
	queue <- parentArticlePath

	corp := corpus.New()

	for corp.Size() < cfg.NumPages {
		if len(queue) == 0 {
			cfg.Logger.Info("Queue is empty, unable to build full corpus", slog.Int("final_size", corp.Size()))
			return &corp
		}

		numStartGoroutines := min(len(queue), cfg.NumConcurrentFetches, cfg.NumPages-corp.Size())

		var wg sync.WaitGroup
		wg.Add(numStartGoroutines)

		for i := numStartGoroutines; i > 0; i-- {
			go func() {
				defer wg.Done()

				path := <-queue
				articles, err := scraper.ScrapeArticlesInWikipediaArticle(path, cfg.Logger)
				if err != nil {
					cfg.Logger.Error("error when scraping "+path, slog.Any("err", err))
					return
				}

				corp.Set(path, articles)

				for _, articlePath := range articles {
					if _, ok := corp.Get(articlePath); ok {
						cfg.Logger.Debug("Skipping "+articlePath, slog.String("reason", "already in corpus"))
						continue
					}
					select {
					case queue <- articlePath:
					default:
						cfg.Logger.Debug("Skipping "+articlePath, slog.String("reason", "queue is full"))
						goto end_goroutine
					}
				}
			end_goroutine:
			}()
		}

		wg.Wait()
	}
	return &corp
}

func printResults(corp *corpus.Corpus, pr map[string]float64) {
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

	prSortedUrls := slices.Clone(urls)
	slices.SortFunc(prSortedUrls, func(a, b string) int {
		rankA := pr[a]
		rankB := pr[b]

		if rankA < rankB {
			return 1
		}
		if rankA > rankB {
			return -1
		}
		return 0
	})

	prSum := 0.0
	for _, v := range pr {
		prSum += v
	}

	fmt.Println()
	fmt.Println(corp.Size(), "pages in the corpus.")
	fmt.Printf("%d cross-references in the corpus.\n", corp.TotalLinks())
	fmt.Println()
	fmt.Println("Top three articles by most cross-references:")
	fmt.Printf("1. %s with %d links, last link: %s\n", urls[0], len(a), a[len(a)-1])
	fmt.Printf("2. %s with %d links, last link: %s\n", urls[1], len(b), b[len(b)-1])
	fmt.Printf("3. %s with %d links, last link: %s\n", urls[2], len(c), c[len(c)-1])
	fmt.Println()
	fmt.Println("Top three articles by PageRank:")
	fmt.Printf("1. %s at %f \n", prSortedUrls[0], pr[prSortedUrls[0]])
	fmt.Printf("1. %s at %f \n", prSortedUrls[1], pr[prSortedUrls[1]])
	fmt.Printf("1. %s at %f \n", prSortedUrls[2], pr[prSortedUrls[2]])
	fmt.Println()
	fmt.Printf("PageRank sums to: %f\n", prSum)
}
