# ADR-011: ArgoCD App of Apps 패턴 도입

## 상태

**채택 (Accepted)**

## 맥락 (Context)

-   애플리케이션의 규모가 커지고 마이크로서비스가 늘어나면서, 각 환경(smb, enterprise, monitoring 등)별 ArgoCD `Application` 리소스를 수동으로 관리하는 부담이 증가했습니다.
-   새로운 애플리케이션이나 환경을 추가할 때마다 ArgoCD에 `Application` CRD를 반복적으로 등록해야 하는 번거로움이 있습니다.
-   전체 애플리케이션들의 배포 상태를 한눈에 파악하고 관리하기 위한 중앙화된 접근 방식이 필요합니다.

## 결정 (Decision)

**ArgoCD의 'App of Apps' 패턴을 도입하여 애플리케이션 관리를 자동화하고 중앙화합니다.**

-   최상위 `root-app`이라는 ArgoCD `Application`을 생성합니다. 이 `root-app`은 Git 저장소의 특정 디렉터리(e.g., `argocd/apps`)를 바라봅니다.
-   개별 애플리케이션(notiflex-smb, notiflex-enterprise, notiflex-monitoring 등)의 ArgoCD `Application` 매니페스트를 `argocd/apps` 디렉터리 안에 각각의 YAML 파일로 정의합니다.
-   `root-app`은 `argocd/apps` 디렉터리의 변경을 감지하고, 그 안에 정의된 모든 `Application`들을 자동으로 ArgoCD에 생성, 동기화, 관리합니다.
-   결과적으로, `root-app` 하나만 관리하면 그에 속한 수많은 자식 애플리케이션들이 자동으로 관리되는 구조가 됩니다.

## 이유 (Rationale)

-   **관리 효율성**: 개별 애플리케이션을 일일이 관리할 필요 없이, Git 저장소의 디렉터리 구조를 통해 전체 애플리케이션 포트폴리오를 선언적으로 관리할 수 있습니다.
-   **자동화 및 확장성**: 새로운 애플리케이션을 추가할 때, Git 저장소에 `Application` YAML 파일을 추가하고 push하기만 하면 ArgoCD에 자동으로 배포가 구성됩니다.
-   **중앙화된 가시성**: ArgoCD UI에서 `root-app`을 통해 모든 하위 애플리케이션들의 배포 상태와 관계를 트리 구조로 한눈에 파악할 수 있습니다.
-   **GitOps 원칙 준수**: 모든 애플리케이션의 구성 정보가 Git에 단일 진실 공급원(Single Source of Truth)으로 관리되어, 버전 관리, 변경 추적, 롤백이 용이해집니다.
