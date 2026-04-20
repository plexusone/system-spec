# Specification Overview

The System Specification (system-spec) is a strongly-typed format for describing system topologies.

## Versions

| Version | Status | Release Date |
|---------|--------|--------------|
| [v0.1.0](v0.1.0.md) | Current | 2026-04-15 |

## Format

System-spec documents are JSON files. YAML is not supported to ensure consistent parsing.

## Schema

Each version has:

- **Human-readable spec** (`spec.md`) - For developers
- **JSON Schema** (`system.schema.json`) - For tools and AI agents

## Design Principles

### Go-first

Go structs are the source of truth. JSON Schema is generated from Go types using `github.com/invopop/jsonschema`.

```
Go structs → jsonschema → system.schema.json
```

### No Polymorphism

Unlike OAM or similar specs, system-spec uses concrete types only:

```go
// ✓ Allowed: typed maps
Services map[string]Service

// ✗ Rejected: interface{}
Properties map[string]interface{}
```

### Schema Validation

Generated schemas pass `schemalint` validation:

```bash
schemalint lint versions/v0.1.0/system.schema.json
```

## Core Objects

```
System
├── Services (map[string]Service)
│   ├── Image (ContainerImage)
│   ├── Repo (GitRepo)
│   ├── Connections (map[string]Connection)
│   └── AWS/GCP/Cloudflare resources
├── Networks (map[string]Network)
└── Deployments
    ├── Helm (map[string]HelmDeployment)
    └── Terraform (map[string]TerraformDeployment)
```

## Graph Model

System-spec is designed for graph operations:

- **Nodes**: Services, databases, queues, storage, AI models
- **Edges**: Connections with protocol and port

This enables:

- Dependency analysis
- Network policy generation
- Visualization
- Impact analysis
