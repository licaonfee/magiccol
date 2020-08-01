package magiccol

import (
	"reflect"
	"time"
)

//Mapper translate sql types to golang types
type Mapper interface {
	//Get typeName should be sql type as is called in sql.ColumnType.DatabaseTypeName()
	Get(typeName string, fallback reflect.Type) reflect.Type
	//Type allow to set alias, extends or fix mapper behaviour
	Type(t reflect.Type, asTypes ...string)
}

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

type LookupMapper struct {
	m map[string]reflect.Type
}

//Get do a map lookup if type is not found return a ScanType itself
func (l LookupMapper) Get(typeName string, fallback reflect.Type) reflect.Type {
	t, ok := l.m[typeName]
	if !ok {
		return fallback
	}
	return t
}

func (l *LookupMapper) Type(t reflect.Type, asType ...string) {
	for _, x := range asType {
		tp := x
		l.m[tp] = t
	}
}

func DefaultMapper() Mapper {
	m := map[string]reflect.Type{
		//Character types
		"CHARACTER":                       stringType,
		"CHAR":                            stringType,
		"CHARACTER VARYING":               stringType,
		"CHAR VARYING":                    stringType,
		"VARCHAR":                         stringType,
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
	return &LookupMapper{m: m}
}
