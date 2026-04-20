# Cytoscape.js Renderer

[Cytoscape.js](https://js.cytoscape.org/) is a graph visualization library for the web.

## Usage

=== "CLI"

    ```bash
    system-spec render system.json --format cytoscape > system.cyto.json
    ```

=== "Go"

    ```go
    renderers := render.NewRenderers()
    output, err := renderers.Cytoscape.Render(g)
    ```

## Output Format

```json
{
  "elements": {
    "nodes": [
      {
        "data": {
          "id": "svc:api",
          "label": "api",
          "kind": "service",
          "image": "myorg/api:v1.0"
        }
      }
    ],
    "edges": [
      {
        "data": {
          "id": "api->backend",
          "source": "svc:api",
          "target": "svc:backend",
          "label": "grpc:8080",
          "protocol": "grpc",
          "port": 8080
        }
      }
    ]
  }
}
```

## Web Integration

```html
<!DOCTYPE html>
<html>
<head>
  <script src="https://unpkg.com/cytoscape/dist/cytoscape.min.js"></script>
</head>
<body>
  <div id="cy" style="width: 100%; height: 600px;"></div>
  <script>
    fetch('system.cyto.json')
      .then(res => res.json())
      .then(data => {
        cytoscape({
          container: document.getElementById('cy'),
          elements: data.elements,
          style: [
            {
              selector: 'node',
              style: {
                'label': 'data(label)',
                'background-color': '#4fc3f7'
              }
            },
            {
              selector: 'node[kind="database"]',
              style: { 'background-color': '#ffb74d' }
            },
            {
              selector: 'edge',
              style: {
                'label': 'data(label)',
                'curve-style': 'bezier',
                'target-arrow-shape': 'triangle'
              }
            }
          ],
          layout: { name: 'cose' }
        });
      });
  </script>
</body>
</html>
```

## Node Properties

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `label` | Display name |
| `kind` | Node type (service, database, etc.) |
| `provider` | Cloud provider (aws, gcp, cloudflare) |
| `image` | Container image (for services) |
| `repo` | Git repository URL |

## Edge Properties

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `source` | Source node ID |
| `target` | Target node ID |
| `label` | Display label (protocol:port) |
| `kind` | Edge type (connection, database, etc.) |
| `protocol` | Connection protocol |
| `port` | Target port |
