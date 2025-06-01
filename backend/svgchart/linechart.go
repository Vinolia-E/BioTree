package svgchart

import (
	"fmt"
	"math"
	"strings"
)

// lineChart implements the Chart interface for line charts
type lineChart struct {
	data    ChartData
	options Options
	minValue float64
	maxValue float64
	padding  int
}

// newLineChart creates a new line chart instance
func newLineChart(data ChartData, options Options) Chart {
	lc := &lineChart{
		data:    data,
		options: options,
		padding: 40,
	}
	
	// Calculate min and max values
	lc.minValue = math.MaxFloat64
	lc.maxValue = -math.MaxFloat64
	for _, d := range data {
		if d.Value < lc.minValue {
			lc.minValue = d.Value
		}
		if d.Value > lc.maxValue {
			lc.maxValue = d.Value
		}
	}
	
	// Add some padding to min/max
	valueRange := lc.maxValue - lc.minValue
	lc.minValue -= valueRange * 0.1
	lc.maxValue += valueRange * 0.1
	
	return lc
}

// Generate produces the SVG string for the line chart
func (lc *lineChart) Generate() string {
	if len(lc.data) == 0 {
		return ""
	}

	width := lc.options.Width
	height := lc.options.Height
	graphWidth := float64(width - (lc.padding * 2))
	graphHeight := float64(height - (lc.padding * 2))

	// Create SVG with styles
	svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<style>
			.axis { stroke: #333; stroke-width: 1; }
			.grid { stroke: #ccc; stroke-width: 0.5; stroke-dasharray: 4 2; }
			.label { font-family: Arial; font-size: 12px; }
			.title { font-family: Arial; font-size: 16px; font-weight: bold; }
			.data-point { fill: #4285f4; }
			.data-line { stroke: #4285f4; stroke-width: 2; fill: none; }
			.data-point:hover { fill: #2962ff; r: 6; }
		</style>`, width, height)

	// Add title if present
	if lc.options.Title != "" {
		svg += fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" class="title">%s</text>`,
			width/2, lc.padding/2, lc.options.Title)
	}

	// Draw axes
	svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" class="axis"/>`, // y-axis
		lc.padding, lc.padding, lc.padding, height-lc.padding)
	svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" class="axis"/>`, // x-axis
		lc.padding, height-lc.padding, width-lc.padding, height-lc.padding)

	// Draw grid and y-axis labels
	numGridLines := 5
	for i := 0; i <= numGridLines; i++ {
		y := float64(lc.padding) + (graphHeight * float64(i) / float64(numGridLines))
		value := lc.maxValue - ((lc.maxValue - lc.minValue) * float64(i) / float64(numGridLines))
		
		// Grid line
		svg += fmt.Sprintf(`<line x1="%d" y1="%f" x2="%d" y2="%f" class="grid"/>`,
			lc.padding, y, width-lc.padding, y)
		
		// Y-axis label
		svg += fmt.Sprintf(`<text x="%d" y="%f" text-anchor="end" alignment-baseline="middle" class="label">%.1f</text>`,
			lc.padding-5, y, value)
	}

	// Generate line path
	points := make([]string, len(lc.data))
	for i, d := range lc.data {
		// Calculate point position
		x := float64(lc.padding) + (float64(i) * graphWidth / float64(len(lc.data)-1))
		heightRatio := (d.Value - lc.minValue) / (lc.maxValue - lc.minValue)
		y := float64(height-lc.padding) - (heightRatio * graphHeight)

		// Add point to path
		points[i] = fmt.Sprintf("%f,%f", x, y)

		// Add data point circle
		svg += fmt.Sprintf(`<circle cx="%f" cy="%f" r="4" class="data-point">
			<title>%s: %.2f</title>
		</circle>`,
			x, y, d.Label, d.Value)

		// X-axis label
		svg += fmt.Sprintf(`<text x="%f" y="%d" text-anchor="middle" transform="rotate(45 %f,%d)" class="label">%s</text>`,
			x, height-lc.padding+5, x, height-lc.padding+5, d.Label)
	}

	// Draw line connecting points
	svg += fmt.Sprintf(`<path d="M%s" class="data-line"/>`, strings.Join(points, " L"))

	// Close SVG
	svg += "</svg>"

	return svg
}
