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
  -fps float
        set fps for display refresh (default 60)
  -logger-force-color
        Force color output even if stderr isn't a terminal
  -logger-no-color
        Prevent colorized output even if stderr is a terminal
  -loglevel level
        log level, one of [Debug Verbose Info Warning Error Critical Fatal] (default Info)
  -profile-cpu file
        write cpu profile to file
  -profile-mem file
        write memory profile to file
  -quiet
        Quiet mode, sets loglevel to Error (quietly) to reduces the output
  -truecolor
        Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set) (default true)
```
