# Notiflex Platform

B2B 알림 SaaS 플랫폼.

## 기술 스택
- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지, 멀티스테이지 빌드
- **인프라**: GKE Standard (Zonal), Spot VM

## GCP 설정
- **프로젝트**: inlaid-vehicle-484800-n8
- **리전**: asia-northeast3
- **존**: asia-northeast3-a
- **Artifact Registry**: asia-northeast3-docker.pkg.dev/inlaid-vehicle-484800-n8/notiflex

## 행동 규칙
- 항상 현재 상태를 확인한 후 실행한다
- 변경 전 현재 상태를 확인한다
- 에러 발생 시 원인을 분석한 후 재시도한다
