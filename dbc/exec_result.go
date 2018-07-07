package dbc

import (
	"context"
	"database/sql"
)

// ExecResult ...
type ExecResult interface {
	Error() error
	LastInsertID() int64
	RowsAffected() int64
	Result() sql.Result
	Then() ExecChain
}

type execResult struct {
	db  executer
	err error
	res sql.Result
}

func (r *execResult) Result() sql.Result {
	return r.res
}

func (r *execResult) Error() error {
	return r.err
}

func (r *execResult) LastInsertID() int64 {
	if r.res == nil {
		return 0
	}
	id, err := r.res.LastInsertId()
	if err != nil {
		return 0
	}
	return id
}

func (r *execResult) RowsAffected() int64 {
	if r.res == nil {
		return 0
	}
	ra, err := r.res.RowsAffected()
	if err != nil {
		return 0
	}
	return ra
}

func (r *execResult) Then() ExecChain {
	if r.err != nil {
		return &execChain{r.db}
	}
	return r.db
}

//
// Chain
//

// ExecChain ...
type ExecChain interface {
	executer
}

type execChain struct {
	executer
}

type brokenChain struct {
	err error
}

func (c *brokenChain) Exec(_ context.Context, _ string, _ ...interface{}) ExecResult {
	return &execResult{err: c.err}
}
