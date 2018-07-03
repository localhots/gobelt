package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	type superEmbed struct {
		B int `toml:"b"`
	}
	type embed struct {
		A     int        `toml:"a"`
		Super superEmbed `toml:"super"`
	}
	type everything struct {
		Int    int             `toml:"int"`
		Int64  int64           `toml:"int64"`
		Uint   uint            `toml:"uint"`
		Uint16 uint16          `toml:"uint16"`
		Str    string          `toml:"str"`
		Slice  []string        `toml:"slice"`
		Map    map[string]uint `toml:"map"`
		Embed  embed           `toml:"embed"`
	}
	var e everything
	Require("everything", &e)
	exp := everything{
		Int:    1,
		Int64:  2,
		Uint:   3,
		Uint16: 4,
		Str:    "cranky",
		Slice:  []string{"alpha", "beta"},
		Map: map[string]uint{
			"a": 10,
			"b": 20,
			"c": 30,
		},
		Embed: embed{
			A: 100,
			Super: superEmbed{
				B: 200,
			},
		},
	}

	type extra struct {
		Foo string `toml:"foo"`
	}
	var x extra
	Require("extra", &x)
	eexp := extra{Foo: "bar"}

	err := Load(testConfig)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if !cmp.Equal(exp, e) {
		t.Errorf("Numbers mismatch: %s", cmp.Diff(exp, e))
	}
	if !cmp.Equal(eexp, x) {
		t.Errorf("Words mismatch: %s", cmp.Diff(eexp, x))
	}
}

const testConfig = `
[everything]
int = 1
int64 = 2
uint = 3
uint16 = 4
str = "cranky"
slice = ["alpha", "beta"]
[everything.map]
a = 10
b = 20
c = 30
[everything.embed]
a = 100
[everything.embed.super]
b = 200

[extra]
foo = "bar"

[too.deep]
answer = 42`
