package dbc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrepareNamedQuery(t *testing.T) {
	q := `SELECT id, "@not_param ", '@name', "i", ` + "`password` " +
		`FROM tbl WHERE name = @name, active = @is_active`
	p, err := newNamedParamsMap(map[string]interface{}{"name": "Bob", "is_active": 1})
	if err != nil {
		t.Fatalf("Failed to create named params map: %v", err)
	}
	q, args, err := prepareNamedQuery(q, p)
	if err != nil {
		t.Fatalf("Failed to prepare named statement: %v", err)
	}
	const expQ = `SELECT id, "@not_param ", '@name', "i", ` + "`password` " +
		`FROM tbl WHERE name = ?, active = ?`
	if q != expQ {
		t.Errorf("Expected query to be\n%s\ngot\n%q", expQ, q)
	}
	expA := []interface{}{"Bob", 1}
	if !cmp.Equal(expA, args) {
		t.Errorf("Returned arguments are different: %s", cmp.Diff(expA, args))
	}
}

func TestNamedParamsMap(t *testing.T) {
	m, err := newNamedParamsMap(map[string]interface{}{
		"num": 1,
		"str": "foo",
	})
	if err != nil {
		t.Fatalf("Failed to create named params map: %v", err)
	}
	testNamedParams(t, m)
}

func TestNamedParamsStruct(t *testing.T) {
	type dummy struct {
		Num int    `db:"num"`
		Str string `db:"str"`
	}
	m, err := newNamedParamsStruct(&dummy{Num: 1, Str: "foo"})
	if err != nil {
		t.Fatalf("Failed to create named params struct: %v", err)
	}
	testNamedParams(t, m)
}

func testNamedParams(t *testing.T, p namedParams) {
	t.Helper()

	v, ok := p.Get("num")
	if !ok {
		t.Error("num was not found")
	}
	if v != 1 {
		t.Errorf("Expected num to equal 1, got %d", v)
	}

	v, ok = p.Get("str")
	if !ok {
		t.Error("num was not found")
	}
	if v != "foo" {
		t.Errorf("Expected num to equal 'foo', got %q", v)
	}

	_, ok = p.Get("missing")
	if ok {
		t.Error("missing value reportedly found")
	}
}

func BenchmarkNamedParamsMap(b *testing.B) {
	m, err := newNamedParamsMap(map[string]interface{}{"foo": 1})
	if err != nil {
		b.Fatalf("Failed to create named params map: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("foo")
	}
}

func BenchmarkNamedParamsStruct(b *testing.B) {
	type dummy struct {
		Foo int `db:"foo"`
	}
	m, err := newNamedParamsStruct(&dummy{Foo: 1})
	if err != nil {
		b.Fatalf("Failed to create named params struct: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("foo")
	}
}

func BenchmarkPrepareNamedOne(b *testing.B) {
	p, _ := newNamedParamsMap(map[string]interface{}{
		"name": "Bob",
	})
	const q = `SELECT * FROM tbl WHER name = @name`

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prepareNamedQuery(q, p)
	}
}

func BenchmarkPrepareNamedQuerySix(b *testing.B) {
	p, _ := newNamedParamsMap(map[string]interface{}{
		"name":      "Bob",
		"is_active": 1,
		"amount":    123.45,
	})
	const q = `SELECT "aaa @false1 bbb" as f1, 'ccc @false2 ddd' as f2, ` +
		"`eee @false3 fff` as f3, name, is_active, amount FROM tbl " +
		`WHERE name = @name AND is_active = @is_active AND amount = @amount`

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prepareNamedQuery(q, p)
	}
}
