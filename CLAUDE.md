# Notiflex Platform

## 프로젝트 개요
Notiflex — B2B 알림 SaaS 플랫폼. 다양한 채널(이메일, SMS, 웹훅)을 통해 알림을 전송하는 서비스.

## 기술 스택
- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **인프라**: GKE Standard (Zonal), Spot VM
- **배포**: GitOps (ArgoCD)

## GCP 설정
- **프로젝트 ID**: inlaid-vehicle-484800-n8
- **리전**: asia-northeast3
- **존**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/inlaid-vehicle-484800-n8/notiflex

## 행동 규칙
- 항상 현재 상태를 확인한 후 실행한다
- 변경 전 기존 상태를 먼저 확인한다
- 에러 발생 시 원인을 파악하고 해결한다
