package magiccol_test

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/licaonfee/magiccol"
)

var errColumnTypes = errors.New("fail ColumnTypes")
var errColumns = errors.New("fail Columns")

type mockRows struct {
	columns []string
	types   []*sql.ColumnType
}

func (f *mockRows) ColumnTypes() ([]*sql.ColumnType, error) { return f.types, nil }

func (f *mockRows) Columns() ([]string, error) { return f.columns, nil }

func (f *mockRows) Err() error { return nil }

func (f *mockRows) Scan(args ...interface{}) error { return nil }

func (f *mockRows) Next() bool { return false }

type failColumns struct {
	mockRows
}

func (f *failColumns) Columns() ([]string, error) {
	return nil, errColumns
}

type failColumnType struct {
	mockRows
}

func (f *failColumnType) ColumnTypes() ([]*sql.ColumnType, error) {
	return nil, errColumnTypes
}

func TestNewScanner(t *testing.T) {
	tests := []struct {
		name    string
		opts    magiccol.Options
		wantErr error
	}{
		{
			name:    "nil rows",
			opts:    magiccol.Options{Rows: nil},
			wantErr: magiccol.ErrNilRows,
		},
		{
			name:    "error Columns()",
			opts:    magiccol.Options{Rows: &failColumns{}},
			wantErr: errColumns,
		},
		{
			name:    "error ColumnType()",
			opts:    magiccol.Options{Rows: &failColumnType{}},
			wantErr: errColumnTypes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := magiccol.NewScanner(tt.opts)
			if err != tt.wantErr || !errors.As(err, &tt.wantErr) {
				t.Errorf("NewScanner() err = %v , wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestScan(t *testing.T) {
	rowError := errors.New("row error")
	tests := []struct {
		name    string
		csvData string
		columns []*sqlmock.Column
		want    []map[string]interface{}
		wantErr error
	}{
		{
			name: "success",
			columns: []*sqlmock.Column{
				sqlmock.NewColumn("name").OfType("VARCHAR", ""),
				sqlmock.NewColumn("age").OfType("INTEGER", int64(0)),
			},
			csvData: strings.Join([]string{`"jhon",35`, `"jeremy",29`}, "\n"),
			want: []map[string]interface{}{
				{"name": "jhon", "age": int64(35)},
				{"name": "jeremy", "age": int64(29)},
			},
			wantErr: nil,
		},
		{
			name: "Rows error",
			columns: []*sqlmock.Column{
				sqlmock.NewColumn("name").OfType("VARCHAR", ""),
			},
			csvData: `"jeimy"`,
			want:    nil,
			wantErr: rowError,
		},
		{
			name: "Scan error",
			columns: []*sqlmock.Column{
				sqlmock.NewColumn("name").OfType("INTEGER", int64(0)),
			},
			csvData: `"jeimy"`,
			want:    nil,
			wantErr: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			r := mock.NewRowsWithColumnDefinition(tt.columns...)
			r.FromCSVString(tt.csvData)
			mock.ExpectQuery("SELECT").WillReturnRows(r)
			if tt.wantErr != nil {
				r.RowError(0, tt.wantErr)
			}
			rows, _ := db.Query("SELECT")
			m, err := magiccol.NewScanner(magiccol.Options{Rows: rows})
			if err != nil {
				t.Errorf("NewScanner() err = %v", err)
			}
			got := make([]map[string]interface{}, 0)
			for m.Scan() {
				got = append(got, m.Value())
			}
			if m.Err() != tt.wantErr {
				e := errors.New("")
				if !(m.Err() != nil && tt.wantErr != nil) {
					t.Errorf("Scan() err = %v , want = %v", m.Err(), tt.wantErr)
				}
				if !errors.As(m.Err(), &e) || m.Err().Error() != tt.wantErr.Error() {
					t.Errorf("Scan() err = %v , want = %v", m.Err(), tt.wantErr)
				}
			}
			if m.Err() != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() got = %v , want = %v", got, tt.want)
			}

		})
	}
}