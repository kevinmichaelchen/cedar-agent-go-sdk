# cedar-agent-go-sdk

[![GoReportCard example](https://goreportcard.com/badge/github.com/kevinmichaelchen/cedar-agent-go-sdk)](https://goreportcard.com/report/github.com/kevinmichaelchen/cedar-agent-go-sdk)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/kevinmichaelchen/cedar-agent-go-sdk)
[![version](https://img.shields.io/github/v/release/kevinmichaelchen/cedar-agent-go-sdk?include_prereleases&label=latest&logo=ferrari)](https://github.com/kevinmichaelchen/cedar-agent-go-sdk/releases/latest)

[Cedar Agent][cedar-agent] is an HTTP Server that runs the [Cedar][cedar] authorization engine.

It's the easiest way to get up and running with Cedar locally, offering a REST API for managing your entities and policies, as well as policy evaluation.

Cedar lets you answer the question: _Is this **user** (principal) allowed to perform this **action** on this **resource**?_

[cedar-agent]: https://github.com/permitio/cedar-agent
[cedar]: https://www.cedarpolicy.com

## Installation

```shell
go get -u github.com/kevinmichaelchen/cedar-agent-go-sdk
```

## Usage

### Creating a client

```go
package main

import (
	"github.com/kevinmichaelchen/cedar-agent-go-sdk/sdk"
	"net/http"
)

func initClient() *sdk.Client {
	c := &http.Client{}

	// The options are entirely ... optional 🙂
	return sdk.NewClient(c,
		sdk.WithBaseURL("http://localhost:8180"),
		sdk.WithParallelizationFactor(3),
	)
}
```

### Performing authorization checks

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
	client := initClient()
	allowed := isAuthorized(ctx, client,
		sdk.CheckRequest{
			Principal: `User::"42"`,
			Action: "viewFoobar",
			Resource: `Foobar::"101"`,
		},
	)
	fmt.Printf("allowed: %t", allowed)
}

func isAuthorized(ctx context.Context, client *sdk.Client, r sdk.CheckRequest) bool {
	res, err := client.Check(ctx, r)
	if err != nil {
		panic(err)
	}
	return res.Decision == "Allow"
}
```
