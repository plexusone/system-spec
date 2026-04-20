# System-Spec Tasks

Go-native system topology specification for modeling SaaS infrastructure.

## Overview

This task list tracks pending features and improvements for system-spec.

---

## Phase 1 - Core Spec ✅ COMPLETE

### Spec Types ✅
- [x] System, Service, GitRepo, ContainerImage
- [x] AWS resources: RDS, DynamoDB, SQS, SNS, S3, Bedrock
- [x] GCP resources: CloudSQL, PubSub, GCS
- [x] Cloudflare resources: Workers, R2
- [x] Service connections with protocol and port
- [x] Deployment mappings: Helm charts, Terraform modules

### Rendering ✅
- [x] Graph intermediate representation
- [x] D2 renderer with styled node shapes
- [x] Mermaid renderer with provider subgraphs
- [x] Cytoscape.js JSON renderer
- [x] Sigma.js JSON renderer
- [x] GraphViz DOT renderer

### CLI ✅
- [x] `validate` - Validate spec files
- [x] `render` - Render to various formats
- [x] `schema` - Generate JSON Schema
- [x] `graph` - Output graph JSON

### Documentation ✅
- [x] Specification document v0.1.0
- [x] JSON Schema generation from Go types
- [x] MkDocs documentation site
- [x] Example payments-system.json

### Integration ✅
- [x] Graphize provider for code graph integration

### AI Agent Skills ✅
- [x] AssistantKit-compatible skill specs
- [x] `create-system-spec` skill - Guide agents to build new specs
- [x] `analyze-system-spec` skill - Validate and visualize existing specs
- [x] Canonical format (converts to Claude, Codex, Gemini, Kiro)

---

## Phase 2 - Cloud Provider Expansion 🔶 PLANNED

### Azure Resources
- [ ] Azure SQL Database
- [ ] Azure Service Bus (queues and topics)
- [ ] Azure Blob Storage
- [ ] Azure Functions
- [ ] Azure Cosmos DB
- [ ] Azure OpenAI Service

### Additional AWS Resources
- [ ] ElastiCache (Redis/Memcached)
- [ ] Lambda functions
- [ ] API Gateway
- [ ] CloudFront
- [ ] ECS/EKS clusters

### Additional GCP Resources
- [ ] Cloud Run
- [ ] Cloud Functions
- [ ] BigQuery
- [ ] Memorystore
- [ ] Cloud CDN

---

## Phase 3 - Import & Discovery 🔶 PLANNED

### Kubernetes Import
- [ ] Import from Kubernetes cluster state
  - [ ] Service discovery from deployments/services
  - [ ] ConfigMap/Secret references
  - [ ] Ingress rules → service connections
- [ ] Import from Helm values.yaml
- [ ] Import from Kustomize overlays

### Terraform Import
- [ ] Import from Terraform state files
  - [ ] Parse tfstate JSON format
  - [ ] Map resources to system-spec types
  - [ ] Extract dependencies as connections
- [ ] Import from Terraform plan output
- [ ] Support for common provider resources

### SBOM Integration
- [ ] CycloneDX import/export
  - [ ] Map components to services
  - [ ] Import vulnerability data
  - [ ] Track component versions
- [ ] SPDX support
- [ ] Dependency graph extraction

---

## Phase 4 - Analysis & Comparison 🔶 PLANNED

### Diff/Compare
- [ ] `system-spec diff <old.json> <new.json>`
  - [ ] Service additions/removals
  - [ ] Connection changes
  - [ ] Resource changes
  - [ ] Breaking change detection
- [ ] JSON patch output for automation
- [ ] Human-readable diff report

### Validation Enhancements
- [ ] Circular dependency detection
- [ ] Orphan resource detection
- [ ] Port conflict detection
- [ ] Cross-system validation

---

## Phase 5 - Agent Integration 🔶 PLANNED

### MCP Server (for querying)
- [ ] `system-spec serve` - Start MCP server
  - [ ] Tool: query_system (list services, connections)
  - [ ] Tool: get_service (service details)
  - [ ] Tool: get_dependencies (service dependencies)
  - [ ] Tool: find_path (trace connection paths)
- [ ] Integration with Claude Code
- [ ] Integration with other AI assistants

Note: For **creating** specs, use the AssistantKit skills (Phase 1) which guide
agents to write JSON directly using the schema. MCP is better suited for
**querying** existing specs.

### API/SDK
- [ ] REST API server option
- [ ] gRPC API for programmatic access
- [ ] Watch mode for spec file changes

---

## Phase 6 - CI/CD Integration ⬜ FUTURE

### GitHub Actions
- [ ] `system-spec-action` for GitHub
  - [ ] Validate spec on PR
  - [ ] Generate diagram artifacts
  - [ ] Diff specs between branches
  - [ ] Post diagram to PR comments
- [ ] GitLab CI template
- [ ] Azure DevOps pipeline template

### Deployment Tracking
- [ ] Link deployments to Git commits
- [ ] Environment-specific specs (dev/staging/prod)
- [ ] Deployment history tracking

---

## Legend

- [x] Implemented
- [ ] Not started
- ✅ Complete
- 🔶 Planned (next phases)
- ⬜ Future
