# Architecture Decision Records

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**날짜**: 2026-04-27 | **상태**: Accepted
- **결정**: ArgoCD 채택 (vs Flux, Jenkins X)
- **이유**: Web UI로 배포 상태 시각화, e2-medium 환경에서 구동 가능, CNCF Graduated

## ADR-002: CI 도구 — GitHub Actions WIF (ch3.4)
**날짜**: 2026-04-27 | **상태**: Accepted
- **결정**: GitHub Actions + Workload Identity Federation 채택
- **이유**: 저장소 네이티브, SA 키 생성 조직 정책 차단 환경에서 WIF가 유일한 선택

## ADR-003: 외부 트래픽 — Gateway API (ch5.2)
**날짜**: 2026-04-27 | **상태**: Accepted
- **결정**: GKE Gateway API (`gke-l7-regional-external-managed`) 채택
- **이유**: GKE 네이티브, 별도 Ingress Controller 불필요, 2.5에서 이미 활성화

## ADR-004: 배포 전략 — Argo Rollouts Blue/Green (ch5.3)
**날짜**: 2026-04-27 | **상태**: Accepted
- **결정**: Argo Rollouts Blue/Green (autoPromotionSeconds: 30) 채택
- **이유**: ArgoCD 동일 생태계, YAML 선언적, preview Pod으로 사전 검증 가능
