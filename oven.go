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
	userAgentDefault = "oven"
)

var (
	// ErrDataNotFound is thrown when query does not find the requested data
	ErrDataNotFound = errors.New("data not found")
)

func New(ctx context.Context, projectID string, opts ...option.ClientOption) *Service {
	s := &Service{
		projectID: projectID,
		options: []option.ClientOption{
			option.WithUserAgent(userAgentDefault),
		},
	}

	s.options = append(s.options, opts...)

	return s
}

// Service provides object representing the inbound HTTP request.
type Service struct {
	client *firestore.Client

	projectID string
	options   []option.ClientOption
}

func isDataNotFoundError(err error) bool {
	return err != nil && status.Code(err) == codes.NotFound
}

// GetClient creates
func (s *Service) GetClient(ctx context.Context) error {
	if s.client == nil {
		c, err := firestore.NewClient(ctx, s.projectID, s.options...)
		if err != nil {
			return errors.Wrapf(err, "error creating Firestore client for project: %s", s.projectID)
		}
		s.client = c
	}

	return nil
}

// GetCollection returns specific store collection by name.
func (s *Service) GetCollection(ctx context.Context, name string) (col *firestore.CollectionRef, err error) {
	if name == "" {
		return nil, errors.New("collection name required")
	}

	if err := s.GetClient(ctx); err != nil {
		return nil, errors.Wrapf(err, "error getting collection %s", name)
	}

	return s.client.Collection(name), nil
}

// Delete deletes specific record by id.
func (s *Service) Delete(ctx context.Context, col, id string) error {
	if id == "" {
		return errors.New("nil id")
	}

	c, err := s.GetCollection(ctx, col)
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
func (s *Service) Get(ctx context.Context, col, id string, in interface{}) error {
	if id == "" {
		return errors.New("id required")
	}

	c, err := s.GetCollection(ctx, col)
	if err != nil {
		return err
	}

	d, err := c.Doc(id).Get(ctx)
	if err != nil {
		if isDataNotFoundError(err) {
			return ErrDataNotFound
		}
		return errors.Wrapf(err, "error getting %s record with id %s", col, id)
	}

	if d == nil || d.Data() == nil {
		return errors.Errorf("record with id %s found in %s collection but has not data", id, col)
	}

	if err := d.DataTo(in); err != nil {
		return errors.Wrapf(err, "data in %s for id %s is in an incorrect format", col, id)
	}

	return nil
}

// Save saves the data to the store.
func (s *Service) Save(ctx context.Context, col, id string, in interface{}) error {
	if in == nil {
		return errors.New("nil in state to save")
	}

	c, err := s.GetCollection(ctx, col)
	if err != nil {
		return err
	}

	_, err = c.Doc(id).Set(ctx, in)
	if err != nil {
		return errors.Wrapf(err, "error saving %s record with id %s", col, id)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, col, id string, args map[string]interface{}) error {
	if col == "" || id == "" {
		return errors.New("nil collection or id  in update")
	}

	c, err := s.GetCollection(ctx, col)
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
