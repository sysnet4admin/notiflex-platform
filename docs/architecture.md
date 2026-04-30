# Architecture

This document outlines the architecture of the Notiflex platform.

## ADRs

- [ADR-001: GitOps 도구로 ArgoCD 채택](./architecture-decisions.md#adr-001-gitops-도구로-argocd-채택-3장)
- [ADR-002: CI 도구로 GitHub Actions 채택](./architecture-decisions.md#adr-002-ci-도구로-github-actions-채택-3장)
- [ADR-003: 메트릭 모니터링으로 Prometheus+Grafana 채택](./architecture-decisions.md#adr-003-메트릭-모니터링으로-prometheusgrafana-채택-4장)
- [ADR-004: 로그 수집으로 Loki+Fluent Bit 채택](./architecture-decisions.md#adr-004-로그-수집으로-lokifluent-bit-채택-4장)
- [ADR-005: 알림은 PrometheusRule + Alertmanager](./architecture-decisions.md#adr-005-알림은-prometheusrule--alertmanager-4장)
- [ADR-006: 외부 트래픽 관리로 Gateway API 채택](./architecture-decisions.md#adr-006-외부-트래픽-관리로-gateway-api-채택-5장)
- [ADR-007: 무중단 배포 전략으로 Argo Rollouts 채택](./architecture-decisions.md#adr-007-무중단-배포-전략으로-argo-rollouts-채택-5장)
- [ADR-008: 캐시는 Valkey](./architecture-decisions.md#adr-008-캐시는-valkey-6장)
- [ADR-009: 시크릿은 GKE Secret Manager CSI + Workload Identity](./architecture-decisions.md#adr-009-시크릿은-gke-secret-manager-csi--workload-identity-6장)
- [ADR-010: 배포 전략은 Canary로 전환](./architecture-decisions.md#adr-010-배포-전략은-canary로-전환-6장)
- [ADR-011: 워크로드별 노드풀 분리](./architecture-decisions.md#adr-011-워크로드별-노드풀-분리-7장)
- [ADR-012: 다중 앱 관리는 App of Apps](./architecture-decisions.md#adr-012-다중-앱-관리는-app-of-apps-7장)
- [ADR-013: 멀티테넌시는 Namespace 분리](./architecture-decisions.md#adr-013-멀티테넌시는-namespace-분리-7장)
- [ADR-014: 메시징은 Strimzi 기반 Kafka](./architecture-decisions.md#adr-014-메시징은-strimzi-기반-kafka-8장)
- [ADR-015: 분산 트레이싱은 Tempo + OpenTelemetry](./architecture-decisions.md#adr-015-분산-트레이싱은-tempo--opentelemetry-8장)
- [ADR-016: 주기 작업은 Kubernetes CronJob](./architecture-decisions.md#adr-016-주기-작업은-kubernetes-cronjob-8장)
