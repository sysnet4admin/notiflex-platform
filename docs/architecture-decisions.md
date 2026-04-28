# Architecture Decision Records

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**시점**: 2026-04 / **결정**: ArgoCD 채택 (vs Flux, Jenkins X)
**이유**:
- Web UI로 배포 상태 시각화
- e2-medium 2노드 환경에서 구동 가능
- CNCF Graduated — 안정성 검증
- automated sync + selfHeal로 GitOps 드리프트 자동 복구

## ADR-002: CI 도구 — GitHub Actions + WIF (ch3.4)
**시점**: 2026-04 / **결정**: GitHub Actions + Workload Identity Federation 채택
**이유**:
- 저장소 네이티브 — 별도 서버 없이 GitHub 저장소와 통합
- SA 키 생성이 조직 정책으로 차단된 환경에서 WIF가 유일한 GCP 인증 수단
- 무료 사용량 확보 (private 저장소 500분/월)
- YAML 선언적 파이프라인

## ADR-003: 메트릭 모니터링 — Prometheus + Grafana (ch4.2)
**시점**: 2026-04 / **결정**: kube-prometheus-stack 채택 (vs Datadog, New Relic)
**이유**:
- GKE 네이티브 — 추가 에이전트 없이 kube-state-metrics + node-exporter 통합
- 오픈소스 — 비용 없음
- kube-prometheus-stack으로 Prometheus + Grafana + Alertmanager 일괄 설치
- Loki, Tempo와 같은 Grafana 생태계로 통합 관측 가능성 구성 가능

## ADR-004: 로깅 — Loki + Fluent Bit (ch4.3)
**시점**: 2026-04 / **결정**: Loki + Fluent Bit 채택 (vs ELK Stack)
**이유**:
- Grafana와 통합 — 메트릭·로그를 동일 대시보드에서 조회
- 인덱싱 없이 로그 저장 — ELK 대비 리소스 절감
- loki-stack 차트로 단순 설치
- e2-medium 2노드 환경에서도 안정적으로 구동

## ADR-005: 알림 — PrometheusRule + Alertmanager (ch4.4)
**시점**: 2026-04 / **결정**: PrometheusRule + Alertmanager 채택 (vs Grafana Alert)
**이유**:
- Prometheus와 네이티브 연동 — PromQL 기반 알림 규칙
- git 버전 관리 — YAML로 선언적 관리
- kube-prometheus-stack에 Alertmanager 포함
- `release: kube-prometheus` 라벨로 Prometheus Operator가 자동 로드

## ADR-006: 외부 트래픽 — GKE Gateway API (ch5.2)
**시점**: 2026-04 / **결정**: GKE Gateway API (`gke-l7-regional-external-managed`) 채택 (vs Ingress, NGINX, Istio)
**이유**:
- GKE 네이티브 — 별도 Ingress Controller 설치 불필요
- ch2.5에서 `--gateway-api=standard`로 이미 활성화
- HTTPRoute + HealthCheckPolicy로 세밀한 트래픽 제어
- Kubernetes 표준 API — 벤더 종속 없음

## ADR-007: 배포 전략 — Argo Rollouts Blue/Green (ch5.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Blue/Green (autoPromotionSeconds: 30) 채택 (vs Flagger, Istio)
**이유**:
- ArgoCD와 동일 생태계 — 별도 서비스 메시 불필요
- YAML 선언적 — `strategy.blueGreen` 필드로 배포 전략을 git에 기록
- preview Pod으로 프로덕션 배포 전 사전 검증 가능
- autoPromotionSeconds로 자동 프로모션

## ADR-008: 캐시 — Valkey (ch6.1)
**시점**: 2026-04 / **결정**: Valkey 채택 (vs Redis, Memcached)
**이유**:
- Redis fork — API 완전 호환, 라이선스 자유(BSD-3)
- bitnami Helm 차트 지원 — 단순 설치, Secret 자동 생성
- INCR 명령으로 분산 ID 카운터 구현
- e2-medium 환경에서 50m CPU로 구동 가능

## ADR-009: 시크릿 관리 — GKE Secret Manager CSI + WI (ch6.2)
**시점**: 2026-04 / **결정**: GKE Secret Manager CSI Driver + Workload Identity 채택 (vs K8s Secret, HashiCorp Vault)
**이유**:
- GCP 네이티브 — Secret Manager와 GKE 완전 통합
- Workload Identity — SA 키 없이 Pod이 GCP API 직접 인증
- 파일 마운트 방식 — 메모리에 Secret 노출 없음
- GKE managed CSI (addon) — 별도 헬름 설치 불필요

## ADR-010: 배포 전략 전환 — Canary (ch6.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Canary (20%→50%→80% 점진) 채택 (Blue/Green 대체)
**이유**:
- 점진적 트래픽 이동 — 문제 발생 시 영향 범위 최소화
- Blue/Green보다 리소스 효율 — 2배 복제 불필요
- 동일 Argo Rollouts 생태계 — 전략 전환 시 컨트롤러 유지
- 30초 간격 단계적 승격

## ADR-011: 노드 스케줄링 — nodeSelector (ch7.2)
**시점**: 2026-04 / **결정**: nodeSelector + `cloud.google.com/gke-nodepool` 채택
**이유**:
- GKE 자동 라벨 기반 — `--node-labels` 불필요
- 단순 선언 — `spec.nodeSelector` 필드 1줄로 완료
- 워크로드 분리 — api/worker/ops 3풀로 비용·성능 최적화
- 커스텀 라벨 금지 — `workload: api` 등 임의 키는 영구 Pending의 원인

## ADR-012: 메시징 — Kafka (ch8.1)
**시점**: 2026-04 / **결정**: Kafka (Strimzi 0.51.0, KRaft 4.1.0) 채택
**이유**:
- 고처리량 + 영속 메시지 — 알림 데이터 유실 없이 비동기 처리
- KRaft 모드 — ZooKeeper 없이 단순화된 구조
- Strimzi Operator — K8s CRD로 Kafka 클러스터 선언적 관리
- worker-pool 배치 — 브로커 전용 노드로 리소스 격리

## ADR-013: 분산 트레이싱 — Tempo (ch8.2)
**시점**: 2026-04 / **결정**: Grafana Tempo 채택
**이유**:
- Grafana 통합 — 메트릭·로그·트레이스를 동일 UI에서 조회
- 단일 바이너리 — ops-pool 1 Pod로 동작
- OTLP 표준 — OTel SDK와 직접 연동, 벤더 종속 없음
- ops-pool 배치 — 운영 도구 전용 노드로 분리
