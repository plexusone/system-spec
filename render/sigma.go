package render

import (
	"encoding/json"
	"math"

	"github.com/plexusone/system-spec/graph"
)

// SigmaRenderer renders graphs to Sigma.js JSON format.
// Optimized for large graphs with force-directed layout hints.
// See: https://www.sigmajs.org/
type SigmaRenderer struct {
	// Width of the layout canvas
	Width float64
	// Height of the layout canvas
	Height float64
}

// NewSigmaRenderer creates a SigmaRenderer with default settings.
func NewSigmaRenderer() *SigmaRenderer {
	return &SigmaRenderer{
		Width:  1000,
		Height: 1000,
	}
}

// Format returns the output format.
func (r *SigmaRenderer) Format() Format {
	return FormatSigma
}

// SigmaOutput is the JSON structure for Sigma.js.
type SigmaOutput struct {
	Nodes []SigmaNode `json:"nodes"`
	Edges []SigmaEdge `json:"edges"`
}

// SigmaNode represents a node in Sigma format.
type SigmaNode struct {
	ID         string            `json:"id"`
	Label      string            `json:"label"`
	X          float64           `json:"x"`
	Y          float64           `json:"y"`
	Size       float64           `json:"size"`
	Color      string            `json:"color"`
	Kind       string            `json:"kind"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// SigmaEdge represents an edge in Sigma format.
type SigmaEdge struct {
	ID       string  `json:"id"`
	Source   string  `json:"source"`
	Target   string  `json:"target"`
	Label    string  `json:"label,omitempty"`
	Size     float64 `json:"size"`
	Color    string  `json:"color"`
	Kind     string  `json:"kind"`
	Protocol string  `json:"protocol,omitempty"`
	Port     int     `json:"port,omitempty"`
}

// Render converts the graph to Sigma.js JSON.
func (r *SigmaRenderer) Render(g *graph.Graph) ([]byte, error) {
	out := SigmaOutput{
		Nodes: make([]SigmaNode, 0, len(g.Nodes)),
		Edges: make([]SigmaEdge, 0, len(g.Edges)),
	}

	// Calculate initial positions (circular layout)
	nodeCount := len(g.Nodes)
	centerX := r.Width / 2
	centerY := r.Height / 2
	radius := math.Min(r.Width, r.Height) * 0.4

	for i, n := range g.Nodes {
		// Circular layout
		angle := 2 * math.Pi * float64(i) / float64(nodeCount)
		x := centerX + radius*math.Cos(angle)
		y := centerY + radius*math.Sin(angle)

		// Calculate node size based on connections
		edgeCount := len(g.EdgesFrom(n.ID)) + len(g.EdgesTo(n.ID))
		size := math.Max(5, math.Min(20, float64(5+edgeCount*2)))

		node := SigmaNode{
			ID:    n.ID,
			Label: n.Label,
			X:     x,
			Y:     y,
			Size:  size,
			Color: r.colorFor(n.Kind),
			Kind:  string(n.Kind),
		}

		// Copy metadata to attributes
		if len(n.Metadata) > 0 {
			node.Attributes = make(map[string]string, len(n.Metadata))
			for k, v := range n.Metadata {
				node.Attributes[k] = v
			}
		}

		out.Nodes = append(out.Nodes, node)
	}

	for _, e := range g.Edges {
		edge := SigmaEdge{
			ID:       e.ID,
			Source:   e.Source,
			Target:   e.Target,
			Label:    e.Label,
			Size:     1,
			Color:    r.edgeColorFor(e.Kind),
			Kind:     string(e.Kind),
			Protocol: e.Protocol,
			Port:     e.Port,
		}
		out.Edges = append(out.Edges, edge)
	}

	return json.MarshalIndent(out, "", "  ")
}

func (r *SigmaRenderer) colorFor(k graph.NodeKind) string {
	switch k {
	case graph.NodeKindService:
		return "#4fc3f7" // light blue
	case graph.NodeKindDatabase:
		return "#ffb74d" // orange
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return "#ba68c8" // purple
	case graph.NodeKindStorage:
		return "#81c784" // green
	case graph.NodeKindAIModel:
		return "#f06292" // pink
	case graph.NodeKindCDN:
		return "#4dd0e1" // cyan
	case graph.NodeKindWorker:
		return "#9575cd" // deep purple
	case graph.NodeKindHelm, graph.NodeKindTerraform:
		return "#a1887f" // brown
	default:
		return "#90a4ae" // blue grey
	}
}

func (r *SigmaRenderer) edgeColorFor(k graph.EdgeKind) string {
	switch k {
	case graph.EdgeKindConnection:
		return "#78909c" // blue grey
	case graph.EdgeKindDatabase:
		return "#ffcc80" // light orange
	case graph.EdgeKindQueue:
		return "#ce93d8" // light purple
	case graph.EdgeKindStorage:
		return "#a5d6a7" // light green
	case graph.EdgeKindDeploys:
		return "#bcaaa4" // light brown
	default:
		return "#b0bec5" // blue grey lighter
	}
}
