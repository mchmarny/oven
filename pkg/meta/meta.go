package meta

import (
	"reflect"

	"github.com/pkg/errors"
)

// Destination represents a meta information derived from an interface representing slice of structs.
type Destination struct {
	list      reflect.Value
	listByPtr bool
	itemType  reflect.Type
}

// New creates a new MetaList from a given list (pointer to []something or []*something).
func New(d interface{}) (*Destination, error) {
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

	return &Destination{
		list:      list,
		listByPtr: itemIsPtr,
		itemType:  itemType,
	}, nil
}

// Item returns the item.
func (d *Destination) Item() interface{} {
	return reflect.New(d.itemType).Interface()
}

// Append appends a new item to the list.
func (d *Destination) Append(item interface{}) {
	var itemVal reflect.Value
	if d.listByPtr {
		itemVal = reflect.ValueOf(item)
	} else {
		itemVal = reflect.ValueOf(item).Elem()
	}

	listVal := reflect.Append(d.list, itemVal)

	d.list.Set(listVal)
}
