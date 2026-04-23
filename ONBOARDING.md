# Notiflex 신규 엔지니어 온보딩

생성일: 2026-04-23 (KST)  
기준 클러스터 컨텍스트: `gke-sysnet4admin_book_gitaiops`

## 1) 클러스터 실제 상태

### 노드풀/머신 타입/용도/워크로드

| 노드풀 | 머신 타입 | 노드 수 | 용도(라벨) | 현재 워크로드(요약) |
|---|---|---:|---|---|
| `api-pool` | `e2-medium` | 1 | `role=api` | node-exporter, gmp collector, GKE system daemon |
| `default-pool` | `e2-medium` | 2 | - | ArgoCD, Argo Rollouts, External Secrets, 모니터링 핵심 컴포넌트 |
| `ops-pool` | `e2-small` | 1 | `role=ops` | Tempo, healthcheck CronJob pod, node-exporter |
| `worker-pool` | `e2-standard-2` | 1 | `role=worker` | `notiflex-api`, Kafka(Strimzi/broker) |

검증 명령:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get nodes -o wide --show-labels
```

### 네임스페이스별 Pod 집계

| Namespace | Pod 수 | Running | Completed(Succeeded) | 역할 |
|---|---:|---:|---:|---|
| `argo-rollouts` | 1 | 1 | 0 | Rollout controller |
| `argocd` | 7 | 7 | 0 | GitOps 컨트롤 플레인 |
| `external-secrets` | 3 | 3 | 0 | Secret 동기화 |
| `kafka` | 3 | 3 | 0 | 메시징 |
| `monitoring` | 11 | 11 | 0 | Prometheus/Grafana/Alertmanager/Tempo |
| `notiflex` | 5 | 2 | 3 | API, Valkey, healthcheck job |
| `kube-system` | 47 | 47 | 0 | GKE 시스템 워크로드 |

검증 명령:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get pods -A -o wide
```

### 네임스페이스 역할 테이블

| Namespace | 역할 | 대표 Pod |
|---|---|---|
| `argocd` | Git 저장소 상태를 클러스터에 동기화 | `argocd-application-controller-0`, `argocd-server-*` |
| `argo-rollouts` | Canary 단계 전환 제어 | `argo-rollouts-*` |
| `notiflex` | 비즈니스 API + 캐시 + 헬스체크 잡 | `notiflex-api-*`, `valkey-primary-0` |
| `monitoring` | 메트릭/로그/트레이스/알림 조회 | `prometheus-*`, `kube-prometheus-grafana-*`, `tempo-0`, `alertmanager-*` |
| `kafka` | 비동기 메시징 | `notiflex-kafka-broker-0`, `strimzi-cluster-operator-*` |
| `external-secrets` | GCP Secret Manager 연동 | `external-secrets-*` |

## 2) 저장소 디렉터리 구조

현재 `notiflex-platform` 기준:

```text
.
├── app/               # Go API 소스 + Dockerfile
├── k8s/               # 워크로드 매니페스트(smb, enterprise, monitoring, kafka)
├── argocd/            # ArgoCD Application(App of Apps 포함)
├── helm-values/       # Helm values(모니터링 스택)
├── .github/workflows/ # CI 파이프라인
├── JOURNEY.md         # 진행 이력
└── ONBOARDING.md      # 이 문서
```

온보딩 표준 확장 디렉터리(현재 저장소에는 미생성):

```text
claude-context/      # 프롬프트 컨텍스트/요약
command-guardrails/  # 운영 명령 가드레일
.claude/             # Claude/Codex 도구 설정
```

## 3) 접근 방법

### ArgoCD UI

```bash
# 1) 포트포워딩
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd \
  port-forward svc/argocd-server 8080:443

# 2) 초기 admin 비밀번호 조회
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd \
  get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 --decode; echo
```

접속: `https://localhost:8080`  
계정: `admin`

### Grafana

```bash
# 1) 포트포워딩
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring \
  port-forward svc/kube-prometheus-grafana 3000:80

# 2) admin 비밀번호 조회
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring \
  get secret kube-prometheus-grafana -o jsonpath='{.data.admin-password}' | base64 --decode; echo
```

접속: `http://localhost:3000`  
계정: `admin`

