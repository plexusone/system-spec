package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/plexusone/system-spec/graph"
)

// D2Renderer renders graphs to D2 diagram language.
// See: https://d2lang.com/
type D2Renderer struct {
	// Direction of the diagram (right, down, left, up)
	Direction string
}

// NewD2Renderer creates a D2Renderer with default settings.
func NewD2Renderer() *D2Renderer {
	return &D2Renderer{
		Direction: "right",
	}
}

// Format returns the output format.
func (r *D2Renderer) Format() Format {
	return FormatD2
}

// Render converts the graph to D2 format.
func (r *D2Renderer) Render(g *graph.Graph) ([]byte, error) {
	var buf bytes.Buffer

	// Direction
	fmt.Fprintf(&buf, "direction: %s\n\n", r.Direction)

	// Style classes for different node kinds
	buf.WriteString("classes: {\n")
	buf.WriteString("  service: {\n")
	buf.WriteString("    shape: rectangle\n")
	buf.WriteString("    style.fill: \"#e1f5fe\"\n")
	buf.WriteString("  }\n")
	buf.WriteString("  database: {\n")
	buf.WriteString("    shape: cylinder\n")
	buf.WriteString("    style.fill: \"#fff3e0\"\n")
	buf.WriteString("  }\n")
	buf.WriteString("  queue: {\n")
	buf.WriteString("    shape: queue\n")
	buf.WriteString("    style.fill: \"#f3e5f5\"\n")
	buf.WriteString("  }\n")
	buf.WriteString("  storage: {\n")
	buf.WriteString("    shape: stored_data\n")
	buf.WriteString("    style.fill: \"#e8f5e9\"\n")
	buf.WriteString("  }\n")
	buf.WriteString("  ai_model: {\n")
	buf.WriteString("    shape: hexagon\n")
	buf.WriteString("    style.fill: \"#fce4ec\"\n")
	buf.WriteString("  }\n")
	buf.WriteString("}\n\n")

	// Nodes
	for _, n := range g.Nodes {
		id := sanitizeD2ID(n.ID)
		shape := r.shapeFor(n.Kind)
		class := r.classFor(n.Kind)

		if class != "" {
			fmt.Fprintf(&buf, "%s: %s { class: %s }\n", id, n.Label, class)
		} else {
			fmt.Fprintf(&buf, "%s: %s { shape: %s }\n", id, n.Label, shape)
		}
	}

	buf.WriteString("\n")

	// Edges
	for _, e := range g.Edges {
		srcID := sanitizeD2ID(e.Source)
		tgtID := sanitizeD2ID(e.Target)

		if e.Label != "" {
			fmt.Fprintf(&buf, "%s -> %s: %s\n", srcID, tgtID, e.Label)
		} else {
			fmt.Fprintf(&buf, "%s -> %s\n", srcID, tgtID)
		}
	}

	return buf.Bytes(), nil
}

func (r *D2Renderer) shapeFor(k graph.NodeKind) string {
	switch k {
	case graph.NodeKindDatabase:
		return "cylinder"
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return "queue"
	case graph.NodeKindStorage:
		return "stored_data"
	case graph.NodeKindCDN:
		return "cloud"
	case graph.NodeKindAIModel:
		return "hexagon"
	case graph.NodeKindVPC, graph.NodeKindSubnet:
		return "rectangle"
	case graph.NodeKindHelm, graph.NodeKindTerraform:
		return "package"
	default:
		return "rectangle"
	}
}

func (r *D2Renderer) classFor(k graph.NodeKind) string {
	switch k {
	case graph.NodeKindService:
		return "service"
	case graph.NodeKindDatabase:
		return "database"
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return "queue"
	case graph.NodeKindStorage:
		return "storage"
	case graph.NodeKindAIModel:
		return "ai_model"
	default:
		return ""
	}
}

// sanitizeD2ID makes an ID safe for D2 (replaces special chars)
func sanitizeD2ID(id string) string {
	// D2 IDs can't have colons, replace with underscores
	return strings.ReplaceAll(id, ":", "_")
}
