# Notiflex 아키텍처 스냅샷 — ch7 완료 시점

이 문서는 AI가 매 대화에서 현재 아키텍처를 빠르게 파악할 수 있도록 한 페이지로 요약한다.

## 3층 지식 구조

| 문서 | 역할 | 업데이트 주기 |
|------|------|------------|
| **CLAUDE.md** | 프로젝트 메타데이터 (GCP 프로젝트, 리전, Artifact Registry) | 초기 설정 시 |
| **claude-context/** | 현재 아키텍처 스냅샷 (토폴로지·컴포넌트·파이프라인) | 챕터 완료 시 |
| **docs/architecture-decisions.md** | 결정 누적 기록 (왜 이 도구를 선택했는가) | 결정 시점마다 |

CLAUDE.md는 매 대화 자동 로드, claude-context는 AI가 참조 요청 시, ADR은 사람·AI가 결정 근거 검토 시 사용한다.

## 클러스터 토폴로지

| 항목 | 값 |
|------|-----|
| 클러스터 | notiflex-cluster |
| 리전/존 | asia-northeast3 / asia-northeast3-a |
| 노드풀 | default-pool(e2-medium×2), api-pool(e2-medium×1), worker-pool(e2-standard-2×1), ops-pool(e2-small×1) |
| GKE 기능 | Gateway API, Workload Identity, Secret Manager CSI |

## 컴포넌트 다이어그램

```
외부 클라이언트
    │
    ▼ HTTP :80
GKE Gateway (35.216.99.80)
  gke-l7-regional-external-managed
    │ HTTPRoute
    ▼
Service: notiflex-api (stable)
Service: notiflex-api-preview (canary)
    │
    ▼
Rollout: notiflex-api (Canary 전략)
  ├── stable ReplicaSet (notiflex/api:v0.2.1)
  └── canary ReplicaSet (배포 시 생성)
    │
    ├── Secret 볼륨 (CSI → Secret Manager)
    │     SecretProviderClass: notiflex-secrets
    │     valkey-password: /mnt/secrets/valkey-password
    │
    └── Valkey StatefulSet (valkey-primary)
          INCR 명령으로 분산 ID 생성
```

## 배포 파이프라인

```
코드 변경 (app/)
    │
    ▼ git push → main
GitHub Actions CI
    │ docker build + push
    ▼
Artifact Registry
  asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api
    │ 매니페스트 이미지 태그 업데이트 → git push
    ▼
ArgoCD (notiflex-smb Application)
  auto-sync + selfHeal
    │
    ▼
Argo Rollouts (Canary)
  20% → 30s → 50% → 30s → 80% → 30s → 100%
```

## 관측 가능성

| 도구 | 역할 | 네임스페이스 |
|------|------|------------|
| Prometheus | 메트릭 수집 (kube-prometheus-stack) | monitoring |
| Grafana | 대시보드·데이터소스 통합 | monitoring |
| Loki | 로그 저장 (SingleBinary, filesystem) | monitoring |
| Fluent Bit | 로그 수집 DaemonSet → Loki | monitoring |
| Alertmanager | 알림 라우팅 (PrometheusRule 연동) | monitoring |

## 주요 네임스페이스

| 네임스페이스 | 주요 워크로드 |
|------------|-------------|
| notiflex | Rollout/notiflex-api(api-pool), valkey-primary(default), notiflex-sa (WI) |
| enterprise | Rollout/notiflex-api(api-pool), notiflex-sa (WI) — enterprise 테넌트 |
| argocd | ArgoCD v3.3.8 (root-app + notiflex-smb + notiflex-enterprise) |
| argo-rollouts | Argo Rollouts 컨트롤러 |
| monitoring | kube-prometheus-stack, Loki, Fluent Bit |
