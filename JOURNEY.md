# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-05 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-05 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-05 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-05 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-05 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-05 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-05 | ArgoCD v3.3.6 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-05 | /notify 엔드포인트 |
| ch3 | 3.4 CI | ✅ | 2026-04-05 | GitHub Actions |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-05 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-05 | kube-prometheus-stack |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-05 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-05 | PrometheusRule |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-05 | GKE Gateway API |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-05 | Argo Rollouts B/G |
| ch6 | 6.1 캐시 | ✅ | 2026-04-05 | Valkey standalone |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-05 | WI + CSI Secret Manager |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-05 | B/G → Canary |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-05 | api/worker/ops 3풀 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-05 | root-app 패턴 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-05 | enterprise 테넌트 |
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
| GitOps | ArgoCD | Flux, Jenkins X | UI, CNCF Graduated |
| 메트릭 | Prometheus+Grafana | Datadog | 무료, K8s 네이티브 |
| 로그 | Loki+Fluent Bit | ELK | 경량, Grafana 통합 |
| 트래픽 | GKE Gateway API | Ingress NGINX, Istio | K8s 표준, 설치 불필요 |
| 배포 | Argo Rollouts | Flagger | ArgoCD 생태계 |
| 캐시 | Valkey | Redis, Memcached | BSD 라이선스, Redis 호환 |
| 시크릿 | CSI + Secret Manager | Sealed Secrets, ESO | GKE 네이티브, WI 통합 |
| 배포 전략 | Canary | B/G 유지 | 점진적 전환, 리소스 효율 |
| 노드 분리 | nodeSelector | taint, affinity | 단순, 직관적 |
| 앱 관리 | App of Apps | ApplicationSet | Git 기반, 직관적 YAML |
| 멀티테넌시 | Namespace + RBAC | vCluster | 경량, K8s 네이티브 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.5.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0 |
| ArgoCD | v3.3.6 | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | 시스템, Valkey |
| api-pool | e2-medium | 1 | notiflex-api (SMB+Enterprise) |
| worker-pool | e2-standard-2 | 1 | (ch8: Kafka) |
| ops-pool | e2-small | 1 | (ch8: 모니터링 이전 예정) |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch6 | CSI DaemonSet으로 CPU 240m 추가 점유 → 대규모 Pending | monitoring CPU 최소화 (Prom/Grafana/AM 5m) |
| ch6 | B/G replicas:2 시 4 Pod 필요 → CPU 부족 | replicas: 1로 축소 |
| ch6 | Canary+Valkey 순환 Pending | old RS scale=0으로 교착 해소 |
