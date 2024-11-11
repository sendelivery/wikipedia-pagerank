package corpus

import (
	"slices"
	"sync"
	"sync/atomic"
)

type Corpus struct {
	syncmap    sync.Map
	size       atomic.Int64
	totalLinks atomic.Int64
}

// Read-only method, we won't be modifying the map so a value receiver is fine.
// func (c Corpus) GetRandomLink(page string) (string, error) {
// 	links, ok := c[page]

// 	if !ok {
// 		return "", errors.New("page does not exist in corpus")
// 	}

// 	i := rand.Intn(len(links))
// 	link := links[i]

// 	if _, ok := c[link]; !ok {
// 		return "", errors.New("link does not exist in corpus")
// 	}

// 	return link, nil
// }

func (c *Corpus) Size() int {
	return int(c.size.Load())
}

func (c *Corpus) TotalLinks() int {
	return int(c.totalLinks.Load())
}

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

func (c *Corpus) Get(key string) (value []string, ok bool) {
	val, ok := c.syncmap.Load(key)
	if !ok {
		return []string{}, false
	}

	v, _ := val.([]string)
	return v, true
}

func (c *Corpus) ForEach(f func(page string, links []string)) {
	c.syncmap.Range(func(key, value any) bool {
		k, _ := key.(string)
		v, _ := value.([]string)

		f(k, v)

		return true
	})
}

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

func New(size int) Corpus {
	return Corpus{}
}
