# Graphize Integration

System-spec integrates with [Graphize](https://github.com/plexusone/graphize) to combine infrastructure topology with code graphs. This enables powerful queries across your entire system - from services to the code that implements them.

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Central Docs Repo                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  system-spec.json          .graphize/                           │
│  ┌─────────────┐           ┌──────────────────────────────┐     │
│  │ services:   │           │ nodes/                        │     │
│  │   api:      │──extract─▶│   svc:api, svc:payments, ... │     │
│  │   payments: │           │   pkg:api/handler, ...        │     │
│  │   ...       │           │ edges/                        │     │
│  └─────────────┘           │   svc:api ──links_to──▶ repo  │     │
│                            └──────────────────────────────┘     │
│                                      ▲                           │
│                                      │                           │
│         ┌────────────────────────────┼────────────────────┐     │
│         │                            │                    │     │
│   graphize add              graphize add           graphize add │
│         │                            │                    │     │
└─────────┼────────────────────────────┼────────────────────┼─────┘
          │                            │                    │
    /local/api-service         /local/payments       /local/notifications
    (cloned repo)              (cloned repo)         (cloned repo)
```

## Multi-Repo Workflow

### Step 1: Create system-spec with service→repo mappings

```json
{
  "name": "platform",
  "description": "E-commerce platform",
  "services": {
    "api": {
      "image": { "name": "api-gateway", "tag": "v1.0" },
      "repo": { "url": "https://github.com/org/api-service" },
      "connections": {
        "payments": { "port": 8080, "protocol": "grpc" }
      }
    },
    "payments": {
      "image": { "name": "payments", "tag": "v1.0" },
      "repo": { "url": "https://github.com/org/payments" },
      "aws": {
        "rds": [{ "name": "payments-db", "engine": "postgres", "port": 5432 }]
      }
    },
    "notifications": {
      "image": { "name": "notifications", "tag": "v1.0" },
      "repo": { "url": "https://github.com/org/notifications" },
      "aws": {
        "sqs": [{ "name": "notification-queue" }]
      }
    }
  }
}
```

### Step 2: Initialize graphize in central repo

```bash
cd /path/to/central-docs-repo
graphize init
```

### Step 3: Add each service's repo (cloned locally)

```bash
# Add repos - they can be anywhere on your filesystem
graphize add /path/to/local/api-service
graphize add /path/to/local/payments
graphize add /path/to/local/notifications

# The system-spec file is auto-detected if present
# graphize add .  # adds current repo with system-spec
```

### Step 4: Analyze all sources

```bash
graphize analyze
```

This extracts:

- **From system-spec**: Service nodes, cloud resources, connections, deployments
- **From code repos**: Packages, functions, types, imports, call graphs

### Step 5: Query across the combined graph

```bash
# List all services
graphize query --type service

# Find what the api service connects to
graphize query svc:api

# Trace path from service to database
graphize path "svc:payments" "rds:payments-db"

# Find all code in a specific repo
graphize query --filter "file:/path/to/payments/*"
```

## Graph Structure

### Node Types from system-spec

| Node Type | ID Format | Description |
|-----------|-----------|-------------|
| `system` | `system:<name>` | Root system node |
| `service` | `svc:<name>` | Service from spec |
| `database` | `rds:<name>`, `cloudsql:<name>` | Database resources |
| `queue` | `sqs:<name>` | Message queues |
| `topic` | `sns:<name>`, `pubsub:<name>` | Pub/sub topics |
| `storage` | `s3:<name>`, `gcs:<name>` | Object storage |
| `helm_chart` | `helm:<name>` | Helm deployments |
| `terraform_module` | `terraform:<name>` | Terraform modules |

### Edge Types from system-spec

| Edge Type | From | To | Description |
|-----------|------|-----|-------------|
| `contains` | system | service | System contains service |
| `connects_to` | service | service | Service-to-service connection |
| `uses` | service | resource | Service uses cloud resource |
| `links_to` | service | repo | Service implemented by repo |
| `deploys` | helm_chart | service | Helm deploys service |
| `manages` | terraform | resource | Terraform manages resource |

### The `links_to` Edge

The `links_to` edge connects services to their source repositories:

```
svc:payments ──links_to──▶ repo:https://github.com/org/payments
                                │
                                ▼
                          (code graph from that repo)
                                │
                          pkg:payments/api
                          func:ProcessPayment
                          type:PaymentRequest
```

This enables queries like "show me all code that implements the payments service."

## Visualization

### Generate multi-page HTML site

For multi-service systems, generate a documentation site with per-service code graphs:

```bash
graphize export htmlsite -o ./site
open ./site/index.html
```

This creates:

```
site/
├── index.html              # System topology overview
└── services/
    ├── api/index.html      # API service code graph
    ├── payments/index.html # Payments service code graph
    └── .../
```

**Features:**

- **Index page**: System topology with service cards and statistics
- **Service pages**: Filtered code graphs for each service's repository
- **Interactive navigation**: Click services in the topology to drill down
- **Dark mode**: Use `--dark` flag for dark theme

```bash
# With dark mode and custom title
graphize export htmlsite -o ./docs --dark --title "Platform Architecture"
```

### Generate single-page HTML visualization

For simpler visualization without per-service breakdown:

```bash
graphize export html -o graph.html
open graph.html
```

The visualization shows:

- Service nodes (from system-spec) in one color
- Code nodes (from repos) in another color
- Infrastructure resources (databases, queues) with icons
- Connections and dependencies as edges

### Generate analysis report

```bash
graphize report -o GRAPH_REPORT.md
```

## Example Queries

### Find services that use a specific database

```bash
graphize query rds:payments-db --dir in
```

### Trace a request path

```bash
graphize path "svc:api" "svc:notifications"
```

### Find all functions in the payments service

```bash
# First find what repo the service links to
graphize query svc:payments

# Then query functions in that repo
graphize query --type function --filter "*/payments/*"
```

### Export for AI agents

```bash
# Token-efficient format
graphize export toon -o graph.toon

# Include in AI context
cat graph.toon
```

## Storing Graphs in Central Repo

The `.graphize/` directory can be committed to git:

```bash
# What to commit (portable)
git add .graphize/manifest.json
git add system-spec.json

# What to regenerate locally (large, derived)
# .graphize/nodes/
# .graphize/edges/
# .graphize/cache/
```

Add to `.gitignore`:

```gitignore
.graphize/nodes/
.graphize/edges/
.graphize/cache/
```

After cloning, regenerate with:

```bash
graphize rebuild
```

## CI/CD Integration

### GitHub Actions example

```yaml
name: Update Graph
on:
  push:
    paths:
      - 'system-spec.json'
      - '.graphize/manifest.json'

jobs:
  rebuild:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # Clone service repos
      - name: Clone service repos
        run: |
          git clone https://github.com/org/api-service ../api-service
          git clone https://github.com/org/payments ../payments

      - name: Install graphize
        run: go install github.com/plexusone/graphize@latest

      - name: Rebuild graph
        run: |
          graphize add ../api-service
          graphize add ../payments
          graphize analyze
          graphize export html -o graph.html

      - name: Upload visualization
        uses: actions/upload-artifact@v4
        with:
          name: graph-visualization
          path: graph.html
```

## Troubleshooting

### "Service not linked to code"

Ensure your system-spec has `repo.url` for each service:

```json
"services": {
  "api": {
    "repo": { "url": "https://github.com/org/api" }
  }
}
```

### "Repo not found"

Add the repo to graphize:

```bash
graphize add /path/to/local/clone
graphize status  # verify it's tracked
```

### "Graph is stale"

Re-analyze after code changes:

```bash
graphize status   # shows which repos have new commits
graphize analyze  # re-extract
```
