# Notiflex Platform

이 저장소는 [AI 시대에 개발자가 알아야 하는 인프라 구성 배포 with 클로드 코드]() 책의 실습 결과물입니다.

Claude Code와 함께 GKE 위에서 SMB → Enterprise 수준의 배포 파이프라인을 구축하고, 그 과정에서 자연스럽게 살아있는 운영 표준(GitAIOps)을 만들어간 결과를 담고 있습니다.

- 실습 가이드 저장소: https://github.com/sysnet4admin/_Book_GitAIOps

---

## 저장소 구조

```
notiflex-platform/
├── CLAUDE.md                ← 프로젝트 규칙 (Claude Code 자동 로드)
├── JOURNEY.md               ← 진행 현황 + 도구 선택 기록
├── ONBOARDING.md            ← 신규 팀원을 위한 온보딩 가이드 (부록A)
├── app/                     ← Go 앱 (main.go, Dockerfile)
├── k8s/
│   ├── smb/                 ← SMB 테넌트 (Rollout, Service, CronJob, Gateway)
│   ├── enterprise/          ← Enterprise 테넌트 (Rollout, Service)
│   ├── kafka/               ← Strimzi CRD (Kafka, KafkaNodePool, KafkaTopic)
│   └── monitoring/          ← PrometheusRule (알림 규칙)
├── argocd/
│   ├── root-app.yaml        ← App of Apps 루트 (ch7)
│   └── apps/                ← 개별 Application (smb, enterprise)
├── helm-values/             ← Helm 차트 values 파일
├── .github/workflows/       ← GitHub Actions CI (WIF 인증)
├── claude-context/          ← AI용 아키텍처 맥락 문서
└── docs/
    └── architecture-decisions.md  ← ADR-001~014
```

---

## 무엇을 만들었나

책의 흐름을 따라 SMB 스타트업에서 Enterprise 플랫폼으로 성장하는 과정을 그대로 구현했습니다.

| 챕터 | 구현 내용 |
|:---:|---------|
| ch2 | GKE 클러스터 + Artifact Registry + 첫 배포 |
| ch3 | ArgoCD GitOps + GitHub Actions CI (WIF) + E2E 자동화 |
| ch4 | Prometheus + Grafana + Loki + Fluent Bit + PrometheusRule |
| ch5 | GKE Gateway API + Argo Rollouts Blue/Green + ADR-001~007 |
| ch6 | Valkey 캐시 + GKE CSI Secret Manager + Canary 배포 전환 |
| ch7 | 역할별 노드풀 4개 + App of Apps + 멀티테넌시(enterprise) |
| ch8 | Kafka(Strimzi KRaft) + Tempo + 헬스체크 CronJob |
| ch9 | 저장소 분석 + 회고 + 온보딩 문서 + GitAIOps |

---

## 아키텍처

```
[인터넷]
  ↓
[GKE Gateway API]  (gke-l7-regional-external-managed)
  ↓ HTTPRoute
[Notiflex API]  (Argo Rollouts Canary, api-pool)
  ├── Valkey  (분산 ID 카운터)
  ├── Kafka   (비동기 이벤트 발행, worker-pool)
  └── Tempo   (분산 트레이싱, ops-pool)

[ArgoCD root-app]  →  notiflex-smb  +  notiflex-enterprise
[GitHub Actions CI]  →  WIF 인증  →  Artifact Registry
```

클러스터는 역할에 따라 4개 노드풀로 분리됩니다.

| 노드풀 | 머신 타입 | 주요 워크로드 |
|--------|----------|-------------|
| default-pool | e2-medium × 2 | ArgoCD, Prometheus, Grafana, Loki |
| api-pool | e2-medium × 1 | notiflex-api (SMB + Enterprise) |
| worker-pool | e2-standard-2 × 1 | Kafka (Strimzi KRaft 4.1.0) |
| ops-pool | e2-small × 1 | Tempo, CronJob |

---

## 살아있는 운영 표준

이 저장소는 단순한 실습 코드가 아닙니다. 책을 진행하는 동안 만들어진 세 가지 층의 지식 구조가 있습니다.

- **JOURNEY.md** — 진행 기록과 도구 선택 이유. AI가 다음 대화에서 참조합니다
- **docs/architecture-decisions.md** — ADR-001~014. 왜 이 기술을 선택했는지 기록
- **claude-context/architecture.md** — AI가 현재 아키텍처를 빠르게 파악하는 스냅샷

새 대화를 시작할 때 AI는 이 세 파일을 읽고 지금까지의 결정을 이어받아 일관성 있게 동작합니다. 이것이 GitAIOps입니다.

---

## 저자
- ✔️ [조 훈](https://github.com/sysnet4admin)
