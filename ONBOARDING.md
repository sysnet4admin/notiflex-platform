# Notiflex Platform 온보딩 가이드

이 문서는 Notiflex 플랫폼의 아키텍처, 사용법, 개발 워크플로우를 안내합니다.

## 1. 아키텍처 및 디렉터리 구조

`notiflex-platform` 저장소는 다음과 같은 구조를 가집니다.

- **`app/`**: Go로 작성된 애플리케이션 소스 코드
- **`k8s/`**: Kubernetes 매니페스트 (Namespace, Rollout, Service 등)
- **`argocd/`**: ArgoCD 애플리케이션 정의 (App of Apps 패턴)
- **`.github/workflows/`**: GitHub Actions CI 워크플로우
- **`docs/`**: 아키텍처 결정 문서 (ADR)
- **`.claude/`, `claude-context/`**: AI 에이전트용 컨텍스트 및 명령어 가이드

## 2. 배포 플로우 (GitOps with AI)

1.  **Git Push**: 로컬에서 코드 변경 후 `main` 브랜치에 push합니다.
2.  **CI (GitHub Actions)**: push를 트리거하여 CI 파이프라인이 실행됩니다.
    - 코드 빌드 및 테스트
    - Docker 이미지 빌드 및 Google Artifact Registry에 푸시
    - `k8s/`의 매니페스트 파일에 새 이미지 태그 업데이트
3.  **CD (ArgoCD)**: Git 저장소의 매니페스트 변경을 감지하고 클러스터에 자동 동기화합니다.
4.  **Canary 배포 (Argo Rollouts)**: 변경 사항은 카나리 방식으로 점진적으로 배포되며, 문제가 발생하면 자동으로 롤백됩니다.

## 3. 주요 서비스 접근 방법

**ArgoCD UI, Grafana 대시보드, API 엔드포인트** 등 주요 서비스 접근 정보는 클러스터 설정에 따라 달라지므로, 담당자에게 문의하거나 관련 스크립트를 통해 확인하세요.

## 4. 자주 묻는 Q&A

**Q: 카나리 배포를 중단하고 싶으면 어떻게 하나요?**
A: Argo Rollouts CLI 또는 ArgoCD UI를 통해 `abort` 명령을 실행하여 롤백할 수 있습니다.

**Q: 특정 서비스의 로그를 어떻게 검색하나요?**
A: Grafana의 Loki 데이터소스를 사용하여 LogQL 쿼리 (예: `{namespace="smb", app="notiflex"}`)로 검색합니다.

**Q: 분산 트레이싱 데이터는 어떻게 확인하나요?**
A: Grafana의 Tempo 데이터소스를 통해 TraceID로 특정 요청의 전체 흐름을 추적할 수 있습니다.

**Q: 새로운 Kafka 토픽을 추가하려면 어떻게 해야 하나요?**
A: `k8s/kafka/` 디렉터리에 Strimzi `KafkaTopic` CRD를 정의하고 Git에 push합니다.

**Q: 새로운 테넌트(고객사)를 추가하고 싶습니다.**
A: `argocd/apps/`에 새 테넌트용 ArgoCD 애플리케이션을 정의하고, `k8s/`에 해당 테넌트의 리소스를 추가합니다.

**Q: 시스템 알림은 어디서 확인하나요?**
A: Prometheus Alertmanager UI에서 현재 발생한 알림과 이력을 확인할 수 있습니다.
