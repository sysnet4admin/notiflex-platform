# Notiflex Architecture Snapshot

> 7장 완료 시점 아키텍처 스냅샷
> 업데이트 기준: 2026-04-30 (KST)
> 클러스터 조회 컨텍스트: `gke-sysnet4admin_book_gitaiops`

## 1) 3층 지식 구조

이 저장소의 운영 지식은 세 층으로 분리한다. `CLAUDE.md`는 프로젝트 메타데이터와 실행 규칙(대화 시작 시 참조되는 기본 맥락)을 제공하고, `claude-context/`는 지금 클러스터가 어떻게 구성되어 동작하는지에 대한 현재 스냅샷을 제공한다. `docs/architecture-decisions.md`는 왜 그런 결정을 내렸는지(대안/트레이드오프 포함)를 누적 기록한다. 즉, 규칙/메모리(`CLAUDE.md`)와 현재 상태(`claude-context`)와 의사결정 이력(ADR)을 분리해 혼선을 줄인다.

## 2) 클러스터 토폴로지

| 항목 | 현재 상태 |
|---|---|
| 클러스터 | `notiflex-cluster` |
| Kubernetes 버전 | `v1.35.1-gke.1396002` (노드 기준) |
| 리전/존 | `asia-northeast3` / `asia-northeast3-a` |
| 노드 수 | 5 |
| 노드풀 | `default-pool`(2), `api-pool`(1), `worker-pool`(1), `ops-pool`(1) |
| 노드 프로비저닝 | `cloud.google.com/gke-provisioning=spot` |
| 외부 트래픽 진입 | Gateway API (`gke-l7-regional-external-managed`) |
| GitOps | ArgoCD App of Apps (`root-app` + `notiflex-smb` + `notiflex-enterprise`) |
| 배포 컨트롤러 | Argo Rollouts (`argo-rollouts` namespace) |
| 시크릿 연동 | Secrets Store CSI Driver (`secrets-store-gke.csi.k8s.io`) + `SecretProviderClass:notiflex-secrets` |
| 캐시 | Valkey StatefulSet (`notiflex/valkey-primary`) |
| 모니터링 스택 | Prometheus + Grafana + Loki + Fluent Bit (`monitoring` namespace) |
| GKE 관리 기능(관찰됨) | Gateway Managed LB, GMP 관련 namespace(`gmp-system`), CSI 드라이버 DaemonSet |

### 노드풀 상세

| 노드 | 존 | 노드풀 | Spot | 인스턴스 타입 |
|---|---|---|---|---|
| `gke-notiflex-cluster-api-pool-6f140622-cwt8` | `asia-northeast3-a` | `api-pool` | `true` | `e2-medium` |
| `gke-notiflex-cluster-default-pool-783b2e9a-13pv` | `asia-northeast3-a` | `default-pool` | `true` | `e2-medium` |
| `gke-notiflex-cluster-default-pool-783b2e9a-uof4` | `asia-northeast3-a` | `default-pool` | `true` | `e2-medium` |
| `gke-notiflex-cluster-ops-pool-5bdb76fd-c956` | `asia-northeast3-a` | `ops-pool` | `true` | `e2-small` |
| `gke-notiflex-cluster-worker-pool-93827d65-wxjm` | `asia-northeast3-a` | `worker-pool` | `true` | `e2-standard-2` |

## 3) 컴포넌트 다이어그램

```text
[Internet Client]
      |
      v
[Gateway: notiflex-gateway (35.216.99.80)]
      |
      v
[HTTPRoute: notiflex-route]
      |
      v
[Service: notiflex-api (stable)] <-------------------------------+
      |                                                         |
      | (Argo Rollouts canary steps: 20 -> 50 -> 80 -> 100)    |
      +--> [Service: notiflex-api-preview (canary)]             |
                     |                                           |
                     v                                           |
             [Rollout: notiflex-api, replicas=2]                |
                     |                                           |
                     v                                           |
             [Pods: notiflex-api-*]                             |
                 |                     \                         |
                 | VALKEY_ADDR          \ secrets mount         |
                 v                       \                       |
       [Valkey: valkey-primary:6379]     [CSI Volume /mnt/secrets]
                                                 |
                                                 v
             [SecretProviderClass: notiflex-secrets]
                                                 |
                                                 v
 [Google Secret Manager: projects/.../secrets/valkey-password/versions/latest]
```

## 4) 배포 파이프라인

