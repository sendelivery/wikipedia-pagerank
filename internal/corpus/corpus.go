package corpus

import (
	"slices"
	"sync"
	"sync/atomic"
)

// Corpus is a map of pages to their links.
type Corpus struct {
	// syncmap is a map of pages to their links, supporting concurrent access via a mutex.
	syncmap sync.Map

	// size is the number of pages in the corpus.
	size atomic.Int64

	// totalLinks is the total number of cross-references in the corpus.
	totalLinks atomic.Int64
}

// Size returns the number of pages in the corpus.
func (c *Corpus) Size() int {
	return int(c.size.Load())
}

// TotalLinks returns the total number of cross-references in the corpus.
func (c *Corpus) TotalLinks() int {
	return int(c.totalLinks.Load())
}

// Set adds a page to the corpus, or updates the links for an existing page.
func (c *Corpus) Set(key string, value []string) {
	oldVal, ok := c.Get(key)
	if !ok {
		c.size.Add(1)
		c.totalLinks.Add(int64(len(value)))
	} else {
		diff := len(value) - len(oldVal)
		c.totalLinks.Add(int64(diff))
	}
	c.syncmap.Store(key, value)
}

// Get returns the links for a page in the corpus.
func (c *Corpus) Get(key string) (value []string, ok bool) {
	val, ok := c.syncmap.Load(key)
	if !ok {
		return []string{}, false
	}

	v, _ := val.([]string)
	return v, true
}

// ForEach iterates over the pages in the corpus, calling the given function for each page.
func (c *Corpus) ForEach(f func(page string, links []string)) {
	c.syncmap.Range(func(key, value any) bool {
		k, _ := key.(string)
		v, _ := value.([]string)

		f(k, v)

		return true
	})
}

// EnforceConsistency removes any links to pages that are not in the corpus.
func (c *Corpus) EnforceConsistency() {
	// A comparison function, for sorting unknown pages to the end of their list.
	sortUnknownPages := func(a, b string) int {
		_, aOk := c.Get(a)
		_, bOk := c.Get(b)

		if aOk && !bOk {
			return -1
		}
		if !aOk && bOk {
			return 1
		}

		return 0
	}

	c.ForEach(func(page string, links []string) {
		slices.SortFunc(links, sortUnknownPages)

		// Count the number of unknown pages.
		linksToRemove := 0
		for _, link := range links {
			if _, ok := c.Get(link); !ok {
				linksToRemove++
			}
		}

		// Adjust the corpus by removing the unknown pages.
		end := len(links) - linksToRemove
		c.Set(page, links[:end])
	})
}

// New creates a new, empty corpus.
func New() Corpus {
	return Corpus{}
}
