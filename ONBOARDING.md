# Notiflex Platform 온보딩 가이드

새로 합류하는 엔지니어를 위한 가이드입니다. 실제 클러스터 상태를 기반으로 작성됩니다.

## 프로젝트 개요

Notiflex는 B2B 알림 SaaS 플랫폼입니다. 기업 고객(테넌트)이 API를 호출하면, 알림을 생성하고 이벤트를 발행합니다. Go 표준 라이브러리로 작성된 단일 바이너리(main.go)가 GKE에서 실행되며, GitOps(ArgoCD)로 배포됩니다.

현재 2개 테넌트가 운영 중입니다.

- **SMB(notiflex 네임스페이스)**: 기본 테넌트. Canary 배포, Valkey 캐시, Kafka 이벤트 발행
- **Enterprise(enterprise 네임스페이스)**: 대형 고객용 격리 테넌트. 별도 Rollout, Secret 관리

## 아키텍처

```
[인터넷]
  ↓
[Gateway API] (notiflex-gateway, 35.216.122.221)
  ↓ HTTPRoute
[Notiflex API] (notiflex namespace, api-pool)
  ├── Valkey (ID 생성, 6379/TCP)
  ├── Kafka (이벤트 발행, 9092/TCP, worker-pool)
  └── Tempo (트레이싱, OTLP 4317/TCP, ops-pool)
```

## 클러스터 구성

```
$ kubectl --context gke-sysnet4admin_book_gitaiops get nodes \
  -o custom-columns=NAME:.metadata.name,POOL:.metadata.labels.cloud\\.google\\.com/gke-nodepool
NAME                                              POOL
gke-notiflex-cluster-api-pool-ffda9966-vv9h       api-pool
gke-notiflex-cluster-default-pool-fdc9c324-kvzo   default-pool
gke-notiflex-cluster-default-pool-fdc9c324-s6bt   default-pool
gke-notiflex-cluster-ops-pool-a023aa47-csz1       ops-pool
gke-notiflex-cluster-worker-pool-2be68d79-v150    worker-pool
```

| 노드풀 | 용도 | 머신 타입 | 워크로드 |
|--------|------|----------|---------|
| default-pool | 플랫폼 인프라 | e2-medium (Spot) × 2 | ArgoCD, Prometheus, Grafana, Loki |
| api-pool | API 서빙 | e2-medium (Spot) × 1 | notiflex-api (SMB, Enterprise) |
| worker-pool | 데이터 처리 | e2-standard-2 (Spot) × 1 | Kafka Broker (KRaft) |
| ops-pool | 운영 도구 | e2-small (Spot) × 1 | Tempo, Health Check CronJob |

모든 노드가 Spot VM입니다. 비용 효율적이지만 선점(preemption)으로 노드가 갑자기 종료될 수 있습니다. Pod이 재스케줄링되므로 서비스 중단은 없지만, 일시적으로 Pod이 Pending 상태가 될 수 있습니다.

## 네임스페이스별 워크로드

| 네임스페이스 | 주요 Pod | 역할 |
|-----------|---------|------|
| notiflex | notiflex-api, valkey-primary, healthcheck-cronjob | API 서빙, 캐시, 헬스체크 |
| enterprise | notiflex-api | Enterprise 테넌트 API |
| kafka | notiflex-kafka-controller | Strimzi 기반 Kafka (KRaft 모드) |
| monitoring | prometheus, grafana, loki, fluent-bit, tempo | 관측 가능성 3요소 (메트릭, 로그, 트레이스) |
| argocd | argocd-server, application-controller, repo-server | GitOps 배포 |
| argo-rollouts | argo-rollouts-controller | Canary/Blue-Green 배포 제어 |

## 저장소 구조

```
notiflex-platform/
├── CLAUDE.md                ← 프로젝트 규칙 (AI가 자동 로드)
├── JOURNEY.md               ← 진행 현황 + 의사결정 기록
├── ONBOARDING.md            ← 이 문서
├── app/                     ← Go 앱 (main.go, Dockerfile)
├── k8s/
│   ├── smb/                 ← SMB 테넌트 (Rollout, Service, CronJob, Gateway)
│   ├── enterprise/          ← Enterprise 테넌트 (Rollout, Service)
│   ├── kafka/               ← Strimzi CRD (Kafka, KafkaNodePool, KafkaTopic)
│   └── monitoring/          ← PrometheusRule (Pod 재시작 알림)
├── argocd/
│   ├── root-app.yaml        ← App of Apps 루트
│   └── apps/                ← 개별 Application (smb, enterprise)
├── helm-values/             ← Helm 차트 values 파일
├── .github/workflows/       ← GitHub Actions CI (WIF 인증)
├── claude-context/          ← AI용 아키텍처 맥락 문서
└── docs/
    └── architecture-decisions.md  ← ADR-001~014 누적
```

