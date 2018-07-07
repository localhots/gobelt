package main

import (
	"context"
	"flag"

	"github.com/localhots/gobelt/dbc"
	"github.com/localhots/gobelt/log"
)

func main() {
	dsn := flag.String("dsn", "", "Database source name")
	flag.Parse()

	ctx := context.Background()
	conn, err := dbc.Connect(ctx, dbc.MySQL, *dsn)
	if err != nil {
		log.Fatal(ctx, "Failed to establish database conneciton", log.F{
			"dsn":   dsn,
			"error": err,
		})
	}

	_ = conn
}
