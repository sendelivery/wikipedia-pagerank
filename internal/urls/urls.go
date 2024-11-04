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
		d := time.Duration(rand.IntN(20) + 20)
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

func GetArticleHTML(path string) (*html.Node, error) {
	fmt.Printf("Fetching link %s\n", path)
	resp, err := http.Get(BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s with error %s", path, err.Error())
	}

	fmt.Println("Parsing response body")
	articleHtml, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html for %s with error %e", path, err)
	}

	return articleHtml, nil
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
