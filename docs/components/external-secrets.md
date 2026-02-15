# External Secrets

**Namespace:** `external-secrets` | **Sync Wave:** 1 | **Chart:** `external-secrets` v0.14.3

The External Secrets Operator syncs secrets from external providers (Vault, AWS Secrets Manager) into Kubernetes Secrets.

## Upstream Chart

- **Chart:** `external-secrets`
- **Version:** 0.14.3
- **Repository:** `https://charts.external-secrets.io`

## Key Configuration

**Base values:**

- CRDs installed
- 1 replica each for webhook and cert-controller

**Per-environment resources:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| Replicas | 1 | 2 | 3 |
| CPU request | 10m | 50m | 100m |
| Memory request | 32Mi | 64Mi | 128Mi |

## How It Works

External Secrets Operator watches for `ExternalSecret` and `SecretStore`/`ClusterSecretStore` custom resources. When it finds one, it:

1. Connects to the configured secret provider (Vault in this repo)
2. Fetches the secret data
3. Creates or updates a Kubernetes Secret with the fetched data
4. Periodically refreshes the secret based on the configured interval

This is used by the [ArgoCD Repo Secret](argocd-repo-secret.md) component to sync the Git SSH key from Vault into a Kubernetes Secret that ArgoCD can use.

## Files

```
components/external-secrets/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
