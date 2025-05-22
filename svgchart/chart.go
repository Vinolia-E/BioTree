package svgchart

type ChartType string

const (
	Line ChartType = "line"
	Bar  ChartType = "bar"
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
	default:
		return newLineChart(chartData, options), nil
	}
}