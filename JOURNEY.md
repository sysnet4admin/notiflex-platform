# Notiflex 여정 기록

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-06 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-06 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-06 | 기존 저장소 초기화 후 재사용 |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-06 | e2-medium 2노드 Spot |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-06 | v0.1.0 |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-06 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-06 | ArgoCD |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-06 | v0.1.1 /version |
| ch3 | 3.4 CI | ✅ | 2026-04-06 | GitHub Actions |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-06 | CI→매니페스트→ArgoCD |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-06 | kube-prometheus-stack |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-06 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-06 | PrometheusRule |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-06 | GKE Gateway API |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-06 | Argo Rollouts B/G |
| ch6 | 6.1 캐시 | ✅ | 2026-04-06 | Valkey |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-06 | GKE CSI + Secret Manager |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-06 | B/G → Canary |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-06 | api/worker/ops-pool |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-06 | root-app |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-06 | enterprise NS |
| ch8 | 8.1 메시징 | ✅ | 2026-04-06 | Strimzi Kafka KRaft |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-06 | Tempo + OTel SDK |
| ch8 | 8.3 CronJob | ✅ | 2026-04-06 | healthcheck 5분 |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-06 | |
| ch9 | 9.2 회고 | ✅ | 2026-04-06 | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-06 | |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-06 | |
| ch9 | 9.5 마무리 | ✅ | 2026-04-06 | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | UI 강점, CNCF Graduated |
| 메트릭 | Prometheus+Grafana | Datadog, CloudWatch | 오픈소스, CRD 기반 |
| 로그 | Loki+Fluent Bit | ELK, CloudWatch Logs | 경량, Grafana 통합 |
| 트래픽 | GKE Gateway API | Ingress NGINX, Istio | GKE 네이티브, 표준 API |
| 배포 | Argo Rollouts | Flagger, K8s native | ArgoCD 통합, B/G→Canary |
| 캐시 | Valkey | Redis, Memcached | Redis 호환, 오픈소스 |
| 시크릿 | GKE CSI Driver | Sealed Secrets, ESO | GKE managed, WI 통합 |
| 배포전략 | Canary | B/G 유지 | 점진적 트래픽 이동, 안전 |
| 노드 배치 | nodeSelector | taint, affinity | 단순, 직관적 |
| 앱 관리 | App of Apps | ApplicationSet | 선언적, 디렉터리 기반 |
| 메시징 | Kafka (Strimzi) | RabbitMQ, NATS | 이벤트 스트리밍, KRaft |
| 트레이싱 | Tempo | Jaeger | 경량, Grafana 통합 |
| 주기 작업 | CronJob | 외부 스케줄러 | K8s 네이티브, 단순 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | ch2.6 |
| Notiflex 이미지 | v0.6.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0→v0.6.0 |
| ArgoCD | 3.3.6 (stable) | ch3.2 |
| Kafka | 4.1.0 (Strimzi 0.51) | ch8.1 |
| OTel SDK | 1.43.0 | ch8.2 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | monitoring, Valkey |
| api-pool | e2-medium | 1 | notiflex-api (2 replicas), enterprise |
| worker-pool | e2-standard-2 | 1 | Strimzi, Kafka broker |
| ops-pool | e2-small | 1 | Tempo, CronJob |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch6.2 | Prometheus 리소스 이름 찾기 | kubectl get prometheus -n monitoring으로 실제 이름 확인 |
