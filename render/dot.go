package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/plexusone/system-spec/graph"
)

// DOTRenderer renders graphs to GraphViz DOT format.
// See: https://graphviz.org/
type DOTRenderer struct {
	// Rankdir: LR, TB, RL, BT
	Rankdir string
}

// NewDOTRenderer creates a DOTRenderer with default settings.
func NewDOTRenderer() *DOTRenderer {
	return &DOTRenderer{
		Rankdir: "LR",
	}
}

// Format returns the output format.
func (r *DOTRenderer) Format() Format {
	return FormatDOT
}

// Render converts the graph to GraphViz DOT format.
func (r *DOTRenderer) Render(g *graph.Graph) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("digraph system {\n")
	fmt.Fprintf(&buf, "    rankdir=%s;\n", r.Rankdir)
	buf.WriteString("    node [fontname=\"Helvetica\", fontsize=10];\n")
	buf.WriteString("    edge [fontname=\"Helvetica\", fontsize=8];\n")
	buf.WriteString("\n")

	// Nodes
	for _, n := range g.Nodes {
		id := sanitizeDOTID(n.ID)
		shape := r.shapeFor(n.Kind)
		color := r.colorFor(n.Kind)
		style := r.styleFor(n.Kind)

		fmt.Fprintf(&buf, "    %s [label=\"%s\", shape=%s, style=\"%s\", fillcolor=\"%s\"];\n",
			id, n.Label, shape, style, color)
	}

	buf.WriteString("\n")

	// Edges
	for _, e := range g.Edges {
		srcID := sanitizeDOTID(e.Source)
		tgtID := sanitizeDOTID(e.Target)

		if e.Label != "" {
			fmt.Fprintf(&buf, "    %s -> %s [label=\"%s\"];\n", srcID, tgtID, e.Label)
		} else {
			fmt.Fprintf(&buf, "    %s -> %s;\n", srcID, tgtID)
		}
	}

	buf.WriteString("}\n")

	return buf.Bytes(), nil
}

func (r *DOTRenderer) shapeFor(k graph.NodeKind) string {
	switch k {
	case graph.NodeKindDatabase:
		return "cylinder"
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return "cds"
	case graph.NodeKindStorage:
		return "folder"
	case graph.NodeKindAIModel:
		return "hexagon"
	case graph.NodeKindCDN:
		return "cloud"
	case graph.NodeKindWorker:
		return "component"
	case graph.NodeKindHelm, graph.NodeKindTerraform:
		return "box3d"
	default:
		return "box"
	}
}

func (r *DOTRenderer) colorFor(k graph.NodeKind) string {
	switch k {
	case graph.NodeKindService:
		return "#e1f5fe"
	case graph.NodeKindDatabase:
		return "#fff3e0"
	case graph.NodeKindQueue, graph.NodeKindTopic:
		return "#f3e5f5"
	case graph.NodeKindStorage:
		return "#e8f5e9"
	case graph.NodeKindAIModel:
		return "#fce4ec"
	case graph.NodeKindCDN:
		return "#e0f7fa"
	case graph.NodeKindWorker:
		return "#ede7f6"
	case graph.NodeKindHelm, graph.NodeKindTerraform:
		return "#efebe9"
	default:
		return "#eceff1"
	}
}

func (r *DOTRenderer) styleFor(k graph.NodeKind) string {
	return "filled,rounded"
}

// sanitizeDOTID makes an ID safe for GraphViz DOT
func sanitizeDOTID(id string) string {
	// DOT IDs need to be quoted if they contain special chars
	// or start with a number
	id = strings.ReplaceAll(id, ":", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	return id
}
