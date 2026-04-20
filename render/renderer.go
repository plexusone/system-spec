// Package render provides multiple output formats for system-spec graphs.
package render

import (
	"github.com/plexusone/system-spec/graph"
)

// Format identifies the output format.
type Format string

const (
	FormatD2        Format = "d2"
	FormatMermaid   Format = "mermaid"
	FormatCytoscape Format = "cytoscape"
	FormatSigma     Format = "sigma"
	FormatDOT       Format = "dot"
	FormatJSON      Format = "json"
)

// Renderer converts a Graph to a specific output format.
type Renderer interface {
	// Format returns the output format identifier.
	Format() Format

	// Render converts the graph to the output format.
	Render(g *graph.Graph) ([]byte, error)
}

// Renderers provides access to all available renderers.
// This is a concrete struct, not a registry with interface{}.
type Renderers struct {
	D2        *D2Renderer
	Mermaid   *MermaidRenderer
	Cytoscape *CytoscapeRenderer
	Sigma     *SigmaRenderer
	DOT       *DOTRenderer
}

// NewRenderers creates a Renderers instance with default settings.
func NewRenderers() *Renderers {
	return &Renderers{
		D2:        NewD2Renderer(),
		Mermaid:   NewMermaidRenderer(),
		Cytoscape: NewCytoscapeRenderer(),
		Sigma:     NewSigmaRenderer(),
		DOT:       NewDOTRenderer(),
	}
}

// Get returns a renderer by format, or nil if not found.
func (r *Renderers) Get(f Format) Renderer {
	switch f {
	case FormatD2:
		return r.D2
	case FormatMermaid:
		return r.Mermaid
	case FormatCytoscape:
		return r.Cytoscape
	case FormatSigma:
		return r.Sigma
	case FormatDOT:
		return r.DOT
	default:
		return nil
	}
}

// SupportedFormats returns all supported output formats.
func SupportedFormats() []Format {
	return []Format{
		FormatD2,
		FormatMermaid,
		FormatCytoscape,
		FormatSigma,
		FormatDOT,
		FormatJSON,
	}
}
