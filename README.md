# System Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/system-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/system-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/system-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/system-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/system-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/system-spec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/system-spec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/system-spec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/system-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/system-spec
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/system-spec
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fsystem-spec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/system-spec
 [repo-url]: https://github.com/plexusone/system-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/system-spec/blob/main/LICENSE

A Go-native system topology specification for modeling microservices, infrastructure, and connectivity.

## Features

- 🔒 **Strongly-typed** - Go structs are source of truth; no `interface{}` or `any`
- 📄 **JSON-first** - JSON format with generated JSON Schema
- 🔗 **Graph-aware** - Services as nodes, connections as edges
- ☁️ **Multi-cloud** - AWS, GCP, and Cloudflare resource bindings
- 🎨 **Multi-render** - D2, Mermaid, Cytoscape.js, Sigma.js, GraphViz

## Installation

### CLI

```bash
go install github.com/plexusone/system-spec/cmd/system-spec@latest
```

### Go SDK

```bash
go get github.com/plexusone/system-spec@latest
```

## Quick Start

Create `system.json`:

```json
{
  "name": "my-system",
  "services": {
    "api": {
      "image": { "name": "myorg/api", "tag": "v1.0" },
      "repo": { "url": "https://github.com/myorg/api" },
      "connections": {
        "database": { "port": 5432, "protocol": "sql" }
      }
    },
    "database": {
      "image": { "name": "postgres", "tag": "15" }
    }
  }
}
```

Validate and render:

```bash
# Validate
system-spec validate system.json

# Render to D2
system-spec render system.json --format d2 > system.d2
d2 system.d2 system.svg

# Render to Mermaid (GitHub-compatible)
system-spec render system.json --format mermaid

# Render to Cytoscape.js (web visualization)
system-spec render system.json --format cytoscape > system.json
```

## Go SDK Usage

```go
import (
    "github.com/plexusone/system-spec/spec"
    "github.com/plexusone/system-spec/graph"
    "github.com/plexusone/system-spec/render"
)

// Load and validate
sys, err := spec.LoadFromFile("system.json")
if err != nil {
    log.Fatal(err)
}

// Convert to graph
g := graph.FromSystem(sys)

// Render
renderers := render.NewRenderers()
output, err := renderers.D2.Render(g)
```

## Specification

| Version | Spec | Schema |
|---------|------|--------|
| v0.1.0 | [spec.md](versions/v0.1.0/spec.md) | [system.schema.json](versions/v0.1.0/system.schema.json) |

The human-readable spec (`spec.md`) is for developers. The JSON Schema is for tools and AI agents.

## Supported Resources

### Cloud Providers

| Provider | Resources |
|----------|-----------|
| AWS | RDS, DynamoDB, SQS, SNS, S3, Bedrock |
| GCP | CloudSQL, PubSub, GCS |
| Cloudflare | Workers, R2 |

### Deployment Tools

| Tool | Support |
|------|---------|
| Helm | Chart references with service mapping |
| Terraform | Module references with resource mapping |

## Rendering Formats

| Format | Output | Use Case |
|--------|--------|----------|
| D2 | Text | Static diagrams, documentation |
| Mermaid | Text | GitHub Markdown, wikis |
| Cytoscape | JSON | Interactive web visualization |
| Sigma | JSON | Large graph visualization |
| DOT | Text | GraphViz, PDF/PNG export |

## Design Principles

1. **Go-first**: Go structs → JSON Schema (not the reverse)
2. **No polymorphism**: `map[string]Service` not `map[string]interface{}`
3. **Schema validation**: Must pass [schemalint](https://github.com/grokify/schemalint)
4. **Graph-native**: Built for visualization and dependency analysis

## Documentation

- [Full Documentation](https://plexusone.github.io/system-spec/) (MkDocs)
- [Getting Started](docs/getting-started/quickstart.md)
- [CLI Reference](docs/getting-started/cli.md)
- [Go SDK](docs/sdk/usage.md)
- [Examples](docs/examples/payments.md)

## Project Structure

```
system-spec/
├── spec/           # Core Go types (System, Service, etc.)
├── graph/          # Graph representation and conversion
├── render/         # Multi-format renderers
├── schema/         # JSON Schema generation
├── cmd/system-spec/# CLI tool
├── versions/       # Versioned specs and schemas
│   └── v0.1.0/
│       ├── spec.md              # Human-readable
│       └── system.schema.json   # Machine-readable
├── docs/           # MkDocs documentation
└── examples/       # Example system specs
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `go test ./...` and `golangci-lint run`
5. Submit a pull request

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

Generated from [CHANGELOG.json](CHANGELOG.json) using [structured-changelog](https://github.com/grokify/structured-changelog).

## License

MIT License - see [LICENSE](LICENSE) for details.
