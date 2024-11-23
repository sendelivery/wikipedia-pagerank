# wikipedia-pagerank

Provide any Wikipedia article path and this program will build a corpus of Wikipedia articles by scraping the hyperlinks of each page. Then, it'll calculate the pagerank of each article in that corpus, printing out the result!

How fun! :)

Built to learn Go.

## Usage

Once [built](#build-the-project), you can run the program with a Wikipedia article path as an argument. For example:

```
$ ./bin/wikipedia-pagerank /wiki/Computer_science
Building corpus...
Calculating PageRank..

3000 pages in the corpus.
277497 cross-references in the corpus.

Top three articles by most cross-references:
1. /wiki/Artificial_intelligence with 2142 links
2. /wiki/Glossary_of_artificial_intelligence with 1284 links
3. /wiki/Glossary_of_computer_science with 699 links

Top three articles by PageRank:
1. /wiki/The_New_York_Times at 0.016329 
1. /wiki/The_Guardian at 0.007479 
1. /wiki/The_Wall_Street_Journal at 0.007251 

PageRank sums to: 1.000000
```

### Build the Project

There are a number of ways to build this project, the simplest being:

```sh
go build -o bin/
```

If you have GNU Make installed, you can also run one of the following:

- Build to `./bin/`:

    ```sh
    make build
    ```

-  Vets the code using `go vet`, `staticcheck`, and `govulncheck`, then builds to ./bin/

    Requires [staticcheck](https://staticcheck.dev/) and [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) to be installed and available in your path.

    ```sh
    make build-vet
    ```

- You can also build and immediately run the program with an example argument.

    This method also requires staticcheck and govulncheck to be installed and available in your path.

    ```sh
    make build-run
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