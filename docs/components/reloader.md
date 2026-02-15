# Reloader

**Namespace:** `reloader` | **Sync Wave:** 2 | **Chart:** `reloader` v2.2.8

Stakater Reloader watches for changes to ConfigMaps and Secrets, and automatically triggers rolling restarts on workloads that reference them.

## Upstream Chart

- **Chart:** `reloader`
- **Version:** 2.2.8
- **Repository:** `https://stakater.github.io/stakater-charts`

## Key Configuration

**Base values:**

- Global watch enabled (monitors all namespaces)

**Per-environment resources:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| CPU request | 10m | 50m | 100m |
| Memory request | 32Mi | 64Mi | 128Mi |
| Memory limit | 128Mi | 256Mi | 512Mi |

## Why Reloader?

When a Secret or ConfigMap changes (e.g., a credential rotation via External Secrets), Kubernetes does not automatically restart pods that mount the old version. Reloader detects these changes and performs rolling updates to ensure workloads always use the latest configuration.

## Files

```
components/reloader/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
