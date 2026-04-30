# Architecture Decision Records

## ADR-001: GitOps 도구로 ArgoCD 채택 (3장)
**시점**: 2026-04 / **결정**: GitOps 도구로 ArgoCD를 채택하고 Flux는 사용하지 않음.
**이유**:
- Web UI를 통해 배포 상태를 실시간으로 확인할 수 있어 학습 과정에 유리합니다.
- Application CRD를 통해 선언적으로 Git 경로와 클러스터 네임스페이스를 관리할 수 있습니다.
- Self-Heal 기능으로 클러스터의 상태를 Git과 일관되게 유지합니다.
- GKE Standard와 호환되며 e2-medium 노드에서 구동 가능합니다.

## ADR-002: CI 도구로 GitHub Actions 채택 (3장)
**시점**: 2026-04 / **결정**: CI 도구로 GitHub Actions를 채택하고 Cloud Build, GitLab CI, Jenkins는 사용하지 않음.
**이유**:
- 코드 저장소인 GitHub에 내장되어 있어 별도의 서버 설치나 관리가 필요 없습니다.
- YAML 파일을 통해 선언적으로 파이프라인을 정의할 수 있습니다.
- 프라이빗 저장소에 월 2,000분의 무료 크레딧을 제공하여 비용 효율적입니다.
- `google-github-actions/auth` 액션을 통해 GCP 인증을 간편하게 처리할 수 있습니다.

## ADR-003: 메트릭 모니터링으로 Prometheus+Grafana 채택 (4장)
**시점**: 2026-04 / **결정**: 메트릭 수집 및 시각화 도구로 kube-prometheus-stack(Prometheus+Grafana)을 채택하고 Datadog, Google Cloud Monitoring은 사용하지 않음.
**이유**:
- Kubernetes 모니터링의 사실상 표준(CNCF Graduated)으로, 방대한 커뮤니티와 자료를 보유하고 있습니다.
- e2-medium 노드에서도 충분히 운영 가능한 가벼운 리소스 요구사항을 가집니다.
- Grafana 대시보드를 통해 로그(Loki), 트레이스(Tempo)와 통합된 뷰를 제공하여 관측 가능성을 높입니다.
- Helm 차트로 제공되어 6개의 관련 컴포넌트를 한 번에 설치하고 관리할 수 있습니다.

## ADR-004: 로그 수집으로 Loki+Fluent Bit 채택 (4장)
**시점**: 2026-04 / **결정**: 로그 수집 및 저장 도구로 Loki와 Fluent Bit을 채택하고 ELK Stack은 사용하지 않음.
**이유**:
- ELK Stack 대비 현저히 낮은 메모리(128Mi)를 사용하여 e2-medium 노드 환경에 적합합니다.
- Prometheus와 동일한 라벨 기반의 쿼리(LogQL)를 사용하여 학습 곡선이 완만합니다.
- Grafana에 통합되어 메트릭과 로그를 동일한 UI에서 함께 조회할 수 있습니다.
- 로그 내용 전체를 인덱싱하는 대신 라벨만 인덱싱하여 저장 비용이 낮습니다.

## ADR-005: 외부 트래픽 관리로 Gateway API 채택 (5장)
**시점**: 2026-04 / **결정**: 외부 트래픽 관리 방식으로 Gateway API를 채택하고 Ingress는 사용하지 않음.
**이유**:
- Kubernetes의 차세대 공식 트래픽 관리 표준 API입니다.
- GKE에서 네이티브로 지원하여 별도의 Controller 설치가 필요 없습니다.
- Gateway(인프라)와 HTTPRoute(애플리케이션)로 역할과 책임이 분리됩니다.
- HTTPRoute의 backendRefs를 통해 Blue/Green, Canary 배포와 쉽게 연동할 수 있습니다.

## ADR-006: 무중단 배포 전략으로 Argo Rollouts 채택 (5장)
**시점**: 2026-04 / **결정**: 무중단 배포 전략으로 Argo Rollouts를 채택하고 Kubernetes 기본 Rolling Update는 사용하지 않음.
**이유**:
- ArgoCD와 통합되어 UI에서 Rollout 진행 상태를 시각적으로 확인할 수 있습니다.
- YAML 선언적 방식으로 Blue/Green, Canary 등 고급 배포 전략을 GitOps 기반으로 관리할 수 있습니다.
- 5장의 Blue/Green 전략에서 6장의 Canary 전략으로 점진적으로 발전시키기 용이합니다.
- `kubectl argo rollouts` 플러그인을 통해 배포 상태를 실시간으로 모니터링할 수 있습니다.

## ADR-007: Valkey 기반 캐시 도입 (6장)
**시점**: 2026-04 / **결정**: Valkey를 인메모리 캐시 솔루션으로 채택.
**이유**:
- **성능**: Valkey는 Redis에서 파생된 고성능 인메모리 데이터 스토어로, 빠른 읽기/쓰기 속도를 제공하여 API 응답 시간 단축에 매우 효과적입니다.
- **단순성**: `INCR`과 같은 간단한 명령어를 지원하여 캐시 로직 구현이 직관적이고 개발 복잡성이 낮습니다.
- **확장성**: 단일 노드로 시작하여 향후 클러스터 모드로 확장이 용이하며, 데이터 샤딩을 통해 대규모 트래픽을 수용할 수 있습니다.
- **생태계**: Redis와 호환되는 클라이언트 라이브러리와 도구를 그대로 사용할 수 있어 기존 개발 생태계를 활용하기 유리합니다.

## ADR-008: Secret Manager CSI Driver for GKE 도입 (6장)
**시점**: 2026-04 / **결정**: Google Secret Manager와 Secret Manager CSI Driver, Workload Identity를 연동하여 시크릿 관리.
**이유**:
- **보안 강화**: 민감한 정보가 Git 저장소나 Kubernetes etcd에 평문으로 저장되지 않으므로 유출 위험이 원천적으로 차단됩니다.
- **중앙화된 관리**: 모든 Secret을 Google Secret Manager에서 통합 관리하므로, 버전 관리, 접근 제어(IAM), 감사 로깅이 용이합니다.
- **자동화 및 편의성**: 애플리케이션 코드 변경 없이 Secret을 파일처럼 사용할 수 있으며, GitOps 워크플로우와 자연스럽게 통합됩니다.
- **GCP 생태계 활용**: Google Cloud의 강력한 IAM 및 보안 기능을 최대로 활용하여 Kubernetes 클러스터의 보안 수준을 높일 수 있습니다.

## ADR-009: Canary 배포 전략으로 전환 (6장)
**시점**: 2026-04 / **결정**: Argo Rollouts의 배포 전략을 Blue/Green에서 Canary로 전환.
**이유**:
- **위험 감소**: 소수의 사용자에게만 신규 버전을 노출시키므로, 장애 발생 시 전체 서비스에 미치는 영향을 최소화할 수 있습니다.
- **신뢰성 있는 검증**: 실제 운영 트래픽을 통해 신규 버전의 성능과 안정성을 객관적인 지표로 검증할 수 있습니다.
- **자동화된 롤백**: 모니터링 시스템과 연동하여 문제 발생 시 사람의 개입 없이 자동으로 이전 버전으로 롤백하여 서비스 안정성을 확보합니다.
- **유연한 트래픽 제어**: 서비스의 특성과 중요도에 따라 트래픽 전환 비율과 검증 시간을 유연하게 조절할 수 있습니다.
