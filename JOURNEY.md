# Notiflex 여정 기록

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-29 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-29 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-29 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-29 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-29 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-29 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-29 | |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-29 | |
| ch3 | 3.4 CI | ✅ | 2026-04-29 | |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-29 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-29 | |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-29 | |
| ch4 | 4.4 알림 | ✅ | 2026-04-29 | |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-29 | |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-29 | |
| ch5 | 5.4 ADR | ✅ | 2026-04-29 | |
| ch6 | 6.1 캐시 | ✅ | 2026-04-29 | |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-29 | |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-29 | |
| ch6 | 6.4 아키텍처 컨텍스트 | ✅ | 2026-04-29 | |
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

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps (ch3.2) | ArgoCD | Flux, Jenkins X | Web UI 배포 상태 시각화, e2-medium 환경 구동 가능, CNCF Graduated |
| CI (ch3.4) | GitHub Actions + WIF | Jenkins, GitLab CI | 저장소 네이티브, SA 키 조직 정책 차단 환경에서 WIF가 유일한 선택 |
| 메트릭 (ch4.2) | Prometheus + Grafana | Datadog, New Relic | GKE 네이티브, 오픈소스, kube-prometheus-stack으로 통합 설치 |
| 로깅 (ch4.3) | Loki + Fluent Bit | ELK Stack, Datadog | Grafana와 통합, 경량, 인덱싱 없이 로그 저장 |
| 알림 (ch4.4) | PrometheusRule + Alertmanager | Grafana Alert | Prometheus와 네이티브 연동, git 버전 관리 |
| 외부 트래픽 (ch5.2) | GKE Gateway API (gke-l7-regional-external-managed) | Ingress, NGINX, Istio | GKE 네이티브, 별도 Ingress Controller 불필요 |
| 배포 전략 (ch5.3) | Argo Rollouts Blue/Green | Flagger, Istio | ArgoCD 동일 생태계, YAML 선언적, preview Pod으로 사전 검증 가능 |
| 캐시 (ch6.1) | Valkey | Redis, Memcached | Redis fork, OSS, bitnami 차트 지원, INCR으로 분산 ID 카운터 구현 |
| 시크릿 관리 (ch6.2) | GKE Secret Manager CSI + WI | K8s Secret, HashiCorp Vault | GCP 네이티브, 키 노출 없음, 파일 마운트 방식 |
| 배포 전략 전환 (ch6.3) | Argo Rollouts Canary | Blue/Green | 점진적 트래픽 이동 (20%→50%→80%), 리소스 효율적 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.1.1 | 초기 버전 |
| ArgoCD | v3.3.8 | |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | 모든 워크로드 |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| | | |
