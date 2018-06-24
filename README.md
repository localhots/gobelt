# Gobelt

Gobelt is a collection of Go tools.

### Thread pool

```go
import "github.com/localhots/gobelt/threadpool"
```

```go
ctx := context.Background()
ctx, cancel = context.WithTimeout(ctx, 30 * time.Second)
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

### File cache

```go
import "github.com/localhots/gobelt/filecache"
```

```go
var val int
filecache.Load(&val, "path/to/cachefile", func() interface{} {
    var items []Item
    err := conn.Query(ctx, "SELECT * FROM items").Load(&items).Error()
    if err != nil {
        log.Fatal("Failed to load items", log.F{"error": err})
    }
    return items
})
```

### Log

```go
import "github.com/localhots/gobelt/log"
```

```go
ctx := context.Background()
ctx = log.ContextWithFields(ctx, log.F{"email": params["email"]})

user, err := signup(ctx, params)
if err != nil {
    log.Errorf(ctx, "Signup failed: %v", err)
    // [ERRO] Signup failed: db: duplicate entry    email=bob@example.com
    return
}

log.Info(ctx, "New user signed up", log.F{"id": user.ID})
// [INFO] New user signed up    email=bob@example.com  id=14 
```