Grafana 데이터소스 사용:

- `Prometheus`(uid: `prometheus`): 메트릭/알람 지표 조회, 예) `rate(http_requests_total[5m])`
- `Loki`(uid: `loki`): 애플리케이션 로그 검색, 예) `{namespace="notiflex"} |= "error"`
- `Tempo`(uid: `tempo`): TraceID 기반 분산추적, service map은 Prometheus와 연결

### API 엔드포인트

Gateway 주소 확인:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops \
  get gateways.gateway.networking.k8s.io -n notiflex
```

현재 주소: `35.216.122.221`

```bash
export GATEWAY_IP=35.216.122.221
curl -s "http://${GATEWAY_IP}/health"
curl -s "http://${GATEWAY_IP}/version"
curl -s "http://${GATEWAY_IP}/id"
```

## 4) 배포 플로우 (Git Push -> Canary)

1. 개발자가 `app/` 변경 후 `main` 브랜치에 push.
2. GitHub Actions(`.github/workflows/ci.yaml`)가 실행.
3. Cloud Build로 이미지 빌드/푸시: `asia-northeast3-docker.pkg.dev/<project>/notiflex/api:sha-<7자리>`.
4. CI가 `k8s/smb/rollout.yaml`의 `image:` 태그를 새 SHA 태그로 자동 갱신 후 커밋/푸시.
5. ArgoCD가 Git 변경을 감지해 `notiflex-smb` 애플리케이션을 자동 Sync.
6. Argo Rollouts가 Canary 단계 적용: `20% -> pause 30s -> 50% -> pause 30s -> 80% -> pause 30s -> promote`.
7. 이상 없으면 stable 서비스(`notiflex-api`)로 최종 전환.

검증 명령:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd get applications.argoproj.io
kubectl --context gke-sysnet4admin_book_gitaiops -n notiflex get rollout notiflex-api
kubectl --context gke-sysnet4admin_book_gitaiops -n notiflex argo rollouts get rollout notiflex-api
```

## 5) 자주 묻는 Q&A

### Q1. Canary 배포를 즉시 중단(Abort)하려면?
A. 아래 명령으로 중단 후 상태를 확인한다.

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n notiflex \
  argo rollouts abort notiflex-api
kubectl --context gke-sysnet4admin_book_gitaiops -n notiflex \
  argo rollouts get rollout notiflex-api
```

### Q2. Loki에서 에러 로그만 빠르게 찾으려면?
A. Grafana Explore에서 데이터소스 `Loki` 선택 후 LogQL 예시를 사용한다.

```logql
{namespace="notiflex", app="notiflex-api"} |= "error"
```

### Q3. Tempo에서 TraceID 추적은 어떻게 하나?
A. API 응답/로그에서 TraceID를 확인한 뒤 Grafana Explore -> `Tempo` 데이터소스에서 TraceID로 조회한다. 이후 Span 상세와 Service Map으로 병목 구간을 본다.

### Q4. Kafka 토픽을 추가하려면?
A. `k8s/kafka/notifications-topic.yaml`을 복제해 새 `KafkaTopic` 리소스를 추가하고 git push하면 ArgoCD가 반영한다.

### Q5. 새 테넌트를 추가하려면?
A. `k8s/enterprise/` 패턴처럼 새 namespace 디렉터리를 만들고 `namespace.yaml`, `service.yaml`, `rollout.yaml`, `service-preview.yaml`를 구성한 뒤 `argocd/` Application에 경로를 등록한다.

### Q6. 알림(Alertmanager) 확인은 어디서 하나?
A. 기본은 Grafana Alerting UI 또는 Alertmanager UI를 사용한다.

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring \
  port-forward svc/kube-prometheus-kube-prome-alertmanager 9093:9093
```

접속: `http://localhost:9093`

---

운영 팁: 이 문서는 클러스터 실측 기반이므로, 클러스터 변경 후 아래 명령으로 수치(노드/파드/Gateway IP)를 먼저 재검증하고 문서를 갱신한다.

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get nodes -o wide --show-labels
kubectl --context gke-sysnet4admin_book_gitaiops get pods -A -o wide
kubectl --context gke-sysnet4admin_book_gitaiops get gateways.gateway.networking.k8s.io -A
```
