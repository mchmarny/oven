package book

import "time"

const (
	// CollectionName is the book collection name.
	CollectionName = "books"
)

// Book represents a simple book document.
type Book struct {
	BookID    string    `json:"book_id" firestore:"book_id"`
	Name      string    `json:"name" firestore:"name"`
	Author    string    `json:"author" firestore:"author"`
	Published time.Time `json:"published" firestore:"published"`
	Pages     int       `json:"pages" firestore:"pages"`
	HardCover bool      `json:"hardcover" firestore:"hardcover"`
}
