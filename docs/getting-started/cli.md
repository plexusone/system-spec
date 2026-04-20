# CLI Reference

## Commands

### validate

Validate a system specification file.

```bash
system-spec validate <file.json>
```

**Example:**

```bash
system-spec validate system.json
```

**Output:**

```
valid: payments-platform (4 services)
  - api-gateway: ghcr.io/org/api-gateway:v2.1.0
      repo: https://github.com/org/api-gateway
      connections: 2
  - payments-service: ghcr.io/org/payments@sha256:abc123
      repo: https://github.com/org/payments
      connections: 1
```

### render

Render a system spec to a diagram format.

```bash
system-spec render <file.json> --format <format>
```

**Formats:**

| Format | Description | Output |
|--------|-------------|--------|
| `d2` | D2 diagram language | Text |
| `mermaid` | Mermaid syntax | Text |
| `cytoscape` | Cytoscape.js JSON | JSON |
| `sigma` | Sigma.js JSON | JSON |
| `dot` | GraphViz DOT | Text |

**Examples:**

```bash
# D2
system-spec render system.json --format d2 > system.d2

# Mermaid
system-spec render system.json --format mermaid > system.mmd

# Cytoscape.js (for web visualization)
system-spec render system.json --format cytoscape > system.cyto.json

# Sigma.js (for large graphs)
system-spec render system.json --format sigma > system.sigma.json

# GraphViz DOT
system-spec render system.json --format dot > system.dot
```

### graph

Output the intermediate graph representation as JSON.

```bash
system-spec graph <file.json>
```

**Example:**

```bash
system-spec graph system.json | jq '.nodes | length'
```

### schema

Output the JSON Schema for system-spec.

```bash
system-spec schema
```

**Example:**

```bash
system-spec schema > system.schema.json
```

### skill

Output AI agent skill instructions. These are prompts that guide AI agents (Claude, Codex, Gemini, etc.) through creating or analyzing system-spec files.

```bash
system-spec skill <name>
```

**Available skills:**

| Skill | Description |
|-------|-------------|
| `create` | Guide for creating a new system-spec |
| `analyze` | Guide for analyzing/validating existing specs |

**Examples:**

```bash
# List available skills
system-spec skill

# Output the create skill (for AI agents to read)
system-spec skill create

# Use in a prompt
SKILL=$(system-spec skill create)
echo "Using skill instructions to help build a spec..."
```

**Use case:** AI agents can fetch skill instructions without needing filesystem installation:

```bash
# AI agent runs this to get guidance
system-spec skill create
```

The output contains structured instructions for:

- Gathering requirements from the user
- Building valid JSON following the schema
- Validating and rendering the result
- Integrating with Graphize for code graph linking

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid file, validation failure, etc.) |

## Environment Variables

None currently. All configuration is via command-line flags.
