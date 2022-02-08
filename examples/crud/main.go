package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mchmarny/oven"
	"github.com/mchmarny/oven/pkg/id"
)

const (
	collection = "books"
)

type Book struct {
	BookID    string    `json:"book_id" firestore:"book_id"`
	Name      string    `json:"name" firestore:"name"`
	Published time.Time `json:"published" firestore:"published"`
	Pages     int       `json:"pages" firestore:"pages"`
	HardCover bool      `json:"hardcover" firestore:"hardcover"`
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("PROJECT_ID not set")
	}

	ctx := context.Background()

	// create a new oven service instance
	service := oven.New(ctx, projectID)

	// create a book
	book1 := &Book{
		BookID:    id.NewID(),
		Name:      "The Hitchhiker's Guide to the Galaxy",
		Published: time.Now().AddDate(-1, -2, -3).UTC(),
		Pages:     100,
		HardCover: true,
	}

	// save book
	if err := service.Save(ctx, collection, book1.BookID, book1); err != nil {
		log.Fatalf("failed to save: %v", err)
	}

	// get the previously saved book
	book2 := &Book{}
	if err := service.GetByID(ctx, collection, book1.BookID, book2); err != nil {
		log.Fatalf("failed to get: %v", err)
	}

	// check book's content
	if fmt.Sprintf("%v", book1) != fmt.Sprintf("%v", book2) {
		log.Fatalf("books are not the same, want %v, got %v", book1, book2)
	}

	// update
	updatedPages := 200
	updates := map[string]interface{}{"pages": updatedPages}
	if err := service.Update(ctx, collection, book2.BookID, updates); err != nil {
		log.Fatalf("failed to update: %v", err)
	}

	book3 := &Book{}
	if err := service.GetByID(ctx, collection, book2.BookID, book3); err != nil {
		log.Fatalf("failed to get: %v", err)
	}

	// check book's content
	if book3.Pages != updatedPages {
		log.Fatalf("book not updated, wanted %d, got %d", updatedPages, book3.Pages)
	}

	// delete
	if err := service.DeleteByID(ctx, collection, book3.BookID); err != nil {
		log.Fatalf("failed to delete: %v", err)
	}

	// no data found error after delete
	book4 := &Book{}
	if err := service.GetByID(ctx, collection, book3.BookID, book4); err != nil {
		if err != oven.ErrDataNotFound {
			log.Fatalf("expected ErrDataNotFound error, got: %v", err)
		}
	}
}
