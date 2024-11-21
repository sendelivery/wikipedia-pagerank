package pagerank

import (
	"math"

	"github.com/sendelivery/wikipedia-pagerank/internal/corpus"
)

const DAMPING_FACTOR = 0.85
const CONVERGENCE_THRESHOLD = 0.0001

// Return a probability distribution over which page to visit next, given a current page.
// With probability `damping_factor`, choose a link at random linked to by `page`.
// With probability `1 - damping_factor`, choose a link at random chosen from all pages in the
// corpus.
func TransitionModel(corp *corpus.Corpus, page string) map[string]float64 {
	numLinks := 0
	links, ok := corp.Get(page)
	if ok {
		numLinks = len(links)
	}

	distribution := make(map[string]float64, corp.Size())

	if numLinks == 0 {
		// Current page links to no pages, evenly distribute probability.
		probability := 1.0 / float64(corp.Size())
		corp.ForEach(func(key string, _ []string) {
			distribution[key] = probability
		})
		return distribution
	}

	linkedPageProbability := DAMPING_FACTOR / float64(numLinks)
	randomPageProbability := (1 - DAMPING_FACTOR) / float64(corp.Size())

	// Calculate probability
	for _, link := range links {
		distribution[link] = linkedPageProbability
	}

	corp.ForEach(func(key string, _ []string) {
		if _, ok := distribution[key]; ok {
			distribution[key] += randomPageProbability
			return
		}
		distribution[key] = randomPageProbability
	})

	return distribution
}

// Return PageRank values for each page by iteratively updating
// PageRank values until convergence.
//
// Return a map where keys are page names, and values are
// their estimated PageRank value (a value between 0 and 1). All
// PageRank values should sum to 1.
func CalculatePagerank(corp *corpus.Corpus) map[string]float64 {
	// Pages that have no links should be interpreted as having one link
	// for every page in the corpus, including itself

	// Work out all pages that link to a given page
	linksTo := make(map[string][]string, corp.Size())
	corp.ForEach(func(page string, _ []string) {
		// If our `page` is linked to by `candidate`, add `candidate` to the `linksTo`
		// list for our `page`
		corp.ForEach(func(candidate string, candidateLinks []string) {
			for _, cl := range candidateLinks {
				if cl == page {
					linksTo[page] = append(linksTo[page], candidate)
				}
			}
		})
	})

	defaultPageRank := 1.0 / float64(corp.Size())

	pageRank := make(map[string]float64, corp.Size())
	pageRankDiff := make(map[string]float64, corp.Size())

	iterate := true
	for iterate {
		corp.ForEach(func(page string, links []string) {
			// Sum the probability of choosing `page` from each page that links to it
			pChosePage := 0.0
			for _, linkingPage := range linksTo[page] {
				prVal, ok := pageRank[linkingPage]
				if !ok {
					prVal = defaultPageRank
				}
				links, _ := corp.Get(linkingPage)
				pChosePage += prVal / float64(len(links))
			}

			// Calculate `page`s new page rank as the probability across
			// all pages + the normalised page rank of all linking pages
			newRank := (1-DAMPING_FACTOR)/float64(corp.Size()) + (DAMPING_FACTOR * pChosePage)
			oldRank, ok := pageRank[page]
			if !ok {
				oldRank = defaultPageRank
			}
			pageRankDiff[page] = math.Abs(newRank - oldRank)
			pageRank[page] = newRank

			// Keep calculating page ranks until none change by more than the convergence threshold
			iterate = false
			for _, diff := range pageRankDiff {
				if diff >= CONVERGENCE_THRESHOLD {
					iterate = true
					break
				}
			}
		})
	}

	// Normalise PageRank values to sum to 1
	prSum := 0.0
	for _, rank := range pageRank {
		prSum += rank
	}
	for page := range pageRank {
		pageRank[page] /= prSum
	}

	return pageRank
}
