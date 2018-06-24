package sqldb

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
)

// Conn represents database connection.
type Conn interface {
	connOrTx
	Begin(context.Context, func(Tx) error) error
	Close() error
	DB() *sql.DB
	Before(BeforeCallback)
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

// Begin executes a transaction.
func (c *dbWrapper) Begin(ctx context.Context, fn func(tx Tx) error) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	err = fn(c.wrapTx(tx))
	if err != nil {
		tx.Rollback()
	}
	return err
}

// DB returns the underlying DB object.
func (c *dbWrapper) DB() *sql.DB {
	return c.conn
}

// Close closes the connection.
func (c *dbWrapper) Close() error {
	return c.conn.Close()
}

// Before adds a callback function that would be called before a query is
// executed.
func (c *dbWrapper) Before(cb BeforeCallback) {
	c.cb.addBefore(cb)
}

// After adds a callback function that would be called after a query was
// executed.
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
