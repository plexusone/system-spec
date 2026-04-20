// Package spec defines the core system-spec types.
// All types are concrete with no interface{} or any fields.
// JSON Schema is generated from these Go structs.
package spec

// System represents a complete system topology.
// It contains services, their dependencies, and infrastructure bindings.
type System struct {
	// Name is the system identifier (e.g., "payments-platform")
	Name string `json:"name"`

	// Description provides context about the system
	Description string `json:"description,omitempty"`

	// Version of this system spec document
	Version string `json:"version,omitempty"`

	// Owners maps owner ID to owner details for team ownership
	Owners map[string]Owner `json:"owners,omitempty"`

	// Services maps service name to service definition
	Services map[string]Service `json:"services"`

	// Networks maps network name to network definition (future)
	Networks map[string]Network `json:"networks,omitempty"`

	// Deployments maps deployment strategies (helm, terraform) (future)
	Deployments *Deployments `json:"deployments,omitempty"`
}

// Owner represents a team or individual responsible for services.
type Owner struct {
	// Name is the display name (e.g., "Payments Team")
	Name string `json:"name"`

	// Email is the team contact email
	Email string `json:"email,omitempty"`

	// Slack is the Slack channel or handle (e.g., "#payments-eng")
	Slack string `json:"slack,omitempty"`

	// Pager is the PagerDuty/Opsgenie service or escalation policy ID
	Pager string `json:"pager,omitempty"`
}

// Service represents a single microservice or workload.
type Service struct {
	// Description of the service
	Description string `json:"description,omitempty"`

	// Owner references an owner ID from system.owners
	Owner string `json:"owner,omitempty"`

	// Repo is the source code repository
	Repo *GitRepo `json:"repo,omitempty"`

	// Image is the container image specification
	Image ContainerImage `json:"image"`

	// Registry is the container registry path (e.g., "ghcr.io/org/repo")
	Registry string `json:"registry,omitempty"`

	// EnvVars documents environment variables for discovering this service
	EnvVars *ServiceEnvVars `json:"envVars,omitempty"`

	// Connections maps target service name to connection details
	Connections map[string]Connection `json:"connections,omitempty"`

	// Exposes defines ports/protocols this service exposes
	Exposes []Endpoint `json:"exposes,omitempty"`

	// AWS contains AWS-specific resource bindings
	AWS *AWSResources `json:"aws,omitempty"`

	// GCP contains GCP-specific resource bindings
	GCP *GCPResources `json:"gcp,omitempty"`

	// Azure contains Azure-specific resource bindings
	Azure *AzureResources `json:"azure,omitempty"`

	// Cloudflare contains Cloudflare resource bindings
	Cloudflare *CloudflareResources `json:"cloudflare,omitempty"`

	// Resources contains provider-agnostic resource bindings
	Resources *GenericResources `json:"resources,omitempty"`
}

// ServiceEnvVars documents environment variables for service discovery.
type ServiceEnvVars struct {
	// URL is the env var name for the full service URL (e.g., "PAYMENTS_URL")
	URL string `json:"url,omitempty"`

	// Host is the env var name for the service host (e.g., "PAYMENTS_SERVICE_HOST")
	Host string `json:"host,omitempty"`

	// Port is the env var name for the service port (e.g., "PAYMENTS_SERVICE_PORT")
	Port string `json:"port,omitempty"`
}

// GitRepo represents a source code repository.
type GitRepo struct {
	// URL is the full repository URL (e.g., "https://github.com/org/repo")
	URL string `json:"url"`

	// Path is the relative path within the repo (for monorepos)
	Path string `json:"path,omitempty"`

	// Ref is the git reference (branch, tag, or commit)
	Ref string `json:"ref,omitempty"`

	// Commit is the specific commit SHA (for pinning)
	Commit string `json:"commit,omitempty"`
}

// ContainerImage represents a container image specification.
type ContainerImage struct {
	// Name is the image name without tag (e.g., "nginx", "ghcr.io/org/app")
	Name string `json:"name"`

	// Tag is the image tag (e.g., "v1.2.3", "latest")
	Tag string `json:"tag,omitempty"`

	// Digest is the image digest for pinning (e.g., "sha256:abc123...")
	Digest string `json:"digest,omitempty"`
}

// FullName returns the complete image reference.
func (i ContainerImage) FullName() string {
	if i.Digest != "" {
		return i.Name + "@" + i.Digest
	}
	if i.Tag != "" {
		return i.Name + ":" + i.Tag
	}
	return i.Name
}

// Connection represents a connection from one service to another.
type Connection struct {
	// Port is the target port number
	Port int `json:"port"`

	// Protocol is the connection protocol (http, grpc, tcp, sql, amqp, redis)
	Protocol string `json:"protocol"`

	// Description explains the purpose of this connection
	Description string `json:"description,omitempty"`
}

// Endpoint represents an exposed port/protocol.
type Endpoint struct {
	// Port number
	Port int `json:"port"`

	// Protocol (http, grpc, tcp)
	Protocol string `json:"protocol"`

	// Description of the endpoint
	Description string `json:"description,omitempty"`
}

// Network represents a network boundary (VPC, subnet, security group).
type Network struct {
	// Description of the network
	Description string `json:"description,omitempty"`

	// CIDR block for the network
	CIDR string `json:"cidr,omitempty"`

	// Rules defines firewall/security group rules
	Rules []NetworkRule `json:"rules,omitempty"`
}

// NetworkRule represents a firewall or security group rule.
type NetworkRule struct {
	// Direction is "inbound" or "outbound"
	Direction string `json:"direction"`

	// FromService is the source service (for inbound rules)
	FromService string `json:"fromService,omitempty"`

	// ToService is the target service (for outbound rules)
	ToService string `json:"toService,omitempty"`

	// Port number
	Port int `json:"port"`

	// Protocol (tcp, udp)
	Protocol string `json:"protocol"`

	// Action is "allow" or "deny"
	Action string `json:"action"`
}

// Deployments contains deployment configuration references.
type Deployments struct {
	// Helm maps helm release name to helm chart reference
	Helm map[string]HelmDeployment `json:"helm,omitempty"`

	// Terraform maps terraform module name to terraform reference
	Terraform map[string]TerraformDeployment `json:"terraform,omitempty"`
}

// HelmDeployment references a Helm chart for deployment.
type HelmDeployment struct {
	// Chart is the chart name or URL
	Chart string `json:"chart"`

	// Version is the chart version
	Version string `json:"version,omitempty"`

	// Repo is the Helm repository URL
	Repo string `json:"repo,omitempty"`

	// ValuesFile is the path to values file (relative to system spec)
	ValuesFile string `json:"valuesFile,omitempty"`

	// Services lists which services this chart deploys
	Services []string `json:"services,omitempty"`
}

// TerraformDeployment references a Terraform module for infrastructure.
type TerraformDeployment struct {
	// Source is the module source (local path, git URL, registry)
	Source string `json:"source"`

	// Version is the module version (for registry modules)
	Version string `json:"version,omitempty"`

	// Path is the path within the source (for git sources)
	Path string `json:"path,omitempty"`

	// Resources lists which infrastructure resources this module manages
	Resources []string `json:"resources,omitempty"`
}
