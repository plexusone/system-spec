---
name: create-system-spec
description: Create a system-spec topology file for SaaS infrastructure
triggers:
  - create-system-spec
  - system-spec
  - create system spec
  - infrastructure spec
dependencies:
  - system-spec
references:
  - ../../versions/v0.1.0/system.schema.json
  - ../../examples/payments-system.json
---

# Create System Spec

You are helping the user create a **system-spec** file - a JSON specification for modeling SaaS infrastructure topology.

## Workflow

1. **Gather Requirements** - Ask the user about their system:
   - System name and description
   - Services (microservices, APIs, workers)
   - Cloud resources (databases, queues, storage)
   - Service connections (which services talk to which)
   - Deployment method (Helm, Terraform, or both)

2. **Build the Spec** - Create a JSON file following the schema:
   - Start with system metadata (name, description, version)
   - Add each service with image, repo, and resources
   - Define connections between services
   - Add deployment configurations

3. **Validate** - Run `system-spec validate <file>` to check the spec

4. **Render** - Generate diagrams with `system-spec render <file> --format <fmt>`

## Schema Reference

The system-spec JSON structure:

```json
{
  "name": "system-name",
  "description": "What this system does",
  "version": "1.0.0",
  "services": {
    "service-name": {
      "image": { "name": "image-name", "tag": "v1.0.0" },
      "registry": "gcr.io/project",
      "repo": { "url": "https://github.com/org/repo", "branch": "main" },
      "description": "What this service does",
      "connections": {
        "target-service": { "port": 8080, "protocol": "grpc" }
      },
      "aws": { ... },
      "gcp": { ... },
      "cloudflare": { ... }
    }
  },
  "deployments": {
    "helm": { ... },
    "terraform": { ... }
  }
}
```

## Cloud Resources

### AWS Resources
```json
"aws": {
  "rds": [{ "name": "main-db", "engine": "postgres", "port": 5432 }],
  "dynamodb": [{ "name": "sessions" }],
  "sqs": [{ "name": "task-queue" }],
  "sns": [{ "name": "notifications" }],
  "s3": [{ "name": "uploads" }],
  "bedrock": [{ "modelId": "anthropic.claude-3-sonnet" }]
}
```

### GCP Resources
```json
"gcp": {
  "cloudSQL": [{ "name": "main-db", "databaseType": "POSTGRES" }],
  "pubsub": [{ "name": "events" }],
  "gcs": [{ "name": "assets" }]
}
```

### Cloudflare Resources
```json
"cloudflare": {
  "workers": [{ "name": "edge-auth", "route": "auth.example.com/*" }],
  "r2Buckets": [{ "name": "static-assets" }]
}
```

## Deployments

### Helm Deployment
```json
"deployments": {
  "helm": {
    "app-chart": {
      "chart": "myorg/app",
      "version": "1.0.0",
      "repo": "https://charts.example.com",
      "services": ["api", "worker"]
    }
  }
}
```

### Terraform Deployment
```json
"deployments": {
  "terraform": {
    "infra": {
      "source": "git::https://github.com/org/terraform-modules",
      "version": "v1.0.0",
      "resources": ["rds:main-db", "sqs:task-queue"]
    }
  }
}
```

## Full Pipeline: Create → Validate → Render → View

### Step 1: Create the spec
Write the JSON file (e.g., `system.json`) following the schema above.

### Step 2: Validate
```bash
system-spec validate system.json
```

### Step 3: Render to diagram format
```bash
# Choose one or more formats:
system-spec render system.json --format d2 > system.d2
system-spec render system.json --format mermaid > system.mmd
system-spec render system.json --format cytoscape > system.cyto.json
system-spec render system.json --format dot > system.dot
```

### Step 4: Generate visual output

**D2 → SVG/PNG** (requires [d2](https://d2lang.com)):
```bash
d2 system.d2 system.svg          # SVG output
d2 system.d2 system.png          # PNG output
open system.svg                   # View the diagram
```

**Mermaid → Image**:
```bash
mmdc -i system.mmd -o system.svg
# Or paste into: https://mermaid.live
```

**Cytoscape → Interactive HTML**:
The JSON output needs an HTML wrapper with Cytoscape.js library.

**GraphViz DOT → Image**:
```bash
dot -Tsvg system.dot -o system.svg
```

### Step 5 (Optional): Add to code graph with Graphize

Link infrastructure topology with your codebase:
```bash
graphize init                    # Initialize graph database
graphize add .                   # Add repo (detects system-spec files)
graphize analyze                 # Extract code + infrastructure graph
graphize export html             # Generate interactive visualization
```

## Render Format Summary

| Format | Output | Tool Required | Best For |
|--------|--------|---------------|----------|
| D2 | `.d2` → `.svg/.png` | `d2` CLI | Clean architecture diagrams |
| Mermaid | `.mmd` → `.svg` | `mmdc` or online | Markdown docs, GitHub |
| Cytoscape | `.json` → HTML | Browser + JS | Interactive exploration |
| Sigma.js | `.json` → HTML | Browser + JS | Very large graphs |
| GraphViz | `.dot` → `.svg/.png` | `dot` CLI | Traditional layouts |

## Example Session

**User:** Create a system spec for my e-commerce platform with an API, worker, and PostgreSQL database on AWS.

**Assistant:** I'll create a system-spec for your e-commerce platform.

```json
{
  "name": "ecommerce-platform",
  "description": "E-commerce platform with API and background workers",
  "version": "1.0.0",
  "services": {
    "api": {
      "image": { "name": "ecommerce-api", "tag": "latest" },
      "description": "Main REST API for e-commerce operations",
      "connections": {
        "worker": { "port": 5672, "protocol": "amqp" }
      },
      "aws": {
        "rds": [{ "name": "ecommerce-db", "engine": "postgres", "port": 5432 }]
      }
    },
    "worker": {
      "image": { "name": "ecommerce-worker", "tag": "latest" },
      "description": "Background job processor",
      "aws": {
        "sqs": [{ "name": "order-queue" }],
        "rds": [{ "name": "ecommerce-db", "engine": "postgres", "port": 5432 }]
      }
    }
  }
}
```

Then validate and render:
```bash
system-spec validate ecommerce.json
system-spec render ecommerce.json --format d2 > ecommerce.d2
```

## Validation Rules

The spec must satisfy:

1. **name** is required (non-empty string)
2. **services** must have at least one service
3. **connections** must reference existing services
4. **ports** must be 1-65535 when specified
5. Each service needs an **image.name**

## Tips

- Use lowercase hyphenated names for services (e.g., `order-api`, not `OrderAPI`)
- Group related cloud resources under the service that uses them
- Use connections to show runtime dependencies, not build dependencies
- Include repo URLs to link services to source code
