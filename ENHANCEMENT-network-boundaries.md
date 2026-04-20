# Enhancement: Network Boundary Rendering from `exposes` Field

## Summary

Renderers (D2, Mermaid, Cytoscape) should use the `exposes` field and/or `networks` definitions to visually group services into network boundary containers (e.g., "Public-Facing" vs "Internal").

## Problem

Services that expose ports to external traffic (load balancers, ingress) are visually indistinguishable from internal-only services in rendered diagrams. Users must manually post-process D2 output to add container boundaries, which breaks the generate-from-spec workflow.

## Proposed Behavior

### Option A: Derive from `exposes`

Add a convention where `exposes` entries with a `visibility` field indicate network exposure:

```json
"ecm": {
  "exposes": [
    { "port": 8080, "protocol": "http", "visibility": "public" }
  ]
}
```

Services with `"visibility": "public"` are grouped into a "Public-Facing" container in D2/Mermaid. All others are "Internal".

### Option B: Use `networks` with service membership

```json
"networks": {
  "public": {
    "description": "Public-facing services (ingress/ALB)",
    "services": ["ecm", "ui-track"]
  },
  "internal": {
    "description": "Internal services (ClusterIP only)"
  }
}
```

Services listed in a network are grouped. Services not listed in any network default to "internal".

### Rendering

**D2**: Wrap services in named containers with distinct border colors.

```d2
public: "Public-Facing" {
  style.stroke: "#4ecca3"
  svc_ecm: ecm
}
internal: "Internal" {
  svc_authms: authms
}
```

**Mermaid**: Use `subgraph` blocks.

**Cytoscape**: Add `parent` field to node data for compound nodes.

## Current Workaround

Post-process D2 output with a script that moves node definitions into container blocks and qualifies edge references. See `eic-atlas` for an example.

## Priority

Medium — the workaround is functional but fragile and breaks when services are added/removed.
