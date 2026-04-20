# Analyze System Spec

You are helping the user analyze, validate, or visualize an existing **system-spec** file.

## Available Commands

### Validate
Check if a system-spec file is valid:
```bash
system-spec validate <file.json>
```

Output shows:
- System name and service count
- Each service with image name
- Repository URLs
- Connection counts

### Render Diagrams
Generate visual diagrams from the spec:

```bash
# D2 diagram (clean architecture style)
system-spec render system.json --format d2 > system.d2

# Mermaid (for markdown/GitHub)
system-spec render system.json --format mermaid > system.mmd

# Cytoscape.js (interactive web)
system-spec render system.json --format cytoscape > system.cyto.json

# Sigma.js (large graphs)
system-spec render system.json --format sigma > system.sigma.json

# GraphViz DOT
system-spec render system.json --format dot > system.dot
```

### Extract Graph
Get the intermediate graph representation:
```bash
system-spec graph system.json
```

### Generate Schema
Output the JSON Schema for system-spec:
```bash
system-spec schema > system.schema.json
```

## Analysis Tasks

When asked to analyze a system-spec, provide:

### 1. System Overview
- Total services count
- Cloud provider usage (AWS, GCP, Cloudflare)
- Deployment methods (Helm, Terraform)

### 2. Service Dependencies
- Which services connect to which
- Protocols used (HTTP, gRPC, AMQP, etc.)
- Shared resources (databases used by multiple services)

### 3. Infrastructure Summary
- Databases (RDS, CloudSQL, DynamoDB)
- Message queues (SQS, PubSub)
- Storage (S3, GCS, R2)
- AI/ML resources (Bedrock)

### 4. Potential Issues
- Orphan services (no connections)
- Missing repo URLs
- Services without descriptions
- Circular dependencies

## Example Analysis

Given a system-spec file, provide analysis like:

```
## System: payments-platform

### Overview
- 4 services
- AWS resources: RDS, SQS, SNS, S3
- Deployments: Helm (1 chart), Terraform (1 module)

### Service Map
- api-gateway → payment-processor (grpc:8080)
- api-gateway → notification-service (http:8081)
- payment-processor → notification-service (amqp:5672)

### Infrastructure
| Resource Type | Count | Services Using |
|---------------|-------|----------------|
| RDS PostgreSQL | 1 | payment-processor |
| SQS Queue | 2 | payment-processor, notification-service |
| S3 Bucket | 1 | api-gateway |

### Observations
- All services have repo URLs defined
- payment-processor is a hub (most connections)
- Consider adding health check endpoints to spec
```

## Comparison Tasks

When asked to compare two specs:
1. Load both specs with `system-spec validate`
2. Identify added/removed services
3. Identify changed connections
4. Identify resource changes
5. Highlight potential breaking changes

## Integration with Graphize

Link infrastructure topology with your codebase using graphize:

### Setup
```bash
# In your repo directory:
graphize init                    # Initialize graph database
graphize add .                   # Add repo (auto-detects system-spec JSON)
graphize analyze                 # Extract code graph + infrastructure
```

### Query the combined graph
```bash
# Find what code implements a service
graphize query svc:api-gateway

# Find all database connections
graphize query --edge-type uses | grep rds

# Trace path from service to database
graphize path "svc:payment-processor" "rds:main-db"
```

### Generate visualizations
```bash
graphize export html             # Interactive Cytoscape.js viewer
graphize export toon             # Token-efficient format for AI agents
graphize report                  # Markdown analysis report
```

### What graphize extracts from system-spec

| Node Type | Examples |
|-----------|----------|
| `system` | The root system |
| `service` | Each service in the spec |
| `database` | RDS, CloudSQL, DynamoDB |
| `queue` | SQS queues |
| `topic` | SNS, PubSub topics |
| `storage` | S3, GCS, R2 buckets |
| `ai_model` | Bedrock models |
| `helm_chart` | Helm deployments |
| `terraform_module` | Terraform modules |

| Edge Type | Meaning |
|-----------|---------|
| `contains` | System contains service |
| `connects_to` | Service-to-service connection |
| `uses` | Service uses a resource |
| `deploys` | Helm chart deploys service |
| `manages` | Terraform manages resource |

This allows queries like "which services use the payments database?" or "what does the api-gateway connect to?"
