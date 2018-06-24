# Gobelt

Gobelt is a collection of Go tools.

### Thread pool

```go
import "github.com/localhots/gobelt/threadpool"
```

```go
ctx := context.Background()
pool := threadpool.New(5)
pool.Enqueue(ctx, func() {
    fmt.Println("Hello")
})
pool.Close()
```

### File cache

```go
import "github.com/localhots/gobelt/filecache"
```

```go
var val int
filecache.Load(&val, "path/to/cachefile", func() interface{} {
    // Expensive calls here
    return 100
})
```

### Log

```go
import "github.com/localhots/gobelt/log"
```

```go
log.Info("New user signed up", log.F{
    "name":  u.Name,
    "email": u.Email,
})
log.Errorf("Invalid database flavor: %s", flavor)
```