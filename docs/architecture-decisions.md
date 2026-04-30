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

## ADR-003: 외부 트래픽 관리로 Gateway API 채택 (5장)
**시점**: 2026-04 / **결정**: 외부 트래픽 관리 방식으로 Gateway API를 채택하고 Ingress는 사용하지 않음.
**이유**:
- Kubernetes의 차세대 공식 트래픽 관리 표준 API입니다.
- GKE에서 네이티브로 지원하여 별도의 Controller 설치가 필요 없습니다.
- Gateway(인프라)와 HTTPRoute(애플리케이션)로 역할과 책임이 분리됩니다.
- HTTPRoute의 backendRefs를 통해 Blue/Green, Canary 배포와 쉽게 연동할 수 있습니다.

## ADR-004: 무중단 배포 전략으로 Argo Rollouts 채택 (5장)
**시점**: 2026-04 / **결정**: 무중단 배포 전략으로 Argo Rollouts를 채택하고 Kubernetes 기본 Rolling Update는 사용하지 않음.
**이유**:
- ArgoCD와 통합되어 UI에서 Rollout 진행 상태를 시각적으로 확인할 수 있습니다.
- YAML 선언적 방식으로 Blue/Green, Canary 등 고급 배포 전략을 GitOps 기반으로 관리할 수 있습니다.
- 5장의 Blue/Green 전략에서 6장의 Canary 전략으로 점진적으로 발전시키기 용이합니다.
- `kubectl argo rollouts` 플러그인을 통해 배포 상태를 실시간으로 모니터링할 수 있습니다.
