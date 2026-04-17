# Notiflex Platform

B2B 알림 SaaS 플랫폼, Kubernetes 운영 환경 구축 실습 프로젝트

## 기술 스택
- **언어**: Go (표준 라이브러리만 사용, 외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **인프라**: GKE Standard (Zonal)
- **CI/CD**: GitHub Actions + ArgoCD (예정)

## GCP 설정
- **프로젝트 ID**: inlaid-vehicle-484800-n8
- **리전**: asia-northeast3 (서울)
- **존**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/inlaid-vehicle-484800-n8/notiflex

## 행동 규칙
1. 명령 실행 전 현재 상태를 확인한다 (`kubectl get`, `gcloud config list` 등)
2. 파일 변경 전 기존 내용을 먼저 읽는다
3. 에러 발생 시 원인을 분석하고 해결 방안을 제시한 뒤 진행한다
4. 매니페스트 작성 시 네임스페이스(`notiflex`)를 명시한다
5. **모든 `kubectl` 명령에 `--context gke-sysnet4admin_book_gitaiops`를 반드시 지정한다**
   (잘못된 클러스터 대상 실행 방지)
6. 리소스 생성/삭제 전에는 영향 범위를 먼저 설명한다
7. 이미지 태그는 `latest`를 쓰지 않고 명시적 버전(`v0.1.0` 등)을 사용한다
8. 토큰, 키 그리고 비밀번호는 코드/매니페스트에 하드코딩하지 않는다
   (환경 변수, GitHub Secrets, Secret Manager 사용)
