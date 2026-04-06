# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-06 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-06 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-06 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-06 | e2-medium×2, Spot, Gateway API |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-06 | v0.1.0, 2 Pods Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-06 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-06 | ArgoCD v3 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-06 | v0.1.1 Rolling Update |
| ch3 | 3.4 CI | ✅ | 2026-04-06 | GitHub Actions |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-06 | CI→매니페스트 자동 업데이트 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-06 | Prometheus + Grafana |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-06 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-06 | PrometheusRule |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-06 | Gateway API |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-06 | Argo Rollouts Blue/Green |
| ch6 | 6.1 캐시 | ✅ | 2026-04-06 | Valkey (Bitnami Helm) |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-06 | GKE Secret Manager CSI |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-06 | B/G → Canary 전환 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-06 | api-pool, worker-pool, ops-pool |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-06 | root-app → smb, enterprise |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-06 | enterprise 네임스페이스 |
| ch8 | 8.1 메시징 | ✅ | 2026-04-06 | Strimzi Kafka (KRaft) |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-06 | Tempo + OTel SDK |
| ch8 | 8.3 CronJob | ✅ | 2026-04-06 | 헬스체크 CronJob |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-06 | |
| ch9 | 9.2 회고 | ✅ | 2026-04-06 | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-06 | |
| ch9 | 9.4 GitAIOps 분석 | ✅ | 2026-04-06 | |
| ch9 | 9.5 마무리 | ✅ | 2026-04-06 | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X | Web UI로 배포 상태 실시간 확인, CNCF Graduated |
| 메트릭 모니터링 | Prometheus + Grafana | Datadog, CloudWatch | 오픈소스, Helm 원클릭, 커뮤니티 생태계 |
| 로그 수집 | Loki + Fluent Bit | ELK, CloudWatch Logs | Grafana 통합, 경량, 인덱스 없는 설계 |
| 트래픽 관리 | Gateway API | Ingress NGINX, Istio | K8s 표준, GKE 네이티브 지원 |
| 무중단 배포 | Argo Rollouts | Flagger, K8s native | ArgoCD와 통합, B/G+Canary 모두 지원 |
| 캐시 | Valkey | Redis, Memcached | Redis 호환 오픈소스, LFDE 라이선스 |
| 시크릿 관리 | GKE Secret Manager CSI | Sealed Secrets, External Secrets | GKE 네이티브, Workload Identity 연동 |
| 배포 전략 | Canary | Blue/Green 유지 | 리소스 효율, 점진적 위험 감소 |
| 노드 스케줄링 | nodeSelector + 전용 풀 | taint/toleration, affinity | 직관적, GKE 노드풀과 자연스러운 매핑 |
| 멀티앱 관리 | App of Apps | ApplicationSet | 선언적, Git 기반 관리 |
| 메시징 | Kafka (Strimzi) | RabbitMQ, NATS | 이벤트 스트리밍, 대규모 처리, KRaft 모드 |
| 트레이싱 | Grafana Tempo | Jaeger, Zipkin | Grafana 통합, 관측성 3축 완성 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | v0.6.0 | v0.1.0→v0.1.1→v0.2.0→v0.3.0→v0.4.0→v0.5.0→v0.6.0 |
| ArgoCD | v3.3.6 | |
| Argo Rollouts | v1.8.1 | |
| Prometheus | kube-prometheus-stack | |
| Grafana | (kube-prometheus 번들) | |
| Loki | 6.x | |
| Fluent Bit | 0.x | |
| Tempo | 2.9.0 | |
| Kafka | 4.1.0 (Strimzi 0.51) | |
| OTel SDK | v1.43.0 | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | Prometheus, Grafana, Loki, Fluent Bit, Valkey |
| api-pool | e2-medium | 1 | Notiflex API (Rollout) |
| worker-pool | e2-medium | 1 | Strimzi Operator, Kafka Broker |
| ops-pool | e2-medium | 1 | ArgoCD, Argo Rollouts, Tempo, CronJob |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch3.2 | ArgoCD NetworkPolicy가 repo-server egress 차단 | NetworkPolicy 삭제 |
| ch3.2 | ArgoCD private repo 인증 실패 | forceHttpBasicAuth: true 설정 |
| ch5.3 | B/G→Canary 전환 시 selfHeal 충돌 | git push 먼저 → Rollout 삭제 → ArgoCD 재생성 |
| ch6.2 | CSI driver/provider 이름 오류 | GKE: secrets-store-gke.csi.k8s.io / provider: gke |
| ch6.2 | Workload Identity 미활성화 | 클러스터+노드풀 모두 활성화 필요 |
| ch7.2 | Grafana strategic merge patch 실패 | JSON patch (op:replace, container index 2) |
| ch9.5 | GKE 삭제 후 orphan 리소스 잔류 | Gateway 리소스 먼저 삭제 후 클러스터 삭제 |
