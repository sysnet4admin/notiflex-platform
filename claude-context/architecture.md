# Notiflex Platform Architecture (7장 완료 시점)

## 1. 3층 지식 구조 (3-Tier Knowledge Structure)

이 프로젝트는 명확한 3단계 지식 계층을 따라 문서화 및 의사결정을 수행합니다. 이는 기술 선택의 근거, 실행 방법, 그리고 결과물에 대한 명확한 추적을 가능하게 합니다.

-   **Decision Guides (`decision-guides/`):** 특정 기술이나 아키텍처 패턴을 왜 선택했는지에 대한 근거를 기록합니다. (e.g., 왜 ArgoCD를 GitOps 도구로 선택했는가?)
-   **Prompt Guardrails (`prompt-guardrails/`):** 결정된 사항을 실제로 구현하기 위한 구체적인 단계와 명령어를 제공합니다. AI 에이전트가 이 가이드라인을 따라 작업을 수행합니다.
-   **Result Templates (`result-templates/`):** 작업 완료 후, 기대하는 결과 상태를 정의하고 검증하기 위한 템플릿입니다.

## 2. 클러스터 토폴로지 (Cluster Topology)

GKE 클러스터는 워크로드 특성에 따라 3개의 노드풀로 구성됩니다. 이를 통해 리소스를 격리하고 비용을 최적화합니다. (ADR-010)

| Node Pool     | Machine Type (예시) | 주요 워크로드                                   | 특징                                     |
| :------------ | :------------------ | :---------------------------------------------- | :--------------------------------------- |
| **api-pool**  | e2-standard-2       | `notiflex-smb`, `notiflex-enterprise` API 서버  | 사용자 트래픽 처리, CPU 중심, 수평적 확장 |
| **worker-pool** | n2-standard-4       | Kafka 메시지 처리, 배치 작업 등 (추후 확장용)    | 백그라운드 작업, Memory/IO 중심          |
| **ops-pool**  | e2-medium           | ArgoCD, Prometheus, Grafana, Loki 등 운영 도구 | 안정적인 운영을 위한 격리된 환경         |

## 3. 컴포넌트 다이어그램 (Component Diagram)

'App of Apps' 패턴과 멀티테넌시 구조가 적용된 컴포넌트 다이어그램입니다.

```
+------------------------------------------------------------------------------------------------------------------+
|  GitHub Repository (notiflex-platform)                                                                           |
|                                                                                                                  |
| +---------------------------+   +----------------------+    +-------------------------------------------------+  |
| |    Application Code (Go)  |-->|  GitHub Actions (CI) |--> |  ArgoCD Application Manifests (/argocd/apps/)   |  |
| +---------------------------+   | - Build, Push Image  |    |   - notiflex-smb.yaml                           |  |
| +---------------------------+   +----------------------+    |   - notiflex-enterprise.yaml                    |  |
| |   k8s Manifests (/k8s/)   |                               |   - notiflex-monitoring.yaml                    |  |
| | - /smb, /enterprise       |                               +-------------------------------------------------+  |
| +---------------------------+                                                                                    |
+----------------------------------------------------------------|-------------------------------------------------+
                                                                 |
                                                                 v
+----------------------------------------------------------------|-------------------------------------------------+
| GKE Cluster (notiflex-cluster)                                 |                                                 |
|                                                                | (Root App)                                      |
| +------------------------------------------------------------+ | +-----------------------------------------------+ |
| | ArgoCD                                                     | | | App of Apps                                   | |
| | - Syncs and manages applications based on git manifests    | | | - notiflex-smb Application                    | |
| +------------------------------------------------------------+ | | - notiflex-enterprise Application             | |
|                                                                | | - notiflex-monitoring Application             | |
|   +--------------------------------------------------------+   | +-----------------------------------------------+ |
|   | Namespace: notiflex-smb (on api-pool)                  |   |                                                 |
|   | +-----------------+    +----------------------------+  |   |                                                 |
|   | | Notiflex SMB    |--> | Valkey                     |  |   |                                                 |
|   | +-----------------+    +----------------------------+  |   |                                                 |
|   +--------------------------------------------------------+   |                                                 |
|   +--------------------------------------------------------+   |                                                 |
|   | Namespace: notiflex-enterprise (on api-pool)           |   |                                                 |
|   | +---------------------+  +--------------------------+  |   |                                                 |
|   | | Notiflex Enterprise |->| Valkey (or other DB)     |  |   |                                                 |
|   | +---------------------+  +--------------------------+  |   |                                                 |
|   +--------------------------------------------------------+   |                                                 |
|   +--------------------------------------------------------+   |                                                 |
|   | Namespace: monitoring (on ops-pool)                    |   |                                                 |
|   | + Prometheus, Grafana, Loki...                         |   |                                                 |
|   +--------------------------------------------------------+   |                                                 |
+------------------------------------------------------------------------------------------------------------------+
```

