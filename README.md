# builder-go
[![Build Status](https://github.com/reevolute/builder-go/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/reevolute/builder-go/actions/workflows/test.yml?query=branch%3Amaster)

## Requirements

- Go 1.15 or later

## Installation

Make sure your project is using Go Modules:

``` sh
go mod init
```

Then, reference builder-go in a Go program with `import`:

``` go
import	"github.com/reevolute/builder-go"
```
Run any of the normal `go` commands. The Go toolchain will resolve and fetch the stripe-go module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/reevolute/builder-go
```

## License ##

This library is distributed under the MIT-style license found in the [LICENSE](./LICENSE)
file.