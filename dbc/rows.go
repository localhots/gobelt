package dbc

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/localhots/gobelt/reflect2"
)

// Rows ...
type Rows interface {
	Error() error
	Load(dest interface{}) error
	Rows() *sql.Rows
}

const tagName = "db"

type rows struct {
	err  error
	rows *sql.Rows
}

func (r *rows) Rows() *sql.Rows {
	return r.rows
}

func (r *rows) Error() error {
	return r.err
}

func (r *rows) Load(dest interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.rows.Close()

	dtyp := reflect.TypeOf(dest)
	if dtyp.Kind() != reflect.Ptr {
		panic("Value must be a pointer")
	}
	dtyp = dtyp.Elem()

	switch dtyp.Kind() {
	case reflect.Struct:
		r.loadStruct(dtyp, dest)
	case reflect.Map:
		r.loadMap(dest.(*map[string]interface{}))
	case reflect.Slice:
		switch dtyp.Elem().Kind() {
		case reflect.Struct:
			r.loadSliceOfStructs(dtyp, dest)
		case reflect.Map:
			r.loadSliceOfMaps(dest.(*[]map[string]interface{}))
		default:
			r.loadSlice(dtyp, dest)
		}
	default:
		r.loadValue(dest)
	}

	if r.err == nil && r.rows.Err() != nil {
		return r.rows.Err()
	}

	return r.err
}

func (r *rows) loadValue(dest interface{}) {
	if r.rows.Next() {
		r.err = r.rows.Scan(dest)
	}
}

func (r *rows) loadSlice(typ reflect.Type, dest interface{}) {
	vSlice := reflect.MakeSlice(typ, 0, 0)
	for r.rows.Next() {
		val := reflect.New(typ.Elem())
		r.err = r.rows.Scan(val.Interface())
		if r.err != nil {
			return
		}
		vSlice = reflect.Append(vSlice, val.Elem())
	}
	reflect.ValueOf(dest).Elem().Set(vSlice)
}

func (r *rows) loadMap(dest *map[string]interface{}) {
	if !r.rows.Next() {
		return
	}

	cols, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return
	}
	colTypes, err := r.rows.ColumnTypes()
	if err != nil {
		r.err = err
		return
	}

	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = newValue(colTypes[i])
	}
	err = r.rows.Scan(vals...)
	if err != nil {
		r.err = err
		return
	}

	if *dest == nil {
		*dest = make(map[string]interface{}, len(cols))
	}
	for i, col := range cols {
		switch tval := vals[i].(type) {
		case *int64:
			(*dest)[col] = *tval
		case *string:
			(*dest)[col] = *tval
		case *bool:
			(*dest)[col] = *tval
		}
	}
}

func (r *rows) loadSliceOfMaps(dest *[]map[string]interface{}) {
	cols, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return
	}
	colTypes, err := r.rows.ColumnTypes()
	if err != nil {
		r.err = err
		return
	}

	if *dest == nil {
		*dest = make([]map[string]interface{}, 0)
	}
	for r.rows.Next() {
		vals := make([]interface{}, len(cols))
		for i := range cols {
			vals[i] = newValue(colTypes[i])
		}
		err = r.rows.Scan(vals...)
		if err != nil {
			r.err = err
			return
		}

		row := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			switch tval := vals[i].(type) {
			case *int64:
				row[col] = *tval
			case *string:
				row[col] = *tval
			case *bool:
				row[col] = *tval
			}
		}
		*dest = append(*dest, row)
	}
}

func (r *rows) loadStruct(typ reflect.Type, dest interface{}) {
	if !r.rows.Next() {
		return
	}

	cols, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return
	}

	val := reflect.ValueOf(dest).Elem()
	vals := make([]interface{}, len(cols))
	tm := reflect2.AssociateColumns(val.Type(), tagName, cols)
	for i := range cols {
		if fi, ok := tm[i]; ok {
			fval := val.Field(fi)
			vals[i] = reflect.New(fval.Type()).Interface()
		} else {
			var dummy interface{}
			vals[i] = &dummy
		}
	}

	if r.err = r.rows.Scan(vals...); r.err != nil {
		return
	}

	for i, fi := range tm {
		fval := val.Field(fi)
		fval.Set(reflect.ValueOf(vals[i]).Elem())
	}
}

func (r *rows) loadSliceOfStructs(typ reflect.Type, dest interface{}) {
	cols, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return
	}

	vSlice := reflect.ValueOf(dest).Elem()
	tSlice := vSlice.Type()
	tElem := tSlice.Elem()
	tm := reflect2.AssociateColumns(tElem, tagName, cols)

	for r.rows.Next() {
		vals := make([]interface{}, len(cols))
		val := reflect.New(tElem).Elem()
		for i := range cols {
			if fi, ok := tm[i]; ok {
				fval := val.Field(fi)
				vals[i] = reflect.New(fval.Type()).Interface()
			} else {
				vals[i] = nopScanner{}
			}
		}

		if r.err = r.rows.Scan(vals...); r.err != nil {
			return
		}

		for i, fi := range tm {
			fval := val.Field(fi)
			fval.Set(reflect.ValueOf(vals[i]).Elem())
		}
		vSlice.Set(reflect.Append(vSlice, val))
	}
}

func (r *rows) withError(err error) Rows {
	r.err = err
	return r
}

func newValue(typ *sql.ColumnType) interface{} {
	switch typ.DatabaseTypeName() {
	case "VARCHAR", "NVARCHAR", "TEXT":
		var s string
		return &s
	case "INT", "BIGINT":
		var i int64
		return &i
	case "BOOL":
		var b bool
		return &b
	default:
		panic(fmt.Errorf("Unsupported MySQL type: %s", typ.DatabaseTypeName()))
	}
}

type nopScanner struct{}

func (s *nopScanner) Scan(interface{}) error { return nil }
