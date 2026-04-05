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
| ch6 | 6.1 캐시 | ✅ | 2026-04-05 | Valkey |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-05 | GKE CSI Driver |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-05 | B/G → Canary |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-05 | api/worker/ops |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-05 | root-app |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-05 | enterprise NS |
| ch8 | 8.1 메시징 | ✅ | 2026-04-05 | Strimzi + Kafka 4.2.0 |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-05 | Tempo + OTel v1.43 |
| ch8 | 8.3 CronJob | ✅ | 2026-04-05 | healthcheck 5분 |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-05 | |
| ch9 | 9.2 회고 | ✅ | 2026-04-05 | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-05 | |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-05 | |
| ch9 | 9.5 마무리 | ✅ | 2026-04-05 | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | UI, CNCF Graduated |
| 메트릭 | Prometheus+Grafana | Datadog | 무료, K8s 네이티브 |
| 로그 | Loki+Fluent Bit | ELK | 경량, Grafana 통합 |
| 트래픽 | GKE Gateway API | Ingress NGINX, Istio | K8s 표준, 설치 불필요 |
| 배포 | Argo Rollouts | Flagger | ArgoCD 생태계 |
| 캐시 | Valkey | Redis, Memcached | 오픈소스, Redis 호환 |
| 시크릿 | GKE CSI Driver | Sealed Secrets | GKE 네이티브, WI 통합 |
| 배포 전략 | Canary | Blue/Green 유지 | 점진적, 리소스 효율 |
| 노드 분리 | nodeSelector | Taint/Toleration | 단순, 직관적 |
| 앱 관리 | App of Apps | ApplicationSet | 직관적, YAML 기반 |
| 메시징 | Kafka(Strimzi) | RabbitMQ, NATS | 대용량, 내구성 |
| 트레이싱 | Tempo | Jaeger | Grafana 통합, 경량 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.5.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0 |
| ArgoCD | v3.3.6 | |
| Strimzi | 0.51.0 | |
| Kafka | 4.2.0 | |
| OTel SDK | v1.43.0 | |
| valkey-go | v1.0.73 | |
| IBM/sarama | v1.47.0 | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | monitoring, argocd |
| api-pool | e2-medium | 1 | notiflex-api |
| worker-pool | e2-standard-2 | 1 | Strimzi, Kafka |
| ops-pool | e2-small | 1 | Tempo, CronJob |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch6 | ArgoCD Sync 정상이나 매니페스트 미반영 (v3 NetworkPolicy가 repo-server egress 차단) | NetworkPolicy 삭제 + 전체 재시작 |
