# Notiflex 프로젝트

## 프로젝트 개요
- **프로젝트명**: Notiflex
- **설명**: B2B 알림 SaaS 플랫폼
- **목표**: 안정적이고 확장 가능한 알림 서비스 제공

## 기술 스택
- **언어**: Go (표준 라이브러리)
- **컨테이너**: scratch 베이스 이미지
- **클라우드**: Google Cloud Platform (GCP)

## GCP 설정
- **Project ID**: project-75fce205-dfa5-4975-a56
- **Region**: asia-northeast3
- **Zone**: asia-northeast3-a

## 주요 리소스
- **Artifact Registry**: `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex`
- **GKE Cluster**: `notiflex-cluster` (asia-northeast3)

## 행동 규칙
- **실행 전 확인**: 모든 변경은 실행 전 동료 리뷰 또는 자동화된 검증을 거칩니다.
- **상태 확인**: 변경 작업 전후로 반드시 시스템의 현재 상태(예: Pod, Service 상태)를 확인하고 기록합니다.
- **문서화**: 중요한 결정이나 아키텍처 변경은 문서로 남깁니다.
