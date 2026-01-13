[![GoDoc](https://godoc.org/fortio.org/tcalc?status.svg)](https://pkg.go.dev/fortio.org/tcalc)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/tcalc)](https://goreportcard.com/report/fortio.org/tcalc)
[![GitHub Release](https://img.shields.io/github/release/fortio/tcalc.svg?style=flat)](https://github.com/fortio/tcalc/releases/)
[![CI Checks](https://github.com/fortio/tcalc/actions/workflows/include.yml/badge.svg)](https://github.com/fortio/tcalc/actions/workflows/include.yml)
[![codecov](https://codecov.io/github/fortio/tcalc/graph/badge.svg?token=Yx6QaeQr1b)](https://codecov.io/github/fortio/tcalc)

# tcalc

tcalc is a bitwise calculator that is run from the terminal. It supports basic variable assignments, and most arithmetic and bitwise operations.

## Install
You can get the binary from [releases](https://github.com/fortio/tcalc/releases)

Or just run
```
CGO_ENABLED=0 go install fortio.org/tcalc@latest  # to install (in ~/go/bin typically) or just
CGO_ENABLED=0 go run fortio.org/tcalc@latest  # to run without install
```

or
```
brew install fortio/tap/tcalc
```

or
```
docker run -ti fortio/tcalc
```


## Usage

```
tcalc help

flags:
```
