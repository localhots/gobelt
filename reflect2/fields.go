package reflect2

import "reflect"

// TagIndex returns a map that associates tag values with field indices.
func TagIndex(typ reflect.Type, tag string) map[string]int {
	tagIndex := map[string]int{}
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get(tag)
		if tag != "" {
			tagIndex[tag] = i
		}
	}
	return tagIndex
}

// AssociateColumns returns a map that associates column indices with fields.
func AssociateColumns(typ reflect.Type, tag string, cols []string) map[int]int {
	tagIndex := TagIndex(typ, tag)
	colFields := map[int]int{}
	for i, col := range cols {
		if fi, ok := tagIndex[col]; ok {
			colFields[i] = fi
		}
	}
	return colFields
}
