package svgchart

import (
	"strings"
	"testing"
)

func TestBarChart(t *testing.T) {
	tests := []struct {
		name    string
		data    ChartData
		options Options
		checks  func(t *testing.T, svg string)
	}{
		{
			name:    "Empty data",
			data:    ChartData{},
			options: DefaultOptions(),
			checks: func(t *testing.T, svg string) {
				if !strings.Contains(svg, "No data available") {
					t.Error("Empty chart should contain 'No data available' message")
				}
			},
		},
		{
			name:    "Single data point",
			data:    ChartData{{Unit: "A", Value: 10.0}},
			options: DefaultOptions(),
			checks: func(t *testing.T, svg string) {
				if !strings.Contains(svg, "<rect") {
					t.Error("Bar chart should contain rect elements")
				}
			},
		},
		{
			name:    "Equal Y values",
			data:    ChartData{{Unit: "A", Value: 10.0}, {Unit: "B", Value: 10.0}},
			options: DefaultOptions(),
			checks: func(t *testing.T, svg string) {
				if !strings.Contains(svg, "<rect") {
					t.Error("Bar chart should contain rect elements")
				}
			},
		},
		{
			name:    "Negative values",
			data:    ChartData{{Unit: "A", Value: -10.0}, {Unit: "B", Value: 20.0}},
			options: DefaultOptions(),
			checks: func(t *testing.T, svg string) {
				if !strings.Contains(svg, "<rect") {
					t.Error("Bar chart should contain rect elements")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := newBarChart(tt.data, tt.options)
			svg := bc.Generate()
			if svg == "" {
				t.Error("Generate() returned empty string")
			}
			if tt.checks != nil {
				tt.checks(t, svg)
			}
		})
	}
}
