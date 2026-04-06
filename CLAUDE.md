# Notiflex Platform

## 프로젝트 개요
Notiflex — B2B 알림 SaaS 플랫폼. 쿠버네티스 위에서 운영되는 REST API 서버.

## 기술 스택
- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (멀티스테이지 빌드)
- **인프라**: GKE Standard (Zonal), Spot VM
- **CI/CD**: GitHub Actions → ArgoCD

## GCP 설정
- **프로젝트 ID**: inlaid-vehicle-484800-n8
- **리전**: asia-northeast3
- **존**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/inlaid-vehicle-484800-n8/notiflex

## 행동 규칙
- 변경 전 반드시 현재 상태를 확인한다
- 확인 후 실행한다
- K8s 리소스 적용 시 namespace를 먼저 생성한다
