package render

import (
	"encoding/json"

	"github.com/plexusone/system-spec/graph"
)

// CytoscapeRenderer renders graphs to Cytoscape.js JSON format.
// See: https://js.cytoscape.org/
type CytoscapeRenderer struct{}

// NewCytoscapeRenderer creates a CytoscapeRenderer.
func NewCytoscapeRenderer() *CytoscapeRenderer {
	return &CytoscapeRenderer{}
}

// Format returns the output format.
func (r *CytoscapeRenderer) Format() Format {
	return FormatCytoscape
}

// CytoscapeOutput is the JSON structure for Cytoscape.js.
type CytoscapeOutput struct {
	Elements CytoscapeElements `json:"elements"`
}

// CytoscapeElements contains nodes and edges.
type CytoscapeElements struct {
	Nodes []CytoscapeNode `json:"nodes"`
	Edges []CytoscapeEdge `json:"edges"`
}

// CytoscapeNode represents a node in Cytoscape format.
type CytoscapeNode struct {
	Data CytoscapeNodeData `json:"data"`
}

// CytoscapeNodeData contains node properties.
type CytoscapeNodeData struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Kind     string `json:"kind"`
	Group    string `json:"group,omitempty"`
	Provider string `json:"provider,omitempty"`
	Image    string `json:"image,omitempty"`
	Repo     string `json:"repo,omitempty"`
}

// CytoscapeEdge represents an edge in Cytoscape format.
type CytoscapeEdge struct {
	Data CytoscapeEdgeData `json:"data"`
}

// CytoscapeEdgeData contains edge properties.
type CytoscapeEdgeData struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Label    string `json:"label,omitempty"`
	Kind     string `json:"kind"`
	Protocol string `json:"protocol,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// Render converts the graph to Cytoscape.js JSON.
func (r *CytoscapeRenderer) Render(g *graph.Graph) ([]byte, error) {
	out := CytoscapeOutput{
		Elements: CytoscapeElements{
			Nodes: make([]CytoscapeNode, 0, len(g.Nodes)),
			Edges: make([]CytoscapeEdge, 0, len(g.Edges)),
		},
	}

	for _, n := range g.Nodes {
		node := CytoscapeNode{
			Data: CytoscapeNodeData{
				ID:    n.ID,
				Label: n.Label,
				Kind:  string(n.Kind),
				Group: n.Group,
			},
		}

		// Copy relevant metadata
		if v, ok := n.Metadata["provider"]; ok {
			node.Data.Provider = v
		}
		if v, ok := n.Metadata["image"]; ok {
			node.Data.Image = v
		}
		if v, ok := n.Metadata["repo"]; ok {
			node.Data.Repo = v
		}

		out.Elements.Nodes = append(out.Elements.Nodes, node)
	}

	for _, e := range g.Edges {
		edge := CytoscapeEdge{
			Data: CytoscapeEdgeData{
				ID:       e.ID,
				Source:   e.Source,
				Target:   e.Target,
				Label:    e.Label,
				Kind:     string(e.Kind),
				Protocol: e.Protocol,
				Port:     e.Port,
			},
		}
		out.Elements.Edges = append(out.Elements.Edges, edge)
	}

	return json.MarshalIndent(out, "", "  ")
}
