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
| ch2 | 2.7 첫 커밋 | ⬜ | | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-29 | ArgoCD 설치 + private GitHub 저장소 연동 완료 |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-29 | `/version` 엔드포인트 추가 및 v0.1.1 롤링 업데이트 완료 |
| ch3 | 3.4 CI | ✅ | 2026-04-29 | GitHub Actions CI 추가 (push to main, app 변경 시 빌드/푸시) |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-29 | CI에서 k8s/smb/deployment.yaml 이미지 태그 자동 갱신 + git push로 ArgoCD 자동 배포 연결 완료 |
| ch4 | 4.2 메트릭 모니터링 | ⬜ | | |
| ch4 | 4.3 로그 수집 | ⬜ | | |
| ch4 | 4.4 알림 | ⬜ | | |
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

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | 2026-04-29: 초기 버전 설정 |
| Notiflex 이미지 | v0.1.1 | 2026-04-29: `/version` 엔드포인트 포함 이미지 업로드 (sha256:2a8ad8a5bfc28e7d758085d05a88887551186375c97effe634fbcb2e8d8d5d88) |
| ArgoCD | quay.io/argoproj/argocd:v3.3.8 | 2026-04-29: gke-sysnet4admin_book_gitaiops 클러스터에 설치 및 notiflex-platform 저장소 연결 |
| Kafka | | |
| OTel SDK | | |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| default-pool | e2-medium | 2 | notiflex-api |

## 트러블슈팅 이력

독자가 겪은 문제와 해결 방법을 기록한다. 같은 문제를 다시 겪지 않도록 한다.

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.6 | Artifact Registry 리포지토리 미존재 | notiflex 리포지토리를 생성한 뒤 이미지 푸시 |
| ch3.2 | ArgoCD Application이 `Sync Unknown` (`Repository not found`) | repo Secret에 GitHub 토큰(`forceHttpBasicAuth: "true"`) 등록 후 ArgoCD 전체 rollout restart |
| ch3.3 | ArgoCD가 최신 커밋(`a08e25a`)을 즉시 반영하지 않음 | `kubectl --context gke-sysnet4admin_book_gitaiops -n argocd annotate application notiflex-smb argocd.argoproj.io/refresh=hard --overwrite`로 강제 재동기화 |
