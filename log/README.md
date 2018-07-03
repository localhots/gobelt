## Log

Package `log` wraps [logrus](https://github.com/sirupsen/logrus) with a set of
convenient short functions. It can also add values to context for future 
logging.

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

log.Info(ctx, "New user signed up", log.F{
    "id":   user.ID,
    "name": user.Name,
})
// [INFO] New user signed up    email=bob@example.com  id=14  name="Bob Fierce"
```
