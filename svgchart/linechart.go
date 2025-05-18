package svgchart

import (
	"fmt"
	"math"
)

// lineChart implements the Chart interface for line charts
type lineChart struct {
	data    ChartData
	options Options
}

// newLineChart creates a new line chart instance
func newLineChart(data ChartData, options Options) *lineChart {
	return &lineChart{
		data:    data,
		options: options,
	}
}

// Generate produces the SVG string for the line chart
func (lc *lineChart) Generate() string {
	if len(lc.data) == 0 {
		return generateEmptyChart(lc.options, "No data available")
	}

	// Calculate chart dimensions
	chartWidth := lc.options.Width - lc.options.Margins.Left - lc.options.Margins.Right
	chartHeight := lc.options.Height - lc.options.Margins.Top - lc.options.Margins.Bottom

	// Get min/max values for scaling
	yMin, yMax := getYMinMax(lc.data)

	if yMin == yMax {
		padding := math.Max(1, yMin*0.1)
		if yMin == 0 {
			padding = 1
		}
		yMin -= padding
		yMax += padding
	}

	sb := NewSVGBuilder(lc.options.Width, lc.options.Height)

	// Background
	sb.AddRect(0, 0, lc.options.Width, lc.options.Height, map[string]string{
		"fill": lc.options.Colors.Background,
	})

	// Add title if provided
	if lc.options.Title != "" {
		sb.AddText(lc.options.Width/2, 20, lc.options.Title, map[string]string{
			"text-anchor":  "middle",
			"font-family":  "Arial",
			"font-size":    "16px",
			"font-weight": "bold",
			"fill":        lc.options.Colors.Title,
		})
	}

	if lc.options.ShowGrid {
		lc.addGrid(sb, chartWidth, chartHeight)
	}

	lc.addAxes(sb, chartWidth, chartHeight, yMin, yMax)

	// Add axis labels if provided
	if lc.options.XLabel != "" {
		sb.AddText(
			lc.options.Margins.Left+chartWidth/2,
			lc.options.Height-10,
			lc.options.XLabel,
			map[string]string{
				"text-anchor": "middle",
				"font-family": "Arial",
				"font-size":   "12px",
				"fill":       lc.options.Colors.Text,
			},
		)
	}

	if lc.options.YLabel != "" {
		sb.AddText(
			15,
			lc.options.Margins.Top+chartHeight/2,
			lc.options.YLabel,
			map[string]string{
				"text-anchor":  "middle",
				"font-family":  "Arial",
				"font-size":    "12px",
				"fill":        lc.options.Colors.Text,
				"transform":   fmt.Sprintf("rotate(-90, 15, %d)", lc.options.Margins.Top+chartHeight/2),
			},
		)
	}

	lc.addLinePath(sb, chartWidth, chartHeight, yMin, yMax)

	return sb.String()
}

func (lc *lineChart) addGrid(sb *SVGBuilder, chartWidth, chartHeight int) {
	numLines := 5
	for i := 0; i <= numLines; i++ {
		yPos := lc.options.Margins.Top + chartHeight - (i * chartHeight / numLines)
		sb.AddLine(
			lc.options.Margins.Left,
			yPos,
			lc.options.Margins.Left+chartWidth,
			yPos,
			map[string]string{
				"stroke":         lc.options.Colors.Grid,
				"stroke-width":   "1",
				"stroke-dasharray": "5,5",
			},
		)
	}
}

func (lc *lineChart) addAxes(sb *SVGBuilder, chartWidth, chartHeight int, yMin, yMax float64) {
	m := lc.options.Margins

	// X-axis
	sb.AddLine(
		m.Left,
		m.Top+chartHeight,
		m.Left+chartWidth,
		m.Top+chartHeight,
		map[string]string{
			"stroke":       lc.options.Colors.Axis,
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
			"stroke":       lc.options.Colors.Axis,
			"stroke-width": "2",
		},
	)

	// X-axis labels
	numLabels := 7
	if len(lc.data) < numLabels {
		numLabels = len(lc.data)
	}

	if numLabels > 0 {
		step := 1
		if len(lc.data) > numLabels {
			step = len(lc.data) / numLabels
		}

		for i := 0; i < len(lc.data); i += step {
			if i >= len(lc.data) {
				break
			}

			var xPos int
			if len(lc.data) == 1 {
				xPos = m.Left + chartWidth/2
			} else {
				xPos = m.Left + (i * chartWidth / (len(lc.data) - 1))
			}

			sb.AddText(
				xPos,
				m.Top+chartHeight+15,
				lc.data[i].Label,
				map[string]string{
					"text-anchor": "middle",
					"font-family": "Arial",
					"font-size":   "10px",
					"fill":       lc.options.Colors.Text,
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
				"fill":       lc.options.Colors.Text,
			},
		)
	}
}

func (lc *lineChart) addLinePath(sb *SVGBuilder, chartWidth, chartHeight int, yMin, yMax float64) {
	var points []struct{ x, y int }
	for i, dp := range lc.data {
		var xPos int
		if len(lc.data) == 1 {
			xPos = lc.options.Margins.Left + chartWidth/2
		} else {
			xPos = lc.options.Margins.Left + (i * chartWidth / (len(lc.data) - 1))
		}

		var yScaled float64
		if yMax == yMin {
			yScaled = float64(chartHeight) / 2
		} else {
			yScaled = ((dp.Value - yMin) / (yMax - yMin)) * float64(chartHeight)
		}

		yPos := lc.options.Margins.Top + (chartHeight - int(yScaled))
		points = append(points, struct{ x, y int }{xPos, yPos})
	}

	// Create line path
	if len(points) > 0 {
		path := fmt.Sprintf("M%d,%d ", points[0].x, points[0].y)
		for i := 1; i < len(points); i++ {
			path += fmt.Sprintf("L%d,%d ", points[i].x, points[i].y)
		}

		sb.AddPath(
			path,
			map[string]string{
				"fill":         "none",
				"stroke":       lc.options.Colors.Line,
				"stroke-width": "2",
			},
		)

		// Add points
		for _, p := range points {
			sb.AddCircle(
				p.x,
				p.y,
				4,
				map[string]string{
					"fill": lc.options.Colors.Line,
				},
			)
		}
	}
}