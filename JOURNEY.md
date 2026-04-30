# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-30 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-30 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-30 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-30 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-30 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-30 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-30 | |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-30 | |
| ch3 | 3.4 CI | ✅ | 2026-04-30 | |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-30 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-30 | |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-30 | |
| ch4 | 4.4 알림 | ✅ | 2026-04-30 | |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-30 | |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-30 | |
| ch6 | 6.1 캐시 | ✅ | 2026-04-30 | Valkey INCR 연동 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-30 | GKE Secret Manager CSI+WI |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-30 | Blue/Green -> Canary |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-30 | api/worker/ops 노드풀 생성 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-30 | Root-app, App of Apps 패턴 적용 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-30 | enterprise 네임스페이스 분리 |
| ch8 | 8.1 메시징 | ✅ | 2026-04-30 | Kafka 4.1.0 KRaft(worker-pool) |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-30 | Tempo OTel(ops-pool) |
| ch8 | 8.3 CronJob | ✅ | 2026-04-30 | 헬스체크(ops-pool) |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-30 | 저장소/커밋/클러스터 분석 완료 |
| ch9 | 9.2 회고 | ✅ | 2026-04-30 | 도구 선택 패턴/리소스 현황 회고 |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-30 | ONBOARDING.md 생성 |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-30 | Git+AI+Ops 루프 분석 완료 |
| ch9 | 9.5 마무리 | ✅ | 2026-04-30 | 다음 단계 제안(보안/스케일링/비용) |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| ch3.2 GitOps | ArgoCD | Flux | Kubernetes 네이티브, 선언적, UI 제공 |
| ch3.4 CI | GitHub Actions | Cloud Build, GitLab CI, Jenkins | GitHub 내장, YAML 선언적, 무료 크레딧 |
| ch4.2 메트릭 모니터링 | Prometheus + Grafana | Datadog, CloudWatch, Google Cloud Monitoring | 오픈소스 표준, 무료, Grafana 시각화 |
| ch4.3 로그 수집 | Loki + Fluent Bit | ELK Stack, CloudWatch Logs | 경량, Grafana 통합, 라벨 기반 인덱싱 |
| ch4.4 알림 | PrometheusRule + Alertmanager | Grafana UI 알림, Cloud Monitoring Alert | GitOps 호환, kube-prometheus-stack 통합 |
| ch5.2 외부 트래픽 관리 | Gateway API | Ingress NGINX, Istio, Traefik | Kubernetes 차세대 표준, GKE 네이티브 |
| ch5.3 무중단 배포 | Argo Rollouts | Flagger, K8s Rolling Update | ArgoCD 통합, CRD 기반, 점진적 전략 진화 |
| ch6.1 캐시 | Valkey | Redis OSS, Memcached | Redis 호환, BSD 라이선스, Bitnami 차트 |
| ch6.2 시크릿 관리 | GKE Secret Manager CSI + WI | Sealed Secrets, HashiCorp Vault | GKE 네이티브, SA 키 불필요 |
| ch6.3 배포 전략 | Argo Rollouts Canary | Blue/Green 유지 | 점진 트래픽 전환, 리스크 단계별 관찰 |
| ch7.2 노드 스케줄링 | GKE 멀티 노드풀 + nodeSelector | taint/toleration, nodeAffinity | 단순, GKE 라벨 자동 부여 |
| ch7.3 멀티 앱 관리 | ArgoCD App of Apps | ApplicationSet | 디렉터리 기반 일괄 관리, sync-wave |
| ch7.4 멀티테넌시 | Namespace 분리 + RBAC | vCluster, 클러스터 분리 | 단일 클러스터 비용 유지, GitOps 확장 |
| ch8.1 메시징 | Strimzi + Kafka 4.1.0 KRaft | RabbitMQ, NATS | 업계 표준, CRD 기반, ZooKeeper 불필요 |
| ch8.2 트레이싱 | Grafana Tempo + OTel | Jaeger, Zipkin | Grafana 스택 통합, 단일 바이너리 |
| ch8.3 배치 | Kubernetes CronJob | Argo Workflows, VM cron | K8s 기본, GitOps 이력 관리 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-04-30: 초기 버전 |
| Notiflex 이미지 | asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:v0.1.1 | ch3.3: /version 엔드포인트 추가 |
| ArgoCD | quay.io/argoproj/argocd:v3.3.8 | ch3.2: 설치 완료 |
| Argo Rollouts | v1.8.4 | ch5.3: Blue/Green, ch6.3: Canary 전환 |
| kube-prometheus-stack | 84.x (Prometheus v3.x + Grafana 13.x) | ch4.2: 설치 완료 |
| Loki | 7.0.0 (SingleBinary) | ch4.3: 설치 완료 |
| Tempo | grafana/tempo (ops-pool) | ch8.2: OTel 트레이싱 설정 |
| Kafka | 4.1.0 KRaft (Strimzi Operator) | ch8.1: worker-pool 배치 |
| OTel SDK | go.opentelemetry.io/otel (Go SDK) | ch8.2: OTLP gRPC exporter |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | ArgoCD, monitoring(Prometheus/Grafana/Loki) |
| api-pool | e2-medium | 1 | notiflex-api (smb/enterprise) |
| worker-pool | e2-standard-2 | 1 | Kafka(Strimzi) |
| ops-pool | e2-small | 1 | Tempo, CronJob |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.6 | Go build failed due to newline in string literal | Fixed the string literal in `main.go` |
