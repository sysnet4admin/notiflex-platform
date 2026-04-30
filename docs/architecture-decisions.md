# Architecture Decision Records

## ADR-014: 메시징 — Kafka (ch8.1)
**시점**: 2026-04 / **결정**: Strimzi Kafka (KRaft 모드, v4.1.0) 채택 (vs RabbitMQ, NATS, Pulsar)
**이유**:
- 고처리량 + 순서 보장 — 알림 이벤트의 파티션별 순서 유지
- Strimzi로 K8s 네이티브 관리 — KafkaNodePool CRD로 브로커 사양 선언적 관리
- KRaft 모드 — ZooKeeper 없이 단일 구성요소로 운영 단순화
- worker-pool 전용 배치 — e2-standard-2로 브로커 성능 격리

## ADR-015: 분산 트레이싱 — Tempo (ch8.2)
**시점**: 2026-04 / **결정**: Tempo (grafana/tempo 단일 바이너리) 채택 (vs Jaeger, Zipkin)
**이유**:
- Grafana 통합 — 메트릭(Prometheus)·로그(Loki)·트레이스(Tempo)를 같은 UI에서 확인
- OTLP gRPC 수신 — OTel SDK와 표준 프로토콜로 연동, 벤더 종속성 없음
- 단일 바이너리 모드 — ops-pool 한 노드에서 경량 운영
- OTel SDK (Go) — 표준 계측으로 향후 다른 백엔드 전환 가능

## ADR-016: 배치 자동화 — K8s CronJob (ch8.3)
**시점**: 2026-04 / **결정**: K8s CronJob 채택 (vs 외부 cron, Argo Workflows)
**이유**:
- 쿠버네티스 네이티브 — 별도 스케줄러 없이 클러스터 내장 CronJob 활용
- ops-pool 배치 — 배치 워크로드를 운영 전용 노드에 격리
- ArgoCD가 매니페스트로 관리 — git에서 스케줄 변경 시 ArgoCD가 자동 반영
- Job 히스토리 보존 — successfulJobsHistoryLimit/failedJobsHistoryLimit으로 실행 이력 추적

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**시점**: 2026-04 / **결정**: ArgoCD v3.3.8 채택 (vs Flux, Jenkins X)
**이유**:
- K8s 네이티브 CRD(Application, AppProject) — 선언적 배포 상태 관리
- Web UI 제공 — 배포 상태·히스토리·diff를 시각적으로 확인
- App of Apps 패턴 지원 — 7장 멀티앱 관리에서 자연스럽게 확장
- automated sync + selfHeal — git이 단일 진실 소스(Single Source of Truth) 역할

## ADR-002: CI 도구 — GitHub Actions (ch3.4)
**시점**: 2026-04 / **결정**: GitHub Actions 채택 (vs Jenkins, GitLab CI)
**이유**:
- GitHub 저장소와 통합 — 별도 서버 운영 불필요
- Workload Identity Federation 지원 — SA 키 없이 GCP 인증
- push 트리거로 이미지 빌드 → 매니페스트 업데이트 → ArgoCD 자동 배포 흐름 완성
- 무료 티어로 학습 환경 구성

## ADR-003: 메트릭 수집 — Prometheus + Grafana (ch4.2)
**시점**: 2026-04 / **결정**: kube-prometheus-stack 채택 (vs Datadog, New Relic)
**이유**:
- 오픈소스 — 라이선스 비용 없음
- ServiceMonitor/PrometheusRule CRD — K8s 네이티브 운영 규칙 관리
- Grafana와 통합 — Loki(로그)·Tempo(트레이스) 데이터소스를 같은 UI에서 확인
- kube-prometheus-stack 차트로 Operator+Prometheus+Grafana+Alertmanager 일괄 설치

## ADR-004: 로그 수집 — Loki + Fluent Bit (ch4.3)
**시점**: 2026-04 / **결정**: Loki + Fluent Bit 채택 (vs ELK, Datadog Logs)
**이유**:
- Grafana 통합 — 메트릭과 로그를 같은 대시보드에서 확인
- 인덱스 없는 레이블 기반 — Elasticsearch 대비 저사양 환경에 유리
- Fluent Bit 경량 DaemonSet — 노드당 최소 자원으로 로그 수집
- SingleBinary 모드로 단일 Pod 운영 — 학습 환경에 적합

## ADR-005: 알림 — PrometheusRule (ch4.4)
**시점**: 2026-04 / **결정**: PrometheusRule CRD 채택 (vs Grafana Alert)
**이유**:
- Prometheus 네이티브 — PromQL 표현식으로 정밀한 조건 정의
- git으로 버전 관리 — 알림 규칙 변경 이력 추적
- kube-prometheus-stack과 자동 연동 — `labels.release` 매칭으로 즉시 활성화
- Alertmanager 라우팅과 결합 — 심각도별 수신 채널 분리

