package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/plexusone/system-spec/graph"
)

// MermaidRenderer renders graphs to Mermaid diagram syntax.
// See: https://mermaid.js.org/
type MermaidRenderer struct {
	// Direction: LR (left-right), TB (top-bottom), RL, BT
	Direction string
}

// NewMermaidRenderer creates a MermaidRenderer with default settings.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		Direction: "LR",
	}
}

// Format returns the output format.
func (r *MermaidRenderer) Format() Format {
	return FormatMermaid
}

// Render converts the graph to Mermaid format.
func (r *MermaidRenderer) Render(g *graph.Graph) ([]byte, error) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "graph %s\n", r.Direction)

	// Subgraphs by provider (optional grouping)
	providers := make(map[string][]graph.Node)
	var services []graph.Node

	for _, n := range g.Nodes {
		if n.Kind == graph.NodeKindService {
			services = append(services, n)
		} else if provider, ok := n.Metadata["provider"]; ok {
			providers[provider] = append(providers[provider], n)
		} else {
			services = append(services, n)
		}
	}

	// Services subgraph
	if len(services) > 0 {
		buf.WriteString("    subgraph Services\n")
		for _, n := range services {
			id := sanitizeMermaidID(n.ID)
			shape := r.shapeFor(n.Kind, id, n.Label)
			fmt.Fprintf(&buf, "        %s\n", shape)
		}
		buf.WriteString("    end\n")
	}

	// Provider subgraphs
	for provider, nodes := range providers {
		fmt.Fprintf(&buf, "    subgraph %s\n", strings.ToUpper(provider))
		for _, n := range nodes {
			id := sanitizeMermaidID(n.ID)
			shape := r.shapeFor(n.Kind, id, n.Label)
			fmt.Fprintf(&buf, "        %s\n", shape)
		}
		buf.WriteString("    end\n")
	}

	buf.WriteString("\n")

	// Edges
	for _, e := range g.Edges {
		srcID := sanitizeMermaidID(e.Source)
		tgtID := sanitizeMermaidID(e.Target)

		if e.Label != "" {
			fmt.Fprintf(&buf, "    %s -->|%s| %s\n", srcID, e.Label, tgtID)
		} else {
			fmt.Fprintf(&buf, "    %s --> %s\n", srcID, tgtID)
		}
	}

	return buf.Bytes(), nil
}

func (r *MermaidRenderer) shapeFor(k graph.NodeKind, id, label string) string {
	switch k {
	case graph.NodeKindDatabase:
		return fmt.Sprintf("%s[(%s)]", id, label) // cylinder
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return fmt.Sprintf("%s>%s]", id, label) // asymmetric
	case graph.NodeKindStorage:
		return fmt.Sprintf("%s[(%s)]", id, label) // cylinder
	case graph.NodeKindAIModel:
		return fmt.Sprintf("%s{{%s}}", id, label) // hexagon
	case graph.NodeKindHelm, graph.NodeKindTerraform:
		return fmt.Sprintf("%s[/%s/]", id, label) // parallelogram
	default:
		return fmt.Sprintf("%s[%s]", id, label) // rectangle
	}
}

// sanitizeMermaidID makes an ID safe for Mermaid
func sanitizeMermaidID(id string) string {
	// Mermaid IDs can't have colons or special chars
	id = strings.ReplaceAll(id, ":", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	return id
}
