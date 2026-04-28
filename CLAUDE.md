# Notiflex Platform

## 프로젝트 개요

Notiflex — B2B 알림 SaaS 플랫폼. 기업 고객이 이메일·SMS·푸시 알림을 API로 발송하는 서비스.

## 기술 스택

- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **인프라**: GKE Standard (Zonal), Spot VM

## GCP 설정

- **프로젝트 ID**: project-75fce205-dfa5-4975-a56
- **리전**: asia-northeast3 (서울)
- **존**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex

## 행동 규칙

1. 항상 현재 상태를 확인한 후 실행한다
2. kubectl 명령에는 반드시 `--context gke-sysnet4admin_book_gitaiops`를 지정한다
3. 파괴적 작업 전 반드시 확인한다
