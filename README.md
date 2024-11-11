# wikipedia-pagerank

Type in any Wikipedia article and this program will build a corpus of Wikipedia articles by scraping the hyperlinks of each page. Then, it'll calculate the pagerank of each article in that corpus, printing out the result!

How fun! :)

Built to learn Go.

## Usage

```
$ make build
$ ./bin/wikipedia-pagerank
Enter a Wikipedia URL: https://en.wikipedia.org/wiki/Albert_Camus
Building corpus.........................
Enforcing corpus consistency.
Calculating pagerank.

Size of corpus: 1000
181428 cross-references in the corpus.

Top three articles by most cross-references:
1. /wiki/Albert_Camus with 1035 links, last link: /wiki/J%C3%BCrgen_Habermas
2. /wiki/Friedrich_Nietzsche with 998 links, last link: /wiki/OCLC_(identifier)
3. /wiki/Max_Weber with 801 links, last link: /wiki/Jean_Baudrillard

Top three articles by PageRank:
1. /wiki/Main_Page at 0.252343 
1. /wiki/ISBN_(identifier) at 0.102934 
1. /wiki/Doi_(identifier) at 0.068566 

PageRank sums to: 0.994319
```

## Development

### Getting Started

This project uses the following dependencies in the [Makefile](Makefile):

- [Staticcheck](https://staticcheck.dev/docs/getting-started/) to perform code-quality checks.
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) to check for vulnerabilities in dependencies.

To install the above, run:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

This will install the latest version of Staticcheck and govulncheck to $GOPATH/bin. To find out where $GOPATH is, run `go env GOPATH`, then be sure to add $GOPATH/bin to your PATH if you haven't already.