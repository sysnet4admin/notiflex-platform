# Notiflex Platform

B2B 알림 SaaS 플랫폼, Kubernetes 운영 환경 구축 실습 프로젝트

## 기술 스택
- **언어**: Go (표준 라이브러리만 사용, 외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지 (최소 크기)
- **인프라**: VM 바닐라 kubeadm (온프레미스) — 책 원고는 GKE 전제이지만 본 환경은 온프레
- **CI/CD**: GitHub Actions + ArgoCD (예정)

## 클러스터 정보 (온프레)
- **kubectl context**: `onprem-sysnet4admin_book_gitaiops`
- **kubeadm cluster (kubeconfig cluster 필드)**: `onprem-notiflex-cluster`
- **노드**: control-plane `cp-k8s` (192.168.1.10), worker `w1-k8s` ~ `w5-k8s` (192.168.1.101~105)
- **CNI**: Calico v3.31.2
- **StorageClass (default)**: `managed-nfs-storage` (CSI NFS)
- **LoadBalancer**: MetalLB v0.15.3 (IP 풀 192.168.1.11-99)
- **Gateway API**: GatewayClass `nginx` (NGINX Gateway Fabric v2.3.0)

## 컨테이너 레지스트리 (예정, ch2.6에서 결정)
- 책 원고는 `asia-northeast3-docker.pkg.dev/<project>/notiflex` (GCP Artifact Registry)
- 온프레 대체: 로컬 `registry:2` 또는 Harbor (ch2.6에서 선택)

## 행동 규칙
1. 명령 실행 전 현재 상태를 확인한다 (`kubectl get` 등)
2. 파일 변경 전 기존 내용을 먼저 읽는다
3. 에러 발생 시 원인을 분석하고 해결 방안을 제시한 뒤 진행한다
4. 매니페스트 작성 시 네임스페이스(`notiflex`)를 명시한다
5. **모든 `kubectl` 명령에 `--context onprem-sysnet4admin_book_gitaiops`를 반드시 지정한다**
   (잘못된 클러스터 대상 실행 방지 — 호스트에 다른 context 다수 존재)
6. 리소스 생성/삭제 전에는 영향 범위를 먼저 설명한다
7. 이미지 태그는 `latest`를 쓰지 않고 명시적 버전(`v0.1.0` 등)을 사용한다
8. 토큰, 키 그리고 비밀번호는 코드/매니페스트에 하드코딩하지 않는다
   (환경 변수, GitHub Secrets, Secret Manager/Vault 사용)
