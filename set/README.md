## Set

Package `set` is a collection of packages implementing a set data type. 
Supported types are:

* `int`, `int8`, `int16`, `int32`, `int64` 
* `uint`, `uint8`, `uint16`, `uint32`, `uint64`
* `string`

All the package names are type names prefixed with "set", e.g. `setuint64`.

> Note: These packages are generated from a template. Instead of modifying each
> package individually change the template and run `make gen`.

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