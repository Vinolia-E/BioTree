package chart

const (
	width       = 400
	height      = 200
	padding     = 40
	axisColor   = "#333"
	lineColor   = "#40e0d0"
	bgColor     = "#ffffff"
	fontSize    = 10
	gridColor   = "#eee"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// GenerateLineChartSVG creates an SVG string for a line chart
func GenerateLineChartSVG(data map[string]float64) string {
	return ""
}