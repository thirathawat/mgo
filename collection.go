package mgo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Entity is a base entity.
type Entity struct {
	// ID is a id of entity.
	ID primitive.ObjectID `bson:"_id"`

	// CreatedAt is a time when entity is created.
	CreatedAt time.Time `bson:"created_at"`

	// UpdatedAt is a time when entity is updated.
	UpdatedAt time.Time `bson:"update_at"`

	// DeletedAt is a time when entity is deleted.
	// If entity is not deleted, DeletedAt is nil.
	DeletedAt *time.Time `bson:"deleted_at"`
}

func newRecord(v any) any {
	return struct {
		E Entity `bson:",inline"`
		V any    `bson:",inline"`
	}{
		E: Entity{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		V: v,
	}
}

func (e Entity) entity() Entity {
	return e
}

type entiter interface {
	entity() Entity
}

// Collecter is a interface for mongo collection.
type Collecter[T any] interface {
	// InsertOne insert one document.
	InsertOne(ctx context.Context, model T) error

	// InsertMany insert many documents.
	InsertMany(ctx context.Context, models []T) error

	// FindOne find one document.
	FindOne(ctx context.Context, options ...Option) (*T, error)

	// FindMany find many documents.
	FindMany(ctx context.Context, options ...Option) ([]T, error)

	// UpdateOne update one document.
	UpdateOne(ctx context.Context, options ...Option) error

	// UpdateMany update many documents.
	UpdateMany(ctx context.Context, options ...Option) error

	// DeleteOne delete one document.
	DeleteOne(ctx context.Context, options ...Option) error

	// DeleteMany delete many documents.
	DeleteMany(ctx context.Context, options ...Option) error

	// SoftDeleteOne soft delete one document by setting deleted_at.
	SoftDeleteOne(ctx context.Context, options ...Option) error

	// SoftDeleteMany soft delete many documents by setting deleted_at.
	SoftDeleteMany(ctx context.Context, options ...Option) error

	// Count count documents.
	Count(ctx context.Context, options ...Option) (int64, error)

	// Aggregate aggregate documents.
	Aggregate(ctx context.Context, options ...Option) ([]T, error)
}

// Option is a function to set option.
type Option func(*option)

// Filter is a function to set filter.
type option struct {
	// Filter is a filter to find.
	Filter bson.M

	// Update is a update to update.
	Update bson.M

	// Sort is a sort to sort.
	Sort bson.M

	// Projection is a projection to projection.
	Projection bson.M

	// Pipeline is a pipeline to aggregate.
	Pipeline []bson.M

	// Skip is a skip to skip.
	Skip *int64

	// Limit is a limit to limit.
	Limit *int64
}

// Filter is a function to set filter.
func (o option) findOptions() *options.FindOptions {
	return &options.FindOptions{
		Sort:       o.Sort,
		Skip:       o.Skip,
		Limit:      o.Limit,
		Projection: o.Projection,
	}
}

func (o *option) setUpdate() {
	if o.Update == nil {
		o.Update = bson.M{}
	}

	if val, ok := o.Update["$set"]; !ok {
		o.Update["$set"] = bson.M{
			"updated_at": time.Now(),
		}
	} else {
		val.(bson.M)["updated_at"] = time.Now()
	}
}

func (o *option) setSoftDelete() {
	if o.Update == nil {
		o.Update = bson.M{}
	}

	o.Update["$set"] = bson.M{
		"deleted_at": time.Now(),
	}
}

// Filter is a function to set filter.
func bindOptions(options ...Option) *option {
	opt := &option{}
	for _, option := range options {
		option(opt)
	}

	return opt
}

// Filter is a function to set filter.
type collection[T entiter] struct {
	*mongo.Collection
}

// NewCollection is a function to create a new collection.
func NewCollection[T entiter](c *mongo.Collection) Collecter[T] {
	return &collection[T]{c}
}

// InsertOne insert one document.
func (c *collection[T]) InsertOne(ctx context.Context, model T) error {
	_, err := c.Collection.InsertOne(ctx, newRecord(model))
	return err
}

// InsertMany insert many documents.
func (c *collection[T]) InsertMany(ctx context.Context, models []T) error {
	docs := make([]any, len(models))
	for i, model := range models {
		docs[i] = newRecord(model)
	}

	_, err := c.Collection.InsertMany(ctx, docs)
	return err
}

// FindOne find one document.
func (c *collection[T]) FindOne(ctx context.Context, options ...Option) (*T, error) {
	opt := bindOptions(options...)
	var model T
	err := c.Collection.FindOne(ctx, opt.Filter).Decode(&model)
	return &model, err
}

// FindMany find many documents.
func (c *collection[T]) FindMany(ctx context.Context, options ...Option) ([]T, error) {
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

// UpdateOne update one document.
func (c *collection[T]) UpdateOne(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	opt.setUpdate()
	_, err := c.Collection.UpdateOne(ctx, opt.Filter, opt.Update)
	return err
}

// UpdateMany update many documents.
func (c *collection[T]) UpdateMany(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	opt.setUpdate()
	_, err := c.Collection.UpdateMany(ctx, opt.Filter, opt.Update)
	return err
}

// DeleteOne delete one document.
func (c *collection[T]) DeleteOne(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	_, err := c.Collection.DeleteOne(ctx, opt.Filter)
	return err
}

// DeleteMany delete many documents.
func (c *collection[T]) DeleteMany(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	_, err := c.Collection.DeleteMany(ctx, opt.Filter)
	return err
}

// SoftDeleteOne soft delete one document by setting deleted_at.
func (c *collection[T]) SoftDeleteOne(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	opt.setSoftDelete()
	_, err := c.Collection.UpdateOne(ctx, opt.Filter, opt.Update)
	return err
}

// SoftDeleteMany soft delete many documents by setting deleted_at.
func (c *collection[T]) SoftDeleteMany(ctx context.Context, options ...Option) error {
	opt := bindOptions(options...)
	opt.setSoftDelete()
	_, err := c.Collection.UpdateMany(ctx, opt.Filter, opt.Update)
	return err
}

// Count count documents.
func (c *collection[T]) Count(ctx context.Context, options ...Option) (int64, error) {
	opt := bindOptions(options...)
	return c.Collection.CountDocuments(ctx, opt.Filter)
}

// Aggregate aggregate documents.
func (c *collection[T]) Aggregate(ctx context.Context, options ...Option) ([]T, error) {
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
