// Package graphize provides integration with the graphize code graph tool.
// It converts system-spec documents into graphfs nodes and edges, allowing
// infrastructure topology to be queried alongside code graphs.
package graphize

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/plexusone/graphfs/pkg/graph"
	"github.com/plexusone/system-spec/spec"
)

// Node types for system-spec entities
const (
	NodeTypeService   = "service"
	NodeTypeDatabase  = "database"
	NodeTypeQueue     = "queue"
	NodeTypeTopic     = "topic"
	NodeTypeStorage   = "storage"
	NodeTypeAIModel   = "ai_model"
	NodeTypeCDN       = "cdn"
	NodeTypeWorker    = "worker"
	NodeTypeHelm      = "helm_chart"
	NodeTypeTerraform = "terraform_module"
	NodeTypeSystem    = "system"
)

// Edge types for system-spec relationships
const (
	EdgeTypeConnectsTo = "connects_to"
	EdgeTypeUses       = "uses"
	EdgeTypeDeploys    = "deploys"
	EdgeTypeManages    = "manages"
	EdgeTypeContains   = "contains"
	EdgeTypeLinksTo    = "links_to" // service -> repo
)

// Provider extracts graphfs nodes and edges from system-spec documents.
type Provider struct{}

// NewProvider creates a new system-spec provider.
func NewProvider() *Provider {
	return &Provider{}
}

// Language returns the provider identifier.
func (p *Provider) Language() string {
	return "system-spec"
}

// Extensions returns file extensions this provider handles.
func (p *Provider) Extensions() []string {
	return []string{".json"}
}

