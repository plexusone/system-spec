# Go SDK Usage

## Installation

```bash
go get github.com/plexusone/system-spec@latest
```

## Packages

| Package | Description |
|---------|-------------|
| `spec` | Core types and validation |
| `graph` | Graph representation and conversion |
| `render` | Multi-format renderers |
| `schema` | JSON Schema generation |

## Loading a System Spec

```go
import "github.com/plexusone/system-spec/spec"

// From file
sys, err := spec.LoadFromFile("system.json")
if err != nil {
    log.Fatal(err)
}

// From JSON bytes
sys, err := spec.LoadFromJSON(jsonData)
```

## Creating a System Programmatically

```go
sys := &spec.System{
    Name: "my-system",
    Services: map[string]spec.Service{
        "api": {
            Image: spec.ContainerImage{
                Name: "myorg/api",
                Tag:  "v1.0",
            },
            Repo: &spec.GitRepo{
                URL: "https://github.com/myorg/api",
            },
            Connections: map[string]spec.Connection{
                "database": {
                    Port:     5432,
                    Protocol: "sql",
                },
            },
        },
    },
}

// Validate
if err := sys.Validate(); err != nil {
    log.Fatal(err)
}

// Serialize to JSON
jsonData, err := sys.ToJSON()
```

## Converting to Graph

```go
import "github.com/plexusone/system-spec/graph"

// Convert system to graph
g := graph.FromSystem(sys)

// Access nodes and edges
for _, node := range g.Nodes {
    fmt.Printf("Node: %s (%s)\n", node.ID, node.Kind)
}

for _, edge := range g.Edges {
    fmt.Printf("Edge: %s -> %s (%s:%d)\n",
        edge.Source, edge.Target, edge.Protocol, edge.Port)
}

// Query edges
edgesFrom := g.EdgesFrom("svc:api")
edgesTo := g.EdgesTo("rds:main-db")
```

## Rendering

```go
import "github.com/plexusone/system-spec/render"

// Create renderers
renderers := render.NewRenderers()

// Render to D2
d2Output, err := renderers.D2.Render(g)

// Render to Mermaid
mermaidOutput, err := renderers.Mermaid.Render(g)

// Render to Cytoscape.js JSON
cytoOutput, err := renderers.Cytoscape.Render(g)

// Get renderer by format
r := renderers.Get(render.FormatD2)
output, err := r.Render(g)
```

## Generating JSON Schema

```go
import "github.com/plexusone/system-spec/schema"

// Generate schema JSON
schemaJSON, err := schema.Generate()
```

## Working with Images

```go
img := spec.ContainerImage{
    Name:   "ghcr.io/org/api",
    Tag:    "v1.0",
    Digest: "sha256:abc123",
}

// Get full reference (digest takes precedence)
ref := img.FullName() // "ghcr.io/org/api@sha256:abc123"
```

## Error Handling

All functions return errors following Go conventions:

```go
sys, err := spec.LoadFromFile("system.json")
if err != nil {
    // Could be:
    // - File not found
    // - JSON parse error
    // - Validation error
    log.Fatalf("Failed to load: %v", err)
}
```
