package svgchart

import (
	"testing"

	"github.com/Vinolia-E/BioTree/backend/util"
)

func TestDataConversion(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		want    int // expected number of data points
		wantErr bool
	}{
		{
			name:    "Map data",
			data:    map[string]float64{"A": 1.0, "B": 2.0, "C": 3.0},
			want:    3,
			wantErr: false,
		},
		{
			name:    "DataPoint slice",
			data:    []util.DataPoint{{Unit: "A", Value: 1.0}, {Unit: "B", Value: 2.0}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Point slice",
			data:    []Point{{X: "A", Y: 1.0}, {X: "B", Y: 2.0}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Interface slice with strings and floats",
			data:    [][]interface{}{{"A", 1.0}, {"B", 2.0}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Interface slice with strings and ints",
			data:    [][]interface{}{{"A", 1}, {"B", 2}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Interface slice with string values",
			data:    [][]interface{}{{"A", "1.0"}, {"B", "2.0"}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Interface slice with invalid values",
			data:    [][]interface{}{{"A", "invalid"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty map",
			data:    map[string]float64{},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Invalid data type",
			data:    "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) != tt.want {
				t.Errorf("ConvertData() got %d data points, want %d", len(got), tt.want)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "Integer value",
			value: 10.0,
			want:  "10",
		},
		{
			name:  "Decimal value",
			value: 10.5,
			want:  "10.5",
		},
		{
			name:  "Decimal with trailing zeros",
			value: 10.50,
			want:  "10.5",
		},
		{
			name:  "Small decimal",
			value: 0.01,
			want:  "0.01",
		},
		{
			name:  "Zero",
			value: 0.0,
			want:  "0",
		},
		{
			name:  "Negative value",
			value: -10.5,
			want:  "-10.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.value)
			if got != tt.want {
				t.Errorf("formatNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
