# Architecture Decision Records

## ADR-001: GitOps 도구로 ArgoCD 채택 (3장)
**시점**: 2026-04 / **결정**: GitOps 도구로 ArgoCD를 채택하고 Flux는 사용하지 않음.
**이유**:
- Web UI를 통해 배포 상태를 실시간으로 확인할 수 있어 학습 과정에 유리합니다.
- Application CRD를 통해 선언적으로 Git 경로와 클러스터 네임스페이스를 관리할 수 있습니다.
- Self-Heal 기능으로 클러스터의 상태를 Git과 일관되게 유지합니다.
- GKE Standard와 호환되며 e2-medium 노드에서 구동 가능합니다.

## ADR-002: CI 도구로 GitHub Actions 채택 (3장)
**시점**: 2026-04 / **결정**: CI는 GitHub Actions 기반으로 구성한다. Cloud Build 중심 방식은 기본 경로로 채택하지 않는다.
**이유**:
- 코드 저장소와 CI가 같은 플랫폼이라 별도 CI 서버 운영이 필요 없습니다.
- YAML 기반 워크플로로 파이프라인을 저장소에서 선언적으로 관리할 수 있습니다.
- 무료 크레딧 범위에서 실습 워크로드를 충분히 소화할 수 있습니다.
- `google-github-actions/auth`로 GCP 인증 연동이 단순합니다.

## ADR-003: 메트릭 모니터링으로 Prometheus+Grafana 채택 (4장)
**시점**: 2026-04 / **결정**: 메트릭 수집/시각화는 kube-prometheus-stack(Prometheus+Grafana+Alertmanager)으로 통합한다. SaaS 모니터링은 도입하지 않는다.
**이유**:
- Kubernetes 모니터링 표준 스택으로 학습/운영 모두에서 재사용성이 높습니다.
- Helm 번들로 검증된 구성요소를 한 번에 설치해 초기 구축 속도가 빠릅니다.
- Grafana 기반 시각화로 이후 로그/트레이싱 도구와 통합이 용이합니다.
- 구독형 SaaS 비용 없이 실습 환경에서 자체 운영 가능합니다.

## ADR-004: 로그 수집으로 Loki+Fluent Bit 채택 (4장)
**시점**: 2026-04 / **결정**: 중앙 로그 수집은 Loki+Fluent Bit를 사용한다. ELK Stack은 채택하지 않는다.
**이유**:
- e2-medium 제약에서 경량 리소스 사용량으로 운영 가능합니다.
- Grafana와 바로 통합되어 메트릭과 로그를 같은 UI에서 분석할 수 있습니다.
- 라벨 기반 인덱싱으로 저장 비용과 운영 복잡도를 낮출 수 있습니다.
- DaemonSet 기반 수집으로 노드 전반 로그를 일관되게 수집할 수 있습니다.

## ADR-005: 알림은 PrometheusRule + Alertmanager (4장)
**시점**: 2026-04 / **결정**: 알림 규칙과 라우팅은 PrometheusRule + Alertmanager로 구성한다. Grafana UI 중심 알림은 기본 경로로 채택하지 않는다.
**이유**:
- 알림 규칙을 YAML로 Git에 버전 관리하여 GitOps 흐름을 유지할 수 있다.
- kube-prometheus-stack에 포함된 구성요소를 그대로 활용해 추가 설치가 필요 없다.
- 실무 표준 스택으로 재현성과 팀 공유성이 높다.
- Alertmanager의 그룹핑/억제/라우팅으로 다단계 알림 정책을 구현할 수 있다.

## ADR-006: 외부 트래픽 관리로 Gateway API 채택 (5장)
**시점**: 2026-04 / **결정**: 외부 트래픽 관리는 GKE Gateway API로 구성한다. Ingress Controller는 사용하지 않는다.
**이유**:
- Kubernetes 공식 차세대 표준으로 Ingress 대비 확장성이 높습니다.
- GKE 네이티브 지원으로 별도 컨트롤러 설치 없이 운영할 수 있습니다.
- Gateway와 HTTPRoute 분리로 인프라/애플리케이션 책임 경계를 명확히 할 수 있습니다.
- HTTPRoute 및 HealthCheckPolicy로 무중단 라우팅 안정성을 확보할 수 있습니다.

