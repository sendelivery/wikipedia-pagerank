package urls

import (
	"fmt"
	"net/http"
)

type Response struct {
	Status int
	Body   string
}

var URLs = [...]string{
	"http://example.com/",
	"https://jsonplaceholder.typicode.com/",
	"https://httpbin.org/",
	"https://api.coindesk.com/v1/bpi/currentprice.json",
	"https://news.ycombinator.com/",
	"https://www.bbc.com/news",
	"https://dev.to/",
	"https://www.reddit.com/r/golang/",
	"https://data.gov/",
	"https://openweathermap.org/current"}

var DummyURLs = [...]string{
	"https://www.url0.com",
	"https://www.url1.com",
	"https://www.url2.com",
	"https://www.url3.com",
	"https://www.url4.com",
	"https://www.url5.com",
	"https://www.url6.com",
	"https://www.url7.com",
	"https://www.url8.com",
	"https://www.url9.com",
}

// UrlDummyResponses is a map of dummy responses for each URL
//
// 200: OK
//
// 4xx: Non-retryable error
//
// 5xx: Retryable error
var UrlDummyResponses map[string]Response = map[string]Response{
	"https://www.url0.com": {200, "url0"},
	"https://www.url1.com": {404, "url1"},
	"https://www.url2.com": {500, "url2"},
	"https://www.url3.com": {400, "url3"},
	"https://www.url4.com": {200, "url4"},
	"https://www.url5.com": {200, "url5"},
	"https://www.url6.com": {500, "url6"},
	"https://www.url7.com": {200, "url7"},
	"https://www.url8.com": {500, "url8"},
	"https://www.url9.com": {200, "url9"},
}

func FetchAndGetStatusCode(url string) int {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error:", url, response, err)
		return -1
	}
	return response.StatusCode
}
