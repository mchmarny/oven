package oven

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
)

const (
	// MaxBatchSize is the maximum number of items that can be batched together.
	MaxBatchSize = 500

	agentName = "oven"
)

var (
	// ErrDataNotFound is thrown when query does not find the requested data.
	ErrDataNotFound = errors.New("data not found")
)

func isDataNotFoundError(err error) bool {
	return err != nil && status.Code(err) == codes.NotFound
}

// GetClient returns instantiated Firestore client.
func GetClient(ctx context.Context, projectID string) (client *firestore.Client, err error) {
	return firestore.NewClient(ctx, projectID, option.WithUserAgent(agentName))
}

// GetClientOrPanic returns instantiated Firestore client or panics on error.
func GetClientOrPanic(ctx context.Context, projectID string) *firestore.Client {
	c, err := firestore.NewClient(ctx, projectID, option.WithUserAgent(agentName))
	if err != nil {
		panic(err)
	}
	return c
}

// GetCollection returns specific store collection by name.
func GetCollection(client *firestore.Client, name string) (col *firestore.CollectionRef, err error) {
	if client == nil {
		return nil, errors.New("nil client")
	}

	if name == "" {
		return nil, errors.New("collection name required")
	}

	return client.Collection(name), nil
}

// Delete deletes specific record by id.
func Delete(ctx context.Context, client *firestore.Client, col, id string) error {
	if client == nil {
		return errors.New("nil client")
	}

	if id == "" {
		return errors.New("nil id")
	}

	c, err := GetCollection(client, col)
	if err != nil {
		return err
	}

	_, err = c.Doc(id).Delete(ctx)

	if isDataNotFoundError(err) {
		return nil
	}

	if err != nil {
		return errors.Wrapf(err, "error deleting %s record with id %s", col, id)
	}

	return nil
}

// Get sets the in parameter to specific store record by id.
func Get[T any](ctx context.Context, client *firestore.Client, col, id string) (*T, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}

	if id == "" {
		return nil, errors.New("id required")
	}

	c, err := GetCollection(client, col)
	if err != nil {
		return nil, err
	}

	d, err := c.Doc(id).Get(ctx)
	if err != nil {
		if isDataNotFoundError(err) {
			return nil, ErrDataNotFound
		}
		return nil, errors.Wrapf(err, "error getting %s record with id %s", col, id)
	}

	if d == nil || d.Data() == nil {
		return nil, errors.Errorf("record with id %s found in %s collection but has not data", id, col)
	}

	var out T
	if err := d.DataTo(&out); err != nil {
		return nil, errors.Wrapf(err, "data in %s for id %s is in an incorrect format", col, id)
	}

	return &out, nil
}

// Save saves the data to the store.
func Save[T any](ctx context.Context, client *firestore.Client, col, id string, in *T) error {
	if client == nil {
		return errors.New("nil client")
	}

	if in == nil {
		return errors.New("nil in state to save")
	}

	c, err := GetCollection(client, col)
	if err != nil {
		return err
	}

	_, err = c.Doc(id).Set(ctx, in)
	if err != nil {
		return errors.Wrapf(err, "error saving %s record with id %s", col, id)
	}
	return nil
}

// Update updates the data in the store.
func Update(ctx context.Context, client *firestore.Client, col, id string, args map[string]interface{}) error {
	if client == nil {
		return errors.New("nil client")
	}

	if col == "" || id == "" {
		return errors.New("nil collection or id  in update")
	}

	c, err := GetCollection(client, col)
	if err != nil {
		return err
	}

	updates := make([]firestore.Update, 0)
	for k, v := range args {
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}

	_, err = c.Doc(id).Update(ctx, updates)
	if err != nil {
		return errors.Wrapf(err, "error updating %s record with id %s", col, id)
	}
	return nil
}

type Identifiable interface {
	GetID() string
}

func BatchSet[T Identifiable](ctx context.Context, client *firestore.Client, col string, items ...T) error {
	if client == nil {
		return errors.New("nil client")
	}

	if col == "" {
		return errors.New("nil collection name")
	}

	batchSize := len(items)
	if batchSize == 0 {
		return nil
	}

	if batchSize > MaxBatchSize {
		return errors.Errorf("batch size %d exceeds max batch size %d", batchSize, MaxBatchSize)
	}

	c, err := GetCollection(client, col)
	if err != nil {
		return errors.Wrap(err, "error getting collection")
	}

	b := client.Batch()
	for _, item := range items {
		b.Set(c.Doc(item.GetID()), item)
	}

	if _, err := b.Commit(ctx); err != nil {
		return errors.Wrapf(err, "error batch setting %d records on %s", batchSize, c.ID)
	}

	return nil
}
