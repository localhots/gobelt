package sqldb

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Result represents query result.
type Result interface {
	// Load decodes rows into provided variable.
	Load(dest interface{}) Result
	// Error returns an error if one happened during query execution.
	Error() error
	// Rows returns original database rows object.
	Rows() *sql.Rows
	// LastInsertID returns the last inserted record ID for results obtained
	// from Exec calls.
	LastInsertID() int64
	// RowsAffected returns the number of rows affected for results obtained
	// from Exec calls.
	RowsAffected() int64
}

type result struct {
	rows *sql.Rows
	res  sql.Result
	err  error
}

func (r result) Rows() *sql.Rows {
	return r.rows
}

func (r result) Error() error {
	return r.err
}

func (r result) LastInsertID() int64 {
	if r.res == nil {
		return 0
	}
	id, err := r.res.LastInsertId()
	if err != nil {
		return 0
	}
	return id
}

func (r result) RowsAffected() int64 {
	if r.res == nil {
		return 0
	}
	ra, err := r.res.RowsAffected()
	if err != nil {
		return 0
	}
	return ra
}

func (r result) Load(dest interface{}) Result {
	if r.err != nil {
		return r
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
		return r.withError(r.rows.Err())
	}

	return r
}

func (r *result) loadValue(dest interface{}) {
	if r.rows.Next() {
		r.err = r.rows.Scan(dest)
	}
}

func (r *result) loadSlice(typ reflect.Type, dest interface{}) {
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

func (r *result) loadMap(dest *map[string]interface{}) {
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

func (r *result) loadSliceOfMaps(dest *[]map[string]interface{}) {
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

func (r *result) loadStruct(typ reflect.Type, dest interface{}) {
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
	tm := tagMap(cols, val.Type())
	for i := range cols {
		if fi, ok := tm[i]; ok {
			fval := val.Field(fi)
			vals[i] = reflect.New(fval.Type()).Interface()
		} else {
			var dummy interface{}
			vals[i] = &dummy
		}
	}

	err = r.rows.Scan(vals...)
	if err != nil {
		r.err = err
		return
	}

	for i, fi := range tm {
		fval := val.Field(fi)
		fval.Set(reflect.ValueOf(vals[i]).Elem())
	}
}

func (r *result) loadSliceOfStructs(typ reflect.Type, dest interface{}) {
	cols, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return
	}

	vSlice := reflect.ValueOf(dest).Elem()
	tSlice := vSlice.Type()
	tElem := tSlice.Elem()
	tm := tagMap(cols, tElem)

	for r.rows.Next() {
		vals := make([]interface{}, len(cols))
		val := reflect.New(tElem).Elem()
		for i := range cols {
			if fi, ok := tm[i]; ok {
				fval := val.Field(fi)
				vals[i] = reflect.New(fval.Type()).Interface()
			} else {
				var dummy interface{}
				vals[i] = &dummy
			}
		}

		err = r.rows.Scan(vals...)
		if err != nil {
			r.err = err
			return
		}

		for i, fi := range tm {
			fval := val.Field(fi)
			fval.Set(reflect.ValueOf(vals[i]).Elem())
		}
		vSlice.Set(reflect.Append(vSlice, val))
	}
}

func (r result) withError(err error) Result {
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

func tagMap(cols []string, typ reflect.Type) map[int]int {
	fieldIndices := map[string]int{}
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag != "" {
			fieldIndices[tag] = i
		}
	}

	colFields := map[int]int{}
	for i, col := range cols {
		if fi, ok := fieldIndices[col]; ok {
			colFields[i] = fi
		}
	}

	return colFields
}
