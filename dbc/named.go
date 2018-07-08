package dbc

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/localhots/gobelt/reflect2"
)

var namedRegexp = regexp.MustCompile("" +
	"`[^`]+`|" +
	`'[^']+'|` +
	`"[^"]+"|` +
	`@[a-zA-Z][a-zA-Z0-9_]*`)

func prepareNamedQuery(query string, p namedParams) (newQuery string, args []interface{}, err error) {
	newQuery = namedRegexp.ReplaceAllStringFunc(query, func(m string) string {
		if !strings.HasPrefix(m, "@") {
			return m
		}
		val, ok := p.Get(m[1:])
		if !ok {
			err = fmt.Errorf("Named parameter %s was not found", m)
		}
		args = append(args, val)
		return "?"
	})
	return
}

//
// Params
//

type namedParams interface {
	Get(name string) (val interface{}, ok bool)
}

var (
	_ namedParams = namedParamsMap{}
	_ namedParams = &namedParamsStruct{}
)

func newNamedParams(val interface{}) (namedParams, error) {
	switch tval := val.(type) {
	case map[string]interface{}:
		return newNamedParamsMap(tval)
	default:
		return newNamedParamsStruct(val)
	}
}

type namedParamsMap struct {
	m map[string]interface{}
}

func newNamedParamsMap(m map[string]interface{}) (namedParamsMap, error) {
	return namedParamsMap{m}, nil
}

func (p namedParamsMap) Get(name string) (val interface{}, ok bool) {
	val, ok = p.m[name]
	return
}

type namedParamsStruct struct {
	s   reflect.Value
	idx map[string]int
}

func newNamedParamsStruct(s interface{}) (*namedParamsStruct, error) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Unsupported named parameters type: %T", s)
	}
	return &namedParamsStruct{
		s:   val,
		idx: reflect2.TagIndex(val.Type(), tagName),
	}, nil
}

func (p *namedParamsStruct) Get(name string) (val interface{}, ok bool) {
	if i, ok := p.idx[name]; ok {
		return p.s.Field(i).Interface(), true
	}
	return nil, false
}
