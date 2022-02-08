package oven

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

const (
	OperationTypeEqual              OperationType = "=="
	OperationTypeNotEqual           OperationType = "!="
	OperationTypeLessThan           OperationType = "<"
	OperationTypeLessThanOrEqual    OperationType = "<="
	OperationTypeGreaterThan        OperationType = ">"
	OperationTypeGreaterThanOrEqual OperationType = ">="
	OperationTypeArrayContains      OperationType = "array-contains"
	OperationTypeArrayContainsAny   OperationType = "array-contains-any"
	OperationTypeIn                 OperationType = "in"
	OperationTypeNotIn              OperationType = "not-in"
)

var (
	// ErrInvalidDestinationType is thrown when the item type is not a pointer to a struct
	ErrInvalidDestinationType = errors.New("destination type must be a non nil pointer")
)

type Query struct {
	Collection string
	Criteria   *Criterion
	OrderBy    string
	Desc       bool
	Limit      int
}

type Criterion struct {
	Path      string
	Operation OperationType
	Value     interface{}
}

type OperationType string

func appendWhere(col *firestore.CollectionRef, criteria ...*Criterion) {
	for _, c := range criteria {
		col.Where(c.Path, string(c.Operation), c.Value)
	}
}

// GetByQuery retreaves access info for all users since last update.
func (s *Service) Query(ctx context.Context, q *Query, d interface{}) error {
	if q == nil || q.Collection == "" {
		return errors.Errorf("valid query required: %+v", q)
	}

	m, err := getDestinationMeta(d)
	if err != nil {
		return err
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

		item := m.new()
		if e := d.DataTo(&item); e != nil {
			return errors.Errorf("error converting data to struct: %v", e)
		}
		m.append(item)
	}

	return nil
}

type destinationMeta struct {
	list      reflect.Value
	listByPtr bool
	itemType  reflect.Type
}

func (d *destinationMeta) new() interface{} {
	return reflect.New(d.itemType).Interface()
}

func (d *destinationMeta) append(item interface{}) {
	var itemVal reflect.Value
	if d.listByPtr {
		itemVal = reflect.ValueOf(item)
	} else {
		itemVal = reflect.ValueOf(item).Elem()
	}

	listVal := reflect.Append(d.list, itemVal)

	d.list.Set(listVal)
}

func getDestinationMeta(d interface{}) (*destinationMeta, error) {
	if d == nil {
		return nil, errors.New("destination type must be a non nil pointer")
	}

	list := reflect.ValueOf(d)
	if !list.IsValid() || (list.Kind() == reflect.Ptr && list.IsNil()) {
		return nil, errors.Errorf("destination must be a non nil pointer")
	}
	if list.Kind() != reflect.Ptr {
		return nil, errors.Errorf("destination must be a pointer, got: %v", list.Type())
	}

	list = list.Elem()
	listType := list.Type()

	if list.Kind() != reflect.Slice {
		return nil, errors.Errorf("destination must be a slice, got: %v", listType)
	}

	itemType := listType.Elem()
	var itemIsPtr bool
	// dereference to value if pointers to struct
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
		itemIsPtr = true
	}

	// make sure slice is empty
	list.Set(list.Slice(0, 0))

	return &destinationMeta{
		list:      list,
		listByPtr: itemIsPtr,
		itemType:  itemType,
	}, nil
}
