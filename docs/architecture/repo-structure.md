# Repository Structure

```
cluster-bootstrap/
├── apps/                          # App of Apps Helm chart
│   ├── Chart.yaml                 # Chart metadata (no dependencies)
│   ├── templates/                 # Application CR templates
│   │   ├── argocd.yaml
│   │   ├── vault.yaml
│   │   ├── external-secrets.yaml
│   │   ├── argocd-repo-secret.yaml
│   │   ├── prometheus-operator-crds.yaml
│   │   ├── kube-prometheus-stack.yaml
│   │   ├── reloader.yaml
│   │   └── trivy-operator.yaml
│   └── values/                    # Per-environment toggle files
│       ├── dev.yaml
│       ├── staging.yaml
│       └── prod.yaml
├── components/                    # Individual component Helm charts
│   ├── argocd/
│   ├── vault/
│   ├── external-secrets/
│   ├── argocd-repo-secret/
│   ├── prometheus-operator-crds/
│   ├── kube-prometheus-stack/
│   ├── reloader/
│   └── trivy-operator/
├── cli/                           # Go CLI tool
│   ├── main.go
│   ├── cmd/                       # Cobra commands
│   ├── internal/                  # Internal packages
│   ├── Taskfile.yml
│   ├── go.mod
│   └── go.sum
├── docs/                          # Documentation (this site)
├── mkdocs.yml                     # MkDocs configuration
├── .gitignore
├── .sops.yaml                     # SOPS encryption rules
└── README.md
```

## `apps/` — App of Apps

The root Helm chart that ArgoCD deploys. Each template in `apps/templates/` is a Kubernetes `Application` custom resource pointing to one component. The `apps/values/` files control which components are enabled per environment.

All environments currently enable the same components, but the structure allows selectively disabling components per environment by setting `<component>: enabled: false`.

## `components/` — Platform Components

Each subdirectory is a standalone Helm chart (or chart wrapper) for one platform component. The common structure is:

```
components/<name>/
├── Chart.yaml          # Declares upstream chart dependency
├── templates/          # Optional custom templates
└── values/
    ├── base.yaml       # Shared defaults
    ├── dev.yaml        # Dev overrides
    ├── staging.yaml    # Staging overrides
    └── prod.yaml       # Prod overrides
```

Most components are thin wrappers around upstream Helm charts — `Chart.yaml` declares the dependency and `values/` files configure it. Some components (like Vault and ArgoCD Repo Secret) include custom templates for additional resources.

## `cli/` — Bootstrap CLI

A Go application that automates cluster bootstrapping. Built with Cobra (CLI framework), it handles SOPS decryption, Helm installation, and Kubernetes resource creation. See the [CLI documentation](../cli/index.md) for details.

## Config Files

| File | Purpose |
|------|---------|
| `.sops.yaml` | SOPS encryption rules — defines which files to encrypt and with which key |
| `.gitignore` | Ignores charts/, secrets, binaries, IDE files, and MkDocs build output |
| `age-key.txt` | Age private key for SOPS decryption (gitignored) |
