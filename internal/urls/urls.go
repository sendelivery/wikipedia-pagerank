package urls

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/html"
)

const BaseURL = "https://en.wikipedia.org"

// InputWikipediaArticlePath prompts the user to input a Wikipedia article path
// and returns that path.
func InputWikipediaArticlePath() string {
	fmt.Printf("Enter a Wikipedia URL: %v", BaseURL)
	time.Sleep(500)
	path := "/wiki/Go_(programming_language)"
	for _, r := range path {
		fmt.Printf("%c", r)
		d := time.Duration(rand.IntN(80) + 20)
		time.Sleep(d * time.Millisecond)
	}
	fmt.Println()
	return path
}

const ValidArticlePathRegex = "^/wiki/[A-Za-z0-9_()%-]+$"

func IsValidWikiPath(url string) bool {
	re := regexp.MustCompile(ValidArticlePathRegex)
	return re.Match([]byte(url))
}

func FetchURL(url string) (*http.Response, error) {
	response, err := http.Get(url)
	return response, err
}

type Corpus struct {
	corpus map[string][]string
}

func (c Corpus) ProcessHTML(doc *html.Node) []string {
	// pageSet := make(map[string]struct{})
	return []string{"Go_(programming_language)"}
}

func (c Corpus) AddPage(page string, links []string) {
	c.corpus[page] = links
}
