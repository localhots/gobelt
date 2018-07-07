package dbc

import (
	"context"
	"database/sql"
	"time"
)

type connOrTx interface {
	executer
	queryPerformer
}

type executer interface {
	// Exec executes a query and does not expect any result.
	Exec(ctx context.Context, query string, args ...interface{}) ExecResult
	ExecNamed(ctx context.Context, query string, arg interface{}) ExecResult
}

type queryPerformer interface {
	// Query executes a query and returns a result object that can later be used
	// to retrieve values.
	Query(ctx context.Context, query string, args ...interface{}) Rows
	QueryNamed(ctx context.Context, query string, arg interface{}) Rows
}

type stdConnOrTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type caller struct {
	db stdConnOrTx
	cb *callbacks
}

func (c *caller) Exec(ctx context.Context, query string, args ...interface{}) ExecResult {
	c.cb.callBefore(ctx, query)
	startedAt := time.Now()
	res, err := c.db.ExecContext(ctx, query, args...)
	c.cb.callAfter(ctx, query, time.Since(startedAt), err)

	return &execResult{
		db:  c,
		err: err,
		res: res,
	}
}

func (c *caller) ExecNamed(ctx context.Context, query string, arg interface{}) ExecResult {
	return nil
}

func (c *caller) Query(ctx context.Context, query string, args ...interface{}) Rows {
	c.cb.callBefore(ctx, query)
	startedAt := time.Now()
	r, err := c.db.QueryContext(ctx, query, args...)
	c.cb.callAfter(ctx, query, time.Since(startedAt), err)
	return &rows{
		err:  err,
		rows: r,
	}
}

func (c *caller) QueryNamed(ctx context.Context, query string, arg interface{}) Rows {
	return nil
}
