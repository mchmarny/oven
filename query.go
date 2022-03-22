package oven

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

const (
	// OperationTypeEqual is equal to.
	OperationTypeEqual OperationType = "=="
	// OperationTypeNotEqual is not equal to.
	OperationTypeNotEqual OperationType = "!="
	// OperationTypeLessThan is less than.
	OperationTypeLessThan OperationType = "<"
	// OperationTypeLessThanOrEqual is less than or equal to.
	OperationTypeLessThanOrEqual OperationType = "<="
	// OperationTypeGreaterThan is greater than.
	OperationTypeGreaterThan OperationType = ">"
	// OperationTypeGreaterThanOrEqual is greater than or equal to.
	OperationTypeGreaterThanOrEqual OperationType = ">="
	// OperationTypeArrayContains is array contains.
	OperationTypeArrayContains OperationType = "array-contains"
	// OperationTypeArrayContainsAny is array contains any.
	OperationTypeArrayContainsAny OperationType = "array-contains-any"
	// OperationTypeIn is in.
	OperationTypeIn OperationType = "in"
	// OperationTypeNotIn is not in.
	OperationTypeNotIn OperationType = "not-in"
)

var (
	// ErrInvalidDestinationType is thrown when the item type is not a pointer to a struct.
	ErrInvalidDestinationType = errors.New("destination type must be a non nil pointer")
)

// Criteria represents a query to be executed against the Firestore.
type Criteria struct {
	Collection string
	Criterions []*Criterion
	OrderBy    string
	Desc       bool
	Limit      int
}

// Criterion represents a single criteria for a query.
type Criterion struct {
	Path      string
	Operation OperationType
	Value     interface{}
}

// OperationType represents the type of operation to be performed.
type OperationType string

// Query retreaves access info for all users since last update.
func Query[T any](ctx context.Context, client *firestore.Client, c *Criteria) ([]*T, error) {
	if c == nil || c.Collection == "" {
		return nil, errors.Errorf("valid query required: %+v", c)
	}

	// collection
	col, err := GetCollection(client, c.Collection)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting collection %s", c.Collection)
	}

	// where
	for _, r := range c.Criterions {
		col.Where(r.Path, string(r.Operation), r.Value)
	}

	// desc?
	dir := firestore.Desc
	if !c.Desc {
		dir = firestore.Asc
	}

	// iterate
	it := col.OrderBy(c.OrderBy, dir).Limit(c.Limit).Documents(ctx)
	return ToStructs[T](it)
}

// ToStructs converst the Firestore document iterator into a slice of structs.
func ToStructs[T any](it *firestore.DocumentIterator) ([]*T, error) {
	if it == nil {
		return nil, errors.Errorf("valid iterator required")
	}

	list := make([]*T, 0)

	for {
		d, e := it.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return nil, e
		}

		var t T
		if e := d.DataTo(&t); e != nil {
			return nil, errors.Errorf("error converting data to struct: %v", e)
		}
		list = append(list, &t)
	}

	return list, nil
}
