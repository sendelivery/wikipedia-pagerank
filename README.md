# wikipedia-pagerank

Type in any Wikipedia article and this program will build a corpus of Wikipedia articles by scraping the hyperlinks of each page. Then, it'll calculate the pagerank of each article in that corpus, printing out the result!

How fun! :)

Built to learn Go.

# Development

## Getting Started

This project uses the following dependencies in the [Makefile](Makefile):

- [Staticcheck](https://staticcheck.dev/docs/getting-started/) to perform code-quality checks.
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) to check for vulnerabilities in dependencies.

To install the above, run:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

This will install the latest version of Staticcheck and govulncheck to $GOPATH/bin. To find out where $GOPATH is, run `go env GOPATH`, then be sure to add $GOPATH/bin to your PATH if you haven't already.