# CLI Tool

The `cluster-bootstrap` CLI automates cluster bootstrapping. Built in Go with [Cobra](https://github.com/spf13/cobra), it handles secret decryption, Helm installation, and Kubernetes resource creation.

## Building

```bash
cd cli
task build
```

This runs `go mod tidy` and builds the `cluster-bootstrap` binary in the `cli/` directory.

### Other Taskfile commands

| Command | Description |
|---------|-------------|
| `task build` | Build the binary |
| `task clean` | Remove the binary |
| `task tidy` | Run `go mod tidy` |
| `task fmt` | Format Go source files |
| `task vet` | Run `go vet` |

## Commands

### `bootstrap <environment>`

Performs the full cluster bootstrap sequence.

```bash
./cli/cluster-bootstrap bootstrap dev
```

**What it does:**

1. Decrypts `secrets.<env>.enc.yaml` using SOPS + age
2. Creates the `argocd` namespace
3. Creates the `repo-ssh-key` Secret with Git SSH credentials
4. Installs ArgoCD via Helm (from `components/argocd/`)
5. Deploys the App of Apps root Application
6. Prints ArgoCD access instructions

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--secrets-file` | `secrets.<env>.enc.yaml` | Path to SOPS-encrypted secrets file |
| `--dry-run` | `false` | Print manifests without applying |
| `--skip-argocd-install` | `false` | Skip the Helm ArgoCD installation |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |
| `--context` | current context | Kubeconfig context to use |
| `--age-key-file` | `SOPS_AGE_KEY_FILE` env | Path to age private key |
| `-v, --verbose` | `false` | Enable verbose output |

### `init`

Interactive setup to create SOPS configuration and encrypted secrets files.

```bash
./cli/cluster-bootstrap init
```

**What it does:**

1. Prompts for SOPS provider (age, AWS KMS, or GCP KMS)
2. Collects the encryption key
3. Generates `.sops.yaml`
4. Interactively collects per-environment secrets (repo URL, target revision, SSH key path)
5. Creates encrypted `secrets.<env>.enc.yaml` files

**Flags:**

| Flag | Description |
|------|-------------|
| `--provider` | SOPS provider: `age`, `aws-kms`, or `gcp-kms` |
| `--age-key-file` | Path to age public key file |
| `--kms-arn` | AWS KMS key ARN |
| `--gcp-kms-key` | GCP KMS key resource ID |
| `--output-dir` | Output directory (default: current directory) |

### `vault-token`

Stores the Vault root token as a Kubernetes Secret.

```bash
./cli/cluster-bootstrap vault-token --token <root-token>
```

**What it does:**

Creates or updates a `vault-root-token` Secret in the `vault` namespace. This is required for non-dev Vault instances after running `vault operator init`.

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--token` | Yes | Vault root token |
| `--kubeconfig` | No | Path to kubeconfig file |
| `--context` | No | Kubeconfig context to use |

## Dependencies

The CLI uses these key libraries:

| Library | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/huh` | Interactive terminal UI |
| `github.com/getsops/sops/v3` | SOPS encryption/decryption |
| `helm.sh/helm/v3` | Helm SDK for chart installation |
| `k8s.io/client-go` | Kubernetes API client |
