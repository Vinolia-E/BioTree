package svgchart

import (
	"fmt"
	"math"
	"strings"
)

type pieChart struct {
	data    []DataPoint
	options Options
}

func newPieChart(data []DataPoint, options Options) Chart {
	return &pieChart{
		data:    data,
		options: options,
	}
}

func (p *pieChart) Generate() string {
	if len(p.data) == 0 {
		return ""
	}

	width := float64(p.options.Width)
	height := float64(p.options.Height)
	radius := math.Min(width, height) * 0.4
	centerX := width / 2
	centerY := height / 2

	var total float64
	for _, d := range p.data {
		total += d.Value
	}

	var paths []string
	var legends []string
	startAngle := -90.0 // Start from top (12 o'clock position)
	colors := []string{"#4285F4", "#34A853", "#FBBC05", "#EA4335", "#673AB7", "#3F51B5", "#2196F3", "#03A9F4"}

	for i, d := range p.data {
		percentage := (d.Value / total) * 100
		endAngle := startAngle + (d.Value/total)*360

		// Calculate path
		startRad := startAngle * math.Pi / 180
		endRad := endAngle * math.Pi / 180

		x1 := centerX + radius*math.Cos(startRad)
		y1 := centerY + radius*math.Sin(startRad)
		x2 := centerX + radius*math.Cos(endRad)
		y2 := centerY + radius*math.Sin(endRad)

		largeArcFlag := 0
		if endAngle-startAngle > 180 {
			largeArcFlag = 1
		}

		color := colors[i%len(colors)]
		path := fmt.Sprintf(`<path d="M %f %f L %f %f A %f %f 0 %d 1 %f %f L %f %f Z" fill="%s" stroke="white" stroke-width="1"/>`,
			centerX, centerY, x1, y1, radius, radius, largeArcFlag, x2, y2, centerX, centerY, color)
		paths = append(paths, path)

		// Add legend
		legend := fmt.Sprintf(`<g transform="translate(%f, %f)">
			<rect width="10" height="10" fill="%s"/>
			<text x="15" y="9" font-size="12">%s (%.1f%%)</text>
		</g>`, width-120, float64(i)*20+40, color, d.Label, percentage)
		legends = append(legends, legend)

		startAngle = endAngle
	}

	// Create SVG with title, chart, and legend
	svg := fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<style>
			text { font-family: Arial, sans-serif; }
		</style>
		<text x="%f" y="30" text-anchor="middle" font-size="16" font-weight="bold">%s</text>
		<g>%s</g>
		<g>%s</g>
	</svg>`,
		p.options.Width, p.options.Height,
		centerX, p.options.Title,
		strings.Join(paths, "\n"),
		strings.Join(legends, "\n"))

	return svg
}
