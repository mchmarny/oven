package oven

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/mchmarny/oven/pkg/meta"
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

// Query represents a query to be executed against the Firestore.
type Query struct {
	Collection string
	Criteria   *Criterion
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

func appendWhere(col *firestore.CollectionRef, criteria ...*Criterion) {
	if col == nil || len(criteria) == 0 {
		return
	}

	for _, c := range criteria {
		col.Where(c.Path, string(c.Operation), c.Value)
	}
}

// Query retreaves access info for all users since last update.
func (s *Service) Query(ctx context.Context, q *Query, d interface{}) error {
	if q == nil || q.Collection == "" {
		return errors.Errorf("valid query required: %+v", q)
	}

	// collection
	col, err := s.GetCollection(ctx, q.Collection)
	if err != nil {
		return errors.Wrapf(err, "error getting collection %s", q.Collection)
	}

	// where
	appendWhere(col, q.Criteria)

	// desc?
	dir := firestore.Desc
	if !q.Desc {
		dir = firestore.Asc
	}

	m, err := meta.New(d)
	if err != nil {
		return err
	}

	// iterate
	it := col.OrderBy(q.OrderBy, dir).Limit(q.Limit).Documents(ctx)
	for {
		d, e := it.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return e
		}

		item := m.Item()
		if e := d.DataTo(&item); e != nil {
			return errors.Errorf("error converting data to struct: %v", e)
		}
		m.Append(item)
	}

	return nil
}
