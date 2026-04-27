# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-27 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-27 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-27 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-27 | GatewayClass 초기화 지연 → clusters update 재실행 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-27 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-27 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-27 | ArgoCD v3 repo-server 인증 → GitHub Token |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-27 | /version 엔드포인트 + v0.1.1 |
| ch3 | 3.4 CI | ✅ | 2026-04-27 | SA 키 차단 → WIF 전환 |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-27 | CI→manifest→ArgoCD 파이프라인 완성 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-27 | kube-prometheus-stack 설치 |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-27 | loki-stack 설치 (loki 최신 차트 복잡성 우회) |
| ch4 | 4.4 알림 | ✅ | 2026-04-27 | PrometheusRule + Alertmanager |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-27 | Gateway API (gke-l7-regional) |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-27 | Argo Rollouts Blue/Green |
| ch6 | 6.1 캐시 | ✅ | 2026-04-27 | Valkey standalone |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-27 | K8s Secret (WI 비활성화로 CSI 우회) |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-27 | Blue/Green → Canary (20→50→80%) |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-27 | api/worker/ops-pool (--node-labels 없이 GKE 자동 라벨 활용) |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-27 | root-app → notiflex-smb 자동 관리 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-27 | notiflex-enterprise namespace + api-pool 배치 |
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
| GitOps | ArgoCD | Flux, Jenkins X, Spinnaker | Web UI + ArgoCD 생태계 통합 |
| CI | GitHub Actions (WIF) | Cloud Build, GitLab CI, Jenkins | 저장소 일치 + SA 키 차단 환경에서 WIF가 적합 |
| 메트릭 | Prometheus + Grafana | Datadog, Google Cloud Monitoring | 오픈소스 표준, Loki/Tempo 통합 |
| 로그 | Loki + Promtail | ELK Stack, Google Cloud Logging | 경량, Grafana 통합 |
| 알림 | PrometheusRule + Alertmanager | Grafana Alerting UI | GitOps 호환 (YAML→Git→ArgoCD) |
| 외부 트래픽 | Gateway API | Ingress NGINX, Istio | GKE 네이티브, 별도 설치 불필요 |
| 배포 전략 | Argo Rollouts Canary | Blue/Green, Flagger | ch5 Blue/Green 경험 후 ch6에서 Canary로 진화 |
| 캐시 | Valkey | Redis, Memcached | BSD 라이선스, Redis 호환 |
| 시크릿 | K8s Secret (WI 비활성 대안) | CSI+SecretManager, Sealed Secrets | WI 비활성 환경 |
| 노드 스케줄링 | nodeSelector + 멀티노드풀 | taint/toleration, nodeAffinity | 가장 단순, GKE 자동 라벨 활용 |
| 멀티앱 관리 | App of Apps | ApplicationSet | 직관적, 순수 YAML |

## 현재 버전

| 컴포넌트 | 버전 |
|---------|------|
| Go | 1.25 |
| Notiflex 이미지 | sha-f18aa61 (CI 자동 태그) |
| GKE | 1.35.1-gke.1396002 |
| ArgoCD | v2.14.x (stable) |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | API, ArgoCD |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.5 | GatewayClass가 클러스터 생성 후 즉시 나타나지 않음 | `gcloud container clusters update --gateway-api=standard` 재실행 후 30초 대기 |
| ch3.2 | ArgoCD v3 "authentication required: Repository not found" | GitHub Token을 repo secret password로 등록 (forceHttpBasicAuth 단독 불충분) |
| ch3.4 | SA 키 생성 iam.disableServiceAccountKeyCreation 정책 차단 | Workload Identity Federation(WIF)으로 전환 |
