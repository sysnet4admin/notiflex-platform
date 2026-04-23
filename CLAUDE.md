# Notiflex Platform

## 프로젝트 개요
- 프로젝트명: Notiflex
- 도메인: B2B 알림 SaaS 플랫폼

## 기술 스택
- 언어: Go (표준 라이브러리 중심)
- 컨테이너: `scratch` 베이스 이미지
- 오케스트레이션: Kubernetes (GKE)

## GCP 설정
- Project ID: `project-75fce205-dfa5-4975-a56`
- Region: `asia-northeast3`
- Zone: `asia-northeast3-a`

## Artifact Registry
- `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex`

## 행동 규칙
1. 실행 전 목표와 영향 범위를 먼저 확인한다.
2. 변경 전 현재 상태(브랜치, 리소스, 설정)를 점검한다.
3. 배포/인프라 명령은 대상 환경을 명시적으로 확인한 뒤 실행한다.
4. 변경 후 검증 절차(빌드/테스트/상태 확인)를 수행한다.
