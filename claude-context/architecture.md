# Notiflex 아키텍처 스냅샷 (ch6 완료 시점)

## 클러스터 토폴로지
- notiflex-cluster / asia-northeast3-a / default-pool (e2-medium × 2)
- Gateway API, Workload Identity, Secret Manager CSI 활성화

## 컴포넌트
Gateway → Canary Rollout → Pod → Valkey + CSI Secret

## 참고
- JOURNEY.md, docs/architecture-decisions.md, CLAUDE.md
