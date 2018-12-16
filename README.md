# retry package

[![GoDoc](https://godoc.org/github.com/semlette/retry?status.svg)](https://godoc.org/github.com/semlette/retry)

`retry` is a simple package for retrying functions if they return an error.

```go
package handler

import "fmt"
import "net/http"
import "github.com/semlette/retry"

func SendHTTP() {
    value, err := retry.Times(3, func() (interface{}, error) {
        return http.Get("/")
    })
    resp := value.(*http.Response)
    if err != nil {
        // failed once + 3 retries
    }
    fmt.Println(resp.StatusCode)
}
```