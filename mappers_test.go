package magiccol_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/magiccol"
)

type col struct {
	name         string
	databaseType string
	scanType     reflect.Type
}

func (c col) Name() string                                         { return c.name }
func (c col) DatabaseTypeName() string                             { return c.databaseType }
func (c col) ScanType() reflect.Type                               { return c.scanType }
func (c col) Nullable() (nullable bool, ok bool)                   { return false, false }
func (c col) DecimalSize() (precision int64, scale int64, ok bool) { return 0, 0, false }
func (c col) Length() (length int64, ok bool)                      { return 0, false }

func TestMapperGet(t *testing.T) {
	tests := []struct {
		name  string
		col   magiccol.ColumnType
		match []magiccol.Matcher
		want  reflect.Type
	}{
		{
			name: "char type",
			col:  col{databaseType: "VARCHAR"},
			want: reflect.TypeOf(""),
		},
		{
			name: "fallback",
			col:  col{databaseType: "MISSING TYPE", scanType: reflect.TypeOf(uint8(0))},
			want: reflect.TypeOf(uint8(0)),
		},
		{
			name: "match priority",
			col:  col{databaseType: "INTEGER", scanType: reflect.TypeOf(uint64(0))},
			match: []magiccol.Matcher{
				func(c magiccol.ColumnType) (reflect.Type, bool) {
					return reflect.TypeOf(float64(0.0)), true
				}},
			want: reflect.TypeOf(float64(0.0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := magiccol.DefaultMapper()
			for i := 0; i < len(tt.match); i++ {
				m.Match(tt.match[i])
			}
			got := m.Get(tt.col)
			if got != tt.want {
				t.Errorf("Get() got = %v , want = %v", got, tt.want)
			}
		})
	}
}

func TestMatchers(t *testing.T) {
	tests := []struct {
		name     string
		positive magiccol.ColumnType
		negative magiccol.ColumnType
		m        magiccol.Matcher
		want     reflect.Type
	}{
		{
			name:     "DatabaseTypeAs",
			positive: col{databaseType: "EXOTIC"},
			negative: col{databaseType: "DOUBLE"},
			m:        magiccol.DatabaseTypeAs("EXOTIC", reflect.TypeOf(float32(0))),
			want:     reflect.TypeOf(float32(0)),
		},
		{
			name:     "ColumnNameAs",
			positive: col{name: "delimited_data"},
			negative: col{name: "name"},
			m:        magiccol.ColumnNameAs("delimited_data", reflect.TypeOf([]string{})),
			want:     reflect.TypeOf([]string{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.m(tt.positive)
			if !ok || got != tt.want {
				t.Errorf("Matcher() got = (%v,%v) , want = (%v, %v)", got, ok, tt.want, true)
			}
			got, ok = tt.m(tt.negative)
			if ok || got != nil {
				t.Errorf("Matcher() got = (%v,%v) , want = (%v, %v)", got, ok, nil, false)
			}

		})
	}
}
