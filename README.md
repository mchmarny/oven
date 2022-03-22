# oven

[Cloud Firestore](https://firebase.google.com/docs/firestore) client helper library

[![test](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml/badge.svg?branch=main)](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml) 
[![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/oven)](https://goreportcard.com/report/github.com/mchmarny/oven) 
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/oven) 
[![codecov](https://codecov.io/gh/mchmarny/oven/branch/main/graph/badge.svg?token=00H8S7GMPP)](https://codecov.io/gh/mchmarny/oven) 
[![GoDoc](https://godoc.org/github.com/mchmarny/oven?status.svg)](https://godoc.org/github.com/mchmarny/oven)


# Overview 

Cloud Firestore is a Serverless document database. It comes with with multi-regional replication, powerful query engine, and seamless integration into the broader Google Cloud. I found myself use it frequently over last few years. While Firestore exposes both HTTP and RPC APIs to which you can make direct calls, most people rely on one of the [client libraries](https://cloud.google.com/firestore/docs/reference/libraries), and the [Go library](https://pkg.go.dev/cloud.google.com/go/firestore) for Firestore is certainly [well documented](https://firebase.google.com/docs/firestore/quickstart). 

Having used Firestore on a few projects though, I did find it wee bit verbose and repetitive in some places. Oven is a wrapper the the standard Firestore client library. It hides some of the complexity (e.g. iterator over resulting documents), and shortens a few of the more verbose, but common, use-cases. 

# Features

* Easy Save, Update, Get, Delete
* Structured criteria Query
* Extends Firestore client

# Install

```shell
go get github.com/mchmarny/oven
```

# Usage

The [examples](./examples) folder includes some of the most common use-cases with two fully functional code:

* [Basic CRUD](examples/crud/main.go) - Save, Get, Update, Delete operations 
* [Structured Query](examples/query/main.go) - With automatic result mapping to a slice

## Save

```go
b := Book{
	ID:     "id-123",
	Title:  "The Hitchhiker's Guide to the Galaxy",
	Author: "Douglas Adams",
}

if err := oven.Save(ctx, client, "books", b.ID, &b); err != nil {
	handleErr(err)
}
```

## Get

```go
b, err := oven.Get[Book](ctx, client, "books", "id-123")
if err != nil {
	handleErr(err)
}
```

## Update

> Using `Title` field with `firestore:"title"` tag

```go
m := map[string]interface{}{
	"title": "Some new title",
}

if err := oven.Update(ctx, client, "books", "id-123", m); err != nil {
	handleErr(err)
}
```

## Delete

```go
if err := oven.Delete(ctx, client, "books", "id-123"); err != nil {
	handleErr(err)
}
```

## Query

```go
q := &oven.Criteria{
	Collection: book.CollectionName,
	Criterions: []*oven.Criterion{
		{
			Path:      "author", // `firestore:"author"`
			Operation: oven.OperationTypeEqual,
			Value:     bookAuthor,
		},
	},
	OrderBy: "published", // `firestore:"published"`
	Desc:    true,
	Limit:   bookCollectionSize,
}

list, err := oven.Query(ctx, client, q)
if err != nil {
	handleErr(err)
}
```

## Iterator 

In case you already have the Firestore iterator and want to just avoid the verbose `for` loop of spooling the documents into a list, `oven` provides the `ToStructs` method

```go
// given it as *firestore.DocumentIterator
it := col.Documents(ctx)

// this returns a slice of books ([]*Book)
list, err := oven.ToStructs[Book](it)
```

# License

See [LICENSE](LICENSE)