package svgchart

import "testing"

func TestOptions(t *testing.T) {
    t.Run("DefaultOptions", func(t *testing.T) {
        opts := DefaultOptions()
        if opts.Width != 500 || opts.Height != 300 {
            t.Error("Default dimensions should be 500x300")
        }
        if !opts.ShowGrid {
            t.Error("Default should show grid")
        }
        if opts.ChartType != Line {
            t.Error("Default chart type should be Line")
        }
    })
    
    t.Run("WithTitle", func(t *testing.T) {
        opts := DefaultOptions()
        WithTitle("Test Title")(&opts)
        if opts.Title != "Test Title" {
            t.Error("Title should be set correctly")
        }
        if opts.Margins.Top != 40 {
            t.Error("Top margin should be adjusted for title")
        }
        
        WithTitle("")(&opts)
        if opts.Margins.Top != 20 {
            t.Error("Top margin should be reduced for empty title")
        }
    })
    
    t.Run("WithDimensions", func(t *testing.T) {
        opts := DefaultOptions()
        WithDimensions(800, 600)(&opts)
        if opts.Width != 800 || opts.Height != 600 {
            t.Error("Dimensions should be set correctly")
        }
    })
    
    t.Run("WithLabels", func(t *testing.T) {
        opts := DefaultOptions()
        WithXLabel("X Axis")(&opts)
        WithYLabel("Y Axis")(&opts)
        if opts.XLabel != "X Axis" || opts.YLabel != "Y Axis" {
            t.Error("Labels should be set correctly")
        }
        if opts.Margins.Bottom != 50 || opts.Margins.Left != 60 {
            t.Error("Margins should be adjusted for labels")
        }
        
        WithXLabel("")(&opts)
        WithYLabel("")(&opts)
        if opts.Margins.Bottom != 30 || opts.Margins.Left != 40 {
            t.Error("Margins should be reduced for empty labels")
        }
    })
    
    t.Run("WithGrid", func(t *testing.T) {
        opts := DefaultOptions()
        WithGrid(false)(&opts)
        if opts.ShowGrid {
            t.Error("ShowGrid should be set to false")
        }
    })
    
    t.Run("WithColors", func(t *testing.T) {
        opts := DefaultOptions()
        newColors := ColorScheme{
            Background: "#000000",
            Axis:       "#ffffff",
        }
        WithColors(newColors)(&opts)
        if opts.Colors.Background != "#000000" || opts.Colors.Axis != "#ffffff" {
            t.Error("Colors should be updated correctly")
        }
    })
}