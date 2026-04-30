# Notiflex Platform 온보딩 가이드

최종 업데이트: 2026-04-30 (Kubernetes 실측 기준)

## 1) 먼저 알아야 할 현재 상태

### 클러스터 노드풀/용도

| 노드풀 | 머신 타입 | 노드 수 | 용도(실측 워크로드) |
|---|---|---:|---|
| `default-pool` | `e2-medium` | 2 | ArgoCD, Prometheus/Loki 일부, Valkey |
| `api-pool` | `e2-medium` | 1 | `notiflex`/`enterprise` API Rollout Pod |
| `worker-pool` | `e2-standard-2` | 1 | Strimzi/Kafka, Grafana |
| `ops-pool` | `e2-small` | 1 | Tempo, CronJob Job Pod |

검증 명령:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get nodes -L cloud.google.com/gke-nodepool,node.kubernetes.io/instance-type,role
```

### 네임스페이스별 Pod 수 (실측)

| Namespace | Pod 수 |
|---|---:|
| `argo-rollouts` | 1 |
| `argocd` | 7 |
| `enterprise` | 1 |
| `kafka` | 3 |
| `monitoring` | 18 |
| `notiflex` | 6 |
| `kube-system` | 57 |

검증 명령:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get pods -A --no-headers \
| awk '{count[$1]++} END {for (ns in count) printf "%s %d\n", ns, count[ns]}' \
| sort
```

### 핵심 워크로드 역할 매핑

| 역할 | 주요 리소스 | Namespace | 주 배치 노드풀 |
|---|---|---|---|
| GitOps 제어면 | ArgoCD, Argo Rollouts | `argocd`, `argo-rollouts` | `default-pool` |
| API 서비스(SMB) | `notiflex-api` Rollout(2 replicas) | `notiflex` | `api-pool` |
| API 서비스(Enterprise) | `notiflex-api` Rollout(1 replica) | `enterprise` | `api-pool` |
| 캐시/카운터 | `valkey-primary` | `notiflex` | `default-pool` |
| 메시징 | Strimzi, Kafka broker | `kafka` | `worker-pool` |
| 관측(메트릭/로그/트레이싱) | Prometheus, Loki, Grafana, Tempo | `monitoring` | `default/worker/ops` |
| 운영 배치 | `notiflex-healthcheck` CronJob | `notiflex` | `ops-pool` |

## 2) 저장소 구조

```text
.
├── app/                 # Go API 소스, Dockerfile
├── k8s/                 # 배포 매니페스트 (smb, enterprise, monitoring, kafka)
├── argocd/              # App of Apps(root-app + apps/)
├── helm-values/         # Helm 커스텀 값
├── docs/                # ADR 등 운영 문서
├── claude-context/      # 현재 아키텍처 스냅샷
├── .claude/commands/    # 자동 문서 업데이트 명령(/update-docs)
├── JOURNEY.md           # 실제 진행 이력/버전/트러블슈팅 원본
└── ONBOARDING.md        # 이 문서
```

참고:
- `command-guardrails/` 디렉터리는 현재 저장소에 없고, 명령 가드레일은 `.claude/commands/`에서 운영한다.

## 3) 접근 방법

### ArgoCD UI

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd port-forward svc/argocd-server 8080:443
```

- URL: `https://localhost:8080`
- ID: `admin`
- 초기 비밀번호:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | base64 -d && echo
```

### Grafana

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring port-forward svc/kube-prometheus-grafana 3000:80
```

- URL: `http://localhost:3000`
- ID: `admin`
- 비밀번호:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring get secret kube-prometheus-grafana \
  -o jsonpath='{.data.admin-password}' | base64 -d && echo
```

데이터소스 사용:
- Prometheus: 지표 조회/알림 검증
- Loki: 로그 검색(LogQL)
- Tempo: Trace 조회(서비스명 `notiflex-api`)

### API 엔드포인트 (Gateway)

현재 Gateway 주소: `35.216.99.80`

```bash
curl -s http://35.216.99.80/health
curl -s http://35.216.99.80/id
curl -s http://35.216.99.80/version
```

## 4) 배포 플로우 (Git push -> Canary)

1. 개발자가 `main`에 `app/**` 변경을 push
2. GitHub Actions(`.github/workflows/ci.yaml`)가 이미지 빌드/푸시
3. CI가 `k8s/smb/rollout.yaml`의 이미지 태그를 `sha-<commit>`로 갱신 후 다시 commit/push
4. ArgoCD가 Git 변경 감지 후 Sync
5. Argo Rollouts Canary 진행:
   - `20% -> pause 30s`
   - `50% -> pause 30s`
   - `80% -> pause 30s`
   - 이후 100% 전환

상태 확인:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get applications.argoproj.io -n argocd
kubectl --context gke-sysnet4admin_book_gitaiops get rollout -n notiflex
kubectl --context gke-sysnet4admin_book_gitaiops get rollout -n enterprise
```

## 5) 주요 아키텍처 결정 요약

- GitOps: ArgoCD + App of Apps
- 배포 전략: Argo Rollouts Canary
- 트래픽 진입: GKE Gateway API
- 캐시/카운터: Valkey
- 시크릿: GKE Secret Manager CSI + Workload Identity
- 관측: Prometheus/Grafana + Loki + Tempo
- 비동기: Strimzi 기반 Kafka

상세 근거는 `docs/architecture-decisions.md`와 `JOURNEY.md`의 도구 선택 기록을 기준으로 본다.

## 6) 트러블슈팅 가이드 (요약)

- ArgoCD `Repository not found`:
  - repo credential Secret(GitHub token, `forceHttpBasicAuth`) 확인 후 `argocd` rollout restart
- Loki 설치 timeout(Pending cache pod):
  - `helm-values/loki.yaml`에서 cache 비활성화 후 재배포
- Kafka NodePool strict decoding:
  - `nodeSelector` 대신 `nodeAffinity` 사용
- Tempo Pending:
  - ops 노드에 `role=ops` 라벨 존재 확인

## 7) 자주 묻는 질문 (FAQ)

### Q1. Canary를 즉시 중단(Abort)하려면?
```bash
kubectl --context gke-sysnet4admin_book_gitaiops argo rollouts abort notiflex-api -n notiflex
```

### Q2. Loki에서 API 로그를 어떻게 검색하나?
Grafana Explore -> Loki에서:
```logql
{namespace="notiflex", pod=~"notiflex-api-.*"} |= "generated_by"
```

### Q3. Tempo에서 TraceID로 추적하려면?
Grafana Explore -> Tempo에서 TraceID 검색 또는 서비스 `notiflex-api` + 시간 범위로 조회 후 span drill-down.

### Q4. Kafka 토픽을 추가하려면?
`k8s/kafka/topic-<name>.yaml`을 추가하고:
```bash
kubectl --context gke-sysnet4admin_book_gitaiops apply -f k8s/kafka/topic-<name>.yaml
```

### Q5. 새 테넌트를 추가하려면?
`k8s/enterprise` 패턴으로 새 네임스페이스/rollout/service를 만들고, `argocd/apps/`에 테넌트 Application을 추가한 뒤 root-app 동기화로 반영한다.

### Q6. Alertmanager 알림 상태는 어디서 보나?
```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring port-forward svc/kube-prometheus-kube-prome-alertmanager 9093:9093
```
브라우저 `http://localhost:9093`에서 firing/resolved 상태를 확인한다.
