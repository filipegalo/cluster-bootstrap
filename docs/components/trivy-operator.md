# Trivy Operator

**Namespace:** `trivy-system` | **Sync Wave:** 3 | **Chart:** `trivy-operator` v0.32.0

Aqua Trivy Operator provides automated vulnerability scanning for container images and Kubernetes resources.

## Upstream Chart

- **Chart:** `trivy-operator`
- **Version:** 0.32.0
- **Repository:** `https://aquasecurity.github.io/helm-charts/`

## Key Configuration

**Base values:**

- Ignore unfixed vulnerabilities: enabled
- Report TTL: 24h
- Concurrent scan jobs: 3

**Per-environment resources:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| Concurrent scan jobs | 2 | 3 | 5 |
| CPU request | 10m | 50m | 100m |
| Memory request | 32Mi | 64Mi | 256Mi |
| CPU limit | 100m | 100m | 250m |
| Memory limit | 128Mi | 256Mi | 512Mi |

## What It Does

Trivy Operator automatically scans:

- Container images in running workloads for known vulnerabilities (CVEs)
- Kubernetes resource configurations for misconfigurations

Scan results are stored as custom resources (`VulnerabilityReport`, `ConfigAuditReport`) that can be queried with kubectl:

```bash
kubectl get vulnerabilityreports -A
kubectl get configauditreports -A
```

## Files

```
components/trivy-operator/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
