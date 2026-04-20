# Sigma.js Renderer

[Sigma.js](https://www.sigmajs.org/) is optimized for rendering large graphs.

## Usage

=== "CLI"

    ```bash
    system-spec render system.json --format sigma > system.sigma.json
    ```

=== "Go"

    ```go
    renderers := render.NewRenderers()
    output, err := renderers.Sigma.Render(g)
    ```

## Output Format

```json
{
  "nodes": [
    {
      "id": "svc:api",
      "label": "api",
      "x": 500,
      "y": 300,
      "size": 10,
      "color": "#4fc3f7",
      "kind": "service",
      "attributes": {
        "image": "myorg/api:v1.0",
        "repo": "https://github.com/myorg/api"
      }
    }
  ],
  "edges": [
    {
      "id": "api->backend",
      "source": "svc:api",
      "target": "svc:backend",
      "label": "grpc:8080",
      "size": 1,
      "color": "#78909c",
      "protocol": "grpc",
      "port": 8080
    }
  ]
}
```

## Features

- **Initial layout**: Circular layout with positions calculated
- **Node sizing**: Based on connection count (more edges = larger node)
- **Color coding**: By node kind

## Color Scheme

| Kind | Color |
|------|-------|
| service | `#4fc3f7` (light blue) |
| database | `#ffb74d` (orange) |
| queue/topic | `#ba68c8` (purple) |
| storage | `#81c784` (green) |
| ai_model | `#f06292` (pink) |

## Web Integration

```html
<!DOCTYPE html>
<html>
<head>
  <script src="https://unpkg.com/sigma/build/sigma.min.js"></script>
  <script src="https://unpkg.com/graphology/dist/graphology.umd.min.js"></script>
</head>
<body>
  <div id="sigma" style="width: 100%; height: 600px;"></div>
  <script>
    fetch('system.sigma.json')
      .then(res => res.json())
      .then(data => {
        const graph = new graphology.Graph();

        data.nodes.forEach(node => {
          graph.addNode(node.id, node);
        });

        data.edges.forEach(edge => {
          graph.addEdge(edge.source, edge.target, edge);
        });

        new Sigma(graph, document.getElementById('sigma'));
      });
  </script>
</body>
</html>
```

## When to Use Sigma.js

- Systems with 100+ services
- When performance matters more than visual polish
- Interactive exploration of large topologies
- Force-directed layout needed

For smaller systems (< 50 services), Cytoscape.js may provide a better experience.
