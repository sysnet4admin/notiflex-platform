# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ⬜ | | |
| ch2 | 2.3 gcloud 설정 | ⬜ | | |
| ch2 | 2.4 GitHub 저장소 | ⬜ | | |
| ch2 | 2.5 GKE 클러스터 | ⬜ | | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-29 | notiflex-api v0.1.0 배포 완료 |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-29 | 초기 커밋 및 `/update-docs` 커스텀 스킬 추가 |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-29 | ArgoCD 설치 + private GitHub 저장소 연동 완료 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-29 | `/version` 엔드포인트 추가 및 v0.1.1 롤링 업데이트 완료 |
| ch3 | 3.4 CI | ✅ | 2026-04-29 | GitHub Actions CI 추가 (push to main, app 변경 시 빌드/푸시) |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-29 | CI에서 GitOps 매니페스트 이미지 태그 자동 갱신 + git push로 ArgoCD 자동 배포 연결 완료 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-29 | kube-prometheus-stack(Helm) 설치 + Notiflex Grafana 대시보드 ConfigMap 적용 |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-29 | Loki + Fluent Bit 설치, Grafana Loki 데이터소스 추가, `{job="fluent-bit",namespace="notiflex"}` 로그 조회 확인 |
| ch4 | 4.4 알림 | ✅ | 2026-04-29 | PrometheusRule(`pod-restart-alert`) 생성/적용 완료, Alertmanager 연동 확인 |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-29 | Gateway API(Gateway/HTTPRoute) + HealthCheckPolicy 적용, 외부 IP 35.216.99.80 응답 확인 |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-29 | Argo Rollouts Blue/Green 적용(rollout + preview service), 30초 auto-promote 설정 완료 |
| ch5 | 5.4 ADR 기록 | ✅ | 2026-04-29 | `docs/architecture-decisions.md` 생성, ch3~ch5 결정 ADR-001~007 누적 기록 |
| ch6 | 6.1 캐시 | ✅ | 2026-04-29 | Valkey standalone 설치 + notiflex-api `/id`를 Valkey INCR 기반으로 전환 완료 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-29 | Workload Identity + GKE Secret Manager CSI 활성화, `valkey-password`를 Google Secret Manager로 이관, SecretProviderClass/파일 마운트 패턴 적용 |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-30 | Argo Rollouts 전략을 Blue/Green에서 Canary(20→50→80→100, 30초 pause)로 전환 |
| ch6 | 6.4 아키텍처 컨텍스트 | ✅ | 2026-04-30 | `claude-context/architecture.md` 스냅샷 생성, 3층 지식 구조와 현재 클러스터 토폴로지 반영 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-30 | api/worker/ops 역할별 노드풀 생성 + notiflex-api를 `api-pool`로 스케줄링 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-30 | `argocd/root-app.yaml` 추가, `argocd/apps/` 하위 Application 관리 + sync-wave 적용 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-30 | `k8s/enterprise` 테넌트 분리(rollout/service/secret) + `argocd/apps/notiflex-enterprise.yaml` 추가, cross-namespace Valkey DNS 적용 |
| ch8 | 8.1 메시징 | ✅ | 2026-04-30 | Strimzi 0.51 + Kafka 4.1.0(KRaft) 설치, `notifications` 토픽 생성, notiflex-api Producer/Consumer 연동 완료 |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-30 | Tempo(OTLP gRPC) 설치 + Notiflex OTel SDK 연동 + Grafana Tempo datasource 구성 |
| ch8 | 8.3 CronJob | ✅ | 2026-04-30 | `k8s/smb/healthcheck-cronjob.yaml` 추가, `*/5 * * * *` 헬스체크 자동화 + ops-pool 배치 확인 |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| ch2.6 컨테이너 베이스 이미지 | scratch + 멀티스테이지 빌드 | alpine, distroless | 최소 공격면과 작은 이미지 크기 |
| ch3.2 GitOps 도구 | ArgoCD | Flux | UI/자동 동기화 기반 학습 흐름에 적합 |
| ch3.4 CI 도구 | GitHub Actions + docker build/push (방식 A) | gcloud builds submit (방식 B) | 권한 구성이 단순하고 학습 흐름에 적합 |
| ch3.5 CI-CD 연결 | GitHub Actions에서 GitOps 매니페스트 자동 갱신 후 ArgoCD 자동 동기화 | CI 내 직접 `kubectl apply`, 별도 CD 파이프라인 분리 | Git 단일 소스 오브 트루스 유지 + 변경 이력 추적 용이 |
| ch4.2 메트릭 모니터링 | kube-prometheus-stack (Prometheus + Grafana + Alertmanager) | Datadog, New Relic, VictoriaMetrics | Helm 기반으로 실습 환경에서 빠르게 설치 가능하고 Grafana 대시보드 연계가 쉬움 |
| ch4.3 로그 수집 | Loki + Fluent Bit | ELK Stack, Cloud Logging | e2-medium 리소스 제약에서 경량 운영 가능하고 Grafana Explore와 즉시 연계 가능 |
| ch4.4 알림 | PrometheusRule + Alertmanager | Grafana Alerting, Cloud Monitoring Alert | GitOps(YAML/PR)로 이력 관리가 가능하고 기존 kube-prometheus-stack과 네이티브로 통합됨 |
| ch5.2 트래픽 관리 | GKE Gateway API(Regional External Managed) + HealthCheckPolicy | Ingress Controller(NGINX), Istio Gateway | GKE 네이티브 L7 관리형 경로를 단순하게 구성하고 `/health:8080` 헬스체크를 명시해 무중단 라우팅 안정성 확보 |
| ch5.3 무중단 배포 | Argo Rollouts Blue/Green (`activeService`/`previewService`, `autoPromotionSeconds: 30`) | Deployment RollingUpdate, Flagger+Istio | 트래픽 전환 시점을 명시적으로 제어하고 preview를 분리해 검증 후 자동 승격할 수 있어 운영 리스크를 낮춤 |
| ch6.1 캐시 | Valkey (`bitnami/valkey`, standalone) | Redis OSS, Memcached | Redis API 호환성으로 앱 변경을 최소화하면서도 오픈 거버넌스(Valkey) 기반으로 캐시/카운터 공유를 단순하게 구성 가능 |
| ch6.2 시크릿 관리 | GKE Secret Manager CSI Driver + Workload Identity | Kubernetes Secret, HashiCorp Vault | Secret 값을 클러스터 밖 Google Secret Manager에서 중앙 관리하고, Pod에는 CSI 파일 마운트로만 주입해 노출면을 줄일 수 있음 |
| ch6.3 배포 전략 전환 | Argo Rollouts Canary (`canaryService`/`stableService`, 20→50→80→100, `pause: 30s`) | Blue/Green 유지 | 점진 트래픽 전환으로 새 버전 리스크를 단계별로 관찰하고 이상 시 중간 단계에서 제어하기 용이 |
| ch7.2 노드 스케줄링 | GKE 멀티 노드풀 + `nodeSelector(cloud.google.com/gke-nodepool)` | taint/toleration, nodeAffinity | 학습 환경에서 가장 단순하게 워크로드 역할 분리를 적용할 수 있고, GKE 노드풀 라벨을 그대로 사용해 설정 오류 가능성을 낮춤 |
| ch7.3 멀티 앱 관리 | ArgoCD App of Apps (`root-app` + `argocd/apps/` + `directory.recurse`) | Application 단건 수동 관리, ApplicationSet | 하위 Application 선언을 Git 디렉터리 기준으로 일괄 관리하고 sync-wave(인프라→플랫폼→앱)로 설치 순서를 통제하기 쉬움 |
| ch7.4 멀티테넌시 | Namespace 기반 테넌트 분리(`enterprise`) + ArgoCD Application(`notiflex-enterprise`) | vCluster, 테넌트별 별도 클러스터 | 단일 클러스터 비용을 유지하면서 RBAC/리소스 단위 격리와 GitOps(App of Apps) 운영 패턴을 그대로 확장 가능 |
| ch8.1 메시징 | Strimzi Operator + Kafka 4.1.0(KRaft) + Sarama Producer/Consumer | RabbitMQ, NATS, Pulsar | 이벤트 스트림 중심 구조에서 토픽 기반 확장성이 좋고, Strimzi로 Kubernetes 운영 일관성을 유지하기 쉬움 |
| ch8.2 분산 트레이싱 | Grafana Tempo(monolithic) + OpenTelemetry Go SDK(OTLP gRPC) | Jaeger, Zipkin | Grafana/Loki/Prometheus와 동일 스택으로 통합 운영이 쉽고, 현재 리소스 제약에서 단일 바이너리 Tempo로 최소 비용 트레이싱 구성이 가능 |
| ch8.3 배치 자동화 | Kubernetes CronJob (`notiflex-healthcheck`, `*/5 * * * *`) | Argo Workflows, Airflow, 외부 VM cron | 단순 주기 헬스체크에는 K8s 기본 리소스가 가장 간단하고 GitOps(YAML/ArgoCD) 이력 관리에 직접 연결됨 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-04-29: 초기 버전 설정 |
| Notiflex 이미지 | `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:sha-20260430083408-tempo` | 2026-04-30: ch8.2 OTel SDK/Tempo exporter 코드 반영 이미지를 빌드/푸시하고 SMB/Enterprise Rollout에 반영 |
| ArgoCD | quay.io/argoproj/argocd:v3.3.8 | 2026-04-29: gke-sysnet4admin_book_gitaiops 클러스터에 설치 및 notiflex-platform 저장소 연결 |
| Argo Rollouts | kubectl-argo-rollouts v1.8.4 / controller quay.io/argoproj/argo-rollouts:v1.9.0 | 2026-04-29: `argo-rollouts` namespace 생성 후 CRD/Controller 설치 완료 |
| Valkey | bitnami/valkey chart 5.5.1 (app 9.0.3) | 2026-04-29: `notiflex` namespace에 standalone 배포(`valkey-primary`), Secret `valkey/valkey-password` 사용 |
| Prometheus | quay.io/prometheus/prometheus:v3.11.3 | 2026-04-29: `kube-prometheus-stack-84.3.0`으로 monitoring namespace에 설치 |
| Grafana | docker.io/grafana/grafana:13.0.1 | 2026-04-29: `kube-prometheus-grafana` 배포, Notiflex 대시보드 ConfigMap(`notiflex-grafana-dashboard`) 추가 |
| Loki | docker.io/grafana/loki:3.6.7 | 2026-04-29: `loki-7.0.0`(SingleBinary) 설치, `loki-datasource` ConfigMap으로 Grafana 데이터소스 등록(`isDefault: false`) |
| Tempo | docker.io/grafana/tempo:2.9.0 (chart `tempo-1.24.4`) | 2026-04-30: `monitoring` 네임스페이스에 Helm 설치, `role=ops` 노드 스케줄링 + OTLP gRPC(4317) 활성화 |
| Fluent Bit | cr.fluentbit.io/fluent/fluent-bit:5.0.3 | 2026-04-30: 클러스터 DaemonSet 기준 실제 실행 이미지를 `cr.fluentbit.io/fluent/fluent-bit:5.0.3`으로 확인 |
| GKE Secret Manager CSI | `secretManagerConfig.enabled=true`, `csi-secrets-store-gke`/`csi-secrets-store-provider-gke` DaemonSet | 2026-04-29: Workload Identity 활성화 후 Secret Manager addon 활성화, `notiflex-secrets` SecretProviderClass로 `valkey-password` 파일 마운트 구성 |
| Kafka | 4.1.0 (Strimzi 0.51.0, KRaft, single broker) | 2026-04-30: `kafka` namespace에 Strimzi operator 설치, `notiflex-kafka` 클러스터 및 `notifications` 토픽 생성 |
| OTel SDK | 1.43.0 | 2026-04-30: `go.opentelemetry.io/otel` + `otlptracegrpc` 추가, `/health` `/id` `/version` span 계측 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium (Spot) | 2 | argocd, monitoring(kube-prometheus-stack + loki + fluent-bit), valkey |
| api-pool | e2-medium (Spot) | 1 | notiflex-api (Rollout nodeSelector) |
| worker-pool | e2-standard-2 (Spot) | 1 | strimzi-cluster-operator, notiflex-kafka-kafka-0, notiflex-kafka-entity-operator |
| ops-pool | e2-small (Spot) | 1 | tempo(statefulset), notiflex-healthcheck CronJob/Job |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.6 | Artifact Registry 리포지토리 미존재 | notiflex 리포지토리를 생성한 뒤 이미지 푸시 |
| ch3.2 | ArgoCD Application이 `Sync Unknown` (`Repository not found`) | repo Secret에 GitHub 토큰(`forceHttpBasicAuth: "true"`) 등록 후 ArgoCD 전체 rollout restart |
| ch3.3 | ArgoCD가 최신 커밋(`a08e25a`)을 즉시 반영하지 않음 | `kubectl --context gke-sysnet4admin_book_gitaiops -n argocd annotate application notiflex-smb argocd.argoproj.io/refresh=hard --overwrite`로 강제 재동기화 |
| ch3.5 | GitHub Actions `if` 조건에서 `secrets.*` 직접 참조 시 워크플로 문법 오류 | 인증 관련 시크릿을 `env`로 옮기고 `if`는 `env` 기반으로 분기 |
| ch3.5 | CI 인증 방식이 환경마다 달라 빌드 실패 (SA Key/WIF 시크릿 키 이름 불일치) | SA Key, 레거시 WIF, GCP WIF 3가지 입력 조합을 모두 지원하도록 워크플로 보완 |
| ch4.3 | Loki 설치 시 기본 cache(chunks/results) Pod가 CPU 부족으로 Pending되어 Helm install timeout | `helm-values/loki.yaml`에서 `chunksCache.enabled=false`, `resultsCache.enabled=false`로 비활성화 후 `helm upgrade --install` 재실행 |
| ch4.4 | `kube_pod_container_status_restarts_total` 기반 경보는 Pod 삭제만으로 즉시 증가하지 않아 Alert가 `firing`되지 않을 수 있음 | 규칙 로드/Alertmanager 연동 확인 후, 실제 재시작 카운트가 증가하는 장애 시나리오(예: CrashLoop)로 추가 검증 |
| ch7.3 | `root-app` 적용 직후 `Sync Unknown (argocd/apps path does not exist)` | 로컬 App of Apps 변경(`argocd/apps/`)을 GitHub `main`에 커밋/푸시한 뒤 ArgoCD 재동기화하면 `Synced`로 전환 |
| ch8.1 | Strimzi `KafkaNodePool`에 `spec.template.pod.nodeSelector`를 넣으면 strict decoding 에러 발생 | `spec.template.pod.affinity.nodeAffinity`로 worker-pool 고정 방식 변경 후 재적용 |
| ch8.2 | Tempo를 `nodeSelector.role=ops`로 설치했지만 `role` 노드 라벨이 없어 Pending 발생 가능 | `kubectl label node <ops-node> role=ops --overwrite`로 라벨을 추가한 뒤 Helm 배포하면 ops-pool에 정상 스케줄링 |
