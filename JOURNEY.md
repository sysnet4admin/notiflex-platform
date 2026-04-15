# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-15 | Claude Code + statusline 설정 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-15 | 프로젝트/존/Docker 인증 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-15 | notiflex-platform 생성 |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-15 | e2-medium x2, Spot, Gateway API |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-15 | v0.1.0, Pod 2개 Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-15 | GitHub 푸시 완료 |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-15 | ArgoCD 설치, notiflex-smb App |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-15 | /version endpoint, v0.1.1 |
| ch3 | 3.4 CI | ✅ | 2026-04-15 | GitHub Actions CI |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-15 | CI 매니페스트 자동 업데이트 |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-15 | kube-prometheus-stack |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-15 | Loki + Fluent Bit |
| ch4 | 4.4 알림 | ✅ | 2026-04-15 | PrometheusRule Pod 재시작 |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-15 | Gateway API, 외부 IP 35.216.96.244 |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-15 | Argo Rollouts Blue/Green |
| ch6 | 6.1 캐시 | ✅ | 2026-04-15 | Valkey standalone, v0.3.0 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-15 | CSI Driver + Secret Manager + Workload Identity |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-15 | Blue/Green -> Canary 전환 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-15 | api/worker/ops-pool 3개 추가 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-15 | root-app이 argocd/ 감시 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-15 | enterprise NS + App of Apps |
| ch8 | 8.1 메시징 | ✅ | 2026-04-15 | Strimzi Kafka KRaft, Sarama producer |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-15 | Tempo + OTel SDK v1.43.0 |
| ch8 | 8.3 CronJob | ✅ | 2026-04-15 | 5분마다 헬스체크, ops-pool |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

독자가 3-프롬프트 패턴(탐색→비교→실행)에서 실제로 선택한 도구와 이유를 기록한다.

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux, Jenkins X, Spinnaker | Web UI, CNCF Graduated, 학습에 최적 |
| CI | GitHub Actions | Cloud Build, GitLab CI, Jenkins | GitHub 네이티브, 무료, 간편 |
| 트래픽 관리 | Gateway API | Ingress NGINX, Istio, Traefik | GKE 네이티브, 설치 불필요 |
| 무중단 배포 | Argo Rollouts | Flagger, K8s Rolling | ArgoCD 통합, Blue/Green+Canary |
| 캐시 | Valkey | Redis, Memcached, DragonflyDB | Redis 호환, BSD 라이선스 |
| Secret 관리 | CSI + Secret Manager | Sealed Secrets, ESO | GKE 네이티브, Workload Identity |
| 배포 전략 진화 | Canary | Blue/Green 유지 | 점진적 전환, 리소스 효율 |
| 노드 배치 | nodeSelector | taint/toleration, affinity | 가장 단순, GKE 자동 라벨 |
| 앱 관리 | App of Apps | ApplicationSet, 수동 | 직관적, Git 디렉터리 감시 |
| 멀티테넌시 | Namespace+RBAC | vCluster, 클러스터 분리 | K8s 기본, 리소스 효율 |
| 메시징 | Kafka (Strimzi) | RabbitMQ, NATS, Pub/Sub | 높은 처리량, KRaft 모드 |
| 트레이싱 | Grafana Tempo | Jaeger, Zipkin | Grafana 통합, 관측성 3축 |
| 배치 자동화 | K8s CronJob | Argo Workflows, Airflow | K8s 내장, 설치 불필요 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | ch2.6 초기 |
| Notiflex 이미지 | v0.6.0 | ch2.6 v0.1.0 -> ch3.3 v0.1.1 -> ch3.5 v0.2.0 -> ch6.1 v0.3.0 -> ch6.2 v0.4.0 -> ch8.1 v0.5.0 -> ch8.2 v0.6.0 |
| ArgoCD | 3.3.6 (stable) | ch3.2 설치 |
| Kafka | 4.2.0 (Strimzi 0.51) | ch8.1 설치 |
| OTel SDK | v1.43.0 | ch8.2 추가 |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | 시스템, Valkey |
| api-pool | e2-medium | 1 | notiflex-api |
| worker-pool | e2-standard-2 | 1 | Kafka broker |
| ops-pool | e2-small | 1 | CronJob, Tempo |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch8.1 | Kafka 4.0.0 UnsupportedKafkaVersionException | Strimzi 0.51은 4.1.0/4.1.1/4.2.0 지원. 4.2.0으로 변경 |
| ch8.1 | Helm zsh escaping | tolerations[0] 패턴이 zsh에서 실패. --set 없이 설치 |
| ch8.2 | git push rejected (CI 충돌) | CI가 SHA 태그로 이미지 업데이트. rebase로 해결 |
| ch8.2 | Kafka DNS 실패 | Spot VM DNS 불안정. graceful degradation으로 앱 정상 동작 |
