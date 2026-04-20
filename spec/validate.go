package spec

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFromFile reads and parses a System from a JSON file.
func LoadFromFile(path string) (*System, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return LoadFromJSON(data)
}

// LoadFromJSON parses a System from JSON data.
func LoadFromJSON(data []byte) (*System, error) {
	var sys System
	if err := json.Unmarshal(data, &sys); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := sys.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &sys, nil
}

// Validate checks the System for logical errors.
func (s *System) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("system name is required")
	}

	if len(s.Services) == 0 {
		return fmt.Errorf("system must have at least one service")
	}

	// Validate each service
	for name, svc := range s.Services {
		if err := validateService(name, &svc, s); err != nil {
			return err
		}
	}

	// Validate deployments reference existing services
	if s.Deployments != nil {
		if err := validateDeployments(s.Deployments, s); err != nil {
			return err
		}
	}

	return nil
}

func validateDeployments(deployments *Deployments, sys *System) error {
	// Validate Helm chart service references
	for helmName, helm := range deployments.Helm {
		for _, svcName := range helm.Services {
			if _, ok := sys.Services[svcName]; !ok {
				return fmt.Errorf("helm chart %q: references undefined service %q", helmName, svcName)
			}
		}
	}

	// Validate Terraform resource references that look like services
	for tfName, tf := range deployments.Terraform {
		for _, res := range tf.Resources {
			// Check if resource is a service reference (svc:name format)
			if len(res) > 4 && res[:4] == "svc:" {
				svcName := res[4:]
				if _, ok := sys.Services[svcName]; !ok {
					return fmt.Errorf("terraform module %q: references undefined service %q", tfName, svcName)
				}
			}
		}
	}

	return nil
}

func validateService(name string, svc *Service, sys *System) error {
	// Image name is required
	if svc.Image.Name == "" {
		return fmt.Errorf("service %q: image name is required", name)
	}

	// Validate connections reference existing services
	for targetName := range svc.Connections {
		if _, ok := sys.Services[targetName]; !ok {
			return fmt.Errorf("service %q: connection to unknown service %q", name, targetName)
		}
	}

	// Validate connection details
	for targetName, conn := range svc.Connections {
		if conn.Port <= 0 || conn.Port > 65535 {
			return fmt.Errorf("service %q: invalid port %d for connection to %q", name, conn.Port, targetName)
		}
		if conn.Protocol == "" {
			return fmt.Errorf("service %q: protocol is required for connection to %q", name, targetName)
		}
	}

	// Validate Git repo URL if specified
	if svc.Repo != nil && svc.Repo.URL == "" {
		return fmt.Errorf("service %q: repo URL is required when repo is specified", name)
	}

	return nil
}

// ToJSON serializes the System to JSON.
func (s *System) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// ToJSONCompact serializes the System to compact JSON.
func (s *System) ToJSONCompact() ([]byte, error) {
	return json.Marshal(s)
}

// Warning represents a lint warning (advisory, not an error).
type Warning struct {
	// Code is a short identifier (e.g., "no-connections", "isolated-service")
	Code string
	// Service is the service name this warning applies to (empty for system-level)
	Service string
	// Message is the human-readable warning description
	Message string
}

// Lint checks the System for potential issues and returns warnings.
// Unlike Validate(), Lint() returns advisory warnings rather than errors.
func (s *System) Lint() []Warning {
	var warnings []Warning

	// Check for unknown owner references
	for name, svc := range s.Services {
		if svc.Owner != "" {
			if _, ok := s.Owners[svc.Owner]; !ok {
				warnings = append(warnings, Warning{
					Code:    "unknown-owner",
					Service: name,
					Message: fmt.Sprintf("references unknown owner %q", svc.Owner),
				})
			}
		}
	}

	// Check each service for graph-related issues
	for name, svc := range s.Services {
		hasConnections := len(svc.Connections) > 0
		hasResources := hasAnyResources(&svc)
		hasIncoming := s.hasIncomingConnections(name)

		if !hasConnections && !hasResources && !hasIncoming {
			warnings = append(warnings, Warning{
				Code:    "isolated-service",
				Service: name,
				Message: "has no connections or resources (will appear as isolated node)",
			})
		} else {
			if !hasConnections && !hasIncoming {
				warnings = append(warnings, Warning{
					Code:    "no-connections",
					Service: name,
					Message: "has no connections to other services",
				})
			}
			if !hasResources {
				warnings = append(warnings, Warning{
					Code:    "no-resources",
					Service: name,
					Message: "has no cloud resources attached",
				})
			}
		}

		// Check for missing exposes when service has incoming connections
		if hasIncoming && len(svc.Exposes) == 0 {
			warnings = append(warnings, Warning{
				Code:    "missing-exposes",
				Service: name,
				Message: "has incoming connections but no exposes defined",
			})
		}
	}

	// Check deployments for references to non-existent services
	if s.Deployments != nil {
		for helmName, helm := range s.Deployments.Helm {
			for _, svcName := range helm.Services {
				if _, ok := s.Services[svcName]; !ok {
					warnings = append(warnings, Warning{
						Code:    "orphan-helm-service",
						Service: helmName,
						Message: fmt.Sprintf("helm chart references undefined service %q", svcName),
					})
				}
			}
		}

		for tfName, tf := range s.Deployments.Terraform {
			for _, res := range tf.Resources {
				// Check if resource looks like a service reference
				if len(res) > 4 && res[:4] == "svc:" {
					svcName := res[4:]
					if _, ok := s.Services[svcName]; !ok {
						warnings = append(warnings, Warning{
							Code:    "orphan-terraform-resource",
							Service: tfName,
							Message: fmt.Sprintf("terraform module references undefined service %q", svcName),
						})
					}
				}
			}
		}
	}

	return warnings
}

