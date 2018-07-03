package csv2

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	type specie struct {
		Name      string  `csv:"name"`
		FavFood   string  `csv:"fav_food"`
		Age       uint16  `csv:"age"`
		Weight    float32 `csv:"weight"`
		Available bool    `csv:"available"`
	}
	body := `name,available,fav_food,weight,age
Alice,false,Bananas,19.22,5
Frank,true,Burrito,14,9
Joel,true,Pesto,32.5,21`
	exp := []specie{
		{Name: "Alice", FavFood: "Bananas", Age: 5, Weight: 19.22, Available: false},
		{Name: "Frank", FavFood: "Burrito", Age: 9, Weight: 14, Available: true},
		{Name: "Joel", FavFood: "Pesto", Age: 21, Weight: 32.5, Available: true},
	}

	csvReader := csv.NewReader(strings.NewReader(body))
	r := NewReader(csvReader)
	r.ColumnNamesInFirstRow = true

	var out []specie
	err := r.Load(&out)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !cmp.Equal(exp, out) {
		t.Errorf("Result value is different: %s", cmp.Diff(exp, out))
	}
}
