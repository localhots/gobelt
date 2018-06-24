package sqldb

import (
	"context"
	"time"
)

type (
	// BeforeCallback is a kind of function that can be called before a query is
	// executed.
	BeforeCallback func(ctx context.Context, query string)
	// AfterCallback is a kind of function that can be called after a query was
	// executed.
	AfterCallback func(ctx context.Context, query string, took time.Duration, err error)
)

type callbacks struct {
	before []BeforeCallback
	after  []AfterCallback
}

func (c *callbacks) addBefore(cb BeforeCallback) {
	c.before = append(c.before, cb)
}

func (c *callbacks) addAfter(cb AfterCallback) {
	c.after = append(c.after, cb)
}

func (c *callbacks) callBefore(ctx context.Context, query string) {
	for _, cb := range c.before {
		cb(ctx, query)
	}
}

func (c *callbacks) callAfter(ctx context.Context, query string, took time.Duration, err error) {
	for _, cb := range c.after {
		cb(ctx, query, took, err)
	}
}
