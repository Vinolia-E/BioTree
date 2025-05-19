package chart

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	width     = 400
	height    = 200
	padding   = 40
	axisColor = "#333"
	lineColor = "#40e0d0"
	bgColor   = "#ffffff"
	fontSize  = 10
	gridColor = "#eee"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// GenerateLineChartSVG creates an SVG string for a line chart
func GenerateLineChartSVG(data map[string]float64) string {
	if len(data) == 0 {
		return `<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200"><text x="10" y="20" font-size="14" fill="red">No data available</text></svg>`
	}

	// Sort the keys (categorical x-axis)
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Compute min/max for scaling
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, v := range data {
		if v < minY {
			minY = v
		}
		if v > maxY {
			maxY = v
		}
	}
	if minY == maxY {
		minY -= 1
		maxY += 1
	}

	scaleX := float64(width-2*padding) / float64(len(keys)-1)
	scaleY := float64(height-2*padding) / (maxY - minY)

	// Generate points
	points := make([]Point, len(keys))
	for i, k := range keys {
		x := float64(padding) + float64(i)*scaleX
		y := float64(height-padding) - (data[k]-minY)*scaleY
		points[i] = Point{x, y}
	}

	// Build SVG parts
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s"/>`, bgColor))

	// Axes
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s"/>`, padding, padding, padding, height-padding, axisColor))              // Y axis
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s"/>`, padding, height-padding, width-padding, height-padding, axisColor)) // X axis

	// Gridlines & labels
	for i, k := range keys {
		x := float64(padding) + float64(i)*scaleX
		y := float64(height - padding)

		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="%s" stroke-dasharray="2,2"/>`, x, padding, x, height-padding, gridColor))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%d" text-anchor="middle">%s</text>`, x, y+fontSize+2, fontSize, k))
	}

	// Y-axis labels
	steps := 5
	for i := 0; i <= steps; i++ {
		val := minY + (maxY-minY)*float64(i)/float64(steps)
		y := float64(height-padding) - (val-minY)*scaleY
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%.1f" font-size="%d" text-anchor="end">%.0f</text>`, padding-5, y+3, fontSize, val))
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="%s" stroke-dasharray="2,2"/>`, padding, y, width-padding, y, gridColor))
	}

	return ""
}
