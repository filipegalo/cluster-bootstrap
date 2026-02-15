# Cluster Bootstrap

GitOps repo for bootstrapping Kubernetes clusters with ArgoCD using the **App of Apps** pattern.

## Documentation

Full documentation is available at the [MkDocs site](docs/index.md). To preview locally:

```bash
pip install mkdocs-material
mkdocs serve
```

Online documentation available at [Cluster Boostrap Docs](https://user-cube.github.io/cluster-bootstrap/)

## Prerequisites

- `kubectl` configured with access to the target cluster
- `helm` (for local template testing)
- `sops` and `age` (for secrets encryption/decryption)
- `go` 1.25+ (to build the CLI)
- `task` (task runner for CLI development)
- SSH private key with read access to this repo

## Quick Start

### 1. Build the CLI

```bash
cd cli
task build
```

### 2. Initialize secrets (first time only)

```bash
./cli/cluster-bootstrap init
```

### 3. Bootstrap the cluster

```bash
./cli/cluster-bootstrap bootstrap dev
```

This will:

1. Decrypt environment secrets using SOPS + age
2. Create the `argocd` namespace and SSH credentials secret
3. Install ArgoCD via Helm
4. Deploy the root **App of Apps** Application

### 4. Access ArgoCD UI

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Get the initial admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d
```

## Architecture

```
CLI bootstrap  →  ArgoCD + App of Apps (root Application)
                        ↓
                   apps/ (Helm chart generating Application CRs)
                        ↓
              components/argocd/  (self-managed ArgoCD)
              components/xxx/    (other components)
```

ArgoCD manages itself — changes pushed to this repo are automatically synced.

## Components

| Component | Namespace | Sync Wave | Description |
|-----------|-----------|-----------|-------------|
| ArgoCD | `argocd` | 0 | Self-managed GitOps controller |
| Vault | `vault` | 1 | Secrets management |
| External Secrets | `external-secrets` | 1 | Syncs external secrets into Kubernetes |
| Prometheus Operator CRDs | `monitoring` | 2 | CRDs for the monitoring stack |
| ArgoCD Repo Secret | `argocd` | 2 | SSH credentials for repo access |
| Reloader | `reloader` | 2 | Restarts pods on ConfigMap/Secret changes |
| Kube Prometheus Stack | `monitoring` | 3 | Prometheus monitoring stack |
| Trivy Operator | `trivy-system` | 3 | Vulnerability scanning |

## CLI Commands

| Command | Description |
|---------|-------------|
| `bootstrap <env>` | Full cluster bootstrap (decrypt secrets, install ArgoCD, deploy App of Apps) |
| `init` | Interactive setup for SOPS config and encrypted secrets files |
| `vault-token` | Store Vault root token as Kubernetes secret |

## Environments

| Environment | Values File | Description |
|-------------|-------------|-------------|
| dev | `apps/values/dev.yaml` | Local/development clusters, minimal resources |
| staging | `apps/values/staging.yaml` | Pre-production, moderate resources |
| prod | `apps/values/prod.yaml` | Production, HA configuration |
