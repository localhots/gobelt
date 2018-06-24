# Gobelt

Gobelt is a collection of Go tools.

### Thread pool

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

### File cache

```go
import "github.com/localhots/gobelt/filecache"
```

```go
ctx := context.Background()
var items []Item
filecache.Load(&items, "tmp/cache/items.json", func() {  
    err := conn.Query(ctx, "SELECT * FROM items").Load(&items).Error()
    if err != nil {
        log.Fatalf("Failed to load items: %v", err)
    }
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

### Set

There is a collection of packages implementing a set data type located inside
the `set` package. Implemented types include:

* `int`, `int8`, `int16`, `int32`, `int64` 
* `uint`, `uint8`, `uint16`, `uint32`, `uint64`
* `string`

All the package names are type names prefixed with "set", e.g. `setuint64`.

```go
import "github.com/localhots/gobelt/set/setstring"
```

```go
s := setstring.New("one", "two")
s.Add("three")
s.Remove("one", "two").Add("four", "five")
fmt.Println("Size:", s.Len()) // 3
fmt.Println("Has one", s.Has("one")) // false
fmt.Println(s.SortedSlice()) // [three four five]
```

### Config

```go
import "github.com/localhots/gobelt/config"
```

Describe configuration structure inside a target package.

```go
package db

var conf struct {
    Flavor string `toml:"flavor"`
    DSN    string `toml:"dsn"`
}

func init() {
    config.Require("db", &conf)
}
```

Load configuration from a `main` function:

```go
package main

func main() {
    config.Load("config/config.toml")
}
```