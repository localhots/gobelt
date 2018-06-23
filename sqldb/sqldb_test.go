package sqldb

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type record struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

var conn *Conn

func TestMain(m *testing.M) {
	dsn := flag.String("dsn", "", "Database source name")
	flag.Parse()
	if *dsn == "" {
		log.Println("Database source name is not provided, skipping package tests")
		os.Exit(0)
	}

	log.Println("Establishing connection to the test database")
	ctx := context.Background()
	var err error
	conn, err = Connect(ctx, MySQL, *dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	log.Println("Seeding database")
	must(conn.Exec(ctx, `
		DROP TABLE IF EXISTS sqldb_test
	`))
	must(conn.Exec(ctx, `
		CREATE TABLE sqldb_test (
			id int(11) UNSIGNED NOT NULL,
			name VARCHAR(10) DEFAULT '',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=ascii
	`))
	must(conn.Exec(ctx, `
		INSERT INTO sqldb_test (id, name) 
		VALUES
			(1, "Alice"),
			(2, "Bob")
	`))

	fmt.Println("Starting test suite")
	exitCode := m.Run()
	log.Println("Test suite finished")
	if err := conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v\n", err)
	}
	os.Exit(exitCode)
}

func mustT(t *testing.T, r Result) Result {
	t.Helper()
	if r.Error() != nil {
		t.Fatalf("Query failed: %v", r.Error())
	}
	return r
}

func must(r Result) Result {
	if r.Error() != nil {
		log.Fatalf("Query failed: %v\n", r.Error())
	}
	return r
}
