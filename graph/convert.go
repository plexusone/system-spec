package graph

import (
	"fmt"

	"github.com/plexusone/system-spec/spec"
)

// FromSystem converts a System spec to a Graph for rendering.
func FromSystem(s *spec.System) *Graph {
	g := NewGraph()

	// Add service nodes
	for name, svc := range s.Services {
		node := Node{
			ID:    serviceNodeID(name),
			Label: name,
			Kind:  NodeKindService,
			Metadata: map[string]string{
				"image": svc.Image.FullName(),
			},
		}

		if svc.Repo != nil {
			node.Metadata["repo"] = svc.Repo.URL
		}
		if svc.Registry != "" {
			node.Metadata["registry"] = svc.Registry
		}

		g.AddNode(node)

		// Add service-to-service connections
		for targetName, conn := range svc.Connections {
			g.AddEdge(Edge{
				ID:       fmt.Sprintf("%s->%s", name, targetName),
				Source:   serviceNodeID(name),
				Target:   serviceNodeID(targetName),
				Label:    fmt.Sprintf("%s:%d", conn.Protocol, conn.Port),
				Kind:     EdgeKindConnection,
				Protocol: conn.Protocol,
				Port:     conn.Port,
			})
		}

		// Add AWS resources
		if svc.AWS != nil {
			addAWSResources(g, name, svc.AWS)
		}

		// Add GCP resources
		if svc.GCP != nil {
			addGCPResources(g, name, svc.GCP)
		}

		// Add Cloudflare resources
		if svc.Cloudflare != nil {
			addCloudflareResources(g, name, svc.Cloudflare)
		}
	}

	// Add deployment nodes (Helm, Terraform)
	if s.Deployments != nil {
		addDeploymentNodes(g, s.Deployments)
	}

	return g
}

func serviceNodeID(name string) string {
	return "svc:" + name
}

func addAWSResources(g *Graph, serviceName string, aws *spec.AWSResources) {
	svcNodeID := serviceNodeID(serviceName)

	// RDS instances
	for _, rds := range aws.RDS {
		nodeID := "rds:" + rds.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: rds.Name,
				Kind:  NodeKindDatabase,
				Metadata: map[string]string{
					"engine":   rds.Engine,
					"provider": "aws",
				},
			})
		}
		port := rds.Port
		if port == 0 {
			port = 3306 // default MySQL
			if rds.Engine == "postgres" || rds.Engine == "aurora-postgresql" {
				port = 5432
			}
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->%s", serviceName, rds.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindDatabase,
			Protocol: "sql",
			Port:     port,
		})
	}

	// DynamoDB tables
	for _, ddb := range aws.DynamoDB {
		nodeID := "dynamodb:" + ddb.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: ddb.Name,
				Kind:  NodeKindDatabase,
				Metadata: map[string]string{
					"provider": "aws",
					"type":     "dynamodb",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->dynamodb:%s", serviceName, ddb.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindDatabase,
			Protocol: "https",
			Port:     443,
		})
	}

	// SQS queues
	for _, sqs := range aws.SQS {
		nodeID := "sqs:" + sqs.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: sqs.Name,
				Kind:  NodeKindQueue,
				Metadata: map[string]string{
					"provider": "aws",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->sqs:%s", serviceName, sqs.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindQueue,
			Protocol: "https",
			Port:     443,
		})
	}

	// SNS topics
	for _, sns := range aws.SNS {
		nodeID := "sns:" + sns.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: sns.Name,
				Kind:  NodeKindTopic,
				Metadata: map[string]string{
					"provider": "aws",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->sns:%s", serviceName, sns.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindQueue,
			Protocol: "https",
			Port:     443,
		})
	}

	// S3 buckets
	for _, s3 := range aws.S3 {
		nodeID := "s3:" + s3.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: s3.Name,
				Kind:  NodeKindStorage,
				Metadata: map[string]string{
					"provider": "aws",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->s3:%s", serviceName, s3.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindStorage,
			Protocol: "https",
			Port:     443,
		})
	}

	// Bedrock models
	for _, bedrock := range aws.Bedrock {
		nodeID := "bedrock:" + bedrock.ModelID
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: bedrock.ModelID,
				Kind:  NodeKindAIModel,
				Metadata: map[string]string{
					"provider": "aws",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->bedrock:%s", serviceName, bedrock.ModelID),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindConnection,
			Protocol: "https",
			Port:     443,
		})
	}
}

