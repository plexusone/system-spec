# Rendering Overview

System-spec supports multiple output formats for visualization and integration.

## Supported Formats

| Format | Type | Use Case |
|--------|------|----------|
| [D2](d2.md) | Text | Static diagrams, documentation |
| [Mermaid](mermaid.md) | Text | GitHub, Markdown embedding |
| [Cytoscape.js](cytoscape.md) | JSON | Interactive web visualization |
| [Sigma.js](sigma.md) | JSON | Large graph visualization |
| DOT | Text | GraphViz, PDF/PNG generation |

## Architecture

```
System → Graph → Renderer → Output
```

1. **System**: The spec document (JSON)
2. **Graph**: Intermediate representation (nodes + edges)
3. **Renderer**: Format-specific output generator
4. **Output**: Text or JSON for the target tool

## CLI Usage

```bash
system-spec render <file.json> --format <format>
```

## Go SDK Usage

```go
// Load and convert
sys, _ := spec.LoadFromFile("system.json")
g := graph.FromSystem(sys)

// Render
renderers := render.NewRenderers()
output, _ := renderers.D2.Render(g)
```

## Node Styling

Each renderer maps node kinds to appropriate shapes:

| Kind | D2 | Mermaid | DOT |
|------|----|---------|----- |
| service | rectangle | `[label]` | box |
| database | cylinder | `[(label)]` | cylinder |
| queue | queue | `>label]` | cds |
| storage | stored_data | `[(label)]` | folder |
| ai_model | hexagon | `{{label}}` | hexagon |

## Edge Labels

Edges display protocol and port:

```
api-gateway -> payments-service: grpc:8080
```
