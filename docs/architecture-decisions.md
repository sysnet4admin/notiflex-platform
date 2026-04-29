# Architecture Decision Records

## ADR-001: GitOps 도구 — ArgoCD (ch3.2)
**시점**: 2026-04 / **결정**: ArgoCD 채택 (vs Flux, Jenkins X)
**이유**:
- Web UI로 배포 상태 시각화 — 어느 버전이 배포되어 있는지 한눈에 확인
- e2-medium 2노드 환경에서 구동 가능 (경량)
- CNCF Graduated — 안정성이 검증된 생태계 표준
- automated sync + selfHeal로 GitOps 드리프트 자동 복구

## ADR-002: CI 도구 — GitHub Actions + WIF (ch3.4)
**시점**: 2026-04 / **결정**: GitHub Actions + Workload Identity Federation 채택
**이유**:
- 저장소 네이티브 — 별도 서버 없이 GitHub 저장소와 통합
- SA 키 생성이 조직 정책(`iam.disableServiceAccountKeyCreation`)으로 차단된 환경에서 WIF가 유일한 GCP 인증 수단
- 무료 사용량 확보 (private 저장소 500분/월)
- YAML 선언적 파이프라인 — ci.yaml도 git으로 버전 관리

## ADR-003: 메트릭 모니터링 — Prometheus + Grafana (ch4.2)
**시점**: 2026-04 / **결정**: kube-prometheus-stack 채택 (vs Datadog, New Relic)
**이유**:
- GKE 네이티브 — 추가 에이전트 없이 kube-state-metrics + node-exporter 통합
- 오픈소스 — 비용 없음, 데이터 외부 전송 없음
- kube-prometheus-stack Helm 차트로 Prometheus + Grafana + Alertmanager 일괄 설치
- Loki, Tempo와 같은 Grafana 생태계로 통합 관측 가능성 구성 가능

## ADR-004: 로깅 — Loki + Fluent Bit (ch4.3)
**시점**: 2026-04 / **결정**: Loki + Fluent Bit 채택 (vs ELK Stack)
**이유**:
- Grafana와 통합 — 메트릭·로그를 동일 대시보드에서 조회
- 인덱싱 없이 로그 저장 — ELK 대비 리소스 사용량 대폭 감소
- loki-stack 차트로 단순 설치 (최신 loki 차트 bucket 설정 불필요)
- e2-medium 2노드 환경에서도 안정적으로 구동

## ADR-005: 알림 — PrometheusRule + Alertmanager (ch4.4)
**시점**: 2026-04 / **결정**: PrometheusRule + Alertmanager 채택 (vs Grafana Alert)
**이유**:
- Prometheus와 네이티브 연동 — PromQL 기반 알림 규칙
- git 버전 관리 — YAML로 선언적 관리, 코드 리뷰 가능
- kube-prometheus-stack에 Alertmanager 포함 — 추가 설치 불필요
- `release: kube-prometheus` 라벨로 Prometheus Operator가 자동 로드

## ADR-006: 외부 트래픽 — GKE Gateway API (ch5.2)
**시점**: 2026-04 / **결정**: GKE Gateway API (`gke-l7-regional-external-managed`) 채택 (vs Ingress, NGINX, Istio)
**이유**:
- GKE 네이티브 — 별도 Ingress Controller 설치 불필요
- ch2.5에서 `--gateway-api=standard`로 이미 활성화 (추가 비용 없음)
- HTTPRoute + HealthCheckPolicy로 세밀한 트래픽 제어 가능
- Kubernetes 표준 API — 벤더 종속 없음

## ADR-007: 배포 전략 — Argo Rollouts Blue/Green (ch5.3)
**시점**: 2026-04 / **결정**: Argo Rollouts Blue/Green (autoPromotionSeconds: 30) 채택 (vs Flagger, Istio)
**이유**:
- ArgoCD와 동일 생태계 — 별도 서비스 메시 불필요
- YAML 선언적 — `strategy.blueGreen` 필드로 배포 전략을 git에 기록
- preview Pod으로 프로덕션 배포 전 사전 검증 가능
- autoPromotionSeconds로 자동 프로모션, 문제 시 즉시 rollback