// hasIncomingConnections checks if any service connects to the given service.
func (s *System) hasIncomingConnections(serviceName string) bool {
	for _, svc := range s.Services {
		if _, ok := svc.Connections[serviceName]; ok {
			return true
		}
	}
	return false
}

// hasAnyResources checks if a service has any cloud resources attached.
func hasAnyResources(svc *Service) bool {
	if svc.AWS != nil {
		if len(svc.AWS.RDS) > 0 || len(svc.AWS.DynamoDB) > 0 ||
			len(svc.AWS.SQS) > 0 || len(svc.AWS.SNS) > 0 ||
			len(svc.AWS.S3) > 0 || len(svc.AWS.Bedrock) > 0 ||
			len(svc.AWS.Lambda) > 0 || len(svc.AWS.ElastiCache) > 0 ||
			len(svc.AWS.Kinesis) > 0 || len(svc.AWS.EventBridge) > 0 ||
			len(svc.AWS.SecretsManager) > 0 || len(svc.AWS.APIGateway) > 0 ||
			len(svc.AWS.CloudFront) > 0 {
			return true
		}
	}
	if svc.GCP != nil {
		if len(svc.GCP.CloudSQL) > 0 || len(svc.GCP.PubSub) > 0 ||
			len(svc.GCP.GCS) > 0 || len(svc.GCP.CloudRun) > 0 ||
			len(svc.GCP.CloudFunctions) > 0 || len(svc.GCP.Memorystore) > 0 ||
			len(svc.GCP.BigQuery) > 0 || len(svc.GCP.SecretManager) > 0 ||
			len(svc.GCP.Firestore) > 0 || len(svc.GCP.Spanner) > 0 {
			return true
		}
	}
	if svc.Azure != nil {
		if len(svc.Azure.CosmosDB) > 0 || len(svc.Azure.ServiceBus) > 0 ||
			len(svc.Azure.BlobStorage) > 0 || len(svc.Azure.Functions) > 0 ||
			len(svc.Azure.RedisCache) > 0 || len(svc.Azure.KeyVault) > 0 ||
			len(svc.Azure.EventHubs) > 0 || len(svc.Azure.SQLDatabase) > 0 {
			return true
		}
	}
	if svc.Cloudflare != nil {
		if len(svc.Cloudflare.Workers) > 0 || len(svc.Cloudflare.R2Buckets) > 0 ||
			len(svc.Cloudflare.KV) > 0 || len(svc.Cloudflare.D1) > 0 ||
			len(svc.Cloudflare.Queues) > 0 || len(svc.Cloudflare.DurableObjects) > 0 {
			return true
		}
	}
	if svc.Resources != nil {
		if len(svc.Resources.Redis) > 0 || len(svc.Resources.Kafka) > 0 ||
			len(svc.Resources.Elasticsearch) > 0 || len(svc.Resources.MongoDB) > 0 ||
			len(svc.Resources.PostgreSQL) > 0 || len(svc.Resources.MySQL) > 0 ||
			len(svc.Resources.RabbitMQ) > 0 || len(svc.Resources.Cassandra) > 0 ||
			len(svc.Resources.Memcached) > 0 || len(svc.Resources.Minio) > 0 ||
			len(svc.Resources.NATS) > 0 || len(svc.Resources.Etcd) > 0 ||
			len(svc.Resources.Vault) > 0 || len(svc.Resources.Consul) > 0 {
			return true
		}
	}
	return false
}
