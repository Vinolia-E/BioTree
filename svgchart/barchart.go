package svgchart

import (
	"fmt"
	"math"
)

// barChart implements the Chart interface for bar charts
type barChart struct {
	data    ChartData
	options Options
}

// newBarChart creates a new bar chart instance
func newBarChart(data ChartData, options Options) *barChart {
	return &barChart{
		data:    data,
		options: options,
	}
}

// Generate produces the SVG string for the bar chart
func (bc *barChart) Generate() string {
	if len(bc.data) == 0 {
		return generateEmptyChart(bc.options, "No data available")
	}

	// Calculate chart dimensions
	chartWidth := bc.options.Width - bc.options.Margins.Left - bc.options.Margins.Right
	chartHeight := bc.options.Height - bc.options.Margins.Top - bc.options.Margins.Bottom

	// Get min/max values for scaling
	yMin, yMax := getYMinMax(bc.data)

	// Handle uniform values
	if yMin == yMax {
		padding := math.Max(1, yMin*0.1)
		if yMin == 0 {
			padding = 1
		}
		yMin -= padding
		yMax += padding
	}

	sb := NewSVGBuilder(bc.options.Width, bc.options.Height)

	// Background
	sb.AddRect(0, 0, bc.options.Width, bc.options.Height, map[string]string{
		"fill": bc.options.Colors.Background,
	})

	// Add title if provided
	if bc.options.Title != "" {
		sb.AddText(bc.options.Width/2, 20, bc.options.Title, map[string]string{
			"text-anchor":  "middle",
			"font-family":  "Arial",
			"font-size":    "16px",
			"font-weight": "bold",
			"fill":        bc.options.Colors.Title,
		})
	}

	// Add grid if enabled
	if bc.options.ShowGrid {
		bc.addGrid(sb, chartWidth, chartHeight)
	}

	// Add axes
	bc.addAxes(sb, chartWidth, chartHeight, yMin, yMax)

	// Add axis labels if provided
	if bc.options.XLabel != "" {
		sb.AddText(
			bc.options.Margins.Left+chartWidth/2,
			bc.options.Height-10,
			bc.options.XLabel,
			map[string]string{
				"text-anchor": "middle",
				"font-family": "Arial",
				"font-size":   "12px",
				"fill":       bc.options.Colors.Text,
			},
		)
	}

	if bc.options.YLabel != "" {
		sb.AddText(
			15,
			bc.options.Margins.Top+chartHeight/2,
			bc.options.YLabel,
			map[string]string{
				"text-anchor":  "middle",
				"font-family":  "Arial",
				"font-size":    "12px",
				"fill":        bc.options.Colors.Text,
				"transform":   fmt.Sprintf("rotate(-90, 15, %d)", bc.options.Margins.Top+chartHeight/2),
			},
		)
	}

	// Generate the bars
	bc.addBars(sb, chartWidth, chartHeight, yMin, yMax)

	return sb.String()
}

func (bc *barChart) addGrid(sb *SVGBuilder, chartWidth, chartHeight int) {
	numLines := 5
	for i := 0; i <= numLines; i++ {
		yPos := bc.options.Margins.Top + chartHeight - (i * chartHeight / numLines)
		sb.AddLine(
			bc.options.Margins.Left,
			yPos,
			bc.options.Margins.Left+chartWidth,
			yPos,
			map[string]string{
				"stroke":         bc.options.Colors.Grid,
				"stroke-width":   "1",
				"stroke-dasharray": "5,5",
			},
		)
	}
}

func (bc *barChart) addAxes(sb *SVGBuilder, chartWidth, chartHeight int, yMin, yMax float64) {
	m := bc.options.Margins

	// X-axis
	sb.AddLine(
		m.Left,
		m.Top+chartHeight,
		m.Left+chartWidth,
		m.Top+chartHeight,
		map[string]string{
			"stroke":       bc.options.Colors.Axis,
			"stroke-width": "2",
		},
	)

	// Y-axis
	sb.AddLine(
		m.Left,
		m.Top,
		m.Left,
		m.Top+chartHeight,
		map[string]string{
			"stroke":       bc.options.Colors.Axis,
			"stroke-width": "2",
		},
	)

	// X-axis labels
	numLabels := 7
	if len(bc.data) < numLabels {
		numLabels = len(bc.data)
	}

	if numLabels > 0 {
		step := 1
		if len(bc.data) > numLabels {
			step = len(bc.data) / numLabels
		}

		for i := 0; i < len(bc.data); i += step {
			if i >= len(bc.data) {
				break
			}

			xPos := m.Left + int((float64(i)+0.5)*(float64(chartWidth)/float64(len(bc.data))))
			sb.AddText(
				xPos,
				m.Top+chartHeight+15,
				bc.data[i].Label,
				map[string]string{
					"text-anchor": "middle",
					"font-family": "Arial",
					"font-size":   "10px",
					"fill":       bc.options.Colors.Text,
				},
			)
		}
	}

	// Y-axis labels
	numYLabels := 5
	for i := 0; i <= numYLabels; i++ {
		yValue := yMin + float64(i)*(yMax-yMin)/float64(numYLabels)
		yPos := m.Top + chartHeight - int(float64(i)*float64(chartHeight)/float64(numYLabels))

		yLabel := formatNumber(yValue)

		sb.AddText(
			m.Left-5,
			yPos+3,
			yLabel,
			map[string]string{
				"text-anchor": "end",
				"font-family": "Arial",
				"font-size":   "10px",
				"fill":       bc.options.Colors.Text,
			},
		)
	}
}

func (bc *barChart) addBars(sb *SVGBuilder, chartWidth, chartHeight int, yMin, yMax float64) {
	barWidth := float64(chartWidth) / float64(len(bc.data)) * 0.8
	gapWidth := float64(chartWidth) / float64(len(bc.data)) * 0.2

	for i, dp := range bc.data {
		// Calculate position
		xPos := float64(bc.options.Margins.Left) + (float64(i) * float64(chartWidth) / float64(len(bc.data))) + gapWidth/2

		// Scale y value
		var barHeight float64
		if yMax == yMin {
			barHeight = float64(chartHeight) / 2
		} else {
			yScaled := ((dp.Value - yMin) / (yMax - yMin)) * float64(chartHeight)
			barHeight = yScaled
		}

		yPos := float64(bc.options.Margins.Top) + float64(chartHeight) - barHeight

		sb.AddRect(
			int(xPos),
			int(yPos),
			int(barWidth),
			int(barHeight),
			map[string]string{
				"fill":   bc.options.Colors.Bar,
				"stroke": "none",
			},
		)
	}
}