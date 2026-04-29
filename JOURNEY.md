# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ⬜ | | |
| ch2 | 2.3 gcloud 설정 | ⬜ | | |
| ch2 | 2.4 GitHub 저장소 | ⬜ | | |
| ch2 | 2.5 GKE 클러스터 | ⬜ | | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-29 | notiflex-api v0.1.0 배포 완료 |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-29 | 초기 커밋 및 `/update-docs` 커스텀 스킬 추가 |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-29 | ArgoCD 설치 + private GitHub 저장소 연동 완료 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-29 | `/version` 엔드포인트 추가 및 v0.1.1 롤링 업데이트 완료 |
| ch3 | 3.4 CI | ✅ | 2026-04-29 | GitHub Actions CI 추가 (push to main, app 변경 시 빌드/푸시) |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-29 | CI에서 k8s/smb/deployment.yaml 이미지 태그 자동 갱신 + git push로 ArgoCD 자동 배포 연결 완료 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-29 | kube-prometheus-stack(Helm) 설치 + Notiflex Grafana 대시보드 ConfigMap 적용 |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-29 | Loki + Fluent Bit 설치, Grafana Loki 데이터소스 추가, `{job="fluent-bit",namespace="notiflex"}` 로그 조회 확인 |
| ch4 | 4.4 알림 | ✅ | 2026-04-29 | PrometheusRule(`pod-restart-alert`) 생성/적용 완료, Alertmanager 연동 확인 |
| ch5 | 5.2 트래픽 관리 | ⬜ | | |
| ch5 | 5.3 무중단 배포 | ⬜ | | |
| ch6 | 6.1 캐시 | ⬜ | | |
| ch6 | 6.2 시크릿 관리 | ⬜ | | |
| ch6 | 6.3 Canary 전환 | ⬜ | | |
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
| ch2.6 컨테이너 베이스 이미지 | scratch + 멀티스테이지 빌드 | alpine, distroless | 최소 공격면과 작은 이미지 크기 |
| ch3.2 GitOps 도구 | ArgoCD | Flux | UI/자동 동기화 기반 학습 흐름에 적합 |
| ch3.4 CI 도구 | GitHub Actions + docker build/push (방식 A) | gcloud builds submit (방식 B) | 권한 구성이 단순하고 학습 흐름에 적합 |
| ch3.5 CI-CD 연결 | GitHub Actions에서 GitOps 매니페스트 자동 갱신 후 ArgoCD 자동 동기화 | CI 내 직접 `kubectl apply`, 별도 CD 파이프라인 분리 | Git 단일 소스 오브 트루스 유지 + 변경 이력 추적 용이 |
| ch4.2 메트릭 모니터링 | kube-prometheus-stack (Prometheus + Grafana + Alertmanager) | Datadog, New Relic, VictoriaMetrics | Helm 기반으로 실습 환경에서 빠르게 설치 가능하고 Grafana 대시보드 연계가 쉬움 |
| ch4.3 로그 수집 | Loki + Fluent Bit | ELK Stack, Cloud Logging | e2-medium 리소스 제약에서 경량 운영 가능하고 Grafana Explore와 즉시 연계 가능 |
| ch4.4 알림 | PrometheusRule + Alertmanager | Grafana Alerting, Cloud Monitoring Alert | GitOps(YAML/PR)로 이력 관리가 가능하고 기존 kube-prometheus-stack과 네이티브로 통합됨 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-04-29: 초기 버전 설정 |
| Notiflex 이미지 | `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:sha-7380ee4` | 2026-04-29: CI가 매니페스트를 자동 갱신해 배포, 현재 실행 이미지 digest `sha256:528b11ab3781e47d58764448a5fbde3044b31352600c3cd03d258fa7016e1e86` |
| ArgoCD | quay.io/argoproj/argocd:v3.3.8 | 2026-04-29: gke-sysnet4admin_book_gitaiops 클러스터에 설치 및 notiflex-platform 저장소 연결 |
| Prometheus | quay.io/prometheus/prometheus:v3.11.3 | 2026-04-29: `kube-prometheus-stack-84.3.0`으로 monitoring namespace에 설치 |
| Grafana | docker.io/grafana/grafana:13.0.1 | 2026-04-29: `kube-prometheus-grafana` 배포, Notiflex 대시보드 ConfigMap(`notiflex-grafana-dashboard`) 추가 |
| Loki | docker.io/grafana/loki:3.6.7 | 2026-04-29: `loki-7.0.0`(SingleBinary) 설치, `loki-datasource` ConfigMap으로 Grafana 데이터소스 등록(`isDefault: false`) |
| Fluent Bit | grafana/fluent-bit-plugin-loki:2.1.0-amd64 | 2026-04-29: `fluent-bit-2.6.0` DaemonSet 설치, Loki Gateway(`/loki/api/v1/push`)로 로그 수집 연동 |
| Kafka | 미설치 | 2026-04-29: ch8 이전 단계라 클러스터에 리소스 없음 |
| OTel SDK | 미설치 | 2026-04-29: ch8.2 이전 단계 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | notiflex-api, argocd, monitoring(kube-prometheus-stack + loki + fluent-bit) |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.6 | Artifact Registry 리포지토리 미존재 | notiflex 리포지토리를 생성한 뒤 이미지 푸시 |
| ch3.2 | ArgoCD Application이 `Sync Unknown` (`Repository not found`) | repo Secret에 GitHub 토큰(`forceHttpBasicAuth: "true"`) 등록 후 ArgoCD 전체 rollout restart |
| ch3.3 | ArgoCD가 최신 커밋(`a08e25a`)을 즉시 반영하지 않음 | `kubectl --context gke-sysnet4admin_book_gitaiops -n argocd annotate application notiflex-smb argocd.argoproj.io/refresh=hard --overwrite`로 강제 재동기화 |
| ch3.5 | GitHub Actions `if` 조건에서 `secrets.*` 직접 참조 시 워크플로 문법 오류 | 인증 관련 시크릿을 `env`로 옮기고 `if`는 `env` 기반으로 분기 |
| ch3.5 | CI 인증 방식이 환경마다 달라 빌드 실패 (SA Key/WIF 시크릿 키 이름 불일치) | SA Key, 레거시 WIF, GCP WIF 3가지 입력 조합을 모두 지원하도록 워크플로 보완 |
| ch4.3 | Loki 설치 시 기본 cache(chunks/results) Pod가 CPU 부족으로 Pending되어 Helm install timeout | `helm-values/loki.yaml`에서 `chunksCache.enabled=false`, `resultsCache.enabled=false`로 비활성화 후 `helm upgrade --install` 재실행 |
| ch4.4 | `kube_pod_container_status_restarts_total` 기반 경보는 Pod 삭제만으로 즉시 증가하지 않아 Alert가 `firing`되지 않을 수 있음 | 규칙 로드/Alertmanager 연동 확인 후, 실제 재시작 카운트가 증가하는 장애 시나리오(예: CrashLoop)로 추가 검증 |
