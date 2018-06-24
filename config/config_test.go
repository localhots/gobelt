package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	type numbers struct {
		Ary []int `toml:"ary"`
		I   int   `toml:"i"`
	}
	var nums numbers
	Require("numbers", &nums)
	numsExp := numbers{
		Ary: []int{1, 2, 3},
		I:   123,
	}

	type words struct {
		Foo  string `toml:"foo"`
		List struct {
			Foo []string `toml:"foo"`
		} `toml:"list"`
	}
	var w words
	Require("words", &w)
	wordsExp := words{
		Foo: "bar",
		List: struct {
			Foo []string `toml:"foo"`
		}{
			Foo: []string{"buzz", "fizz"},
		},
	}

	err := Load("example.toml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if !cmp.Equal(numsExp, nums) {
		t.Errorf("Numbers mismatch: %s", cmp.Diff(numsExp, nums))
	}
	if !cmp.Equal(wordsExp, w) {
		t.Errorf("Words mismatch: %s", cmp.Diff(wordsExp, w))
	}
}
