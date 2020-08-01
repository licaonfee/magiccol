package magiccol

import (
	"database/sql"
	"reflect"
)

type Scanner struct {
	o        Options
	columns  []string
	pointers []interface{}
	values   []reflect.Value
	m        map[string]interface{}
	err      error
}

type Options struct {
	Rows   *sql.Rows
	Mapper Mapper
}

//Inspect set scanner to read a column set
func NewScanner(o Options) (*Scanner, error) {
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

func (s *Scanner) SetMap(value map[string]interface{}) {
	for i := 0; i < len(s.columns); i++ {
		value[s.columns[i]] = reflect.Indirect(s.values[i]).Interface()
	}
}

func (s *Scanner) Value() map[string]interface{} {
	value := make(map[string]interface{}, len(s.columns))
	s.SetMap(value)
	return value
}

func (s *Scanner) Err() error {
	return s.err
}