## ADR-007: 무중단 배포 전략으로 Argo Rollouts 채택 (5장)
**시점**: 2026-04 / **결정**: 무중단 배포 전략으로 Argo Rollouts를 채택하고 Kubernetes 기본 Rolling Update는 사용하지 않는다.
**이유**:
- ArgoCD와 통합되어 UI에서 Rollout 진행 상태를 시각적으로 확인할 수 있습니다.
- YAML 선언적 방식으로 Blue/Green, Canary 등 고급 배포 전략을 GitOps 기반으로 관리할 수 있습니다.
- 5장의 Blue/Green 전략에서 6장의 Canary 전략으로 점진적으로 발전시키기 용이합니다.
- `kubectl argo rollouts` 플러그인을 통해 배포 상태를 실시간으로 모니터링합니다.

## ADR-008: 캐시는 Valkey (6장)
**시점**: 2026-04 / **결정**: 애플리케이션 공유 카운터/캐시는 Valkey를 채택한다. Redis OSS와 Memcached는 기본 경로로 채택하지 않는다.
**이유**:
- Redis API 호환으로 앱 코드 변경을 최소화할 수 있습니다.
- standalone 배포로 e2-medium 실습 환경에서도 운영 복잡도를 낮출 수 있습니다.
- 오픈 거버넌스 기반(Valkey)으로 장기 운영 리스크를 줄일 수 있습니다.
- `/id` 카운터 시나리오를 빠르게 검증하기에 적합합니다.

## ADR-009: 시크릿은 GKE Secret Manager CSI + Workload Identity (6장)
**시점**: 2026-04 / **결정**: 시크릿 주입은 GKE Secret Manager CSI Driver와 Workload Identity를 조합해 사용한다. Kubernetes Secret 직접 저장과 Vault 도입은 현재 단계에서 채택하지 않는다.
**이유**:
- 시크릿 원문을 클러스터 외부(Google Secret Manager)에서 중앙 관리할 수 있습니다.
- Pod에는 CSI 파일 마운트로만 노출되어 환경 변수 직접 주입 대비 노출면이 작습니다.
- Workload Identity로 정적 키 배포 없이 권한을 부여할 수 있습니다.
- GKE 네이티브 기능 조합으로 초기 운영 복잡도와 통합 비용이 낮습니다.

## ADR-010: 배포 전략은 Canary로 전환 (6장)
**시점**: 2026-04 / **결정**: 배포 전략을 Argo Rollouts Blue/Green에서 Canary로 전환한다. Blue/Green은 기본 전략으로 유지하지 않는다.
**이유**:
- 20→50→80→100 단계별 트래픽 전환으로 리스크를 점진적으로 관찰할 수 있습니다.
- 각 단계 `pause`로 품질 지표/알림 확인 시간을 확보할 수 있습니다.
- Argo Rollouts/ArgoCD 생태계를 유지하면서 전략만 변경해 운영 일관성을 확보합니다.
- 이상 징후 발생 시 100% 전환 전에 제어하기 유리합니다.

## ADR-011: 워크로드별 노드풀 분리 (7장)
**시점**: 2026-04 / **결정**: GKE 멀티 노드풀(`api-pool`, `worker-pool`, `ops-pool`)과 `nodeSelector(cloud.google.com/gke-nodepool)`로 워크로드를 역할별로 분리한다.
**이유**:
- GKE 노드풀 생성 시 자동 부여되는 라벨을 그대로 사용해 설정 오류 가능성을 낮출 수 있다.
- 학습 환경에서 nodeSelector만으로 충분하며 taint보다 이해하기 쉽다.
- API/Worker/Ops 워크로드를 물리적으로 분리해 리소스 간섭을 줄일 수 있다.
- 이후 taint/toleration으로 확장 시 기반 토폴로지를 재사용할 수 있다.

