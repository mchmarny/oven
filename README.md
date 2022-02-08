# oven

[Cloud Firestore](https://firebase.google.com/docs/firestore) wrapper client library

[![Test](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml/badge.svg?branch=main)](https://github.com/mchmarny/oven/actions/workflows/test-on-push.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/oven)](https://goreportcard.com/report/github.com/mchmarny/oven) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/oven) [![codecov](https://codecov.io/gh/mchmarny/oven/branch/main/graph/badge.svg?token=00H8S7GMPP)](https://codecov.io/gh/mchmarny/oven)


# Overview 

Cloud Firestore is a Serverless document database that's easy to use. It comes with with multi-region replication, powerful query engine, and seamless integration into the broader Google Cloud. I found myself use it frequently over last few years. 

While Firestore exposes both HTTP and RPC APIs to which you can make direct calls, most people rely on one of the [client libraries](https://cloud.google.com/firestore/docs/reference/libraries). The [Go client library](https://pkg.go.dev/cloud.google.com/go/firestore) for Firestore is [well documented](https://firebase.google.com/docs/firestore/quickstart) and rich set of features. Having build a few services with Firestore Go library I did find it wee bit verbose and repetitive in places. It's mostly due to the fact that most of the libraries are auto-generated. Oven is a standalone library on top of the client that hides some of the complexity and shortens some of the more verbose use-cases. 

# Features

* Easy Save, Update, Get, Delete
* Query support using structured criteria 
* Exposes native Firestore client for when you need to extend it

# Install

```shell
go get github.com/mchmarny/oven
```

# Usage

The [examples](./examples) folder includes some of the most common use-cases

* [Basic CRUD Operations](examples/crud) - Save, Get, Update, Delete methods. 
* [Structured Query](examples/query) - Automatic mapping of multiple document results to a Go slice of structs (e.g. `[]Book{}` or `[]*Book{}`)

Simple Example

```go
	var list []*book.Book
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

	if err := service.Query(ctx, q, &list); err != nil {
		panic(err)
	}

	for i, b := range list {
		fmt.Printf("book[%d]: %+v\n", i, b)
	}
```

# License

See [LICENSE](LICENSE)