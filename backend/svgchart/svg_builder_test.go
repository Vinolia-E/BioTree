package svgchart

import (
    "strings"
    "testing"
)

func TestSVGBuilder(t *testing.T) {
    t.Run("Basic SVG generation", func(t *testing.T) {
        sb := NewSVGBuilder(100, 100)
        svg := sb.String()
        
        if !strings.Contains(svg, `<svg width="100" height="100"`) {
            t.Error("SVG should contain correct width and height")
        }
        if !strings.Contains(svg, `</svg>`) {
            t.Error("SVG should be properly closed")
        }
    })
    
    t.Run("Add elements", func(t *testing.T) {
        sb := NewSVGBuilder(100, 100)
        sb.AddRect(10, 10, 80, 80, map[string]string{"fill": "blue"})
        sb.AddCircle(50, 50, 25, map[string]string{"fill": "red"})
        sb.AddLine(10, 10, 90, 90, map[string]string{"stroke": "black"})
        sb.AddText(50, 50, "Test", map[string]string{"text-anchor": "middle"})
        sb.AddPath("M10,10 L90,90", map[string]string{"stroke": "green"})
        
        svg := sb.String()
        
        if !strings.Contains(svg, "<rect") || !strings.Contains(svg, `fill="blue"`) {
            t.Error("SVG should contain rect element with fill=blue")
        }
        
        if !strings.Contains(svg, "<circle") || !strings.Contains(svg, `fill="red"`) {
            t.Error("SVG should contain circle element with fill=red")
        }
        
        if !strings.Contains(svg, "<line") || !strings.Contains(svg, `stroke="black"`) {
            t.Error("SVG should contain line element with stroke=black")
        }
        
        if !strings.Contains(svg, "<text") || !strings.Contains(svg, `text-anchor="middle"`) {
            t.Error("SVG should contain text element with text-anchor=middle")
        }
        
        if !strings.Contains(svg, "<path") || !strings.Contains(svg, `stroke="green"`) {
            t.Error("SVG should contain path element with stroke=green")
        }
        
        if !strings.Contains(svg, ">Test<") {
            t.Error("SVG should contain text with content 'Test'")
        }
    })
}