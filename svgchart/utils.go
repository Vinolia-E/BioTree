package svgchart

import (
	"fmt"
	"math"
	"strings"
)

// getYMinMax returns the minimum and maximum Y values from the data
func getYMinMax(data ChartData) (float64, float64) {
	if len(data) == 0 {
		return 0, 0
	}

	min, max := data[0].Value, data[0].Value
	for _, dp := range data {
		if dp.Value < min {
			min = dp.Value
		}
		if dp.Value > max {
			max = dp.Value
		}
	}
	return min, max
}

// formatNumber formats a float64 for display
func formatNumber(value float64) string {
	if value == math.Floor(value) {
		return fmt.Sprintf("%.0f", value)
	}
	
	formatted := fmt.Sprintf("%.2f", value)
	// Remove trailing zeros
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return formatted
}

// generateEmptyChart creates an SVG with a message for empty data
func generateEmptyChart(options Options, message string) string {
	sb := NewSVGBuilder(options.Width, options.Height)
	sb.AddRect(0, 0, options.Width, options.Height, map[string]string{
		"fill": options.Colors.Background,
	})
	sb.AddText(
		options.Width/2,
		options.Height/2,
		message,
		map[string]string{
			"text-anchor":  "middle",
			"font-family":  "Arial",
			"font-size":    "14px",
			"fill":        "#666666",
		},
	)
	return sb.String()
}