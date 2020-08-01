package magiccol_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/magiccol"
)

func TestMapperGet(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		fallback reflect.Type
		want     reflect.Type
	}{
		{
			name:     "char type",
			typeName: "VARCHAR",
			fallback: nil,
			want:     reflect.TypeOf(""),
		},
		{
			name:     "fallback",
			typeName: "MISSING TYPE",
			fallback: reflect.TypeOf(uint8(0)),
			want:     reflect.TypeOf(uint8(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := magiccol.DefaultMapper()
			got := m.Get(tt.typeName, tt.fallback)
			if got != tt.want {
				t.Errorf("Get() got = %v , want = %v", got, tt.want)
			}
		})
	}
}

func TestMapperType(t *testing.T) {
	tests := []struct {
		name  string
		rType reflect.Type
		as    []string
	}{
		{
			name:  "new type",
			rType: reflect.TypeOf(false),
			as:    []string{"MY_BOOL", "FLAGTYPE"},
		},
		{
			name:  "overwrite",
			rType: reflect.TypeOf(int64(0)),
			as:    []string{"VARCHAR"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := magiccol.DefaultMapper()
			m.Type(tt.rType, tt.as...)
			for _, n := range tt.as {
				if tp := m.Get(n, nil); tp != tt.rType {
					t.Errorf("Type(%s) not set %v", n, tt.rType)
				}
			}
		})
	}
}
