---
description: notiflex-platform JOURNEY.md + CLAUDE.md를 최근 진행에 맞춰 갱신
---

# /update-docs

각 장 마지막에 호출하여 다음을 자동 갱신:

## 1. JOURNEY.md 갱신

`notiflex-platform/JOURNEY.md`에 다음 섹션을 갱신/추가:

- **진행 현황 표**: 각 서브챕터 ⬜→✅, 완료일 기록
- **도구 선택 기록**: 3-prompt 패턴에서 독자가 실제 선택한 도구 + 선택 이유
- **현재 버전**: notiflex-api 이미지 tag, 주요 도구 버전 (ArgoCD, Helm chart 등)
- **현재 리소스**: 노드풀 구성, namespace 목록, 활성 helm release
- **트러블슈팅 이력**: 가드레일에 없는 새로운 문제 + 해결 (있을 시)

## 2. CLAUDE.md 갱신 (메타데이터만)

CLAUDE.md의 환경 섹션 수치값 갱신:
- GCP project, region, zone (변경 없을 시 그대로)
- kubectl context 이름
- Artifact Registry URL
- 새로 추가된 도구·리소스 정보

## 3. 검증

- `git diff JOURNEY.md CLAUDE.md` 출력
- 사용자에게 commit/push 여부 확인
