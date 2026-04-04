# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-04 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-04 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-04 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-04 | e2-medium ×2, gateway-api=standard |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-04 | v0.1.0 배포, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-04 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-04 | ArgoCD 설치, repo secret + preemptive restart |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-04 | /notify 엔드포인트 (v0.1.1) |
| ch3 | 3.4 CI | ✅ | 2026-04-04 | GitHub Actions CI pipeline |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-04 | CI manifest auto-update → ArgoCD sync |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-04 | kube-prometheus-stack (Prometheus + Grafana) |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-04 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-04 | Grafana alerting (CPU/memory rules) |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-04 | Gateway API (GKE native) |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-04 | Argo Rollouts Blue/Green (v0.2.0) |
| ch6 | 6.1 캐시 | ✅ | 2026-04-04 | Valkey + /id 엔드포인트 (v0.3.0) |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-04 | GKE CSI Secret Manager (file-based) |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-04 | B/G → Canary (20→50→80%, 30s pause) |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-04 | api-pool, worker-pool, ops-pool + nodeSelector |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-04 | root-app → notiflex-smb, notiflex-enterprise |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-04 | enterprise 테넌트 (namespace + rollout + service) |
| ch8 | 8.1 메시징 | ✅ | 2026-04-04 | Strimzi Kafka 4.1.0 on worker-pool |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-04 | Tempo + OTel SDK (v0.5.0) |
| ch8 | 8.3 CronJob | ✅ | 2026-04-04 | healthcheck CronJob (5분 간격) |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-04 | 10 commits, 637 lines, 28 files |
| ch9 | 9.2 회고 | ✅ | 2026-04-04 | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-04 | |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-04 | |
| ch9 | 9.5 마무리 | ✅ | 2026-04-04 | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | Web UI, CNCF Graduated, GKE 호환 |
| 모니터링 | Prometheus+Grafana | Datadog, CloudWatch | 오픈소스, Helm 원클릭 |
| 로깅 | Loki+Fluent Bit | ELK, CloudWatch Logs | 경량, Grafana 통합 |
| 트래픽 | Gateway API | Ingress NGINX, Istio | GKE 네이티브, K8s 표준 |
| 배포 전략 | Argo Rollouts | Flagger, K8s native | ArgoCD 시너지, CRD 기반 |
| 캐시 | Valkey | Redis, Memcached | Redis 포크, 오픈소스 |
| 시크릿 | CSI Secret Manager | Sealed Secrets, External Secrets | GKE 네이티브, 관리형 |
| 메시징 | Strimzi Kafka | RabbitMQ, NATS | 이벤트 스트리밍 표준 |
| 트레이싱 | Tempo | Jaeger | Grafana 통합, 경량 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25.2 | 초기 설정 |
| Notiflex 이미지 | v0.5.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0 |
| ArgoCD | latest | ch3에서 설치 |
| Kafka | 4.1.0 (Strimzi) | ch8에서 설치 |
| OTel SDK | 1.43.0 | ch8에서 추가 |
| Tempo | latest | ch8에서 설치 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | kube-system, monitoring, argocd |
| api-pool | e2-medium | 1 | notiflex-api (SMB + enterprise) |
| worker-pool | e2-medium | 1 | Kafka (Strimzi), Valkey |
| ops-pool | e2-medium | 1 | Tempo |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch8 | Strimzi latest → Kafka 3.8.0 미지원 | Kafka version 4.1.0으로 변경 |
| ch8 | entity-operator user-operator OOM 재시작 | userOperator 제거 (불필요) |
