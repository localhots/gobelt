package threadpool

import (
	"context"
	"log"
	"sync"
)

// ThreadPool implements a thread pool model. It allocates a pool of threads
// ready to perform tasks concurrently.
type ThreadPool struct {
	Logger interface {
		Printf(f string, args ...interface{})
	}
	queue chan func()
	wg    sync.WaitGroup
}

// New creates a thread pool with a given number of workers.
func New(size int) *ThreadPool {
	tp := &ThreadPool{
		Logger: &log.Logger{},
		queue:  make(chan func(), size),
	}

	tp.wg.Add(size)
	for i := 0; i < size; i++ {
		go tp.worker()
	}

	return tp
}

// Enqueue adds a task to queue.
func (tp *ThreadPool) Enqueue(ctx context.Context, task func()) {
	select {
	case tp.queue <- task:
	case <-ctx.Done():
	}
}

// Close waits for all currently accepted tasks to be processed and returns.
// Attempts to enqueue a task after calling Close would result in a panic.
func (tp *ThreadPool) Close() {
	close(tp.queue)
	tp.wg.Wait()
}

func (tp *ThreadPool) worker() {
	defer tp.wg.Done()
	for task := range tp.queue {
		tp.perform(task)
	}
}

func (tp *ThreadPool) perform(task func()) {
	defer func() {
		if err := recover(); err != nil {
			tp.Logger.Printf("Thread pool task recovered from panic: %v", err)
		}
	}()
	task()
}
