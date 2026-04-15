# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-15 | Claude Code + statusline 설정 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-15 | 프로젝트/존/Docker 인증 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-15 | notiflex-platform 생성 |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-15 | e2-medium x2, Spot, Gateway API |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-15 | v0.1.0, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-15 | GitHub 푸시 완료 |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-15 | ArgoCD 설치 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-15 | /version, v0.1.1 Rolling |
| ch3 | 3.4 CI | ✅ | 2026-04-15 | GitHub Actions |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-15 | CI→ArgoCD 자동 배포 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-15 | kube-prometheus-stack |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-15 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-15 | PrometheusRule |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-15 | Gateway API, IP 35.216.1.32 |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-15 | Argo Rollouts Blue/Green |
| ch6 | 6.1 캐시 | ✅ | 2026-04-15 | Valkey standalone, v0.3.0 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-15 | CSI + Secret Manager, v0.4.0 |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-15 | B/G→Canary, v0.5.0 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-15 | api/worker/ops-pool |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-15 | root-app + directory.recurse |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-15 | enterprise 네임스페이스 |
| ch8 | 8.1 메시징 | ✅ | 2026-04-15 | Strimzi Kafka, v0.6.0 |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-15 | Tempo + OTel SDK, v0.7.0 |
| ch8 | 8.3 CronJob | ✅ | 2026-04-15 | 헬스체크 5분마다, ops-pool |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | UI, CRD, 생태계 |
| CI | GitHub Actions | GitLab CI, CircleCI | GitHub 네이티브 |
| 모니터링 | kube-prometheus-stack | Datadog, Elastic | 오픈소스, Helm 통합 |
| 로깅 | Loki + Fluent Bit | EFK, Datadog | 경량, Grafana 통합 |
| 알림 | Alertmanager | PagerDuty, Opsgenie | Prometheus 내장 |
| 트래픽 | Gateway API | Ingress, Istio | K8s 표준, GKE 네이티브 |
| 배포 전략 | Argo Rollouts | Flagger, Istio | CRD 기반, ArgoCD 호환 |
| 캐시 | Valkey | Redis, Memcached | BSD 라이선스, Redis 호환 |
| 시크릿 | CSI + Secret Manager | Vault, SOPS | GKE 네이티브, 무설치 |
| 메시징 | Kafka (Strimzi) | RabbitMQ, NATS | 업계 표준, CRD 관리 |
| 트레이싱 | Tempo | Jaeger, Zipkin | Grafana 통합, 경량 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | ch2.6 초기 |
| Notiflex 이미지 | v0.7.0 | ch2.6 v0.1.0 → ch8.2 v0.7.0 |
| ArgoCD | 2.14 | ch3.2 설치 |
| Kafka | 4.1.0 (Strimzi 0.51) | ch8.1 설치 |
| OTel SDK | 1.43.0 | ch8.2 설치 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | 관측성, Valkey |
| api-pool | e2-medium | 1 (Spot) | notiflex-api |
| worker-pool | e2-medium | 1 (Spot) | Kafka |
| ops-pool | e2-medium | 1 (Spot) | ArgoCD, Argo Rollouts |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| | | |