// CanExtract returns true if the file is a system-spec document.
func (p *Provider) CanExtract(path string) bool {
	if !strings.HasSuffix(path, ".json") {
		return false
	}

	// Quick check: try to detect if it's a system-spec file
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	// Check for system-spec signature (has "name" and "services" at root)
	var probe struct {
		Name     string         `json:"name"`
		Services map[string]any `json:"services"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}

	return probe.Name != "" && probe.Services != nil
}

// ExtractFile extracts nodes and edges from a system-spec JSON file.
func (p *Provider) ExtractFile(path, baseDir string) ([]*graph.Node, []*graph.Edge, error) {
	sys, err := spec.LoadFromFile(path)
	if err != nil {
		return nil, nil, err
	}

	return p.ExtractSystem(sys, path)
}

// ExtractSystem extracts nodes and edges from a System struct.
func (p *Provider) ExtractSystem(sys *spec.System, sourcePath string) ([]*graph.Node, []*graph.Edge, error) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	relPath := sourcePath
	if abs, err := filepath.Abs(sourcePath); err == nil {
		relPath = abs
	}

	// Create system node
	systemNodeID := "system:" + sys.Name
	nodes = append(nodes, &graph.Node{
		ID:    systemNodeID,
		Type:  NodeTypeSystem,
		Label: sys.Name,
		Attrs: map[string]string{
			"source_file": relPath,
			"description": sys.Description,
			"version":     sys.Version,
		},
	})

	// Track created resource nodes to avoid duplicates
	resourceNodes := make(map[string]bool)

	// Process services
	for name, svc := range sys.Services {
		svcNodes, svcEdges := p.extractService(name, &svc, systemNodeID, resourceNodes)
		nodes = append(nodes, svcNodes...)
		edges = append(edges, svcEdges...)
	}

	// Process deployments
	if sys.Deployments != nil {
		deplNodes, deplEdges := p.extractDeployments(sys.Deployments, systemNodeID)
		nodes = append(nodes, deplNodes...)
		edges = append(edges, deplEdges...)
	}

	return nodes, edges, nil
}

func (p *Provider) extractService(name string, svc *spec.Service, systemNodeID string, resourceNodes map[string]bool) ([]*graph.Node, []*graph.Edge) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	svcNodeID := "svc:" + name

	// Service node
	attrs := map[string]string{
		"image": svc.Image.FullName(),
	}
	if svc.Registry != "" {
		attrs["registry"] = svc.Registry
	}
	if svc.Repo != nil {
		attrs["repo"] = svc.Repo.URL
		if svc.Repo.Ref != "" {
			attrs["ref"] = svc.Repo.Ref
		}
		if svc.Repo.Commit != "" {
			attrs["commit"] = svc.Repo.Commit
		}
	}
	if svc.Description != "" {
		attrs["description"] = svc.Description
	}

	nodes = append(nodes, &graph.Node{
		ID:    svcNodeID,
		Type:  NodeTypeService,
		Label: name,
		Attrs: attrs,
	})

	// System contains service
	edges = append(edges, &graph.Edge{
		From:       systemNodeID,
		To:         svcNodeID,
		Type:       EdgeTypeContains,
		Confidence: graph.ConfidenceExtracted,
	})

	// Service links to its source repository
	// This enables graphize to connect service nodes to code graphs
	if svc.Repo != nil && svc.Repo.URL != "" {
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         "repo:" + svc.Repo.URL,
			Type:       EdgeTypeLinksTo,
			Confidence: graph.ConfidenceExtracted,
			Attrs: map[string]string{
				"url": svc.Repo.URL,
			},
		})
	}

	// Service connections
	for targetName, conn := range svc.Connections {
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         "svc:" + targetName,
			Type:       EdgeTypeConnectsTo,
			Confidence: graph.ConfidenceExtracted,
			Attrs: map[string]string{
				"protocol": conn.Protocol,
				"port":     itoa(conn.Port),
			},
		})
	}

	// AWS resources
	if svc.AWS != nil {
		awsNodes, awsEdges := p.extractAWSResources(svcNodeID, svc.AWS, resourceNodes)
		nodes = append(nodes, awsNodes...)
		edges = append(edges, awsEdges...)
	}

	// GCP resources
	if svc.GCP != nil {
		gcpNodes, gcpEdges := p.extractGCPResources(svcNodeID, svc.GCP, resourceNodes)
		nodes = append(nodes, gcpNodes...)
		edges = append(edges, gcpEdges...)
	}

	// Cloudflare resources
	if svc.Cloudflare != nil {
		cfNodes, cfEdges := p.extractCloudflareResources(svcNodeID, svc.Cloudflare, resourceNodes)
		nodes = append(nodes, cfNodes...)
		edges = append(edges, cfEdges...)
	}

	return nodes, edges
}

func (p *Provider) extractAWSResources(svcNodeID string, aws *spec.AWSResources, resourceNodes map[string]bool) ([]*graph.Node, []*graph.Edge) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	// RDS
	for _, rds := range aws.RDS {
		nodeID := "rds:" + rds.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeDatabase,
				Label: rds.Name,
				Attrs: map[string]string{
					"provider": "aws",
					"engine":   rds.Engine,
					"port":     itoa(rds.Port),
				},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
			Attrs:      map[string]string{"protocol": "sql"},
		})
	}

	// DynamoDB
	for _, ddb := range aws.DynamoDB {
		nodeID := "dynamodb:" + ddb.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeDatabase,
				Label: ddb.Name,
				Attrs: map[string]string{
					"provider": "aws",
					"type":     "dynamodb",
				},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// SQS
	for _, sqs := range aws.SQS {
		nodeID := "sqs:" + sqs.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeQueue,
				Label: sqs.Name,
				Attrs: map[string]string{"provider": "aws"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// SNS
	for _, sns := range aws.SNS {
		nodeID := "sns:" + sns.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeTopic,
				Label: sns.Name,
				Attrs: map[string]string{"provider": "aws"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// S3
	for _, s3 := range aws.S3 {
		nodeID := "s3:" + s3.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeStorage,
				Label: s3.Name,
				Attrs: map[string]string{"provider": "aws"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// Bedrock
	for _, bedrock := range aws.Bedrock {
		nodeID := "bedrock:" + bedrock.ModelID
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeAIModel,
				Label: bedrock.ModelID,
				Attrs: map[string]string{"provider": "aws"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	return nodes, edges
}

func (p *Provider) extractGCPResources(svcNodeID string, gcp *spec.GCPResources, resourceNodes map[string]bool) ([]*graph.Node, []*graph.Edge) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	// CloudSQL
	for _, sql := range gcp.CloudSQL {
		nodeID := "cloudsql:" + sql.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeDatabase,
				Label: sql.Name,
				Attrs: map[string]string{
					"provider":      "gcp",
					"database_type": sql.DatabaseType,
				},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
			Attrs:      map[string]string{"protocol": "sql"},
		})
	}

	// PubSub
	for _, ps := range gcp.PubSub {
		nodeID := "pubsub:" + ps.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeTopic,
				Label: ps.Name,
				Attrs: map[string]string{"provider": "gcp"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// GCS
	for _, gcs := range gcp.GCS {
		nodeID := "gcs:" + gcs.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeStorage,
				Label: gcs.Name,
				Attrs: map[string]string{"provider": "gcp"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	return nodes, edges
}

func (p *Provider) extractCloudflareResources(svcNodeID string, cf *spec.CloudflareResources, resourceNodes map[string]bool) ([]*graph.Node, []*graph.Edge) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	// Workers
	for _, w := range cf.Workers {
		nodeID := "cf-worker:" + w.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeWorker,
				Label: w.Name,
				Attrs: map[string]string{
					"provider": "cloudflare",
					"route":    w.Route,
				},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	// R2
	for _, r2 := range cf.R2Buckets {
		nodeID := "r2:" + r2.Name
		if !resourceNodes[nodeID] {
			nodes = append(nodes, &graph.Node{
				ID:    nodeID,
				Type:  NodeTypeStorage,
				Label: r2.Name,
				Attrs: map[string]string{"provider": "cloudflare"},
			})
			resourceNodes[nodeID] = true
		}
		edges = append(edges, &graph.Edge{
			From:       svcNodeID,
			To:         nodeID,
			Type:       EdgeTypeUses,
			Confidence: graph.ConfidenceExtracted,
		})
	}

	return nodes, edges
}

func (p *Provider) extractDeployments(deployments *spec.Deployments, systemNodeID string) ([]*graph.Node, []*graph.Edge) {
	var nodes []*graph.Node
	var edges []*graph.Edge

	// Helm charts
	for name, helm := range deployments.Helm {
		nodeID := "helm:" + name
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeHelm,
			Label: name,
			Attrs: map[string]string{
				"chart":   helm.Chart,
				"version": helm.Version,
				"repo":    helm.Repo,
			},
		})

		// System contains helm chart
		edges = append(edges, &graph.Edge{
			From:       systemNodeID,
			To:         nodeID,
			Type:       EdgeTypeContains,
			Confidence: graph.ConfidenceExtracted,
		})

		// Helm deploys services
		for _, svcName := range helm.Services {
			edges = append(edges, &graph.Edge{
				From:       nodeID,
				To:         "svc:" + svcName,
				Type:       EdgeTypeDeploys,
				Confidence: graph.ConfidenceExtracted,
			})
		}
	}

	// Terraform modules
	for name, tf := range deployments.Terraform {
		nodeID := "terraform:" + name
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeTerraform,
			Label: name,
			Attrs: map[string]string{
				"source":  tf.Source,
				"version": tf.Version,
			},
		})

		// System contains terraform module
		edges = append(edges, &graph.Edge{
			From:       systemNodeID,
			To:         nodeID,
			Type:       EdgeTypeContains,
			Confidence: graph.ConfidenceExtracted,
		})

		// Terraform manages resources
		for _, res := range tf.Resources {
			edges = append(edges, &graph.Edge{
				From:       nodeID,
				To:         res,
				Type:       EdgeTypeManages,
				Confidence: graph.ConfidenceExtracted,
			})
		}
	}

	return nodes, edges
}

func itoa(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}
