package svgchart

import (
    "strings"
    "testing"
)

func TestLineChart(t *testing.T) {
    tests := []struct {
        name    string
        data    ChartData
        options Options
        checks  func(t *testing.T, svg string)
    }{
        {
            name: "Empty data",
            data: ChartData{},
            options: DefaultOptions(),
            checks: func(t *testing.T, svg string) {
                if !strings.Contains(svg, "No data available") {
                    t.Error("Empty chart should contain 'No data available' message")
                }
            },
        },
        {
            name: "Single data point",
            data: ChartData{{Label: "A", Value: 10.0}},
            options: DefaultOptions(),
            checks: func(t *testing.T, svg string) {
                if !strings.Contains(svg, "circle") {
                    t.Error("Line chart with single point should contain circle element")
                }
            },
        },
        {
            name: "Equal Y values",
            data: ChartData{{Label: "A", Value: 10.0}, {Label: "B", Value: 10.0}},
            options: DefaultOptions(),
            checks: func(t *testing.T, svg string) {
                if !strings.Contains(svg, "path") {
                    t.Error("Line chart should contain path element")
                }
            },
        },
        {
            name: "With title and labels",
            data: ChartData{{Label: "A", Value: 10.0}, {Label: "B", Value: 20.0}},
            options: func() Options {
                o := DefaultOptions()
                o.Title = "Test Title"
                o.XLabel = "X Label"
                o.YLabel = "Y Label"
                return o
            }(),
            checks: func(t *testing.T, svg string) {
                if !strings.Contains(svg, "Test Title") {
                    t.Error("Chart should contain title")
                }
                if !strings.Contains(svg, "X Label") {
                    t.Error("Chart should contain X label")
                }
                if !strings.Contains(svg, "Y Label") {
                    t.Error("Chart should contain Y label")
                }
            },
        },
        {
            name: "Without grid",
            data: ChartData{{Label: "A", Value: 10.0}, {Label: "B", Value: 20.0}},
            options: func() Options {
                o := DefaultOptions()
                o.ShowGrid = false
                return o
            }(),
            checks: func(t *testing.T, svg string) {
                if strings.Contains(svg, "stroke-dasharray") {
                    t.Error("Chart without grid should not contain dashed lines")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lc := newLineChart(tt.data, tt.options)
            svg := lc.Generate()
            if svg == "" {
                t.Error("Generate() returned empty string")
            }
            if tt.checks != nil {
                tt.checks(t, svg)
            }
        })
    }
}