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

// An unknown page in our corpus is a link for which there is no corresponding entry in the corpus.
// I.e. page A might link to pages B, C, and D. B might link to C, and C to A and D. However, D does not
// exist as a key in our corpus. I.e. we do not know which pages D links to! In this case, D would be an
// unknown page. We need to decide how to deal with such pages.

// Read-only method, we won't be modifying the map so a value receiver is fine.
// func (c Corpus) GetRandomLink(page string) (string, error) {
// 	links, ok := c[page]

// 	if !ok {
// 		return "", errors.New("page does not exist in corpus")
// 	}

// 	i := rand.Intn(len(links))
// 	link := links[i]

// 	// Todo: should we ensure that our corpus is valid? (i.e. no unknown pages)
// 	// Alternatively, we could loop until we find a known link,
// 	// or have our pagerank algorithm deal with unknown pages...
// 	if _, ok := c[link]; !ok {
// 		return "", errors.New("link does not exist in corpus")
// 	}

// 	return link, nil
// }

func (c *Corpus) Size() int {
	return int(c.size.Load())
}

func (c *Corpus) Set(key string, value []string) {
	c.syncmap.Store(key, value)
	c.size.Add(1)
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

func MakeCorpus(size int) Corpus {
	return Corpus{}
}
