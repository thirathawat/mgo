# Package mgo

The `mgo` package provides a MongoDB client for Go applications.

## Usage

Import the `mgo` package into your Go file:

```go
import "github.com/thirathawat/mgo"
```

### Configuration

The `mgo` package reads its configuration from environment variables. Here are the supported environment variables:

- `ADDRS`: A comma-separated list of host:port addresses for the MongoDB servers.
- `NAME`: The name of the MongoDB database.
- `AUTH_SOURCE`: The name of the MongoDB authentication database.
- `USER`: The name of the MongoDB user.
- `PASSWORD`: The password for the MongoDB user.

You can use the `readConfig` function to read the configuration from the environment. It returns a `config` struct that holds the configuration values.

### Creating a MongoDB Client

To create a new MongoDB client and establish a connection, use the `New` function:

```go
db, cleanup, err := mgo.New()
if err != nil {
    // handle error
}
defer cleanup()
```

The `New` function returns a `mongo.Database` object that represents the MongoDB database. It also returns a cleanup function that should be deferred to close the connection when you're done with the database.

### Database Operations

Once you have a `mongo.Database` object, you can perform various operations on the database, such as inserting, updating, and querying documents. Refer to the [MongoDB Go Driver documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo) for more information on working with the MongoDB database.

### Example

Here's an example of how to use the `mgo` package:

```go
package main

import (
    "github.com/tirathawat/mgo"
)

func main() {
    db, cleanup, err := mgo.New()
    if err != nil {
        // handle error
    }
    defer cleanup()

    // Perform database operations...
}
```

## Conclusion

The `mgo` package provides a convenient way to connect to a MongoDB database and perform database operations in your Go applications. It simplifies the process of configuring the client and establishing a connection, allowing you to focus on building your application's functionality.
