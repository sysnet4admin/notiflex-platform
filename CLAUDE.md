# Notiflex Platform

## 프로젝트 개요
Notiflex는 B2B 알림 SaaS 플랫폼이다. 고객사 시스템에서 이벤트가 발생하면 이메일, SMS, 웹훅으로 알림을 전달한다.

## 기술 스택
- 언어: Go (표준 라이브러리만 사용, 외부 프레임워크 없음)
- 컨테이너: scratch 베이스 이미지
- 인프라: GKE Standard (Zonal)

## GCP 설정
- 프로젝트 ID: inlaid-vehicle-484800-n8
- 리전: asia-northeast3
- 존: asia-northeast3-a

## Artifact Registry
- asia-northeast3-docker.pkg.dev/inlaid-vehicle-484800-n8/notiflex

## 행동 규칙
- 항상 현재 상태를 확인한 후 실행한다
- 변경 전 현재 상태를 먼저 확인한다
- 모든 kubectl 명령에 --context gke-notiflex를 지정한다
