package corpus

import (
	"golang.org/x/sync/syncmap"
)

type Corpus = syncmap.Map

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

func MakeCorpus(size int) Corpus {
	return Corpus{}
}
