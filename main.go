package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/sendelivery/wikipedia-pagerank/internal/config"
	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
	"github.com/sendelivery/wikipedia-pagerank/internal/pagerank"
	"github.com/sendelivery/wikipedia-pagerank/internal/queue"
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

	cfg := config.DefaultConfig(slog.LevelInfo)
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

	sem := make(chan struct{}, cfg.NumConcurrentFetches) // A "semaphore" to limit concurrent fetches
	defer close(sem)

	sizedQueue := queue.NewSizedQueue(cfg.NumPages)
	defer sizedQueue.Close()
	sizedQueue.Enqueue(parentArticlePath)

	corp := corpus.New()

	var wg sync.WaitGroup

	// Initial parent pass to populate the queue
	cfg.Logger.Debug("Starting parent pass")
	wg.Add(1)
	sem <- struct{}{} // Acquire a slot
	go processPathInQueue(ctx, sizedQueue, &corp, sem, &wg)
	wg.Wait()

	// Main NumPages pass
	cfg.Logger.Debug("Starting main pass", slog.Int("num_pages", cfg.NumPages))
	for i := 0; i < cfg.NumPages-1; i++ {
		wg.Add(1)
		sem <- struct{}{} // Acquire a slot
		go processPathInQueue(ctx, sizedQueue, &corp, sem, &wg)

		// This will prevent the program from hanging if the queue is empty and all goroutines are
		// waiting for a slot.
		if sizedQueue.Empty() {
			cfg.Logger.Debug(
				"Queue is empty, waiting for remaining goroutines",
				slog.Int("num_pages", cfg.NumPages),
			)
			wg.Wait()
			if sizedQueue.Empty() {
				cfg.Logger.Debug("Queue is still empty, returning incomplete corpus", slog.Int("num_pages", cfg.NumPages))
				return &corp
			}
		}
	}
	wg.Wait()

	// Error correction pass
	cfg.Logger.Debug("Starting error correction pass")
	for i := 1; ; i++ {
		numGoroutines := min(cfg.NumConcurrentFetches, cfg.NumPages-corp.Size(), sizedQueue.Length())
		if numGoroutines == 0 {
			break
		}

		cfg.Logger.Debug("Error correction #"+strconv.Itoa(i), slog.Int("num_goroutines", numGoroutines))

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			sem <- struct{}{} // Acquire a slot
			go processPathInQueue(ctx, sizedQueue, &corp, sem, &wg)
		}
		wg.Wait()
	}

	return &corp
}

func processPathInQueue(ctx context.Context, queue *queue.SizedQueue, corp *corpus.Corpus, sem chan struct{}, wg *sync.WaitGroup) {
	cfg, ok := config.ConfigFromContext(ctx)
	if !ok {
		log.Fatal("Failed to get config from context.")
	}

	defer wg.Done()          // Decrement the wait group
	defer func() { <-sem }() // Release the slot

	path, ok := queue.Dequeue()
	if !ok {
		cfg.Logger.Debug("Queue is empty", slog.String("reason", "no more elements"))
		return
	}

	articles, err := scraper.ScrapeArticlesInWikipediaArticle(path, cfg.Logger)
	if err != nil {
		cfg.Logger.Debug("error when scraping "+path, slog.Any("err", err))
		return
	}

	cfg.Logger.Debug("Setting links for "+path+" in corpus", slog.Int("num_links", len(articles)))
	corp.Set(path, articles)

	cfg.Logger.Debug("Iterating through links, adding to queue", slog.Int("num_links", len(articles)))
	for _, articlePath := range articles {
		if queue.Full() {
			cfg.Logger.Debug("Skipping "+articlePath, slog.String("reason", "queue is full"))
			break
		}
		if ok := queue.Enqueue(articlePath); !ok {
			cfg.Logger.Debug("Skipping "+articlePath, slog.String("reason", "already in queue"))
			continue
		}
		cfg.Logger.Debug("Added to the queue", slog.String("link", articlePath))
	}
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
	fmt.Printf("1. %s with %d links\n", urls[0], len(a))
	fmt.Printf("2. %s with %d links\n", urls[1], len(b))
	fmt.Printf("3. %s with %d links\n", urls[2], len(c))
	fmt.Println()
	fmt.Println("Top three articles by PageRank:")
	fmt.Printf("1. %s at %f \n", prSortedUrls[0], pr[prSortedUrls[0]])
	fmt.Printf("1. %s at %f \n", prSortedUrls[1], pr[prSortedUrls[1]])
	fmt.Printf("1. %s at %f \n", prSortedUrls[2], pr[prSortedUrls[2]])
	fmt.Println()
	fmt.Printf("PageRank sums to: %f\n", prSum)
}
