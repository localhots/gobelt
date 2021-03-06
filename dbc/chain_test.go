package dbc

import (
	"context"
	"testing"
)

func TestCallChain(t *testing.T) {
	requireConn(t)
	ctx := context.Background()
	mustExec(t, conn.
		Exec(ctx, "INSERT INTO sqldb_test (id, name) VALUES (3, 'Fred')").Then().
		Exec(ctx, "UPDATE sqldb_test SET name = 'Wilson' WHERE id = 3").Then().
		Exec(ctx, "DELETE FROM sqldb_test WHERE id = 3"),
	)
}
