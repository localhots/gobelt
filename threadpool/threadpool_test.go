package threadpool

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

func TestThreadPool(t *testing.T) {
	const n = 100
	var s int64
	ctx := context.Background()

	pool := New(n / 10)
	for i := 0; i < n; i++ {
		pool.Enqueue(ctx, func() { atomic.AddInt64(&s, 1) })
	}
	pool.Close()

	if s != n {
		t.Errorf("Thread pool result doesn't match: expected %d, got %d", n, s)
	}
}

func TestThreadPoolPanicHandling(t *testing.T) {
	logger := &bufLogger{buf: bytes.NewBuffer(nil)}
	ctx := context.Background()

	pool := New(1)
	pool.Logger = logger
	pool.Enqueue(ctx, func() { panic("oh no!") })
	pool.Close()

	out := logger.buf.String()
	exp := "Thread pool task recovered from panic: oh no!"
	if out != exp {
		t.Errorf("Expected logger to receive message %q, got %q", exp, out)
	}
}

type bufLogger struct {
	buf *bytes.Buffer
}

func (b *bufLogger) Printf(f string, args ...interface{}) {
	b.buf.WriteString(fmt.Sprintf(f, args...))
}
