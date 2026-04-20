# System Spec

A Go-native system topology specification for modeling microservices, infrastructure, and connectivity.

<div class="grid cards" markdown>

- :material-code-json: **JSON-first**

    ---

    Strongly-typed JSON format with generated JSON Schema. No polymorphism, no `interface{}`.

- :material-graph: **Graph-aware**

    ---

    Services as nodes, connections as edges. Built for visualization and analysis.

- :material-cloud: **Multi-cloud**

    ---

    First-class support for AWS, GCP, and Cloudflare resources.

- :material-language-go: **Go SDK**

    ---

    Idiomatic Go types with full IDE support and compile-time safety.

</div>

## Quick Example

```json
{
  "name": "payments-platform",
  "services": {
    "api-gateway": {
      "image": { "name": "ghcr.io/org/api-gateway", "tag": "v2.0" },
      "repo": { "url": "https://github.com/org/api-gateway" },
      "connections": {
        "payments-service": { "port": 8080, "protocol": "grpc" }
      }
    },
    "payments-service": {
      "image": { "name": "ghcr.io/org/payments", "digest": "sha256:abc123" },
      "aws": {
        "rds": [{ "name": "payments-db", "engine": "aurora-mysql" }]
      }
    }
  }
}
```

## Render to Multiple Formats

```bash
# D2 diagrams
system-spec render system.json --format d2 > system.d2

# Mermaid for GitHub/Markdown
system-spec render system.json --format mermaid > system.mmd

# Cytoscape.js for interactive web
system-spec render system.json --format cytoscape > system.json

# Sigma.js for large graphs
system-spec render system.json --format sigma > system.json
```

## Design Principles

| Principle | Description |
|-----------|-------------|
| **Go-first** | Go structs are source of truth; JSON Schema is generated |
| **No polymorphism** | All types are concrete; `map[string]Service` not `map[string]interface{}` |
| **Graph-native** | System topology modeled as directed graph |
| **Deployment-aware** | Links to Helm charts and Terraform modules |

## Installation

```bash
go install github.com/plexusone/system-spec/cmd/system-spec@latest
```

## Resources

- [Specification v0.1.0](spec/v0.1.0.md) - Human-readable spec
- [JSON Schema](https://github.com/plexusone/system-spec/blob/main/versions/v0.1.0/system.schema.json) - Machine-readable schema
- [Go SDK](sdk/usage.md) - Go package documentation
- [Examples](examples/payments.md) - Real-world examples
