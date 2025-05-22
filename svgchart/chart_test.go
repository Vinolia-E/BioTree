package svgchart

import (
    "strings"
    "testing"
)

func TestNew(t *testing.T) {
    tests := []struct {
        name      string
        data      interface{}
        chartType ChartType
        options   []Option
        wantErr   bool
    }{
        {
            name:      "Valid map data with line chart",
            data:      map[string]float64{"A": 1.0, "B": 2.0},
            chartType: Line,
            wantErr:   false,
        },
        {
            name:      "Valid map data with bar chart",
            data:      map[string]float64{"A": 1.0, "B": 2.0},
            chartType: Bar,
            wantErr:   false,
        },
        {
            name:      "Empty map data",
            data:      map[string]float64{},
            chartType: Line,
            wantErr:   false,
        },
        {
            name:      "DataPoint slice",
            data:      []DataPoint{{Label: "A", Value: 1.0}, {Label: "B", Value: 2.0}},
            chartType: Line,
            wantErr:   false,
        },
        {
            name:      "Point slice",
            data:      []Point{{X: "A", Y: 1.0}, {X: "B", Y: 2.0}},
            chartType: Line,
            wantErr:   false,
        },
        {
            name:      "Interface slice",
            data:      [][]interface{}{{"A", 1.0}, {"B", 2.0}},
            chartType: Line,
            wantErr:   false,
        },
        {
            name:      "Invalid data type",
            data:      "invalid",
            chartType: Line,
            wantErr:   true,
        },
        {
            name:      "With custom options",
            data:      map[string]float64{"A": 1.0, "B": 2.0},
            chartType: Line,
            options:   []Option{WithTitle("Test"), WithDimensions(800, 600)},
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := New(tt.data, tt.chartType, tt.options...)
            if (err != nil) != tt.wantErr {
                t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                svg := got.Generate()
                if svg == "" {
                    t.Error("Generate() returned empty string")
                }
                if !strings.HasPrefix(svg, "<svg") {
                    t.Error("Generated SVG should start with <svg tag")
                }
            }
        })
    }
}