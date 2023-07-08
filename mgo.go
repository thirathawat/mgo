// Package mgo provides a MongoDB client.
package mgo

import (
	"context"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// config holds the configuration for the MongoDB client.
type config struct {
	// Addrs is a comma separated list of host:port addresses for the MongoDB
	Addrs string `envconfig:"ADDRS"`

	// Database is the name of the MongoDB database.
	Name string `envconfig:"NAME"`

	// AuthSource is the name of the MongoDB authentication database.
	AuthSource string `envconfig:"AUTH_SOURCE"`

	// User is the name of the MongoDB user.
	User string `envconfig:"USER"`

	// Password is the password for the MongoDB user.
	Password string `envconfig:"PASSWORD"`
}

// readConfig reads the configuration from the environment.
func readConfig() *config {
	cfg := new(config)
	envconfig.MustProcess("MGO", cfg)
	return cfg
}

// client is a wrapper around the mongo.Client type.
type client struct {
	*mongo.Client
	dbName string
}

// Database returns a handle for a given database.
func (c *client) Database() *mongo.Database {
	return c.Client.Database(c.dbName)
}

// New creates a new MongoDB client and establishes a connection.
func New() (db *mongo.Database, cleanup func(), err error) {
	c, err := connect()
	if err != nil {
		return nil, nil, err
	}

	return c.Database(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		c.Disconnect(ctx)
	}, nil
}

// connect creates a new MongoDB client and establishes a connection.
func connect() (*client, error) {
	cfg := readConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Client().
		SetHosts(strings.Split(cfg.Addrs, ",")).
		SetConnectTimeout(5 * time.Second).
		SetSocketTimeout(5 * time.Second).
		SetAuth(options.Credential{
			AuthSource: cfg.AuthSource,
			Username:   cfg.User,
			Password:   cfg.Password,
		})

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := c.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &client{
		Client: c,
		dbName: cfg.Name,
	}, nil
}
