# Adding a Component

This guide walks through adding a new platform component to the stack.

## 1. Create the component chart

Create a new directory under `components/`:

```
components/my-component/
├── Chart.yaml
└── values/
    ├── base.yaml
    ├── dev.yaml
    ├── staging.yaml
    └── prod.yaml
```

### Chart.yaml

Declare the upstream Helm chart as a dependency:

```yaml
apiVersion: v2
name: my-component
description: My new component
version: 0.1.0
type: application
dependencies:
  - name: upstream-chart-name
    version: "1.2.3"
    repository: https://charts.example.com
```

### Values files

Create `values/base.yaml` with shared configuration:

```yaml
upstream-chart-name:
  someKey: someValue
```

!!! note
    Values must be nested under the dependency name (e.g., `upstream-chart-name:`). This is how Helm routes values to subchart dependencies.

Create environment-specific overrides in `values/dev.yaml`, `values/staging.yaml`, and `values/prod.yaml`:

```yaml
upstream-chart-name:
  replicas: 1
  resources:
    requests:
      cpu: 10m
      memory: 32Mi
```

## 2. Create the ArgoCD Application template

Add a new template in `apps/templates/my-component.yaml`:

```yaml
{{- if .Values.myComponent.enabled }}
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-component
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "3"
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: {{ .Values.repo.url }}
    targetRevision: {{ .Values.repo.targetRevision }}
    path: components/my-component
    helm:
      valueFiles:
        - values/base.yaml
        - values/{{ .Values.environment }}.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: my-namespace
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
{{- end }}
```

Choose the sync wave based on dependencies:

| Wave | When to use |
|------|-------------|
| 0 | Core infrastructure (ArgoCD) |
| 1 | Secrets infrastructure (Vault, External Secrets) |
| 2 | CRDs, credentials, utilities |
| 3 | Application-level components |

## 3. Enable in environment values

Add the component toggle to each environment file in `apps/values/`:

```yaml
# apps/values/dev.yaml
myComponent:
  enabled: true
```

Repeat for `staging.yaml` and `prod.yaml`.

## 4. Verify with helm template

Before pushing, verify the Application template renders correctly:

```bash
helm template apps/ -f apps/values/dev.yaml
```

Check that:

- The Application CR is generated with correct metadata
- The sync wave is set appropriately
- The values file paths are correct
- The namespace matches your component's expectation

## 5. Push and sync

Commit and push. ArgoCD will detect the new Application in the App of Apps and deploy it.

```bash
git add components/my-component/ apps/templates/my-component.yaml apps/values/
git commit -m "feat: add my-component"
git push
```
