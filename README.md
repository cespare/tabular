# tabular

[![Go Reference](https://pkg.go.dev/badge/github.com/cespare/tabular.svg)](https://pkg.go.dev/github.com/cespare/tabular)

tabular is a small Go package for printing tabular text.

I originally wrote this package as an alternative to text/tabwriter because I
needed to omit the trailing padding in each row that that package's algorithm
creates. This resulting package is a little simpler than text/tabwriter and is
better-suited to my typical needs. See the doc comments for more comparisons
between tabular and text/tabwriter.
