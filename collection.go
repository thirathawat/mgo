package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collecter is a interface for mongo collection
type Collecter[T any] interface {
	// InsertOne insert one document
	InsertOne(ctx context.Context, model T) error

	// InsertMany insert many documents
	InsertMany(ctx context.Context, models []T) error

	// FindOne find one document
	FindOne(ctx context.Context, options ...SetOption) (*T, error)

	// FindMany find many documents
	FindMany(ctx context.Context, options ...SetOption) ([]T, error)

	// UpdateOne update one document
	UpdateOne(ctx context.Context, options ...SetOption) error

	// UpdateMany update many documents
	UpdateMany(ctx context.Context, options ...SetOption) error

	// DeleteOne delete one document
	DeleteOne(ctx context.Context, options ...SetOption) error

	// DeleteMany delete many documents
	DeleteMany(ctx context.Context, options ...SetOption) error

	// Count count documents
	Count(ctx context.Context, options ...SetOption) (int64, error)

	// Aggregate aggregate documents
	Aggregate(ctx context.Context, options ...SetOption) ([]T, error)
}

// SetOption is a function to set option
type SetOption func(*Option)

// Filter is a function to set filter
type Option struct {
	// Filter is a filter to find
	Filter bson.M

	// Update is a update to update
	Update bson.M

	// Sort is a sort to sort
	Sort bson.M

	// Projection is a projection to projection
	Projection bson.M

	// Pipeline is a pipeline to aggregate
	Pipeline []bson.M

	// Skip is a skip to skip
	Skip *int64

	// Limit is a limit to limit
	Limit *int64
}

// Filter is a function to set filter
func (o Option) findOptions() *options.FindOptions {
	return &options.FindOptions{
		Sort:       o.Sort,
		Skip:       o.Skip,
		Limit:      o.Limit,
		Projection: o.Projection,
	}
}

// Filter is a function to set filter
func bindOptions(options ...SetOption) *Option {
	opt := &Option{}
	for _, option := range options {
		option(opt)
	}

	return opt
}

// Filter is a function to set filter
type collection[T any] struct {
	*mongo.Collection
}

// NewCollection is a function to create a new collection
func NewCollection[T any](c *mongo.Collection) Collecter[T] {
	return &collection[T]{c}
}

// InsertOne insert one document
func (c *collection[T]) InsertOne(ctx context.Context, model T) error {
	_, err := c.Collection.InsertOne(ctx, model)
	return err
}

// InsertMany insert many documents
func (c *collection[T]) InsertMany(ctx context.Context, models []T) error {
	docs := make([]interface{}, len(models))
	for i, model := range models {
		docs[i] = model
	}

	_, err := c.Collection.InsertMany(ctx, docs)
	return err
}

// FindOne find one document
func (c *collection[T]) FindOne(ctx context.Context, options ...SetOption) (*T, error) {
	opt := bindOptions(options...)
	var model T
	err := c.Collection.FindOne(ctx, opt.Filter).Decode(&model)
	return &model, err
}

// FindMany find many documents
func (c *collection[T]) FindMany(ctx context.Context, options ...SetOption) ([]T, error) {
	opt := bindOptions(options...)

	var models []T
	cursor, err := c.Collection.Find(ctx, opt.Filter, opt.findOptions())
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// UpdateOne update one document
func (c *collection[T]) UpdateOne(ctx context.Context, options ...SetOption) error {
	opt := bindOptions(options...)
	_, err := c.Collection.UpdateOne(ctx, opt.Filter, opt.Update)
	return err
}

// UpdateMany update many documents
func (c *collection[T]) UpdateMany(ctx context.Context, options ...SetOption) error {
	opt := bindOptions(options...)
	_, err := c.Collection.UpdateMany(ctx, opt.Filter, opt.Update)
	return err
}

// DeleteOne delete one document
func (c *collection[T]) DeleteOne(ctx context.Context, options ...SetOption) error {
	opt := bindOptions(options...)
	_, err := c.Collection.DeleteOne(ctx, opt.Filter)
	return err
}

// DeleteMany delete many documents
func (c *collection[T]) DeleteMany(ctx context.Context, options ...SetOption) error {
	opt := bindOptions(options...)
	_, err := c.Collection.DeleteMany(ctx, opt.Filter)
	return err
}

// Count count documents
func (c *collection[T]) Count(ctx context.Context, options ...SetOption) (int64, error) {
	opt := bindOptions(options...)
	return c.Collection.CountDocuments(ctx, opt.Filter)
}

// Aggregate aggregate documents
func (c *collection[T]) Aggregate(ctx context.Context, options ...SetOption) ([]T, error) {
	opt := bindOptions(options...)

	cursor, err := c.Collection.Aggregate(ctx, opt.Pipeline)
	if err != nil {
		return nil, err
	}

	var models []T
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}
