# Prerequisites

Before bootstrapping a cluster, ensure the following tools are installed.

## Required Tools

| Tool | Purpose | Installation |
|------|---------|-------------|
| `kubectl` | Kubernetes CLI | [Install kubectl](https://kubernetes.io/docs/tasks/tools/) |
| `helm` | Helm package manager | [Install Helm](https://helm.sh/docs/intro/install/) |
| `git-crypt` | Encrypted secrets in Git | [Install git-crypt](https://github.com/AGWA/git-crypt) |
| `go` | Required to build the CLI tool | [Install Go](https://go.dev/doc/install) (1.25+) |
| `task` | Task runner for CLI development | [Install Task](https://taskfile.dev/installation/) |

## Cluster Access

You need a running Kubernetes cluster with `kubectl` configured to access it. Any conformant cluster works — local (kind, minikube, k3s) or cloud-managed (EKS, GKE, AKS).

Verify access:

```bash
kubectl cluster-info
kubectl get nodes
```

## SSH Key

An SSH key pair with read access to the Git repository is required. ArgoCD uses this key to pull manifests from the repo.

```bash
ssh-keygen -t ed25519 -f repo-ssh-key.pem -N ""
```

Add the public key (`repo-ssh-key.pem.pub`) as a deploy key in your repository settings.

## git-crypt

Secrets files (`secrets.*.yaml`) are encrypted in the repo using git-crypt. After cloning, run `git-crypt unlock` (with your key or GPG) so the CLI can read them. See [Secrets Management](../guides/secrets-management.md) for setup.

## MkDocs (optional)

To preview this documentation site locally:

```bash
pip install mkdocs-material
mkdocs serve
```