func addGCPResources(g *Graph, serviceName string, gcp *spec.GCPResources) {
	svcNodeID := serviceNodeID(serviceName)

	// Cloud SQL instances
	for _, sql := range gcp.CloudSQL {
		nodeID := "cloudsql:" + sql.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: sql.Name,
				Kind:  NodeKindDatabase,
				Metadata: map[string]string{
					"provider":     "gcp",
					"databaseType": sql.DatabaseType,
				},
			})
		}
		port := 3306
		if sql.DatabaseType == "postgres" {
			port = 5432
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->cloudsql:%s", serviceName, sql.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindDatabase,
			Protocol: "sql",
			Port:     port,
		})
	}

	// Pub/Sub topics
	for _, ps := range gcp.PubSub {
		nodeID := "pubsub:" + ps.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: ps.Name,
				Kind:  NodeKindTopic,
				Metadata: map[string]string{
					"provider": "gcp",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->pubsub:%s", serviceName, ps.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindQueue,
			Protocol: "https",
			Port:     443,
		})
	}

	// GCS buckets
	for _, gcs := range gcp.GCS {
		nodeID := "gcs:" + gcs.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: gcs.Name,
				Kind:  NodeKindStorage,
				Metadata: map[string]string{
					"provider": "gcp",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->gcs:%s", serviceName, gcs.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindStorage,
			Protocol: "https",
			Port:     443,
		})
	}
}

func addCloudflareResources(g *Graph, serviceName string, cf *spec.CloudflareResources) {
	svcNodeID := serviceNodeID(serviceName)

	// Workers
	for _, w := range cf.Workers {
		nodeID := "cf-worker:" + w.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: w.Name,
				Kind:  NodeKindWorker,
				Metadata: map[string]string{
					"provider": "cloudflare",
					"route":    w.Route,
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->cf-worker:%s", serviceName, w.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindConnection,
			Protocol: "https",
			Port:     443,
		})
	}

	// R2 buckets
	for _, r2 := range cf.R2Buckets {
		nodeID := "r2:" + r2.Name
		if g.NodeByID(nodeID) == nil {
			g.AddNode(Node{
				ID:    nodeID,
				Label: r2.Name,
				Kind:  NodeKindStorage,
				Metadata: map[string]string{
					"provider": "cloudflare",
				},
			})
		}
		g.AddEdge(Edge{
			ID:       fmt.Sprintf("%s->r2:%s", serviceName, r2.Name),
			Source:   svcNodeID,
			Target:   nodeID,
			Kind:     EdgeKindStorage,
			Protocol: "https",
			Port:     443,
		})
	}
}

func addDeploymentNodes(g *Graph, deployments *spec.Deployments) {
	// Helm deployments
	for name, helm := range deployments.Helm {
		nodeID := "helm:" + name
		g.AddNode(Node{
			ID:    nodeID,
			Label: name,
			Kind:  NodeKindHelm,
			Metadata: map[string]string{
				"chart":   helm.Chart,
				"version": helm.Version,
			},
		})

		// Link to services this helm chart deploys
		for _, svcName := range helm.Services {
			g.AddEdge(Edge{
				ID:     fmt.Sprintf("helm:%s->%s", name, svcName),
				Source: nodeID,
				Target: serviceNodeID(svcName),
				Kind:   EdgeKindDeploys,
			})
		}
	}

	// Terraform deployments
	for name, tf := range deployments.Terraform {
		nodeID := "terraform:" + name
		g.AddNode(Node{
			ID:    nodeID,
			Label: name,
			Kind:  NodeKindTerraform,
			Metadata: map[string]string{
				"source":  tf.Source,
				"version": tf.Version,
			},
		})

		// Link to resources this terraform module manages
		for _, res := range tf.Resources {
			g.AddEdge(Edge{
				ID:     fmt.Sprintf("terraform:%s->%s", name, res),
				Source: nodeID,
				Target: res,
				Kind:   EdgeKindDeploys,
			})
		}
	}
}
