package svgchart

import "fmt"

type ChartType string

const (
	Line ChartType = "line"
	Bar  ChartType = "bar"
	Pie  ChartType = "pie"
)

type Chart interface {
	Generate() string
}

// New creates a new chart based on the specified type
func New(data interface{}, chartType ChartType, opts ...Option) (Chart, error) {
	chartData, err := ConvertData(data)
	if err != nil {
		return nil, err
	}

	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	switch chartType {
	case Bar:
		return newBarChart(chartData, options), nil
	case Pie:
		return newPieChart(chartData, options), nil
	case Line:
		return newLineChart(chartData, options), nil
	default:
		return nil, fmt.Errorf("unsupported chart type: %s", chartType)
	}
}
