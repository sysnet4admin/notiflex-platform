# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-30 | |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-30 | |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-30 | |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-30 | |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-30 | |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-30 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-30 | |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-30 | |
| ch3 | 3.4 CI | ✅ | 2026-04-30 | |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-30 | |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-04-30 | |
| ch4 | 4.3 로그 수집 | ✅ | 2026-04-30 | |
| ch4 | 4.4 알림 | ✅ | 2026-04-30 | |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-04-30 | |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-30 | |
| ch6 | 6.1 캐시 | ✅ | 2026-04-30 | Valkey INCR 연동 |
| ch6 | 6.2 시크릿 관리 | ✅ | 2026-04-30 | GKE Secret Manager CSI+WI |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-30 | Blue/Green -> Canary |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-30 | api/worker/ops 노드풀 생성 |
| ch7 | 7.3 App of Apps | ✅ | 2026-04-30 | Root-app, App of Apps 패턴 적용 |
| ch7 | 7.4 멀티테넌시 | ✅ | 2026-04-30 | enterprise 네임스페이스 분리 |
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
| GitOps | ArgoCD | Flux | Kubernetes 네이티브, 선언적, 다양한 기능 |
| CI | GitHub Actions | Cloud Build, GitLab CI, Jenkins | GitHub 내장, YAML 선언적, 무료 크레딧 |
| 메트릭 모니터링 | Prometheus + Grafana | Datadog, CloudWatch, Google Cloud Monitoring | 오픈소스 표준, 무료, Grafana 시각화 |
| 로그 수집 | Loki + Fluent Bit | ELK Stack, CloudWatch Logs, Google Cloud Logging | 경량, Grafana 통합, 라벨 기반 인덱싱 |
| 외부 트래픽 관리 | Gateway API | Ingress NGINX, Istio, Traefik | Kubernetes 차세대 표준, GKE 네이티브, 역할 분리 |
| 무중단 배포 | Argo Rollouts | Flagger, K8s native Rolling Update | ArgoCD 통합, CRD 기반, 점진적 전략 진화 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | |
| Notiflex 이미지 | asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:v0.1.1 | ch3.2 |
| ArgoCD | quay.io/argoproj/argocd:v3.3.8 | ch3.2 |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | notiflex-api, argocd |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.6 | Go build failed due to newline in string literal | Fixed the string literal in `main.go` |
