package dbc

import (
	"context"
	"flag"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/localhots/gobelt/log"
)

type record struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

var conn Conn

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn := flag.String("dsn", "", "Database source name")
	flag.Parse()
	if *dsn == "" {
		log.Warn(ctx, "Database source name is not provided, some tests would be skipped")
	} else {
		connect(ctx, *dsn)
		seed(ctx)
	}

	log.Info(ctx, "Starting test suite")
	exitCode := m.Run()
	if exitCode == 0 {
		log.Info(ctx, "Test suite completed successfully")
	} else {
		log.Error(ctx, "Test suite failed")
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Errorf(ctx, "Failed to close connection: %v\n", err)
		}
	}
	os.Exit(exitCode)
}

func connect(ctx context.Context, dsn string) {
	log.Info(ctx, "Establishing connection to the test database")
	var err error
	conn, err = Connect(ctx, MySQL, dsn)
	if err != nil {
		log.Fatalf(ctx, "Failed to connect: %v\n", err)
	}
}

func seed(ctx context.Context) {
	log.Info(ctx, "Seeding database")
	mustExecMain(conn.Exec(ctx, `DROP TABLE IF EXISTS sqldb_test`))
	mustExecMain(conn.Exec(ctx, `
		CREATE TABLE sqldb_test (
			id int(11) UNSIGNED NOT NULL,
			name VARCHAR(10) DEFAULT '',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=ascii
	`))
	mustExecMain(conn.Exec(ctx, `
		INSERT INTO sqldb_test (id, name) 
		VALUES
			(1, "Alice"),
			(2, "Bob")
	`))
}

func mustExec(t *testing.T, r ExecResult) ExecResult {
	t.Helper()
	if r.Error() != nil {
		t.Fatalf("Exec failed: %v", r.Error())
	}
	return r
}

func mustExecMain(r ExecResult) ExecResult {
	if r.Error() != nil {
		log.Fatalf(context.Background(), "Query failed: %v\n", r.Error())
	}
	return r
}

func mustQuery(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
}

func requireConn(t *testing.T) {
	t.Helper()
	if conn == nil {
		t.Skip("Connection required")
	}
}
