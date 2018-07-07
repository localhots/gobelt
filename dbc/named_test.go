package dbc

import "testing"

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
