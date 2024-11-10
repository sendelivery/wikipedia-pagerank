package corpus

import (
	"slices"
	"sync"
	"sync/atomic"
)

type Corpus struct {
	syncmap sync.Map
	size    atomic.Int64
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

func (c *Corpus) Set(key string, value []string) {
	if _, ok := c.syncmap.Load(key); !ok {
		c.size.Add(1)
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

func (c *Corpus) ForEach(f func(string, []string)) {
	c.syncmap.Range(func(key, value any) bool {
		k, _ := key.(string)
		v, _ := value.([]string)

		f(k, v)

		return true
	})
}

func (c *Corpus) EnsureConsistency() {
	c.ForEach(func(page string, links []string) {
		linksToRemove := 0

		// Sort the links so that unknown pages are at the end.
		slices.SortFunc(links, func(a, b string) int {
			_, aOk := c.Get(a)
			_, bOk := c.Get(b)

			if aOk && !bOk {
				return -1
			}
			if !aOk && bOk {
				return 1
			}

			return 0
		})

		// Count the number of unknown pages.
		for _, link := range links {
			if _, ok := c.Get(link); !ok {
				linksToRemove++
			}
		}

		// Adjust the corpus by removing the unknown pages.
		end := len(links) - linksToRemove
		c.syncmap.Store(page, links[:end])
	})

}

func (c *Corpus) CheckConsistency() {
	unknownPages := make(map[string]struct{})

	// Find all unknown pages.
	c.ForEach(func(page string, links []string) {
		for _, link := range links {
			if _, ok := c.Get(link); !ok {
				unknownPages[link] = struct{}{}
			}
		}
	})

	// Print unknown pages.
	for page := range unknownPages {
		println("Unknown page:", page)
	}
}

func (c *Corpus) GetTotalLinks() int {
	total := 0
	c.ForEach(func(_ string, links []string) {
		total += len(links)
	})
	return total
}

func New(size int) Corpus {
	return Corpus{}
}