1. 개발자가 `main` 브랜치로 push (`app/**` 변경 시 CI 트리거).
2. GitHub Actions (`.github/workflows/ci.yaml`)가 GCP 인증(SA Key 또는 WIF) 후 Docker 빌드/푸시.
3. 이미지가 Artifact Registry에 업로드됨: `asia-northeast3-docker.pkg.dev/<PROJECT_ID>/notiflex/api:sha-<7자리>`.
4. 같은 워크플로가 `k8s/smb/rollout.yaml`의 `image:`를 새 SHA 태그로 갱신하고 커밋/푸시.
5. ArgoCD `root-app`이 `argocd/apps/` 하위 Application들을 동기화하고, `notiflex-smb`/`notiflex-enterprise`가 각 경로를 배포한다.
6. Argo Rollouts가 `notiflex-api`를 Canary 전략으로 배포:
   - `setWeight: 20` -> `pause 30s`
   - `setWeight: 50` -> `pause 30s`
   - `setWeight: 80` -> `pause 30s`
   - 최종 100% 전환

### 현재 배포 설정 핵심값

| 키 | 값 |
|---|---|
| Rollout 이름 | `notiflex-api` |
| replicas | `2` |
| stable service | `notiflex-api` |
| canary service | `notiflex-api-preview` |
| 현재 이미지 | `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:sha-5a5e0c6` |
| probe | `/health` (readiness/liveness) |

## 5) 관측 가능성

| 도구 | 배치 위치 | 역할 |
|---|---|---|
| Prometheus | `monitoring` (`prometheus-kube-prometheus-kube-prome-prometheus`) | 메트릭 수집/저장 |
| Grafana | `monitoring` (`kube-prometheus-grafana`) | 대시보드/시각화 |
| Loki | `monitoring` (`loki`, `loki-gateway`) | 로그 저장/조회 |
| Fluent Bit | `monitoring` DaemonSet (`fluent-bit`) | 노드/파드 로그 수집 후 Loki 전달 |
| Alertmanager | `monitoring` (`alertmanager-kube-prometheus-kube-prome-alertmanager`) | 알림 라우팅 |
| Argo Rollouts Metrics | `argo-rollouts` svc (`argo-rollouts-metrics`) | 배포 관련 지표 노출 |

현재 클러스터 조회 결과 기준으로 Tempo 워크로드는 배포되어 있지 않다(Tracing은 향후/별도 단계 관리).

## 6) 주요 네임스페이스

| Namespace | 주요 워크로드 | 역할 |
|---|---|---|
| `notiflex` | `Rollout/notiflex-api`, `Service/notiflex-api`, `Service/notiflex-api-preview`, `StatefulSet/valkey-primary`, `Gateway/notiflex-gateway`, `HTTPRoute/notiflex-route`, `SecretProviderClass/notiflex-secrets` | SMB 테넌트 애플리케이션 본체, 캐시, 트래픽 진입점 |
| `enterprise` | `Rollout/notiflex-api`, `Service/notiflex-api`, `Secret/notiflex-api-secret` | Enterprise 테넌트 워크로드 분리 네임스페이스 |
| `argocd` | `Application/root-app`, `Application/notiflex-smb`, `Application/notiflex-enterprise`, `argocd-server`, `argocd-repo-server`, `argocd-application-controller` | App of Apps 기반 GitOps 동기화/선언적 배포 |
| `argo-rollouts` | `deployment/argo-rollouts` | Canary/Blue-Green 배포 컨트롤러 |
| `monitoring` | `prometheus-*`, `grafana`, `loki*`, `fluent-bit`, `alertmanager-*` | 메트릭/로그/알림 |
| `kube-system` | `csi-secrets-store-gke`, `csi-secrets-store-provider-gke`, `kube-dns`, `metrics-server` | 클러스터 핵심 애드온 및 CSI 드라이버 |
| `gmp-system` | `gmp-operator`, `rule-evaluator` | GKE Managed Prometheus 관련 시스템 워크로드 |

## 참고 조회 명령 (재생성용)

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get app -n argocd \
  -o custom-columns=NAME:.metadata.name,SYNC:.status.sync.status,HEALTH:.status.health.status

kubectl --context gke-sysnet4admin_book_gitaiops get ns
kubectl --context gke-sysnet4admin_book_gitaiops get rollout,deploy,sts -A
kubectl --context gke-sysnet4admin_book_gitaiops get gateway,httproute -A
kubectl --context gke-sysnet4admin_book_gitaiops get ds -A
kubectl --context gke-sysnet4admin_book_gitaiops get svc -A
```
