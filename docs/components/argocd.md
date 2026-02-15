# ArgoCD

**Namespace:** `argocd` | **Sync Wave:** 0 | **Chart:** `argo-cd` v7.8.13

ArgoCD is the GitOps controller that manages all components, including itself. It is the first component deployed (wave 0) and is self-managed after the initial bootstrap.

## Upstream Chart

- **Chart:** `argo-cd`
- **Version:** 7.8.13
- **Repository:** `https://argoproj.github.io/argo-helm`

## Key Configuration

**Base values:**

- CRDs installed and kept on deletion
- 1 replica for server, controller, repo-server, and applicationset-controller

**Per-environment resources:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| Server replicas | 1 | 2 | 3 |
| Controller replicas | 1 | 1 | 2 |
| Repo Server replicas | 1 | 2 | 3 |
| ApplicationSet replicas | 1 | 1 | 2 |
| CPU request | 50m | 100m | 250m |
| Memory request | 128Mi | 256Mi | 512Mi |
| CPU limit | 200m | 500m | 500m |
| Memory limit | 256Mi | 512Mi | 2Gi |

## Self-Management

After the CLI performs the initial Helm install, ArgoCD takes over its own management:

1. The App of Apps creates an ArgoCD Application pointing to `components/argocd/`
2. ArgoCD syncs this Application and manages its own deployment
3. Upgrades are done by updating the chart version in `Chart.yaml` and pushing to Git

## Files

```
components/argocd/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
