# Architecture Decision Records

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**시점**: 2026-04 / **결정**: ArgoCD 채택 (vs Flux, Jenkins X)
**이유**: Web UI 배포 상태 시각화, e2-medium 환경 구동 가능, CNCF Graduated, automated sync + selfHeal

## ADR-002: CI 도구 — GitHub Actions + WIF (ch3.4)
**시점**: 2026-04 / **결정**: GitHub Actions + Workload Identity Federation 채택
**이유**: 저장소 네이티브, SA 키 조직 정책 차단 환경에서 WIF가 유일한 GCP 인증 수단

## ADR-003: 메트릭 모니터링 — Prometheus + Grafana (ch4.2)
**시점**: 2026-04 / **결정**: kube-prometheus-stack 채택 (vs Datadog, New Relic)
**이유**: GKE 네이티브, 오픈소스, kube-prometheus-stack으로 Prometheus + Grafana + Alertmanager 일괄 설치

## ADR-004: 로깅 — Loki + Fluent Bit (ch4.3)
**시점**: 2026-04 / **결정**: Loki + Fluent Bit 채택 (vs ELK Stack)
**이유**: Grafana와 통합, 인덱싱 없이 로그 저장, loki-stack 차트로 단순 설치

## ADR-005: 알림 — PrometheusRule + Alertmanager (ch4.4)
**시점**: 2026-04 / **결정**: PrometheusRule + Alertmanager 채택 (vs Grafana Alert)
**이유**: Prometheus와 네이티브 연동, git 버전 관리, release 라벨로 Prometheus Operator 자동 로드

## ADR-006: 외부 트래픽 — GKE Gateway API (ch5.2)
**시점**: 2026-04 / **결정**: GKE Gateway API (`gke-l7-regional-external-managed`) 채택
**이유**: GKE 네이티브, 별도 Ingress Controller 불필요, HTTPRoute + HealthCheckPolicy로 세밀한 트래픽 제어

## ADR-007: 배포 전략 — Argo Rollouts Blue/Green (ch5.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Blue/Green (autoPromotionSeconds: 30) 채택
**이유**: ArgoCD와 동일 생태계, YAML 선언적, preview Pod으로 사전 검증 가능

## ADR-008: 캐시 — Valkey (ch6.1)
**시점**: 2026-04 / **결정**: Valkey 채택 (vs Redis, Memcached)
**이유**: Redis fork, bitnami 차트, INCR으로 분산 ID 카운터, 50m CPU로 구동 가능

## ADR-009: 시크릿 관리 — GKE Secret Manager CSI + WI (ch6.2)
**시점**: 2026-04 / **결정**: GKE Secret Manager CSI + Workload Identity 채택
**이유**: GCP 네이티브, SA 키 없이 인증, 파일 마운트, GKE managed addon

## ADR-010: 배포 전략 전환 — Canary (ch6.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Canary 채택 (Blue/Green 대체)
**이유**: 점진적 트래픽 이동, 리소스 효율, 동일 Argo Rollouts 생태계

## ADR-011: 노드 스케줄링 — nodeSelector (ch7.2)
**시점**: 2026-04 / **결정**: nodeSelector + `cloud.google.com/gke-nodepool` 채택
**이유**: GKE 자동 라벨 기반, 단순 선언, 워크로드 분리, 커스텀 라벨 금지

## ADR-012: 멀티앱 관리 — App of Apps (ch7.3)
**시점**: 2026-04 / **결정**: ArgoCD App of Apps + sync-wave 채택 (vs ApplicationSet)
**이유**: 단순 구조, 설치 순서 보장, 기존 ArgoCD 확장, directory.recurse: true 필수

## ADR-013: 메시징 — Kafka (ch8.1)
**시점**: 2026-04 / **결정**: Kafka (Strimzi 0.51.0, KRaft 4.1.0) 채택
**이유**: 고처리량 + 영속 메시지, KRaft 모드, Strimzi Operator, worker-pool 배치

## ADR-014: 분산 트레이싱 — Tempo (ch8.2)
**시점**: 2026-04 / **결정**: Grafana Tempo 채택
**이유**: Grafana 통합, 단일 바이너리, OTLP 표준, ops-pool 배치
