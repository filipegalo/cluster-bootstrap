# Secrets Management

This repo uses a multi-layer secrets architecture: git-crypt for encryption at rest in Git, Vault for runtime secrets storage, and External Secrets Operator for syncing secrets into Kubernetes.

## Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Developer Machine                      │
│                                                           │
│  secrets.<env>.yaml (git-crypt: decrypt with git-crypt   │
│  unlock; files are encrypted in the Git repo)             │
│                                                           │
│  CLI bootstrap reads secrets and:                         │
│    1. Creates repo-ssh-key Secret in argocd namespace     │
│    2. Vault seed job copies SSH key into Vault KV         │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                      │
│                                                           │
│  Vault (KV Store)                                         │
│    └── SSH private key                                    │
│           │                                               │
│           ▼                                               │
│  External Secrets Operator                                │
│    └── ExternalSecret (watches Vault)                     │
│           │                                               │
│           ▼                                               │
│  Kubernetes Secret (argocd namespace)                     │
│    └── ArgoCD uses for Git repo access                    │
└─────────────────────────────────────────────────────────┘
```

## git-crypt

[git-crypt](https://github.com/AGWA/git-crypt) transparently encrypts files in the Git repository. Files matching the patterns in `.gitattributes` are stored encrypted in the repo and decrypted in your working tree when you run `git-crypt unlock`.

### Configuration

`.gitattributes` defines which files are encrypted:

```
secrets.*.yaml filter=git-crypt diff=git-crypt
```

Any file matching `secrets.*.yaml` will be encrypted by git-crypt when committed.

### Secrets file structure

Each environment has a `secrets.<env>.yaml` file containing:

```yaml
repo:
  url: git@github.com:filipegalo/cluster-bootstrap.git
  targetRevision: main
  sshPrivateKey: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    ...
    -----END OPENSSH PRIVATE KEY-----
```

### Working with encrypted files

```bash
# Initialize git-crypt in the repo (first time only)
git-crypt init

# Add a GPG user or symmetric key (see git-crypt documentation)
git-crypt add-gpg-user USER_ID
# or: git-crypt export-key ./key && git-crypt unlock ./key

# Unlock the repo so secrets.*.yaml are decrypted on disk
git-crypt unlock

# Lock the repo (secrets appear encrypted again in working tree)
git-crypt lock
```

### Initialize with the CLI

The `init` command sets up `.gitattributes` for git-crypt and creates per-environment secrets files interactively:

```bash
./cli/cluster-bootstrap init --output-dir .
```

If the repo is not yet using git-crypt, run `git-crypt init` in the repo root. Then add your key (GPG user or symmetric key) and commit. Secrets files will be encrypted when you commit.

## Vault Integration

[Vault](https://www.hashicorp.com/products/vault) provides runtime secrets storage in the cluster.

### Bootstrap flow

1. The CLI creates an initial `repo-ssh-key` Kubernetes Secret during bootstrap
2. Vault starts and the seed job copies the SSH key from the Kubernetes Secret into Vault's KV store
3. The config job sets up Kubernetes authentication in Vault

### Non-dev environments

For staging and production, Vault requires initialization:

```bash
# After Vault pods are running
kubectl exec -n vault vault-0 -- vault operator init

# Store the root token
./cli/cluster-bootstrap vault-token --token <root-token>
```

## External Secrets Operator

The [External Secrets Operator](https://external-secrets.io/) bridges Vault and Kubernetes Secrets.

### Components

- **SecretStore** — configures the connection to Vault (address, auth method)
- **ExternalSecret** — defines what to fetch from Vault and where to store it in Kubernetes

### ArgoCD Repo Secret flow

The `argocd-repo-secret` component creates:

1. A `SecretStore` pointing to Vault with Kubernetes auth
2. An `ExternalSecret` that fetches the SSH key from Vault KV
3. The operator creates a Kubernetes Secret in the `argocd` namespace
4. ArgoCD uses this Secret for Git repository access

This closes the loop — after bootstrap, credential rotation flows through Vault and External Secrets without manual intervention.
