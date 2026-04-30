# Notiflex Platform 온보딩 가이드

## 빠른 시작

### 필수 도구

| 도구 | 버전 | 역할 |
|------|------|------|
| gcloud CLI | 529+ | GCP 인증 및 클러스터 접근 |
| kubectl | 1.35+ | K8s 클러스터 운영 |
| helm | 3.x | Helm 차트 관리 |
| gh | 2.x | GitHub 저장소 관리 |

### 클러스터 접근

```bash
gcloud container clusters get-credentials notiflex-cluster \
  --zone=asia-northeast3-a \
  --project=project-75fce205-dfa5-4975-a56
kubectl config rename-context \
  $(kubectl config current-context) \
  gke-sysnet4admin_book_gitaiops
kubectl --context gke-sysnet4admin_book_gitaiops get nodes
```

### 주요 리소스 확인

```bash
# ArgoCD Application 상태
kubectl --context gke-sysnet4admin_book_gitaiops get app -n argocd

# 전체 Pod 상태
kubectl --context gke-sysnet4admin_book_gitaiops get pods -A

# Gateway IP (외부 접근)
kubectl --context gke-sysnet4admin_book_gitaiops get gateway -n notiflex

# API 테스트
curl http://35.216.99.80/health
curl http://35.216.99.80/id
```

## 아키텍처 요약

`claude-context/architecture.md` 참조 — ch9 완료 시점 스냅샷.

핵심 흐름:
1. 코드 변경 → GitHub → ArgoCD → Argo Rollouts (Canary)
2. /id 호출 → Valkey INCR → Kafka 이벤트 → OTel 트레이스 → Tempo
3. PrometheusRule → Alertmanager 알림

## 의사결정 이력

`docs/architecture-decisions.md` — ADR-001~016 16개 결정 기록

| 챕터 | ADR 범위 | 핵심 결정 |
|------|----------|----------|
| ch3 | 001~002 | ArgoCD, GitHub Actions |
| ch4 | 003~005 | Prometheus+Grafana, Loki+Fluent Bit, PrometheusRule |
| ch5 | 006~007 | Gateway API, Blue/Green |
| ch6 | 008~010 | Valkey, Secret Manager CSI, Canary |
| ch7 | 011~013 | 노드풀 분리, App of Apps, 멀티테넌시 |
| ch8 | 014~016 | Kafka, Tempo, CronJob |

## 트러블슈팅

자주 발생하는 문제:
- ArgoCD Sync 실패 → `argocd.argoproj.io/refresh=hard` 어노테이션
- Valkey 연결 실패 → `kubectl get pod valkey-primary-0 -n notiflex`
- CSI Secret 마운트 실패 → WI 바인딩 확인
- Kafka entity-operator CrashLoop → userOperator 섹션 제거
