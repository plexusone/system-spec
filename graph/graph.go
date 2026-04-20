// Package graph provides an intermediate graph representation
// for rendering system-spec to various diagram formats.
package graph

// Graph is the intermediate representation for rendering.
// It decouples the spec structure from renderer output formats.
type Graph struct {
	// Nodes in the graph
	Nodes []Node `json:"nodes"`

	// Edges connecting nodes
	Edges []Edge `json:"edges"`
}

// Node represents a vertex in the graph.
type Node struct {
	// ID is the unique node identifier
	ID string `json:"id"`

	// Label is the display name
	Label string `json:"label"`

	// Kind categorizes the node for styling
	Kind NodeKind `json:"kind"`

	// Group for clustering (e.g., VPC name, namespace)
	Group string `json:"group,omitempty"`

	// Metadata contains additional key-value data for renderers
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NodeKind categorizes nodes for rendering and styling.
type NodeKind string

const (
	NodeKindService    NodeKind = "service"
	NodeKindDatabase   NodeKind = "database"
	NodeKindQueue      NodeKind = "queue"
	NodeKindTopic      NodeKind = "topic"
	NodeKindStorage    NodeKind = "storage"
	NodeKindCDN        NodeKind = "cdn"
	NodeKindWorker     NodeKind = "worker"
	NodeKindAIModel    NodeKind = "ai_model"
	NodeKindVPC        NodeKind = "vpc"
	NodeKindSubnet     NodeKind = "subnet"
	NodeKindExternal   NodeKind = "external"
	NodeKindHelm       NodeKind = "helm"
	NodeKindTerraform  NodeKind = "terraform"
)

// Edge represents a directed edge between nodes.
type Edge struct {
	// ID is the unique edge identifier
	ID string `json:"id"`

	// Source node ID
	Source string `json:"source"`

	// Target node ID
	Target string `json:"target"`

	// Label for display
	Label string `json:"label,omitempty"`

	// Kind categorizes the edge
	Kind EdgeKind `json:"kind"`

	// Protocol (http, grpc, tcp, sql, amqp, redis, etc.)
	Protocol string `json:"protocol,omitempty"`

	// Port number
	Port int `json:"port,omitempty"`

	// Metadata contains additional key-value data
	Metadata map[string]string `json:"metadata,omitempty"`
}

// EdgeKind categorizes edges for rendering and styling.
type EdgeKind string

const (
	EdgeKindConnection EdgeKind = "connection"   // service to service
	EdgeKindDatabase   EdgeKind = "database"     // service to database
	EdgeKindQueue      EdgeKind = "queue"        // service to queue
	EdgeKindStorage    EdgeKind = "storage"      // service to storage
	EdgeKindDeploys    EdgeKind = "deploys"      // helm/terraform to service
)

// NewGraph creates an empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(n Node) {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string)
	}
	g.Nodes = append(g.Nodes, n)
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(e Edge) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	g.Edges = append(g.Edges, e)
}

// NodeByID finds a node by ID, returns nil if not found.
func (g *Graph) NodeByID(id string) *Node {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i]
		}
	}
	return nil
}

// EdgesFrom returns all edges originating from the given node ID.
func (g *Graph) EdgesFrom(nodeID string) []Edge {
	var result []Edge
	for _, e := range g.Edges {
		if e.Source == nodeID {
			result = append(result, e)
		}
	}
	return result
}

// EdgesTo returns all edges pointing to the given node ID.
func (g *Graph) EdgesTo(nodeID string) []Edge {
	var result []Edge
	for _, e := range g.Edges {
		if e.Target == nodeID {
			result = append(result, e)
		}
	}
	return result
}