## 4. 배포 파이프라인 (Deployment Pipeline)

1.  **Commit & Push:** 개발자가 `notiflex-platform` 저장소의 `main` 브랜치에 코드를 커밋하고 푸시합니다.
2.  **CI (GitHub Actions):** 푸시된 코드는 GitHub Actions 워크플로우를 트리거합니다.
    -   Go 애플리케이션을 빌드하고 테스트합니다.
    -   Docker 이미지를 빌드하여 GCR(Google Container Registry)에 푸시합니다.
    -   (필요시) `k8s/` 디렉터리의 Kubernetes 매니페스트를 업데이트합니다.
3.  **CD (ArgoCD):**
    -   **App of Apps Pattern (ADR-011):** 최상위 `root-app`이 `argocd/apps` 디렉터리를 감시합니다.
    -   이 디렉터리에 정의된 `Application` 리소스(e.g., `notiflex-smb`, `notiflex-enterprise`)가 변경되면, ArgoCD는 해당 자식 애플리케이션들을 자동으로 동기화합니다.
    -   각 자식 `Application`은 `k8s/` 아래의 각 테넌트별 매니페스트를 참조하여 클러스터에 배포합니다.
    -   **Argo Rollouts**를 사용하여 점진적 배포(Canary)를 수행하여 안정성을 확보합니다.
4.  **Deployment:** ArgoCD에 의해 각 애플리케이션이 지정된 네임스페이스(`notiflex-smb`, `notiflex-enterprise`)와 노드풀에 배포됩니다.

## 5. 관측 가능성 (Observability)

`monitoring` 네임스페이스와 `ops-pool` 노드풀을 중심으로 관측 가능성 스택이 구축되어 있습니다.

| Tool         | Role                  | Description                                            |
| :----------- | :-------------------- | :----------------------------------------------------- |
| **Prometheus** | Metrics Monitoring    | 클러스터와 애플리케이션의 시계열 메트릭을 수집, 저장, 쿼리합니다. |
| **Grafana**    | Visualization & Dashboards | Prometheus와 Loki의 데이터를 시각화하여 대시보드를 제공합니다.  |
| **Loki**       | Log Aggregation       | 클러스터의 모든 파드에서 생성되는 로그를 수집하고 저장합니다.    |
| **Fluent-bit** | Log Shipper           | 각 노드에서 DaemonSet으로 실행되며, 로그를 Loki로 전송합니다.    |
| **Alertmanager**| Alerting              | Prometheus에서 정의된 임계치 기반의 경고를 관리하고 라우팅합니다.|

## 6. 주요 네임스페이스 (Key Namespaces)

| Namespace              | Description                                                               |
| :--------------------- | :------------------------------------------------------------------------ |
| `argocd`               | GitOps를 위한 ArgoCD 관련 컴포넌트들이 배포된 공간입니다.                  |
| `argo-rollouts`        | 점진적 배포 전략(Canary)을 관리하는 Argo Rollouts 컨트롤러가 위치합니다. |
| `monitoring`           | Prometheus, Grafana, Loki 등 관측 가능성 스택이 배포된 공간입니다.        |
| `notiflex-smb`         | SMB 고객용 Notiflex 애플리케이션과 관련 리소스가 배포됩니다. (ADR-012)     |
| `notiflex-enterprise`  | Enterprise 고객용 Notiflex 애플리케이션이 배포됩니다. (ADR-012)           |
| `kube-system`          | Kubernetes 시스템 자체의 핵심 컴포넌트(DNS, Proxy 등)가 실행되는 공간입니다.|
| `gmp-system`           | Google Cloud Managed Service for Prometheus 관련 리소스가 위치합니다.     |