# Quick Start

This guide walks through bootstrapping a cluster from scratch.

## 1. Build the CLI

```bash
cd cli
task build
```

This produces the `cluster-bootstrap` binary in the `cli/` directory.

## 2. Initialize secrets (first time only)

Run the interactive init command to set up SOPS encryption and create per-environment secrets files:

```bash
./cli/cluster-bootstrap init
```

This will:

1. Prompt you to choose a SOPS provider (age, AWS KMS, or GCP KMS)
2. Collect your encryption key
3. Generate a `.sops.yaml` configuration
4. Interactively collect secrets for each environment (repo URL, target revision, SSH private key)
5. Create encrypted `secrets.<env>.enc.yaml` files

## 3. Bootstrap the cluster

Run the bootstrap command with your target environment:

```bash
./cli/cluster-bootstrap bootstrap dev
```

This performs the following steps:

1. Decrypts environment secrets using SOPS + your age key
2. Creates the `argocd` namespace
3. Creates the `repo-ssh-key` Secret with your Git SSH credentials
4. Installs ArgoCD via Helm (using `components/argocd/` chart and values)
5. Deploys the App of Apps root Application
6. Prints ArgoCD access instructions

### Common flags

```bash
# Use a specific secrets file
./cli/cluster-bootstrap bootstrap dev --secrets-file ./my-secrets.enc.yaml

# Use a specific kubeconfig or context
./cli/cluster-bootstrap bootstrap dev --kubeconfig ~/.kube/my-config --context my-cluster

# Specify age key location
./cli/cluster-bootstrap bootstrap dev --age-key-file ./age-key.txt

# Dry run — print manifests without applying
./cli/cluster-bootstrap bootstrap dev --dry-run

# Skip ArgoCD Helm install (if already installed)
./cli/cluster-bootstrap bootstrap dev --skip-argocd-install
```

## 4. Access ArgoCD

After bootstrap completes, access the ArgoCD UI:

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Get the initial admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

Open [https://localhost:8080](https://localhost:8080) and log in with `admin` and the password above.

## 5. Store Vault token (non-dev environments)

For staging and production, after Vault initializes you need to store the root token:

```bash
./cli/cluster-bootstrap vault-token --token <vault-root-token>
```

This creates a `vault-root-token` Secret in the `vault` namespace, which the Vault configuration and seed jobs use.

## What happens next?

Once ArgoCD is running with the App of Apps deployed, it will automatically sync all enabled components in sync wave order:

1. **Wave 0**: ArgoCD (self-manages)
2. **Wave 1**: Vault, External Secrets
3. **Wave 2**: Prometheus Operator CRDs, ArgoCD Repo Secret, Reloader
4. **Wave 3**: Kube Prometheus Stack, Trivy Operator

All components use automated sync with pruning and self-healing enabled.
