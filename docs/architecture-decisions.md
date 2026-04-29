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
