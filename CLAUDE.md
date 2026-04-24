# Notiflex Platform

B2B 알림 SaaS 플랫폼, Kubernetes 운영 환경 구축 실습 프로젝트

## 환경
- **GCP Project**: project-75fce205-dfa5-4975-a56
- **Region**: asia-northeast3
- **Zone**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex
- **kubectl context**: gke-sysnet4admin_book_gitaiops (ch2.5에서 생성 예정)

## 진행 현황 (run-53)

각 장 [별도] Hints & Tips를 포함하여 책 전체 흐름 검증 중.

## 안전 규칙

모든 kubectl 명령에 `--context gke-sysnet4admin_book_gitaiops`를 명시.

## 행동 규칙 (run-53 ch3 [별도])

1. **ArgoCD Application 변경 후 sync 상태 확인 필수**
   - 변경 commit/push 후 `kubectl annotate app ... refresh=normal` 또는 UI 확인
   - Sync Status=Synced 이전엔 다음 작업 금지

2. **모든 K8s 매니페스트는 `k8s/` 디렉터리 하위에만**
   - `k8s/smb/`, `k8s/enterprise/`, `k8s/monitoring/` 등 namespace별 분리
   - 루트 또는 다른 디렉터리에 .yaml 매니페스트 금지

3. **이미지 태그는 명시적 버전만**
   - `sha-XXXXXXX` (CI 빌드) 또는 `vX.Y.Z` (수동 릴리스) 만 허용
   - `:latest` 절대 금지 (재현성·롤백 불가)
