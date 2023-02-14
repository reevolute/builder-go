# builder-go
[![Build Status](https://github.com/reevolute/builder-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/reevolute/builder-go/actions/workflows/test.yml?query=branch%3Amain)

## Requirements ##

- Go 1.15 or later

## Installation ##

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

## Usage ##

### Create a client ###

Based on API key and tenant id. Assuming the env var `API_KEY` contains your api key.
```go
tenantID := "my_tenant_1234"
client := builder.New(os.Getenv("API_KEY"), tenantID)
```

### Add execution ###
```go
parameters := map[string]interface{}{
		"color": "red",
}

treeID:= "01G5PGEHAPPJZ8WE14E37M721Q"
response, err := client.AddExecution(treeID, "production", parameters)
```

## License ##

This library is distributed under the MIT-style license found in the [LICENSE](./LICENSE)
file.