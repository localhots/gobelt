package csv2

import (
	"encoding/csv"
	"io"
	"reflect"
	"strconv"

	"github.com/juju/errors"
	"github.com/localhots/gobelt/reflect2"
)

// Reader is a wrapper for standard library CSV Reader that allows unmarshalling
// into a slice of structs.
type Reader struct {
	TagName               string
	ColumnNamesInFirstRow bool
	Reader                *csv.Reader
	cols                  []string
}

const defaultTagName = "csv"

// NewReader creates a new reader from a standard CSV Reader.
func NewReader(r *csv.Reader) *Reader {
	return &Reader{
		TagName:               defaultTagName,
		ColumnNamesInFirstRow: true,
		Reader:                r,
	}
}

// SetColumnNames assigns column names. Use this function if column names are
// not provided in the first row.
func (r *Reader) SetColumnNames(cols []string) {
	if r.ColumnNamesInFirstRow {
		panic("Should not assign column names when they are expected in the first row")
	}
	r.cols = cols
}

// Load reads CSV contents and unmarshals them into given destination.
// Destination value must be a pointer to a slice of structs.
func (r *Reader) Load(dest interface{}) error {
	destT := reflect.TypeOf(dest)
	if destT.Kind() != reflect.Ptr ||
		destT.Elem().Kind() != reflect.Slice ||
		destT.Elem().Elem().Kind() != reflect.Struct {
		return errors.New("Destination must be a pointer to a slice of structs")
	}

	if r.cols == nil {
		if r.ColumnNamesInFirstRow {
			cols, err := r.Reader.Read()
			if err != nil {
				return errors.Annotate(err, "Failed to read column names from first row")
			}
			r.cols = cols
		} else {
			return errors.New("Column names are not defined")
		}
	}

	destV := reflect.ValueOf(dest).Elem()
	valT := destT.Elem().Elem()
	colIndex := reflect2.AssociateColumns(valT, r.TagName, r.cols)

	for {
		row, err := r.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Annotate(err, "Failed to read CSV row")
		}

		val := reflect.New(valT).Elem()
		for iCol, iField := range colIndex {
			err := unmarshal(row[iCol], val.Field(iField))
			if err != nil {
				return errors.Annotate(err, "Failed to process CSV row")
			}
		}

		destV.Set(reflect.Append(destV, val))
	}
	return nil
}

func unmarshal(v string, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.String:
		dest.SetString(v)
	case reflect.Bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return unmarshalError(err, v, dest.Kind())
		}
		dest.SetBool(b)
	case reflect.Int, reflect.Int64:
		return unmarshalInt(v, dest, 64)
	case reflect.Int8:
		return unmarshalInt(v, dest, 8)
	case reflect.Int16:
		return unmarshalInt(v, dest, 16)
	case reflect.Int32:
		return unmarshalInt(v, dest, 32)
	case reflect.Uint, reflect.Uint64:
		return unmarshalUint(v, dest, 64)
	case reflect.Uint8:
		return unmarshalUint(v, dest, 8)
	case reflect.Uint16:
		return unmarshalUint(v, dest, 16)
	case reflect.Uint32:
		return unmarshalUint(v, dest, 32)
	case reflect.Float32:
		return unmarshalFloat(v, dest, 32)
	case reflect.Float64:
		return unmarshalFloat(v, dest, 64)
	}

	return nil
}

func unmarshalInt(v string, dest reflect.Value, bitSize int) error {
	i, err := strconv.ParseInt(v, 10, bitSize)
	if err != nil {
		return unmarshalError(err, v, dest.Kind())
	}
	dest.SetInt(i)
	return nil
}

func unmarshalUint(v string, dest reflect.Value, bitSize int) error {
	i, err := strconv.ParseUint(v, 10, bitSize)
	if err != nil {
		return unmarshalError(err, v, dest.Kind())
	}
	dest.SetUint(i)
	return nil
}

func unmarshalFloat(v string, dest reflect.Value, bitSize int) error {
	f, err := strconv.ParseFloat(v, bitSize)
	if err != nil {
		return unmarshalError(err, v, dest.Kind())
	}
	dest.SetFloat(f)
	return nil
}

func unmarshalError(err error, v string, k reflect.Kind) error {
	return errors.Annotatef(err, "Can't unmarshal %q into value of type %s", v, k)
}
