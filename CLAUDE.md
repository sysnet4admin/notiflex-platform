# Notiflex Platform

## 프로젝트 개요

Notiflex — B2B 알림 SaaS 플랫폼. 기업 고객별 알림 채널(이메일, SMS, Slack)을 관리한다.

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

- 항상 현재 상태를 확인한 후 작업을 진행한다
- kubectl 명령은 반드시 `--context gke-sysnet4admin_book_gitaiops`를 지정한다
- 변경 전 현재 리소스 상태를 먼저 확인한다
