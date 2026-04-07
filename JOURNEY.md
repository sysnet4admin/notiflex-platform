# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-07 | Claude Code 정상 동작 확인 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-07 | SDK 529.0.0, 프로젝트/존 설정 완료 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-07 | notiflex-platform 저장소 구성 |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-07 | e2-medium×2, Spot, Gateway API |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-07 | v0.1.0, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-07 | GitHub 푸시 완료 |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-07 | ArgoCD 설치, Application 생성 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-07 | /version 엔드포인트, v0.1.1 |
| ch3 | 3.4 CI | ✅ | 2026-04-07 | GitHub Actions + Cloud Build |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-07 | CI가 매니페스트 자동 업데이트, v0.2.0 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-07 | kube-prometheus-stack + Grafana 대시보드 |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-07 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-07 | PrometheusRule (Pod 재시작 알림) |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-07 | Gateway API + HTTPRoute + HealthCheckPolicy |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-07 | Argo Rollouts Blue/Green |
| ch6 | 6.1 캐시 | ✅ | 2026-04-07 | Valkey standalone, /id 엔드포인트, v0.3.0 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-07 | CSI Secret Store + Workload Identity, v0.4.0 |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-07 | Blue/Green → Canary (20→50→80→100), v0.5.0 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-07 | api-pool, worker-pool, ops-pool 추가 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-07 | root-app + smb/enterprise Application |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-07 | enterprise 네임스페이스, 크로스NS Valkey |
| ch8 | 8.1 메시징 | ✅ | 2026-04-07 | Kafka (Strimzi KRaft), v0.6.0 |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-07 | Tempo + OTel SDK, v0.7.0 |
| ch8 | 8.3 CronJob | ✅ | 2026-04-07 | 5분 주기 헬스체크 CronJob |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-07 | 31개 파일, 1159줄, 코드<매니페스트 |
| ch9 | 9.2 회고 | ✅ | 2026-04-07 | 도구 선택 13개 영역 종합 |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-07 | AI 기반 자동 생성 |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-07 | Git+AI+Ops 결합 분석 |
| ch9 | 9.5 마무리 | ✅ | 2026-04-07 | 프로덕션 전환 로드맵 |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | Web UI, CNCF Graduated, GKE 호환 |
| 메트릭 | Prometheus + Grafana | Datadog, CloudWatch | 오픈소스, GKE 네이티브 |
| 로그 | Loki + Fluent Bit | ELK, CloudWatch Logs | 경량, Grafana 통합 |
| 알림 | PrometheusRule | Grafana Alerting | PromQL 기반, GitOps 관리 |
| 트래픽 | Gateway API | Ingress NGINX, Istio | K8s 공식, GKE 네이티브 |
| 무중단 배포 | Argo Rollouts | Flagger, K8s native | Canary+B/G, ArgoCD 연동 |
| 캐시 | Valkey | Redis, Memcached | Redis 호환 오픈소스, CNCF |
| 시크릿 | CSI Secret Store | Sealed Secrets, ESO | GKE addon, 파일 마운트 |
| 배포 전략 | Canary | Blue/Green 유지 | 점진적 트래픽, 리스크 최소 |
| 노드 분리 | nodeSelector | Taint, Affinity | 단순명료 |
| 앱 관리 | App of Apps | ApplicationSet | 직관적, YAML 기반 |
| 멀티테넌시 | Namespace 격리 | vCluster | 리소스 효율적 |
| 메시징 | Kafka (Strimzi) | RabbitMQ, NATS | 이벤트 스트리밍 표준, KRaft |
| 트레이싱 | Tempo | Jaeger | Grafana 통합, 경량 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | ch2.6에서 설정 |
| Notiflex 이미지 | v0.7.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0→v0.6.0→v0.7.0 |
| ArgoCD | v3.3.6 | ch3.2에서 설치 |
| Argo Rollouts | v1.7 | ch5.3에서 설치 |
| Kafka (Strimzi) | 0.45.0 (KRaft) | ch8.1에서 설치 |
| OTel SDK | v1.43.0 | ch8.2에서 추가 |
| Tempo | 1.32 | ch8.2에서 설치 |
| kube-prometheus-stack | - | ch4.2에서 설치 |
| Loki | - | ch4.3에서 설치 |
| Fluent Bit | - | ch4.3에서 설치 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium (2 vCPU, 4GB) | 2 | 시스템, 모니터링, ArgoCD |
| api-pool | e2-medium (2 vCPU, 4GB) | 1 | notiflex-api Pod |
| worker-pool | e2-standard-2 (2 vCPU, 8GB) | 1 | Kafka 브로커 |
| ops-pool | e2-small (0.5 vCPU, 2GB) | 1 | Tempo, CronJob |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch6.1 | ArgoCD 동기화 지연 | hard refresh 어노테이션 패치 |
| ch6.2 | Valkey Pod Pending (CPU 부족) | 모니터링 스택 CPU requests 1m으로 축소 |
| ch6.2 | Workload Identity 활성화 지연 | 클러스터 + 노드풀 양쪽 순차 활성화 |
| ch8.2 | v0.7.0 빌드 실패 (go.sum) | go mod tidy 재실행 후 재빌드 |
| 매 장 | CI SHA 태그 충돌 | git pull → sed 버전 복원 → push |