## 배포 방법

코드를 변경하고 Git push하면 자동 배포됩니다.

1. `git push origin main`
2. GitHub Actions CI가 이미지를 빌드하고 Artifact Registry에 push
3. CI가 `k8s/smb/rollout.yaml`의 이미지 태그를 업데이트하고 커밋
4. ArgoCD가 Git 변경을 감지하고 Canary 배포를 시작
5. 20% → 50% → 80% → 100% (각 단계 30초 pause)

배포 상태 확인:

```bash
kubectl --context gke-sysnet4admin_book_gitaiops argo rollouts status notiflex-api -n notiflex
```

## 접근 방법

**ArgoCD UI**
```bash
kubectl --context gke-sysnet4admin_book_gitaiops port-forward svc/argocd-server -n argocd 8080:443
# 초기 비밀번호 조회
kubectl --context gke-sysnet4admin_book_gitaiops get secret argocd-initial-admin-secret \
  -n argocd -o jsonpath='{.data.password}' | base64 -d
```
https://localhost:8080 (admin / 위 비밀번호)

**Grafana (메트릭 + 로그 + 트레이스)**
```bash
kubectl --context gke-sysnet4admin_book_gitaiops port-forward svc/kube-prometheus-grafana \
  -n monitoring 3000:80
```
http://localhost:3000 (admin / prom-operator)

- **Prometheus**: 메트릭 조회 (PromQL). 예: `rate(http_requests_total[5m])`
- **Loki**: 로그 조회 (LogQL). 예: `{namespace="notiflex"}`
- **Tempo**: 트레이스 조회. Trace ID로 검색하거나 Service Name으로 필터링

**API 엔드포인트**
```bash
$ kubectl --context gke-sysnet4admin_book_gitaiops get gateway -n notiflex --no-headers
notiflex-gateway   gke-l7-regional-external-managed   35.216.122.221   True

$ curl -s http://35.216.122.221/health
ok

$ curl -s http://35.216.122.221/id
{"id":1,"pod":"notiflex-api-xxx-yyy"}
```

## 자주 묻는 Q&A

**Q: Canary 배포 중 에러율이 올라가면 어떻게 해?**

즉시 abort합니다.
```bash
kubectl --context gke-sysnet4admin_book_gitaiops argo rollouts abort notiflex-api -n notiflex
```
stable 버전으로 즉시 복원됩니다. 원인을 분석하고 수정한 뒤 다시 배포합니다.

**Q: 로그에서 특정 에러를 찾으려면?**

Grafana Explore에서 Loki를 선택하고 LogQL을 사용합니다.
```
{namespace="notiflex"} |= "error"
{namespace="notiflex", container="notiflex-api"} | json | level="error"
```
또는 CLI:
```bash
kubectl --context gke-sysnet4admin_book_gitaiops logs -l app=notiflex-api -n notiflex --tail=100
```

**Q: 요청이 느린데 어디서 병목인지 모르겠어?**

API 응답의 `trace_id`를 Grafana Explore → Tempo에서 검색합니다. API → Valkey → Kafka 각 구간의 소요 시간이 표시됩니다.

**Q: Kafka 토픽을 추가하려면?**

`k8s/kafka/` 디렉터리에 KafkaTopic YAML을 추가하고 Git push합니다.
```yaml
apiVersion: kafka.strimzi.io/v1
kind: KafkaTopic
metadata:
  name: new-topic
  namespace: kafka
  labels:
    strimzi.io/cluster: notiflex-kafka
spec:
  partitions: 3
  replicas: 1
```

**Q: 새 테넌트를 추가하려면?**

Enterprise 패턴을 따릅니다.
1. `k8s/<tenant-name>/` 디렉터리에 Rollout, Service YAML 생성
2. `argocd/apps/<tenant-name>.yaml` Application 추가 (sync-wave: "2", CreateNamespace=true)
3. Git push하면 App of Apps가 자동으로 관리

**Q: 알림은 어떻게 확인해?**

Pod이 5분 내에 3회 이상 재시작되면 PrometheusRule이 감지합니다. 현재 알림 상태:
```bash
kubectl --context gke-sysnet4admin_book_gitaiops port-forward \
  svc/kube-prometheus-kube-prome-alertmanager -n monitoring 9093:9093
```
http://localhost:9093 에서 활성 알림을 확인합니다.
