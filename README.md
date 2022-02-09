# oven

[Cloud Firestore](https://firebase.google.com/docs/firestore) wrapper client library

[![test](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml/badge.svg?branch=main)](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/oven)](https://goreportcard.com/report/github.com/mchmarny/oven) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/oven) [![codecov](https://codecov.io/gh/mchmarny/oven/branch/main/graph/badge.svg?token=00H8S7GMPP)](https://codecov.io/gh/mchmarny/oven) [![publish](https://github.com/mchmarny/oven/actions/workflows/publish-on-tag.yaml/badge.svg?branch=main)](https://github.com/mchmarny/oven/actions/workflows/publish-on-tag.yaml)


# Overview 

Cloud Firestore is a Serverless document database. It comes with with multi-regional replication, powerful query engine, and seamless integration into the broader Google Cloud. I found myself use it frequently over last few years. While Firestore exposes both HTTP and RPC APIs to which you can make direct calls, most people rely on one of the [client libraries](https://cloud.google.com/firestore/docs/reference/libraries). The [Go client library](https://pkg.go.dev/cloud.google.com/go/firestore) for Firestore is [well documented](https://firebase.google.com/docs/firestore/quickstart) and has a rich set of features. 

Having used Firestore on a few projects I did find it wee bit verbose and repetitive in some places. It's mostly due to the fact that it's auto-generated. Oven is a standalone library, wrapping the Firestore client, that hides some of the complexity and shortens some of the more verbose use-cases. 

# Features

* Easy Save, Update, Get, Delete
* Query support using structured criteria 
* Exposes native Firestore client for when you need to extend it

# Install

```shell
go get github.com/mchmarny/oven
```

# Usage

The [examples](./examples) folder includes some of the most common use-cases with two fully functional code:

* [Basic CRUD](examples/crud/main.go) - Save, Get, Update, Delete operations 
* [Structured Query](examples/query/main.go) - With automatic result mapping to a slice

## Service 

```go
package main

import (
	"context"

	"github.com/mchmarny/oven"
)

func main() {
	ctx := context.Background()
	s := oven.New(ctx, "my-project-id")
}
```

You can also use an existing Firestore client: 

```go
// c is an existing *firestore.Client
s := oven.NewWithClient(c) 
```

## Save

```go
b := Book{
	ID:     "id-123",
	Title:  "The Hitchhiker's Guide to the Galaxy",
	Author: "Douglas Adams",
}

if err := s.Save(ctx, "books", b.ID, &b); err != nil {
	handleErr(err)
}
```

## Get

```go
b := &Book{}
if err := s.Get(ctx, "books", "id-123", b); err != nil {
	handleErr(err)
}
```

## Update

> Using `Title` field with `firestore:"title"` tag

```go
m := map[string]interface{}{
	"title": "Some new title",
}

if err := s.Update(ctx, "books", "id-123", m); err != nil {
	handleErr(err)
}
```

## Delete

```go
if err := s.Delete(ctx, "books", "id-123"); err != nil {
	handleErr(err)
}
```

## Query

```go
q := &oven.Query{
	Collection: "books",
	Criteria: &oven.Criterion{
		Path:      "author", // `firestore:"author"`
		Operation: oven.OperationTypeEqual,
		Value:     "Douglas Adams",
	},
	OrderBy: "published", // `firestore:"published"`
	Desc:    true,
	Limit:   10,
}

var list []*Book
if err := s.Query(ctx, q, &list); err != nil {
	handleErr(err)
}
```

# License

See [LICENSE](LICENSE)