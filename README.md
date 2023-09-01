# cedar-agent-go-sdk

## Installation

```shell
go get -u github.com/kevinmichaelchen/cedar-agent-go-sdk
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/kevinmichaelchen/cedar-agent-go-sdk/sdk"
	"net/http"
)

func main() {
	ctx := context.Background()
	c := &http.Client{}
	client := sdk.NewClient(c)
	res, err := client.Check(ctx, sdk.CheckRequest{
		Principal: `User::"42"`,
		Action: "viewFoobar",
		Resource: `Foobar::"101"`,
	})
	if err != nil {
		panic(err)
    }
	fmt.Printf("allowed: %t", res.Decision == "Allow")
}
```