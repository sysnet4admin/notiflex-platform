# Notiflex Platform

「AI 시대에 개발자가 알아야 하는 인프라 구성 배포 with 클로드 코드」 실습 저장소.

B2B 알림 SaaS 플랫폼을 GKE 위에 GitOps + AI로 구축하는 전 과정을 담았다.

## 구성

| 디렉터리 | 내용 |
|---------|------|
| `app/` | Notiflex Go API (Valkey INCR, Kafka Producer, OTel 트레이싱) |
| `k8s/smb/` | SMB 테넌트 매니페스트 (Rollout, Service, Gateway, CronJob) |
| `k8s/enterprise/` | Enterprise 테넌트 매니페스트 |
| `k8s/kafka/` | Strimzi Kafka 클러스터 (KRaft, v4.1.0) |
| `k8s/monitoring/` | PrometheusRule |
| `argocd/` | ArgoCD App of Apps (root-app + apps/) |
| `helm-values/` | Helm 차트 설정값 |
| `docs/` | ADR-001~016 아키텍처 결정 기록 |
| `claude-context/` | AI 참조용 아키텍처 스냅샷 |

## GitAIOps 3층 구조

```
CLAUDE.md          → AI에게 프로젝트 메타데이터 제공 (매 대화 자동 로드)
claude-context/    → 현재 아키텍처 스냅샷 (AI 참조용)
docs/ADR           → 팀의 결정 누적 기록 (사람 + AI가 함께 검토)
```

Git이 인프라의 단일 진실 소스, AI가 운영 표준의 살아있는 저자.

## 빠른 시작

`ONBOARDING.md` 참조.
