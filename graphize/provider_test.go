package graphize

import (
	"testing"

	"github.com/plexusone/system-spec/spec"
)

func TestExtractSystem(t *testing.T) {
	sys := &spec.System{
		Name:        "test-system",
		Description: "Test system",
		Services: map[string]spec.Service{
			"api": {
				Image: spec.ContainerImage{Name: "api", Tag: "v1.0"},
				Repo: &spec.GitRepo{
					URL: "https://github.com/org/api",
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
					SQS: []spec.SQSQueue{
						{Name: "events"},
					},
				},
			},
		},
		Deployments: &spec.Deployments{
			Helm: map[string]spec.HelmDeployment{
				"app-chart": {
					Chart:    "myorg/app",
					Version:  "1.0.0",
					Services: []string{"api", "backend"},
				},
			},
		},
	}

	p := NewProvider()
	nodes, edges, err := p.ExtractSystem(sys, "system.json")
	if err != nil {
		t.Fatalf("ExtractSystem failed: %v", err)
	}

	// Count expected nodes:
	// 1 system + 2 services + 1 RDS + 1 SQS + 1 Helm = 6
	expectedNodes := 6
	if len(nodes) != expectedNodes {
		t.Errorf("expected %d nodes, got %d", expectedNodes, len(nodes))
		for _, n := range nodes {
			t.Logf("  node: %s (%s)", n.ID, n.Type)
		}
	}

	// Count expected edges:
	// 1 system->api + 1 system->backend + 1 api->backend (connects)
	// + 1 api->repo (links_to)
	// + 1 backend->rds + 1 backend->sqs
	// + 1 system->helm + 2 helm->services = 9
	expectedEdges := 9
	if len(edges) != expectedEdges {
		t.Errorf("expected %d edges, got %d", expectedEdges, len(edges))
		for _, e := range edges {
			t.Logf("  edge: %s -> %s (%s)", e.From, e.To, e.Type)
		}
	}

	// Verify service node
	var apiNode *struct{ ID, Type string }
	for _, n := range nodes {
		if n.ID == "svc:api" {
			apiNode = &struct{ ID, Type string }{n.ID, n.Type}
			break
		}
	}
	if apiNode == nil {
		t.Error("expected node 'svc:api'")
	}

	// Verify connection edge
	var connEdge bool
	for _, e := range edges {
		if e.From == "svc:api" && e.To == "svc:backend" && e.Type == EdgeTypeConnectsTo {
			connEdge = true
			if e.Attrs["protocol"] != "grpc" {
				t.Errorf("expected protocol 'grpc', got %q", e.Attrs["protocol"])
			}
			break
		}
	}
	if !connEdge {
		t.Error("expected connection edge from api to backend")
	}
}

func TestCanExtract(t *testing.T) {
	p := NewProvider()

	// Non-JSON file
	if p.CanExtract("file.go") {
		t.Error("should not extract .go files")
	}

	if p.CanExtract("nonexistent.json") {
		t.Error("should not extract nonexistent files")
	}
}

func TestProviderLanguage(t *testing.T) {
	p := NewProvider()

	if p.Language() != "system-spec" {
		t.Errorf("expected language 'system-spec', got %q", p.Language())
	}

	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != ".json" {
		t.Errorf("expected extensions ['.json'], got %v", exts)
	}
}
