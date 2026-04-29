# Architecture Decision Records

## ADR-001: 배포 자동화는 ArgoCD (3장)
**시점**: 2026-04 / **결정**: GitOps 도구로 ArgoCD를 채택한다. Flux는 채택하지 않는다.
**이유**:
- Web UI로 배포 상태를 실시간 시각화할 수 있어 학습 흐름에 적합하다
- Application CRD로 Git 경로와 배포 대상을 선언적으로 관리할 수 있다
- selfHeal로 클러스터 수동 변경 드리프트를 자동 복구할 수 있다
- GKE Standard 및 e2-medium 리소스 환경에서 안정적으로 운영 가능하다

## ADR-002: CI 도구는 GitHub Actions (3장)
**시점**: 2026-04 / **결정**: CI는 GitHub Actions 기반으로 구성한다. Cloud Build 중심 방식은 기본 경로로 채택하지 않는다.
**이유**:
- 코드 저장소와 CI가 같은 플랫폼이라 별도 CI 서버 운영이 필요 없다
- YAML 기반 워크플로로 파이프라인을 저장소에서 선언적으로 관리할 수 있다
- 무료 크레딧 범위에서 실습 워크로드를 충분히 소화할 수 있다
- `google-github-actions/auth`로 GCP 인증 연동이 단순하다

## ADR-003: 메트릭은 Prometheus와 Grafana (4장)
**시점**: 2026-04 / **결정**: 메트릭 수집/시각화는 kube-prometheus-stack(Prometheus+Grafana+Alertmanager)으로 통합한다. SaaS 모니터링은 도입하지 않는다.
**이유**:
- Kubernetes 모니터링 표준 스택으로 학습/운영 모두에서 재사용성이 높다
- Helm 번들로 검증된 구성요소를 한 번에 설치해 초기 구축 속도가 빠르다
- Grafana 기반 시각화로 이후 로그/트레이싱 도구와 통합이 용이하다
- 구독형 SaaS 비용 없이 실습 환경에서 자체 운영 가능하다

## ADR-004: 로그는 Loki와 Fluent Bit (4장)
**시점**: 2026-04 / **결정**: 중앙 로그 수집은 Loki+Fluent Bit를 사용한다. ELK Stack은 채택하지 않는다.
**이유**:
- e2-medium 제약에서 경량 리소스 사용량으로 운영 가능하다
- Grafana와 바로 통합되어 메트릭과 로그를 같은 UI에서 분석할 수 있다
- 라벨 기반 인덱싱으로 저장 비용과 운영 복잡도를 낮출 수 있다
- DaemonSet 기반 수집으로 노드 전반 로그를 일관되게 수집할 수 있다

## ADR-005: 알림은 PrometheusRule과 Alertmanager (4장)
**시점**: 2026-04 / **결정**: 알림 규칙과 라우팅은 PrometheusRule+Alertmanager로 구성한다. Grafana UI 중심 알림은 기본 경로로 채택하지 않는다.
**이유**:
- 알림 규칙을 YAML로 Git에 버전 관리하여 GitOps 흐름을 유지할 수 있다
- kube-prometheus-stack에 포함된 구성요소를 그대로 활용해 추가 설치가 필요 없다
- 실무 표준 스택으로 재현성과 팀 공유성이 높다
- Alertmanager의 그룹핑/억제/라우팅으로 다단계 알림 정책을 구현할 수 있다

## ADR-006: 외부 진입점은 Gateway API (5장)
**시점**: 2026-04 / **결정**: 외부 트래픽 관리는 GKE Gateway API로 구성한다. Ingress Controller는 사용하지 않는다.
**이유**:
- Kubernetes 공식 차세대 표준으로 Ingress 대비 확장성이 높다
- GKE 네이티브 지원으로 별도 컨트롤러 설치 없이 운영할 수 있다
- Gateway와 HTTPRoute 분리로 인프라/애플리케이션 책임 경계를 명확히 할 수 있다
- HTTPRoute 및 HealthCheckPolicy로 무중단 라우팅 안정성을 확보할 수 있다

## ADR-007: 무중단 배포는 Argo Rollouts Blue/Green (5장)
**시점**: 2026-04 / **결정**: 무중단 배포 전략은 Argo Rollouts Blue/Green으로 운영한다. 기본 Deployment RollingUpdate와 Canary는 현 단계에서 채택하지 않는다.
**이유**:
- active/preview 분리로 트래픽 전환 시점을 명시적으로 통제할 수 있다
- preview 검증 후 자동 승격(`autoPromotionSeconds`)으로 운영 리스크를 낮출 수 있다
- ArgoCD와 같은 생태계여서 GitOps 운영 및 상태 가시성이 높다
- 현재 환경에서는 Blue/Green의 리소스 오버헤드가 감당 가능하다
