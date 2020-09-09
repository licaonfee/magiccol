package magiccol

import (
	"reflect"
	"time"
)

//ColumnType is identical as defined in sql.ColumnType struct
type ColumnType interface {
	Name() string
	DatabaseTypeName() string
	ScanType() reflect.Type
	Nullable() (nullable bool, ok bool)
	DecimalSize() (precision int64, scale int64, ok bool)
	Length() (length int64, ok bool)
}

//Matcher return a type and true if column definition match
//on a negative match reflect.Type should be null but is not mandatory
type Matcher func(ColumnType) (reflect.Type, bool)

var (
	stringType   = reflect.TypeOf("")
	intType      = reflect.TypeOf(int64(0))
	dateType     = reflect.TypeOf(time.Time{})
	bytesType    = reflect.TypeOf([]byte{})
	floatType    = reflect.TypeOf(float64(0))
	complexType  = reflect.TypeOf(complex64(0))
	boolType     = reflect.TypeOf(false)
	durationType = reflect.TypeOf(time.Duration(0))
)

//Mapper translate sql types to golang types
type Mapper struct {
	m     map[string]reflect.Type
	match []Matcher
}

//Get do a map lookup if type is not found return a ScanType itself
func (l *Mapper) Get(col ColumnType) reflect.Type {
	for _, m := range l.match {
		t, ok := m(col)
		if ok {
			return t
		}
	}
	t, ok := l.m[col.DatabaseTypeName()]
	if !ok {
		return col.ScanType()
	}
	return t
}

//Match method allow to set custom types as scanneable types
//if m is nil then is a no-op
func (l *Mapper) Match(m ...Matcher) {
	for i := 0; i < len(m); i++ {
		if m[i] != nil {
			l.match = append(l.match, m[i])
		}
	}
}

func DatabaseTypeAs(databaseTypeName string, t reflect.Type) Matcher {
	return func(col ColumnType) (reflect.Type, bool) {
		if col.DatabaseTypeName() == databaseTypeName {
			return t, true
		}
		return nil, false
	}
}

func ColumnNameAs(columnName string, t reflect.Type) Matcher {
	return func(col ColumnType) (reflect.Type, bool) {
		if col.Name() == columnName {
			return t, true
		}
		return nil, false
	}
}

//DefaultMapper provides a mapping for most common sql types
//type list reference used is:
//http://jakewheat.github.io/sql-overview/sql-2011-foundation-grammar.html#predefined-type
func DefaultMapper() *Mapper {
	m := map[string]reflect.Type{
		//Character types
		"CHARACTER":                       stringType,
		"CHAR":                            stringType,
		"CHARACTER VARYING":               stringType,
		"CHAR VARYING":                    stringType,
		"VARCHAR":                         stringType,
		"TEXT":                            stringType,
		"CHARACTER LARGE OBJECT":          stringType,
		"CHAR LARGE OBJECT":               stringType,
		"CLOB":                            stringType,
		"NATIONAL CHARACTER":              stringType,
		"NATIONAL CHAR":                   stringType,
		"NCHAR":                           stringType,
		"NATIONAL CHARACTER VARYING":      stringType,
		"NATIONAL CHAR VARYING":           stringType,
		"NCHAR VARYING":                   stringType,
		"NATIONAL CHARACTER LARGE OBJECT": stringType,
		"NCHAR LARGE OBJECT":              stringType,
		"NCLOB":                           stringType,
		//Binary types
		"BINARY":              bytesType,
		"BINARY VARYING":      bytesType,
		"VARBINARY":           bytesType,
		"BINARY LARGE OBJECT": bytesType,
		"BLOB":                bytesType,
		//exact numeric type
		"NUMERIC":  floatType,
		"DECIMAL":  floatType,
		"DEC":      floatType,
		"SMALLINT": intType,
		"INTEGER":  intType,
		"INT":      intType,
		"BIGINT":   intType,
		//approximate numeric type
		"FLOAT":            floatType,
		"REAL":             complexType,
		"DOUBLE PRECISION": floatType,
		//boolean type
		"BOOLEAN": boolType,
		//datetime type
		"DATE":      dateType,
		"TIME":      dateType,
		"TIMESTAMP": dateType,
		//Interval type
		"INTERVAL": durationType,
	}
	return &Mapper{m: m}
}