## ADR-012: 다중 앱 관리는 App of Apps (7장)
**시점**: 2026-04 / **결정**: ArgoCD `root-app`과 `argocd/apps/` 디렉터리 기반 App of Apps 패턴을 채택한다. ApplicationSet은 현 단계 기본 경로로 채택하지 않는다.
**이유**:
- 하위 Application 선언을 Git 디렉터리 기준으로 일괄 관리할 수 있다.
- sync-wave를 통해 인프라/플랫폼/앱의 설치 순서를 명시적으로 통제할 수 있다.
- 기존 ArgoCD/GitOps 흐름을 바꾸지 않고 구조만 확장할 수 있다.
- ApplicationSet보다 단순해 학습 초기 단계에 적합하다.

## ADR-013: 멀티테넌시는 Namespace 분리 (7장)
**시점**: 2026-04 / **결정**: Namespace 기반 테넌트 분리(`enterprise`)와 전용 ArgoCD Application을 채택한다. vCluster 및 테넌트별 별도 클러스터는 현재 채택하지 않는다.
**이유**:
- 추가 클러스터 없이 비용을 통제하면서 테넌트 격리를 적용할 수 있다.
- App of Apps 구조에서 테넌트별 배포 단위를 독립적으로 운영할 수 있다.
- cross-namespace DNS로 공유 Valkey에 접근하는 패턴을 학습할 수 있다.
- 학습 단계에서 운영 복잡도를 크게 늘리지 않고 멀티테넌시 패턴을 경험할 수 있다.

## ADR-014: 메시징은 Strimzi 기반 Kafka (8장)
**시점**: 2026-04 / **결정**: 이벤트 스트리밍은 Strimzi Operator + Kafka 4.1.0(KRaft)을 채택한다. RabbitMQ, NATS는 현재 채택하지 않는다.
**이유**:
- 이벤트 스트림 기반 토픽 모델이 Notiflex 알림 워크로드 확장에 적합하다.
- Strimzi로 Kubernetes 네이티브 선언형 운영(CRD/Helm)을 유지할 수 있다.
- KRaft 단일 브로커 구성으로 실습 환경 리소스 제약에서 운영 복잡도를 낮출 수 있다.
- 업계 표준 도구로 학습 투자 대비 효과가 높다.

## ADR-015: 분산 트레이싱은 Tempo + OpenTelemetry (8장)
**시점**: 2026-04 / **결정**: 트레이싱 스택은 Grafana Tempo(monolithic)와 OpenTelemetry Go SDK(OTLP gRPC)를 채택한다. Jaeger와 Zipkin은 현재 채택하지 않는다.
**이유**:
- Grafana/Loki/Prometheus와 동일한 관측 스택으로 통합 운영이 쉽다.
- 단일 바이너리 Tempo 구성이 학습 클러스터 리소스에서 운영 부담이 작다.
- OTel 표준 계측으로 향후 벤더/백엔드 변경 시 이식성을 확보할 수 있다.
- OTLP gRPC 경로를 사용해 표준 방식으로 트레이스를 전송할 수 있다.

## ADR-016: 주기 작업은 Kubernetes CronJob (8장)
**시점**: 2026-04 / **결정**: 정기 헬스체크 배치는 Kubernetes CronJob으로 운영한다. Argo Workflows, 외부 VM cron은 현재 채택하지 않는다.
**이유**:
- 단순 반복 작업에는 K8s 기본 리소스만으로 가장 빠르게 구현할 수 있다.
- GitOps(YAML/ArgoCD) 흐름 안에서 정의/배포/이력 추적을 일관되게 유지할 수 있다.
- `ops-pool` 스케줄링으로 운영성 워크로드를 애플리케이션 노드와 분리할 수 있다.
- 실패 재시도/히스토리 보존 등 기본 배치 제어 기능으로 요구사항을 충족한다.
