## File cache

Package `file_cache` implements a helper function that would cache results of
an expensive operation to a file. It would use file contents when the file exist
and won't call the function again unless the file is deleted.

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
