package svgchart

import (
	"fmt"
	"math"
)

// barChart implements the Chart interface for bar charts
type barChart struct {
	data    ChartData
	options Options
	// Calculated values
	minValue float64
	maxValue float64
	padding  int
	barWidth float64
}

// newBarChart creates a new bar chart instance
func newBarChart(data ChartData, options Options) Chart {
	bc := &barChart{
		data:    data,
		options: options,
		padding: 40,
	}

	// Calculate min and max values
	bc.minValue = math.MaxFloat64
	bc.maxValue = -math.MaxFloat64
	for _, d := range data {
		if d.Value < bc.minValue {
			bc.minValue = d.Value
		}
		if d.Value > bc.maxValue {
			bc.maxValue = d.Value
		}
	}

	// Add some padding to min/max
	valueRange := bc.maxValue - bc.minValue
	bc.minValue = math.Max(0, bc.minValue-valueRange*0.1) // Don't go below 0 for bar charts
	bc.maxValue += valueRange * 0.1

	// Calculate bar width
	graphWidth := float64(options.Width - (bc.padding * 2))
	bc.barWidth = (graphWidth / float64(len(data))) * 0.8 // Leave 20% gap between bars

	return bc
}

// Generate produces the SVG string for the bar chart
func (b *barChart) Generate() string {
	if len(b.data) == 0 {
		return ""
	}

	width := b.options.Width
	height := b.options.Height
	graphWidth := float64(width - (b.padding * 2))
	graphHeight := float64(height - (b.padding * 2))

	// Create SVG with styles
	svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<style>
			.axis { stroke: #333; stroke-width: 1; }
			.grid { stroke: #ccc; stroke-width: 0.5; stroke-dasharray: 4 2; }
			.label { font-family: Arial; font-size: 12px; }
			.title { font-family: Arial; font-size: 16px; font-weight: bold; }
			.bar { fill: #4285f4; transition: fill 0.3s; }
			.bar:hover { fill: #2962ff; }
			.value-label { font-family: Arial; font-size: 12px; fill: #333; }
		</style>`, width, height)

	// Add title if present
	if b.options.Title != "" {
		svg += fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" class="title">%s</text>`,
			width/2, b.padding/2, b.options.Title)
	}

	// Draw axes
	svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" class="axis"/>`, // y-axis
		b.padding, b.padding, b.padding, height-b.padding)
	svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" class="axis"/>`, // x-axis
		b.padding, height-b.padding, width-b.padding, height-b.padding)

	// Draw grid and y-axis labels
	numGridLines := 5
	for i := 0; i <= numGridLines; i++ {
		y := float64(b.padding) + (graphHeight * float64(i) / float64(numGridLines))
		value := b.maxValue - ((b.maxValue - b.minValue) * float64(i) / float64(numGridLines))

		// Grid line
		svg += fmt.Sprintf(`<line x1="%d" y1="%f" x2="%d" y2="%f" class="grid"/>`,
			b.padding, y, width-b.padding, y)

		// Y-axis label
		svg += fmt.Sprintf(`<text x="%d" y="%f" text-anchor="end" alignment-baseline="middle" class="label">%.1f</text>`,
			b.padding-5, y, value)
	}

	// Generate bars
	spacing := (graphWidth - (float64(len(b.data)) * b.barWidth)) / float64(len(b.data)+1)
	for i, d := range b.data {
		// Calculate bar position and height
		x := float64(b.padding) + spacing + (float64(i) * (b.barWidth + spacing))
		heightRatio := (d.Value - b.minValue) / (b.maxValue - b.minValue)
		barHeight := heightRatio * float64(graphHeight)
		y := float64(height-b.padding) - barHeight

		// Draw bar
		svg += fmt.Sprintf(`<rect x="%f" y="%f" width="%f" height="%f" class="bar">
			<title>%s: %.2f</title>
		</rect>`,
			x, y, b.barWidth, barHeight, d.Label, d.Value)

		// Add value label on top of bar
		svg += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" class="value-label">%.1f</text>`,
			x+b.barWidth/2, y-5, d.Value)

		// X-axis label
		svg += fmt.Sprintf(`<text x="%f" y="%d" text-anchor="middle" transform="rotate(45 %f,%d)" class="label">%s</text>`,
			x+b.barWidth/2, height-b.padding+5, x+b.barWidth/2, height-b.padding+5, d.Label)
	}

	// Close SVG
	svg += "</svg>"

	return svg
}
