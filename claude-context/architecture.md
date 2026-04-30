# Notiflex Platform Architecture

## 1. 3층 지식 구조 (3-Tier Knowledge Structure)

이 프로젝트는 명확한 3단계 지식 계층을 따라 문서화 및 의사결정을 수행합니다. 이는 기술 선택의 근거, 실행 방법, 그리고 결과물에 대한 명확한 추적을 가능하게 합니다.

- **Decision Guides (`decision-guides/`):** 특정 기술이나 아키텍처 패턴을 왜 선택했는지에 대한 근거를 기록합니다. (e.g., 왜 ArgoCD를 GitOps 도구로 선택했는가?)
- **Prompt Guardrails (`prompt-guardrails/`):** 결정된 사항을 실제로 구현하기 위한 구체적인 단계와 명령어를 제공합니다. AI 에이전트가 이 가이드라인을 따라 작업을 수행합니다.
- **Result Templates (`result-templates/`):** 작업 완료 후, 기대하는 결과 상태를 정의하고 검증하기 위한 템플릿입니다.

## 2. 클러스터 토폴로지 (Cluster Topology)

현재 GKE(Google Kubernetes Engine) 클러스터는 2개의 노드로 구성되어 있으며, `us-central1` 리전에 배포되어 있습니다.

| Node Name                                           | Status | Roles  | Age | Version               | Internal IP   | External IP  | OS Image                           |
| --------------------------------------------------- | ------ | ------ | --- | --------------------- | ------------- | ------------ | ---------------------------------- |
| gke-notiflex-cluster-default-pool-b9c12066-tnxe | Ready  | <none> | 25m | v1.35.1-gke.1396002   | 10.178.15.228 | 34.50.24.85  | Container-Optimized OS from Google |
| gke-notiflex-cluster-default-pool-b9c12066-tvg7 | Ready  | <none> | 22m | v1.35.1-gke.1396002   | 10.178.15.229 | 34.50.52.202 | Container-Optimized OS from Google |

## 3. 컴포넌트 다이어그램 (Component Diagram)

주요 컴포넌트 간의 상호작용은 다음과 같이 요약할 수 있습니다.

```
+------------------------------------------------------------------------------------------------+
|  GitHub Repository (notiflex-platform)                                                         |
|    +--------------------------+      +---------------------------+     +---------------------+   |
|    | Application Code (Go)    |----->| GitHub Actions (CI)       |---->| ArgoCD (GitOps)     |   |
|    | k8s Manifests (YAML)     |      | - Build Docker Image      |     | - Sync Manifests    |   |
|    +--------------------------+      | - Push to GCR             |     | - Deploy to Cluster |   |
|                                      +---------------------------+     +---------------------+   |
+-----------------------------------------|--------------------------------------------------------+
                                          |
                                          v
+-----------------------------------------|--------------------------------------------------------+
|  GKE Cluster (notiflex-cluster)         |                                                        |
|                                         |                                                        |
|  +---------------------------+          v  +---------------------------+                         |
|  | Observability (monitoring)|<--------+  | Core Application (notiflex)|                         |
|  | - Prometheus (Metrics)    |             | - Notiflex API (Go)       |                         |
|  | - Grafana (Dashboard)     |             | - Valkey (Cache)          |                         |
|  | - Loki (Logging)          |             +---------------------------+                         |
|  | - Fluent-bit (Log Agent)  |                                                                  |
|  +---------------------------+                                                                  |
+------------------------------------------------------------------------------------------------+

```

## 4. 배포 파이프라인 (Deployment Pipeline)

1.  **Commit & Push:** 개발자가 `notiflex-platform` 저장소의 `main` 브랜치에 코드를 커밋하고 푸시합니다.
2.  **CI (GitHub Actions):** 푸시된 코드는 GitHub Actions 워크플로우를 트리거합니다.
    - Go 애플리케이션을 빌드하고 테스트합니다.
    - Docker 이미지를 빌드하여 GCR(Google Container Registry)에 푸시합니다.
    - (필요시) Kubernetes 매니페스트를 업데이트합니다.
3.  **CD (ArgoCD):** ArgoCD가 `notiflex-platform` 저장소의 `k8s/` 디렉터리를 지속적으로 감시합니다.
    - 매니페스트 변경이 감지되면, ArgoCD는 자동으로 클러스터의 상태를 Git 저장소의 상태와 동기화합니다.
    - **Argo Rollouts**를 사용하여 점진적 배포(Blue/Green, Canary)를 수행하여 안정성을 확보합니다.
4.  **Deployment:** ArgoCD에 의해 애플리케이션(Notiflex API)이 `notiflex` 네임스페이스에 배포됩니다.

## 5. 관측 가능성 (Observability)

`monitoring` 네임스페이스를 중심으로 관측 가능성 스택이 구축되어 있습니다.

| Tool         | Role                  | Description                                            |
|--------------|-----------------------|--------------------------------------------------------|
| **Prometheus** | Metrics Monitoring    | 클러스터와 애플리케이션의 시계열 메트릭을 수집, 저장, 쿼리합니다. |
| **Grafana**    | Visualization & Dashboards | Prometheus와 Loki의 데이터를 시각화하여 대시보드를 제공합니다.  |
| **Loki**       | Log Aggregation       | 클러스터의 모든 파드에서 생성되는 로그를 수집하고 저장합니다.    |
| **Fluent-bit** | Log Shipper           | 각 노드에서 DaemonSet으로 실행되며, 로그를 Loki로 전송합니다.    |
| **Alertmanager**| Alerting              | Prometheus에서 정의된 임계치 기반의 경고를 관리하고 라우팅합니다.|

## 6. 주요 네임스페이스 (Key Namespaces)

| Namespace       | Description                                                               |
|-----------------|---------------------------------------------------------------------------|
| `argocd`        | GitOps를 위한 ArgoCD 관련 컴포넌트들이 배포된 공간입니다.                  |
| `argo-rollouts` | 점진적 배포 전략(Blue/Green, Canary)을 관리하는 Argo Rollouts 컨트롤러가 위치합니다. |
| `monitoring`    | Prometheus, Grafana, Loki 등 관측 가능성 스택이 배포된 공간입니다.        |
| `notiflex`      | 핵심 애플리케이션인 Notiflex API와 관련 리소스(Valkey 등)가 배포됩니다.    |
| `kube-system`   | Kubernetes 시스템 자체의 핵심 컴포넌트(DNS, Proxy 등)가 실행되는 공간입니다.|
| `gmp-system`    | Google Cloud Managed Service for Prometheus 관련 리소스가 위치합니다.     |

