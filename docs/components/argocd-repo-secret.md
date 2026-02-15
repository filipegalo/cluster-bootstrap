# ArgoCD Repo Secret

**Namespace:** `argocd` | **Sync Wave:** 2 | **Chart:** Custom (no upstream dependency)

This component creates the Git repository credentials that ArgoCD uses to pull manifests. It uses an ExternalSecret to sync the SSH private key from Vault into a Kubernetes Secret.

## Chart Type

This is a custom chart with no upstream dependency — it only contains custom templates.

## Key Configuration

**Base values:**

- Provider type: `vault`
- Vault KV path for SSH key
- Kubernetes auth configuration for Vault
- Repository URL: `git@github.com:user-cube/cluster-bootstrap.git`

**Per-environment:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| Vault address | `http://vault.vault.svc:8200` | `http://vault.vault.svc:8200` | `http://vault.vault.svc:8200` |
| Refresh interval | default | default | 30m |

## Templates

- **`external-secret.yaml`** — ExternalSecret that fetches the SSH key from Vault and creates a Kubernetes Secret in the `argocd` namespace with the correct labels for ArgoCD to recognize it as a repository credential
- **`secret-store.yaml`** — SecretStore that configures the connection to Vault using Kubernetes auth

## Secret Flow

```
Vault KV Store
  └── SSH private key
        │
        ▼
ExternalSecret (watches Vault)
        │
        ▼
Kubernetes Secret (argocd namespace)
        │
        ▼
ArgoCD (uses secret for Git access)
```

## Files

```
components/argocd-repo-secret/
├── Chart.yaml
├── templates/
│   ├── external-secret.yaml
│   └── secret-store.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
