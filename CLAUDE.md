# Notiflex Platform 작업 가이드

## 프로젝트 개요
- Notiflex는 B2B 알림 SaaS 플랫폼이다.

## 기술 스택
- 백엔드: Go (표준 라이브러리 중심)
- 컨테이너: `scratch` 베이스 이미지

## GCP 설정
- Project ID: `project-75fce205-dfa5-4975-a56`
- Region: `asia-northeast3`
- Zone: `asia-northeast3-a`

## Artifact Registry
- `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex`

## 행동 규칙
- 실행 전에 목적과 영향을 먼저 확인한다.
- 변경 전 현재 상태(리소스, 설정, Git 상태)를 먼저 점검한다.
- 위험한 작업(삭제/강제 반영)은 사전 확인 후에만 진행한다.
