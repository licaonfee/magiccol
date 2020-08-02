package magiccol

import (
	"database/sql"
	"errors"
	"reflect"
)

var ErrNilRows = errors.New("nil *sql.Rows as argument")

type Scanner struct {
	o        Options
	columns  []string
	pointers []interface{}
	values   []reflect.Value
	err      error
}

type Options struct {
	//Rows must be a valid sql.Rows object
	Rows *sql.Rows
	//Mapper can be nil, if so DefaultMapper is used
	Mapper Mapper
}

//Inspect set scanner to read a column set
func NewScanner(o Options) (*Scanner, error) {
	if o.Rows == nil {
		return nil, ErrNilRows
	}
	tp, err := o.Rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	cols, err := o.Rows.Columns()
	if err != nil {
		return nil, err
	}
	pointers := make([]interface{}, len(cols))
	values := make([]reflect.Value, len(cols))
	for i := 0; i < len(cols); i++ {
		t := tp[i]
		refType := o.Mapper.Get(t.DatabaseTypeName(), t.ScanType())
		v := reflect.New(refType)
		pointers[i] = v.Interface()
		values[i] = v
	}

	return &Scanner{columns: cols, pointers: pointers, values: values, o: o}, nil
}

//Scan return true if there are rows in queue and false if
//there is no more rows or an error occured. To distinguish
//between error or no more rows Err() method should be consulted
func (s *Scanner) Scan() bool {
	if !s.o.Rows.Next() {
		if s.o.Rows.Err() != nil {
			s.err = s.o.Rows.Err()
		}
		return false
	}
	if err := s.o.Rows.Scan(s.pointers...); err != nil {
		s.err = err
		return false
	}
	return true
}

//SetMap read values from current row and load it in a given map[string]interface{}
//this allow to set default values, or reutilize same map in multiple iterations
//SetMap does not clear map object and any preexistent key will be preserved
func (s *Scanner) SetMap(value map[string]interface{}) {
	for i := 0; i < len(s.columns); i++ {
		value[s.columns[i]] = reflect.Indirect(s.values[i]).Interface()
	}
}

//Value returns a new map object with all values from current row
//successives calls to Value without call Scan returns always same values
//in a new allocated map. Call Value() before Scan return all values as Zero
func (s *Scanner) Value() map[string]interface{} {
	value := make(map[string]interface{}, len(s.columns))
	s.SetMap(value)
	return value
}

//Err return last error in Scanner
func (s *Scanner) Err() error {
	return s.err
}
