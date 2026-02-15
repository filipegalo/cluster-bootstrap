# Prometheus Operator CRDs

**Namespace:** `monitoring` | **Sync Wave:** 2 | **Chart:** `prometheus-operator-crds` v27.0.0

Installs the Custom Resource Definitions (CRDs) required by the Prometheus Operator. Deployed separately from the monitoring stack to ensure CRDs are available before any resources that depend on them.

## Upstream Chart

- **Chart:** `prometheus-operator-crds`
- **Version:** 27.0.0
- **Repository:** `https://prometheus-community.github.io/helm-charts`

## Key Configuration

No custom values — uses upstream defaults. The chart only installs CRDs.

## Why Separate?

The Kube Prometheus Stack chart can install CRDs itself, but deploying them separately provides:

- **Ordering guarantees** — CRDs exist before the stack tries to create CR instances (wave 2 before wave 3)
- **Independent lifecycle** — CRD updates don't require a full stack redeployment
- **Safer upgrades** — CRD changes can be reviewed and applied independently

## Sync Configuration

The ArgoCD Application for this component uses special sync options:

- **`ServerSideApply=true`** — required for large CRD resources that exceed annotation size limits
- **`Replace=true`** — ensures CRDs are fully replaced rather than patched

## Files

```
components/prometheus-operator-crds/
├── Chart.yaml
└── values/
```
