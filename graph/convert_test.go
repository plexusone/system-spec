package graph

import (
	"testing"

	"github.com/plexusone/system-spec/spec"
)

func TestFromSystem(t *testing.T) {
	sys := &spec.System{
		Name: "test-system",
		Services: map[string]spec.Service{
			"api": {
				Image: spec.ContainerImage{Name: "nginx", Tag: "1.25"},
				Repo: &spec.GitRepo{
					URL: "https://github.com/example/api",
				},
				Connections: map[string]spec.Connection{
					"backend": {Port: 8080, Protocol: "grpc"},
				},
			},
			"backend": {
				Image: spec.ContainerImage{Name: "backend", Tag: "v1.0"},
				AWS: &spec.AWSResources{
					RDS: []spec.RDSInstance{
						{Name: "main-db", Engine: "postgres", Port: 5432},
					},
				},
			},
		},
	}

	g := FromSystem(sys)

	// Check node count: 2 services + 1 RDS
	if len(g.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(g.Nodes))
	}

	// Check edge count: 1 service-to-service + 1 service-to-db
	if len(g.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(g.Edges))
	}

	// Verify service node
	apiNode := g.NodeByID("svc:api")
	if apiNode == nil {
		t.Fatal("expected node 'svc:api'")
	}
	if apiNode.Kind != NodeKindService {
		t.Errorf("expected kind 'service', got %q", apiNode.Kind)
	}
	if apiNode.Metadata["image"] != "nginx:1.25" {
		t.Errorf("expected image 'nginx:1.25', got %q", apiNode.Metadata["image"])
	}

	// Verify database node
	dbNode := g.NodeByID("rds:main-db")
	if dbNode == nil {
		t.Fatal("expected node 'rds:main-db'")
	}
	if dbNode.Kind != NodeKindDatabase {
		t.Errorf("expected kind 'database', got %q", dbNode.Kind)
	}
}

func TestEdgesFrom(t *testing.T) {
	g := NewGraph()
	g.AddNode(Node{ID: "a", Label: "A", Kind: NodeKindService})
	g.AddNode(Node{ID: "b", Label: "B", Kind: NodeKindService})
	g.AddNode(Node{ID: "c", Label: "C", Kind: NodeKindService})

	g.AddEdge(Edge{ID: "a->b", Source: "a", Target: "b", Kind: EdgeKindConnection})
	g.AddEdge(Edge{ID: "a->c", Source: "a", Target: "c", Kind: EdgeKindConnection})
	g.AddEdge(Edge{ID: "b->c", Source: "b", Target: "c", Kind: EdgeKindConnection})

	edges := g.EdgesFrom("a")
	if len(edges) != 2 {
		t.Errorf("expected 2 edges from 'a', got %d", len(edges))
	}

	edges = g.EdgesTo("c")
	if len(edges) != 2 {
		t.Errorf("expected 2 edges to 'c', got %d", len(edges))
	}
}
