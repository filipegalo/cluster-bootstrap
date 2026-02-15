# Kube Prometheus Stack

**Namespace:** `monitoring` | **Sync Wave:** 3 | **Chart:** `kube-prometheus-stack` v82.0.0

Deploys the Prometheus monitoring stack for cluster observability. Configured for metrics collection and remote write — Grafana and Alertmanager are disabled.

## Upstream Chart

- **Chart:** `kube-prometheus-stack`
- **Version:** 82.0.0
- **Repository:** `https://prometheus-community.github.io/helm-charts`

## Key Configuration

**Base values:**

- CRD installation disabled (provided by [Prometheus Operator CRDs](prometheus-operator-crds.md))
- Grafana: disabled
- Alertmanager: disabled
- Node Exporter: enabled
- Kube State Metrics: enabled
- Scrape interval: 30s
- Evaluation interval: 30s
- Remote write endpoint configured
- All control plane component monitoring enabled (kubelet, API server, etcd, scheduler, controller-manager, CoreDNS, kube-proxy)

**Per-environment resources:**

| Setting | Dev | Staging | Prod |
|---------|-----|---------|------|
| Prometheus replicas | 1 | 1 | 2 |
| Retention | 2h | 4h | 6h |
| CPU request | 100m | 200m | 500m |
| Memory request | 256Mi | 512Mi | 1Gi |
| CPU limit | 200m | 500m | 1000m |
| Memory limit | 512Mi | 1Gi | 2Gi |

## Sync Configuration

The ArgoCD Application ignores diff on status fields for Prometheus and Alertmanager CRs, which are frequently updated by the operator and would otherwise cause perpetual out-of-sync states.

## Dependencies

- **Prometheus Operator CRDs** (wave 2) must be installed before this component (wave 3)

## Files

```
components/kube-prometheus-stack/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```
