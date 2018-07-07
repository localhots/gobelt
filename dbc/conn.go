package dbc

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
)

// Conn represents database connection.
type Conn interface {
	connOrTx
	// Begin executes a transaction.
	Begin(context.Context, func(Tx) error) error
	// BeginCustom executes a transaction with provided options.
	BeginCustom(context.Context, func(Tx) error, *sql.TxOptions) error
	// Close closes the connection.
	Close() error
	// DB returns the underlying DB object.
	DB() *sql.DB
	// Before adds a callback function that would be called before a query is
	// executed.
	Before(BeforeCallback)
	// After adds a callback function that would be called after a query was
	// executed.
	After(AfterCallback)
}

// Tx represents database transacation.
type Tx interface {
	connOrTx
	Commit() error
	Rollback() error
}

type dbWrapper struct {
	conn *sql.DB
	*caller
}

// Flavor defines a kind of SQL database.
type Flavor string

const (
	// MySQL is the MySQL SQL flavor.
	MySQL Flavor = "mysql"
	// PostgreSQL is the PostgreSQL SQL flavor.
	PostgreSQL Flavor = "postgresql"
)

// Connect establishes a new database connection.
func Connect(ctx context.Context, f Flavor, dsn string) (Conn, error) {
	conn, err := sql.Open(string(f), dsn)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to establish connection")
	}
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "Connection is not responding")
	}
	return &dbWrapper{
		conn: conn,
		caller: &caller{
			db: conn,
			cb: &callbacks{},
		},
	}, nil
}

func (c *dbWrapper) Begin(ctx context.Context, fn func(tx Tx) error) error {
	return c.BeginCustom(ctx, fn, nil)
}

func (c *dbWrapper) BeginCustom(ctx context.Context, fn func(tx Tx) error, opts *sql.TxOptions) error {
	if opts == nil {
		opts = &sql.TxOptions{}
	}
	tx, err := c.conn.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	err = fn(c.wrapTx(tx))
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (c *dbWrapper) DB() *sql.DB {
	return c.conn
}

func (c *dbWrapper) Close() error {
	return c.conn.Close()
}

func (c *dbWrapper) Before(cb BeforeCallback) {
	c.cb.addBefore(cb)
}

func (c *dbWrapper) After(cb AfterCallback) {
	c.cb.addAfter(cb)
}

func (c *dbWrapper) wrapTx(tx *sql.Tx) Tx {
	return &txWrapper{
		tx: tx,
		connOrTx: &caller{
			db: tx,
			cb: c.cb,
		},
	}
}

type txWrapper struct {
	tx *sql.Tx
	connOrTx
}

func (w *txWrapper) Commit() error {
	return w.tx.Commit()
}

func (w *txWrapper) Rollback() error {
	return w.tx.Rollback()
}
