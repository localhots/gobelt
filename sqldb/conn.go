package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/juju/errors"
)

// Conn represents database connection.
type Conn struct {
	db *sql.DB

	beforeCallbacks []BeforeCallback
	afterCallbacks  []AfterCallback
}

type (
	// BeforeCallback is a kind of function that can be called before a query is
	// executed.
	BeforeCallback func(ctx context.Context, query string)
	// AfterCallback is a kind of function that can be called after a query was
	// executed.
	AfterCallback func(ctx context.Context, query string, took time.Duration, err error)
)

// Flavor defines a kind of SQL database.
type Flavor string

const (
	// MySQL is the MySQL SQL flavor.
	MySQL Flavor = "mysql"
	// PostgreSQL is the PostgreSQL SQL flavor.
	PostgreSQL Flavor = "postgresql"
)

// Connect establishes a new database connection.
func Connect(ctx context.Context, f Flavor, dsn string) (*Conn, error) {
	conn, err := sql.Open(string(f), dsn)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to establish connection")
	}
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "Connection is not responding")
	}
	return &Conn{db: conn}, nil
}

// Exec executes a query and does not expect any result.
func (c *Conn) Exec(ctx context.Context, query string, args ...interface{}) Result {
	c.callBefore(ctx, query)
	startedAt := time.Now()
	res, err := c.db.ExecContext(ctx, query, args...)
	c.callAfter(ctx, query, time.Since(startedAt), err)
	if err != nil {
		return result{err: err}
	}
	return result{res: res}
}

// Query executes a query and returns a result object that can later be used to
// retrieve values.
func (c *Conn) Query(ctx context.Context, query string, args ...interface{}) Result {
	var r result
	c.callBefore(ctx, query)
	startedAt := time.Now()
	r.rows, r.err = c.db.QueryContext(ctx, query, args...)
	c.callAfter(ctx, query, time.Since(startedAt), r.err)
	return r
}

// DB returns the underlying DB object.
func (c *Conn) DB() *sql.DB {
	return c.db
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.db.Close()
}

// Before adds a callback function that would be called before a query is
// executed.
func (c *Conn) Before(cb BeforeCallback) {
	c.beforeCallbacks = append(c.beforeCallbacks, cb)
}

// After adds a callback function that would be called after a query was
// executed.
func (c *Conn) After(cb AfterCallback) {
	c.afterCallbacks = append(c.afterCallbacks, cb)
}

func (c *Conn) callBefore(ctx context.Context, query string) {
	for _, cb := range c.beforeCallbacks {
		cb(ctx, query)
	}
}

func (c *Conn) callAfter(ctx context.Context, query string, took time.Duration, err error) {
	for _, cb := range c.afterCallbacks {
		cb(ctx, query, took, err)
	}
}
