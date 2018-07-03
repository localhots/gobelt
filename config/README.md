## Config

Package config allows other packages to require a certain section of TOML 
configuration file to be parsed into its internal configuration structure. Once
config file is processed its values are distributed by the packages that 
required them.

```go
import "github.com/localhots/gobelt/config"
```

Describe configuration structure inside a target package then call 
`config.Require` from the init function of a package.

> Note: subgroups are not currently supported. If a package is called `s3` and
> located in `aws` parent package, call config section `aws_s3` instead of 
> `aws.s3`. 

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

Load configuration from a `main` function of the app:

```go
package main

func main() {
    config.Load("config/config.toml")
}
```
