# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-30 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-30 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-30 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-30 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-30 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-30 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-30 | |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-30 | |
| ch3 | 3.4 CI | ✅ | 2026-04-30 | |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-30 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-30 | |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-30 | |
| ch4 | 4.4 알림 | ✅ | 2026-04-30 | |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-30 | |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-30 | |
| ch5 | 5.4 ADR 기록 | ✅ | 2026-04-30 | |
| ch6 | 6.1 캐시 | ✅ | 2026-04-30 | |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-30 | |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-30 | |
| ch6 | 6.4 아키텍처 스냅샷 | ✅ | 2026-04-30 | |
| ch7 | 7.2 멀티 노드풀 | ⬜ | | |
| ch7 | 7.3 App of Apps | ⬜ | | |
| ch7 | 7.4 멀티테넌시 | ⬜ | | |
| ch8 | 8.1 메시징 | ⬜ | | |
| ch8 | 8.2 트레이싱 | ⬜ | | |
| ch8 | 8.3 CronJob | ⬜ | | |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps (ch3.2) | ArgoCD v3.3.8 | Flux, Jenkins X | K8s 네이티브, Web UI, App of Apps 지원 |
| CI (ch3.4) | GitHub Actions | Jenkins, GitLab CI | 저장소 통합, WIF 지원, 무료 |
| 메트릭 (ch4.2) | Prometheus + Grafana | Datadog, New Relic | 오픈소스, kube-prometheus-stack 통합 |
| 로깅 (ch4.3) | Loki + Fluent Bit | ELK, Datadog Logs | Grafana 통합, 경량 인덱싱 |
| 알림 (ch4.4) | PrometheusRule | Grafana Alert | Prometheus 네이티브, PromQL 표현식 |
| 트래픽 관리 (ch5.2) | Gateway API (gke-l7-regional-external-managed) | Ingress, NGINX | GKE 네이티브, K8s 표준, HealthCheckPolicy |
| 배포 전략 (ch5.3) | Argo Rollouts Blue/Green | Flagger, Istio | ArgoCD 통합, 즉각 롤백, autoPromotion |
| 캐시 (ch6.1) | Valkey (Bitnami, standalone) | Redis, Memcached | Redis fork, BSD-3 라이선스, INCR 분산 ID |
| 시크릿 관리 (ch6.2) | GKE Secret Manager CSI + WI | K8s Secret, Vault | GKE 네이티브, SA 키 불필요, 파일 마운트 |
| 배포 전략 전환 (ch6.3) | Argo Rollouts Canary | Blue/Green 유지 | 트래픽 점진 이동, 운영 위험 최소화 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.2.1 | v0.1.0→v0.1.1(롤링)→v0.2.0(Valkey)→v0.2.1(CSI) |
| ArgoCD | v3.3.8 | |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | notiflex-api |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| | | |
