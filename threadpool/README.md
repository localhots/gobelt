## Thread pool

Package `thread_pool` implements a pool of threads, duh.

```go
import "github.com/localhots/gobelt/threadpool"
```

```go
ctx := context.Background()
ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
defer cancel()

pool := threadpool.New(10)
defer pool.Close()
for i := 0; i < 1000000; i++ {
    i := i
    pool.Enqueue(ctx, func() {
        fmt.Printf("The number is %d\n", i)
    })
}
```
