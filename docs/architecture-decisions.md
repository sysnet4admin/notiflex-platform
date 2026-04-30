# Architecture Decision Records

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**시점**: 2026-04 / **결정**: ArgoCD v3.3.8 채택 (vs Flux, Jenkins X)
**이유**:
- K8s 네이티브 CRD(Application, AppProject) — 선언적 배포 상태 관리
- Web UI 제공 — 배포 상태·히스토리·diff를 시각적으로 확인
- App of Apps 패턴 지원 — 7장 멀티앱 관리에서 자연스럽게 확장
- automated sync + selfHeal — git이 단일 진실 소스(Single Source of Truth) 역할

## ADR-002: CI 도구 — GitHub Actions (ch3.4)
**시점**: 2026-04 / **결정**: GitHub Actions 채택 (vs Jenkins, GitLab CI)
**이유**:
- GitHub 저장소와 통합 — 별도 서버 운영 불필요
- Workload Identity Federation 지원 — SA 키 없이 GCP 인증
- push 트리거로 이미지 빌드 → 매니페스트 업데이트 → ArgoCD 자동 배포 흐름 완성
- 무료 티어로 학습 환경 구성

## ADR-003: 메트릭 수집 — Prometheus + Grafana (ch4.2)
**시점**: 2026-04 / **결정**: kube-prometheus-stack 채택 (vs Datadog, New Relic)
**이유**:
- 오픈소스 — 라이선스 비용 없음
- ServiceMonitor/PrometheusRule CRD — K8s 네이티브 운영 규칙 관리
- Grafana와 통합 — Loki(로그)·Tempo(트레이스) 데이터소스를 같은 UI에서 확인
- kube-prometheus-stack 차트로 Operator+Prometheus+Grafana+Alertmanager 일괄 설치

## ADR-004: 로그 수집 — Loki + Fluent Bit (ch4.3)
**시점**: 2026-04 / **결정**: Loki + Fluent Bit 채택 (vs ELK, Datadog Logs)
**이유**:
- Grafana 통합 — 메트릭과 로그를 같은 대시보드에서 확인
- 인덱스 없는 레이블 기반 — Elasticsearch 대비 저사양 환경에 유리
- Fluent Bit 경량 DaemonSet — 노드당 최소 자원으로 로그 수집
- SingleBinary 모드로 단일 Pod 운영 — 학습 환경에 적합

## ADR-005: 알림 — PrometheusRule (ch4.4)
**시점**: 2026-04 / **결정**: PrometheusRule CRD 채택 (vs Grafana Alert)
**이유**:
- Prometheus 네이티브 — PromQL 표현식으로 정밀한 조건 정의
- git으로 버전 관리 — 알림 규칙 변경 이력 추적
- kube-prometheus-stack과 자동 연동 — `labels.release` 매칭으로 즉시 활성화
- Alertmanager 라우팅과 결합 — 심각도별 수신 채널 분리

## ADR-006: 외부 트래픽 — Gateway API (ch5.2)
**시점**: 2026-04 / **결정**: GKE Gateway API(gke-l7-regional-external-managed) 채택 (vs Ingress, NGINX)
**이유**:
- GKE 네이티브 — 별도 Ingress Controller 설치 없이 Google 관리형 L7 LB 사용
- HTTPRoute CRD — 트래픽 분할·헤더 기반 라우팅 등 세밀한 제어
- HealthCheckPolicy — `/health` 경로·포트를 직접 지정하여 정확한 헬스체크
- K8s 표준 API — 벤더 종속성 최소화

## ADR-008: 캐시 — Valkey (ch6.1)
**시점**: 2026-04 / **결정**: Valkey (Bitnami standalone) 채택 (vs Redis, Memcached)
**이유**:
- Redis fork — API 완전 호환, 라이선스 자유(BSD-3)
- bitnami Helm 차트 지원 — 단순 설치, Secret 자동 생성
- INCR 명령으로 분산 ID 카운터 구현
- e2-medium 환경에서 50m CPU로 구동 가능

## ADR-009: 시크릿 관리 — GKE Secret Manager CSI + WI (ch6.2)
**시점**: 2026-04 / **결정**: GKE managed CSI + Workload Identity 채택 (vs K8s Secret, Vault)
**이유**:
- 서비스 계정 키 불필요 — Workload Identity로 GCP IAM 직접 연동
- 파일 마운트 패턴 — 환경변수 노출 없이 /mnt/secrets/로 안전하게 전달
- Secret Manager 버전 관리 — 비밀번호 교체 시 새 버전 추가
- GKE managed addon — 오픈소스 CSI 설치 없이 클러스터 업데이트만으로 활성화

## ADR-010: 배포 전략 전환 — Canary (ch6.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Canary로 전환 (Blue/Green에서)
**이유**:
- 트래픽 점진 이동(20%→50%→80%→100%) — 운영 위험 최소화
- canaryService/stableService 분리 — 문제 발생 시 즉각 rollback
- pause 단계 — 각 단계에서 모니터링 후 다음 단계로
- Prometheus 메트릭 기반 자동 판단 확장 가능

## ADR-007: 배포 전략 — Blue/Green (ch5.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Blue/Green 채택 (vs Flagger, Istio)
**이유**:
- 즉각 롤백 — active/preview 서비스 전환으로 이전 버전 즉시 복귀
- ArgoCD 통합 — Rollout 상태를 ArgoCD UI에서 확인
- autoPromotionSeconds: 30 — 자동 승격으로 학습 환경에서 빠른 검증
- Deployment API 호환 — kubectl argo rollouts 플러그인으로 상태 확인
