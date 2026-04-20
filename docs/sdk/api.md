# API Reference

## Package `spec`

### Types

#### System

```go
type System struct {
    Name        string                `json:"name"`
    Description string                `json:"description,omitempty"`
    Version     string                `json:"version,omitempty"`
    Services    map[string]Service    `json:"services"`
    Networks    map[string]Network    `json:"networks,omitempty"`
    Deployments *Deployments          `json:"deployments,omitempty"`
}
```

#### Service

```go
type Service struct {
    Description string                   `json:"description,omitempty"`
    Repo        *GitRepo                 `json:"repo,omitempty"`
    Image       ContainerImage           `json:"image"`
    Registry    string                   `json:"registry,omitempty"`
    Connections map[string]Connection    `json:"connections,omitempty"`
    Exposes     []Endpoint               `json:"exposes,omitempty"`
    AWS         *AWSResources            `json:"aws,omitempty"`
    GCP         *GCPResources            `json:"gcp,omitempty"`
    Cloudflare  *CloudflareResources     `json:"cloudflare,omitempty"`
}
```

#### ContainerImage

```go
type ContainerImage struct {
    Name   string `json:"name"`
    Tag    string `json:"tag,omitempty"`
    Digest string `json:"digest,omitempty"`
}

func (i ContainerImage) FullName() string
```

#### GitRepo

```go
type GitRepo struct {
    URL    string `json:"url"`
    Path   string `json:"path,omitempty"`
    Ref    string `json:"ref,omitempty"`
    Commit string `json:"commit,omitempty"`
}
```

#### Connection

```go
type Connection struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    Description string `json:"description,omitempty"`
}
```

### Functions

```go
func LoadFromFile(path string) (*System, error)
func LoadFromJSON(data []byte) (*System, error)
func (s *System) Validate() error
func (s *System) ToJSON() ([]byte, error)
func (s *System) ToJSONCompact() ([]byte, error)
```

---

## Package `graph`

### Types

#### Graph

```go
type Graph struct {
    Nodes []Node `json:"nodes"`
    Edges []Edge `json:"edges"`
}

func NewGraph() *Graph
func (g *Graph) AddNode(n Node)
func (g *Graph) AddEdge(e Edge)
func (g *Graph) NodeByID(id string) *Node
func (g *Graph) EdgesFrom(nodeID string) []Edge
func (g *Graph) EdgesTo(nodeID string) []Edge
```

#### Node

```go
type Node struct {
    ID       string            `json:"id"`
    Label    string            `json:"label"`
    Kind     NodeKind          `json:"kind"`
    Group    string            `json:"group,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

#### Edge

```go
type Edge struct {
    ID       string            `json:"id"`
    Source   string            `json:"source"`
    Target   string            `json:"target"`
    Label    string            `json:"label,omitempty"`
    Kind     EdgeKind          `json:"kind"`
    Protocol string            `json:"protocol,omitempty"`
    Port     int               `json:"port,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

#### NodeKind

```go
type NodeKind string

const (
    NodeKindService   NodeKind = "service"
    NodeKindDatabase  NodeKind = "database"
    NodeKindQueue     NodeKind = "queue"
    NodeKindTopic     NodeKind = "topic"
    NodeKindStorage   NodeKind = "storage"
    NodeKindCDN       NodeKind = "cdn"
    NodeKindWorker    NodeKind = "worker"
    NodeKindAIModel   NodeKind = "ai_model"
    NodeKindVPC       NodeKind = "vpc"
    NodeKindSubnet    NodeKind = "subnet"
    NodeKindExternal  NodeKind = "external"
    NodeKindHelm      NodeKind = "helm"
    NodeKindTerraform NodeKind = "terraform"
)
```

### Functions

```go
func FromSystem(s *spec.System) *Graph
```

---

## Package `render`

### Types

#### Format

```go
type Format string

const (
    FormatD2        Format = "d2"
    FormatMermaid   Format = "mermaid"
    FormatCytoscape Format = "cytoscape"
    FormatSigma     Format = "sigma"
    FormatDOT       Format = "dot"
    FormatJSON      Format = "json"
)
```

#### Renderer

```go
type Renderer interface {
    Format() Format
    Render(g *graph.Graph) ([]byte, error)
}
```

#### Renderers

```go
type Renderers struct {
    D2        *D2Renderer
    Mermaid   *MermaidRenderer
    Cytoscape *CytoscapeRenderer
    Sigma     *SigmaRenderer
    DOT       *DOTRenderer
}

func NewRenderers() *Renderers
func (r *Renderers) Get(f Format) Renderer
func SupportedFormats() []Format
```

---

## Package `schema`

### Functions

```go
func Generate() ([]byte, error)
```
