# Notiflex 아키텍처 스냅샷 (ch6 완료 시점)

## 3층 지식 구조

- **CLAUDE.md** — AI에게 프로젝트 메타데이터 제공 (매 대화 자동 로드)
- **claude-context/** — 현재 아키텍처 스냅샷 (AI 참조용 한눈보기)
- **docs/architecture-decisions.md** — 결정 ADR 누적 (사람+AI 공동 검토)

## 클러스터 토폴로지

| 항목 | 값 |
|------|-----|
| 클러스터 | notiflex-cluster |
| 리전/존 | asia-northeast3 / asia-northeast3-a |
| 노드풀 | default-pool (e2-medium × 2, Spot) |
| GKE 기능 | Gateway API standard, Workload Identity, Secret Manager CSI |

## 컴포넌트 다이어그램

```
외부 요청
   ↓
GKE Gateway (gke-l7-regional-external-managed)
   ↓ HTTPRoute
notiflex-api Service (ClusterIP)
   ↓ Argo Rollouts Canary (20%→50%→80%→100%)
notiflex-api Pod (v0.1.1, Go)
   ├─→ Valkey (valkey-primary.notiflex:6379)
   └─→ Secret Manager CSI (GKE managed, WI 인증)
```

## 참고 문서

- `JOURNEY.md` — 진행 기록·도구 선택 기록
- `docs/architecture-decisions.md` — ADR-001~007
- `CLAUDE.md` — AI 행동 규칙
