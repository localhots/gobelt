package sqldb

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadSingleValue(t *testing.T) {
	ctx := context.Background()
	exp := int(1)
	var out int
	mustT(t, conn.Query(ctx, "SELECT 1").Load(&out))
	if exp != out {
		t.Errorf("Value doesn't match: expected %d, got %d", exp, out)
	}
}

func TestLoadSlice(t *testing.T) {
	ctx := context.Background()
	exp := []int{1, 2}
	var out []int
	mustT(t, conn.Query(ctx, "SELECT id FROM sqldb_test").Load(&out))
	if !cmp.Equal(exp, out) {
		t.Errorf("Values dont't match: %s", cmp.Diff(exp, out))
	}
}

func TestLoadMap(t *testing.T) {
	ctx := context.Background()
	exp := map[string]interface{}{"id": int64(1), "name": "Alice"}
	var out map[string]interface{}
	mustT(t, conn.Query(ctx, "SELECT * FROM sqldb_test WHERE id = 1").Load(&out))
	if !cmp.Equal(exp, out) {
		t.Errorf("Record doesn't match: %s", cmp.Diff(exp, out))
	}
}

func TestLoadSliceOfMaps(t *testing.T) {
	ctx := context.Background()
	exp := []map[string]interface{}{
		{"id": int64(1), "name": "Alice"},
		{"id": int64(2), "name": "Bob"},
	}
	var out []map[string]interface{}
	mustT(t, conn.Query(ctx, "SELECT * FROM sqldb_test ORDER BY id ASC").Load(&out))
	if !cmp.Equal(exp, out) {
		t.Errorf("Records don't match: %s", cmp.Diff(exp, out))
	}
}

func TestLoadStruct(t *testing.T) {
	ctx := context.Background()
	exp := record{ID: 1, Name: "Alice"}
	var out record
	mustT(t, conn.Query(ctx, "SELECT * FROM sqldb_test WHERE id = 1").Load(&out))
	if !cmp.Equal(exp, out) {
		t.Errorf("Record doesn't match: %s", cmp.Diff(exp, out))
	}
}

func TestLoadSliceOfStructs(t *testing.T) {
	ctx := context.Background()
	exp := []record{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	var out []record
	mustT(t, conn.Query(ctx, "SELECT * FROM sqldb_test").Load(&out))
	if !cmp.Equal(exp, out) {
		t.Errorf("Records don't match: %s", cmp.Diff(exp, out))
	}
}
