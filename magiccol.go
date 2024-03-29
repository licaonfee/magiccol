// Package magiccol allows to scan sql rows into an arbitrary map
package magiccol

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

// ErrNilRows a nil Rows interface is provided
var ErrNilRows = errors.New("nil *sql.Rows as argument")

// ErrInvalidDataType columnd declared with a type incompatible with the actual value
var ErrInvalidDataType = errors.New("data type not valid for column")

// Scanner read data from an sql.Rows object into a map
type Scanner struct {
	o        Options
	columns  []string
	pointers []any
	values   []reflect.Value
	err      error
}

// Options for Scanner
type Options struct {
	// Rows must be a valid sql.Rows object
	Rows Rows
	// Mapper can be nil, if so DefaultMapper is used
	Mapper *Mapper
}

// Rows allow to mock sql.Rows object
type Rows interface {
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Next() bool
	Err() error
	Scan(...any) error
}

// NewScanner create a new Scanner object, return an error if
// a nil Rows interface is provided or any error is returned by its
func NewScanner(o Options) (*Scanner, error) {
	if o.Rows == nil {
		return nil, ErrNilRows
	}
	if o.Mapper == nil {
		o.Mapper = DefaultMapper()
	}
	tp, err := o.Rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	cols, err := o.Rows.Columns()
	if err != nil {
		return nil, err
	}
	pointers := make([]any, len(cols))
	values := make([]reflect.Value, len(cols))
	for i := 0; i < len(cols); i++ {
		t := tp[i]
		refType := o.Mapper.Get(t)
		v := reflect.New(refType)
		pointers[i] = v.Interface()
		values[i] = v
	}

	return &Scanner{columns: cols, pointers: pointers, values: values, o: o}, nil
}

// Scan return true if there are rows in queue and false if
// there is no more rows or an error occurred. To distinguish
// between error or no more rows Err() method should be consulted
func (s *Scanner) Scan() bool {
	if !s.o.Rows.Next() {
		if s.o.Rows.Err() != nil {
			s.err = s.o.Rows.Err()
		}
		return false
	}
	if err := s.o.Rows.Scan(s.pointers...); err != nil {
		s.err = fmt.Errorf("%w %s", ErrInvalidDataType, err)
		return false
	}
	return true
}

// SetMap read values from current row and load it in a given map[string]any
// this allow to set default values, or reutilize same map in multiple iterations
// SetMap does not clear map object and any preexistent key will be preserved
func (s *Scanner) SetMap(value map[string]any) {
	for i := 0; i < len(s.columns); i++ {
		v := reflect.Indirect(s.values[i])
		for v.Kind() == reflect.Ptr && !v.IsNil() {
			v = reflect.Indirect(v)
		}
		value[s.columns[i]] = v.Interface()
	}
}

// Value returns a new map object with all values from current row
// successives calls to Value without call Scan returns always same values
// in a new allocated map. Call Value() before Scan return all values as Zero
func (s *Scanner) Value() map[string]any {
	value := make(map[string]any, len(s.columns))
	s.SetMap(value)
	return value
}

// Err return last error in Scanner
func (s *Scanner) Err() error {
	return s.err
}
