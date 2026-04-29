# Notiflex 아키텍처 스냅샷 (ch9 완료 시점)

## 3층 지식 구조

- **CLAUDE.md** — AI에게 프로젝트 메타데이터 제공 (매 대화 자동 로드)
- **claude-context/** — 현재 아키텍처 스냅샷 (AI 참조용 한눈보기). 이 파일이 그 역할
- **docs/architecture-decisions.md** — 결정 ADR 누적 (ADR-001~014, 사람+AI 공동 검토)

## 클러스터 토폴로지

| 항목 | 값 |
|------|-----|
| 클러스터 | notiflex-cluster |
| 리전/존 | asia-northeast3 / asia-northeast3-a |
| 노드풀 | default-pool (e2-medium × 2), api-pool (e2-medium × 1), worker-pool (e2-standard-2 × 1), ops-pool (e2-small × 1) |
| GKE 기능 | Gateway API standard, Workload Identity, Secret Manager CSI |

## 컴포넌트 다이어그램

```
외부 요청
   ↓
GKE Gateway (gke-l7-regional-external-managed)
   ↓ HTTPRoute
notiflex-api Service (ClusterIP)
   ↓ Argo Rollouts Canary (20%→50%→80%→100%)
notiflex-api Pod (v0.1.1, Go) [api-pool]
   ├─→ Valkey (valkey-primary.notiflex:6379) — ID 카운터
   ├─→ Secret Manager CSI (GKE managed, WI 인증) — valkey-password
   └─→ Kafka (notiflex-kafka-kafka-bootstrap.kafka:9092) — notifications 토픽 [worker-pool]

enterprise/notiflex-api Pod [api-pool] — 멀티테넌시
CronJob notiflex-healthcheck [ops-pool] — 5분마다 헬스체크
```

## 배포 파이프라인

```
app/ 변경 → git push
   ↓ GitHub Actions CI (WIF 인증)
Artifact Registry (sha-태그 이미지)
   ↓ rollout.yaml 자동 업데이트
ArgoCD auto-sync via root-app (selfHeal: true)
   └─ Apps: notiflex-smb, notiflex-enterprise
   ↓ Argo Rollouts
Canary 점진 배포 (30초 간격)
```

## 관측 가능성

| 도구 | 역할 | 위치 |
|------|------|------|
| Prometheus | 메트릭 수집 (5m CPU) | monitoring ns |
| Grafana | 대시보드 | monitoring ns |
| Alertmanager | 알림 라우팅 | monitoring ns |
| Loki + Fluent Bit | 로그 수집 | monitoring ns |
| Tempo | 분산 트레이싱 (OTLP gRPC 4317) | monitoring ns / ops-pool |
| PrometheusRule | Pod 재시작 알림 | monitoring/pod-restart-alert |
| CronJob | 5분마다 API 헬스체크 | notiflex ns / ops-pool |

## 주요 네임스페이스

| 네임스페이스 | 주요 워크로드 |
|-----------|------------|
| notiflex | notiflex-api (Rollout/Canary), valkey-primary, notiflex-healthcheck (CronJob) |
| enterprise | notiflex-api (Rollout/Canary) — 멀티테넌시 |
| argocd | ArgoCD 전체 스택 (root-app → notiflex-smb, notiflex-enterprise) |
| argo-rollouts | Argo Rollouts 컨트롤러 |
| kafka | notiflex-kafka (KRaft 4.1.0, Strimzi 0.51.0) [worker-pool] |
| monitoring | Prometheus, Grafana, Alertmanager, Loki, Fluent Bit, Tempo [ops-pool] |
| kube-system | CSI secrets-store-gke DaemonSet |

## 참고 문서

- `JOURNEY.md` — 진행 기록·도구 선택 기록 (ch2~ch9 전체)
- `docs/architecture-decisions.md` — ADR-001~014 누적
- `CLAUDE.md` — AI 행동 규칙
