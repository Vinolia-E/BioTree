package svgchart

import (
	"fmt"
	"strings"
)

// SVGBuilder efficiently constructs SVG content
type SVGBuilder struct {
	builder strings.Builder
}

// NewSVGBuilder creates a new SVGBuilder with initialized SVG tag
func NewSVGBuilder(width, height int) *SVGBuilder {
	sb := &SVGBuilder{}
	sb.builder.Grow(2048)
	sb.builder.WriteString(fmt.Sprintf(
		`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`,
		width, height,
	))
	return sb
}

// AddElement adds an SVG element with attributes and content
func (sb *SVGBuilder) AddElement(tag string, attrs map[string]string, content string) {
	sb.builder.WriteString("<")
	sb.builder.WriteString(tag)
	sb.writeAttrs(attrs)
	if content == "" {
		sb.builder.WriteString("/>")
	} else {
		sb.builder.WriteString(">")
		sb.builder.WriteString(content)
		sb.builder.WriteString("</")
		sb.builder.WriteString(tag)
		sb.builder.WriteString(">")
	}
}

// AddRect adds a rectangle element
func (sb *SVGBuilder) AddRect(x, y, width, height int, attrs map[string]string) {
	attrs["x"] = fmt.Sprintf("%d", x)
	attrs["y"] = fmt.Sprintf("%d", y)
	attrs["width"] = fmt.Sprintf("%d", width)
	attrs["height"] = fmt.Sprintf("%d", height)
	sb.AddElement("rect", attrs, "")
}

// AddText adds a text element
func (sb *SVGBuilder) AddText(x, y int, text string, attrs map[string]string) {
	attrs["x"] = fmt.Sprintf("%d", x)
	attrs["y"] = fmt.Sprintf("%d", y)
	sb.AddElement("text", attrs, text)
}

// AddLine adds a line element
func (sb *SVGBuilder) AddLine(x1, y1, x2, y2 int, attrs map[string]string) {
	attrs["x1"] = fmt.Sprintf("%d", x1)
	attrs["x2"] = fmt.Sprintf("%d", x2)
	attrs["y1"] = fmt.Sprintf("%d", y1)
	attrs["y2"] = fmt.Sprintf("%d", y2)
	sb.AddElement("line", attrs, "")
}

// AddPath adds a path element
func (sb *SVGBuilder) AddPath(d string, attrs map[string]string) {
	attrs["d"] = d
	sb.AddElement("path", attrs, "")
}

// AddCircle adds a circle element
func (sb *SVGBuilder) AddCircle(cx, cy, r int, attrs map[string]string) {
	attrs["cx"] = fmt.Sprintf("%d", cx)
	attrs["cy"] = fmt.Sprintf("%d", cy)
	attrs["r"] = fmt.Sprintf("%d", r)
	sb.AddElement("circle", attrs, "")
}

// String returns the complete SVG content
func (sb *SVGBuilder) String() string {
	sb.builder.WriteString("</svg>")
	return sb.builder.String()
}

func (sb *SVGBuilder) writeAttrs(attrs map[string]string) {
	for k, v := range attrs {
		sb.builder.WriteString(" ")
		sb.builder.WriteString(k)
		sb.builder.WriteString(`="`)
		sb.builder.WriteString(v)
		sb.builder.WriteString(`"`)
	}
}
