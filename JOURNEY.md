# Notiflex 여정 기록

이 파일은 실제 진행 상태를 저장소 기준으로 기록한다.  
기준 시각: 2026-04-23 (KST)

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-04-23 | codex run-01 시작 완료 |
| ch2 | 2.3 gcloud 설정 | ✅ | 2026-04-23 | GCP 기반 CI/이미지 경로 사용 중 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-04-23 | `sysnet4admin/notiflex-platform` 연결 |
| ch2 | 2.5 GKE 클러스터 | ✅ | 2026-04-23 | GKE nodepool 라벨 사용 중 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-04-23 | 초기 앱/매니페스트 생성 |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-04-23 | `e99d548` |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-04-23 | `48d281b` |
| ch3 | 3.3 기능 추가 | ✅ | 2026-04-23 | `/version` + 롤링 반영 |
| ch3 | 3.4 CI | ✅ | 2026-04-23 | GitHub Actions + Cloud Build |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-04-23 | CI가 `k8s/smb/rollout.yaml` 이미지 태그 갱신 |
| ch4 | 4.2 메트릭 모니터링 | ⚠️ | 2026-04-23 | 파일 존재(워크트리), 커밋 반영 전 |
| ch4 | 4.3 로그 수집 | ⚠️ | 2026-04-23 | 파일 존재(워크트리), 커밋 반영 전 |
| ch4 | 4.4 알림 | ⚠️ | 2026-04-23 | 파일 존재(워크트리), 커밋 반영 전 |
| ch5 | 5.2 트래픽 관리 | ⚠️ | 2026-04-23 | Gateway/HealthCheckPolicy 파일 존재(워크트리) |
| ch5 | 5.3 무중단 배포 | ✅ | 2026-04-23 | Argo Rollouts(Blue/Green→Canary) |
| ch6 | 6.1 캐시 | ✅ | 2026-04-23 | Valkey 연동 (`05e5106`) |
| ch6 | 6.2 시크릿 관리 | ⚠️ | 2026-04-23 | ExternalSecret/ClusterSecretStore 파일 존재(워크트리) |
| ch6 | 6.3 Canary 전환 | ✅ | 2026-04-23 | `k8s/smb/rollout.yaml` canary 단계 존재 |
| ch7 | 7.2 멀티 노드풀 | ✅ | 2026-04-23 | `api-pool`/`ops-pool` 사용 |
| ch7 | 7.3 App of Apps | ⚠️ | 2026-04-23 | `argocd/root-app.yaml` 존재(워크트리) |
| ch7 | 7.4 멀티테넌시 | ⚠️ | 2026-04-23 | `k8s/enterprise/*` 존재(워크트리) |
| ch8 | 8.1 메시징 | ⚠️ | 2026-04-23 | Kafka 매니페스트 존재, 앱 코드 연동 미확인 |
| ch8 | 8.2 트레이싱 | ✅ | 2026-04-23 | OTel + Tempo exporter (`785eb72`) |
| ch8 | 8.3 CronJob | ⚠️ | 2026-04-23 | `k8s/smb/healthcheck-cronjob.yaml` 존재(워크트리) |
| ch9 | 9.1 저장소 분석 | ✅ | 2026-04-23 | 구조/히스토리/JOURNEY 대조 완료 |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ✅ | 2026-04-23 | `ONBOARDING.md` 생성(클러스터 실측 반영) |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| GitOps | ArgoCD | Flux | 책 표준 흐름과 실습 연계성 |
| 배포 전략 | Argo Rollouts Canary | Rolling, Blue/Green 유지 | 점진 트래픽 전환 검증 |
| 캐시/공유 ID | Valkey | Redis | 경량 구성 + INCR 사용 단순성 |
| 관측성 트레이싱 | Tempo + OTel | Jaeger | Grafana 스택과 결합 용이 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25.0 | `app/go.mod` |
| Notiflex 이미지 | `v0.1.2`(manifest), `v0.1.3`(appVersion) | 이미지/코드 버전 불일치 점검 필요 |
| ArgoCD | 미기록 | 추후 클러스터 조회 필요 |
| Kafka | 4.1.0 (Strimzi) | `k8s/kafka/kafka-cluster.yaml` |
| OTel SDK | 1.43.0 | `app/go.mod` |

## 현재 리소스

| 노드풀 | 머신 타입 | 노드 수 | 주요 워크로드 |
|--------|----------|---------|-------------|
| api-pool | 미기록 | 미기록 | notiflex-api Rollout |
| ops-pool | 미기록 | 미기록 | healthcheck CronJob |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch3.5 | CI 결과 이미지 태그와 매니페스트 동기화 필요 | workflow에서 `sed`로 rollout image 자동 갱신 |
