package svgchart

// Option is a function that modifies chart options
type Option func(*Options)

// Options contains all configurable chart options
type Options struct {
	Width      int
	Height     int
	Title      string
	XLabel     string
	YLabel     string
	ShowGrid   bool
	ChartType  ChartType
	Colors     ColorScheme
	Margins    Margins
}

// ColorScheme defines the color palette for the chart
type ColorScheme struct {
	Background string
	Axis       string
	Grid       string
	Line       string
	Text       string
	Title      string
	Bar        string
}

// Margins defines the chart margins
type Margins struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// DefaultOptions returns the default chart options
func DefaultOptions() Options {
	return Options{
		Width:     500,
		Height:    300,
		ShowGrid:  true,
		ChartType: Line,
		Colors: ColorScheme{
			Background: "#ffffff",
			Axis:       "#333333",
			Grid:       "#dddddd",
			Line:       "#3366cc",
			Bar:        "#3366cc",
			Text:       "#333333",
			Title:      "#000000",
		},
		Margins: Margins{
			Top:    40,
			Right:  20,
			Bottom: 50,
			Left:   60,
		},
	}
}

// WithTitle sets the chart title
func WithTitle(title string) Option {
	return func(o *Options) {
		o.Title = title
		if title != "" {
			o.Margins.Top = 40
		} else {
			o.Margins.Top = 20
		}
	}
}

// WithDimensions sets the chart dimensions
func WithDimensions(width, height int) Option {
	return func(o *Options) {
		o.Width = width
		o.Height = height
	}
}

// WithXLabel sets the X-axis label
func WithXLabel(label string) Option {
	return func(o *Options) {
		o.XLabel = label
		if label != "" {
			o.Margins.Bottom = 50
		} else {
			o.Margins.Bottom = 30
		}
	}
}

// WithYLabel sets the Y-axis label
func WithYLabel(label string) Option {
	return func(o *Options) {
		o.YLabel = label
		if label != "" {
			o.Margins.Left = 60
		} else {
			o.Margins.Left = 40
		}
	}
}

// WithGrid toggles grid display
func WithGrid(show bool) Option {
	return func(o *Options) {
		o.ShowGrid = show
	}
}

// WithColors sets the color scheme
func WithColors(colors ColorScheme) Option {
	return func(o *Options) {
		o.Colors = colors
	}
}

// WithMargins sets custom margins
func WithMargins(m Margins) Option {
	return func(o *Options) {
		o.Margins = m
	}
}