## ADR-006: 외부 트래픽 — Gateway API (ch5.2)
**시점**: 2026-04 / **결정**: GKE Gateway API(gke-l7-regional-external-managed) 채택 (vs Ingress, NGINX)
**이유**:
- GKE 네이티브 — 별도 Ingress Controller 설치 없이 Google 관리형 L7 LB 사용
- HTTPRoute CRD — 트래픽 분할·헤더 기반 라우팅 등 세밀한 제어
- HealthCheckPolicy — `/health` 경로·포트를 직접 지정하여 정확한 헬스체크
- K8s 표준 API — 벤더 종속성 최소화

## ADR-011: 노드 스케줄링 — nodeSelector (ch7.2)
**시점**: 2026-04 / **결정**: GKE 자동 라벨 `cloud.google.com/gke-nodepool` 기반 nodeSelector 채택 (vs nodeAffinity, Taint/Toleration)
**이유**:
- GKE가 노드풀 이름을 자동으로 라벨로 부여 — 커스텀 라벨 관리 불필요
- 단순한 YAML 표현 — `cloud.google.com/gke-nodepool: api-pool` 한 줄로 배치 결정
- api-pool/worker-pool/ops-pool 역할 분리 — 워크로드 특성에 맞는 머신 타입 배치
- Spot VM 비용 절감 — 역할별 노드풀에 적합한 VM 크기 선택

## ADR-012: 멀티앱 관리 — App of Apps (ch7.3)
**시점**: 2026-04 / **결정**: Argo CD App of Apps 패턴 채택 (vs 개별 Application 수동 관리, ApplicationSet)
**이유**:
- root-app이 argocd/apps/ 디렉터리를 감시 — 새 Application 파일 추가만으로 자동 배포
- Sync Wave 어노테이션 — 인프라(0)→플랫폼(1)→앱(2) 순서 보장
- git이 Application 목록의 단일 진실 소스 — kubectl apply 없이 PR만으로 앱 추가/삭제
- ArgoCD selfHeal — git 상태와 클러스터 상태 자동 동기화

## ADR-013: 멀티테넌시 — Namespace 분리 (ch7.4)
**시점**: 2026-04 / **결정**: Namespace 분리 + per-tenant Rollout 채택 (vs 단일 namespace + 라벨 격리, vCluster)
**이유**:
- 강한 격리 — RBAC, NetworkPolicy, ResourceQuota를 네임스페이스 단위로 독립 적용
- ArgoCD App of Apps와 자연 결합 — 테넌트별 Application을 apps/ 디렉터리에 추가만 하면 됨
- 테넌트별 독립 배포 — smb/enterprise가 각자의 Rollout으로 독립적인 배포 이력 유지
- 운영 가시성 — kubectl get pods -n enterprise로 테넌트별 상태 즉시 확인

## ADR-008: 캐시 — Valkey (ch6.1)
**시점**: 2026-04 / **결정**: Valkey (Bitnami standalone) 채택 (vs Redis, Memcached)
**이유**:
- Redis fork — API 완전 호환, 라이선스 자유(BSD-3)
- bitnami Helm 차트 지원 — 단순 설치, Secret 자동 생성
- INCR 명령으로 분산 ID 카운터 구현
- e2-medium 환경에서 50m CPU로 구동 가능

## ADR-009: 시크릿 관리 — GKE Secret Manager CSI + WI (ch6.2)
**시점**: 2026-04 / **결정**: GKE managed CSI + Workload Identity 채택 (vs K8s Secret, Vault)
**이유**:
- 서비스 계정 키 불필요 — Workload Identity로 GCP IAM 직접 연동
- 파일 마운트 패턴 — 환경변수 노출 없이 /mnt/secrets/로 안전하게 전달
- Secret Manager 버전 관리 — 비밀번호 교체 시 새 버전 추가
- GKE managed addon — 오픈소스 CSI 설치 없이 클러스터 업데이트만으로 활성화

## ADR-010: 배포 전략 전환 — Canary (ch6.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Canary로 전환 (Blue/Green에서)
**이유**:
- 트래픽 점진 이동(20%→50%→80%→100%) — 운영 위험 최소화
- canaryService/stableService 분리 — 문제 발생 시 즉각 rollback
- pause 단계 — 각 단계에서 모니터링 후 다음 단계로
- Prometheus 메트릭 기반 자동 판단 확장 가능

## ADR-007: 배포 전략 — Blue/Green (ch5.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Blue/Green 채택 (vs Flagger, Istio)
**이유**:
- 즉각 롤백 — active/preview 서비스 전환으로 이전 버전 즉시 복귀
- ArgoCD 통합 — Rollout 상태를 ArgoCD UI에서 확인
- autoPromotionSeconds: 30 — 자동 승격으로 학습 환경에서 빠른 검증
- Deployment API 호환 — kubectl argo rollouts 플러그인으로 상태 확인
