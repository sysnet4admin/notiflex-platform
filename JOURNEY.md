# Notiflex 여정 기록 (onprem run-01)

책의 GKE 전제와 달리 본 여정은 **VM 바닐라 kubeadm 온프레 클러스터** (`onprem-sysnet4admin_book_gitaiops`)에서 진행 중. 플랫폼 커버리지 테스트 목적.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| pre  | cluster bootstrap | ✅ | 2026-04-22 | Vagrant + kubeadm (cp + w1~w5) |
| ch2  | 2.2 설치 확인 | ✅ | 2026-04-22 | 플랫폼 중립, 이미 세팅됨 |
| ch2  | 2.3 gcloud 설정 | ⚠️ skip | 2026-04-22 | GKE 특화 전체 건너뜀 |
| ch2  | 2.4 GitHub 저장소 | ✅ | 2026-04-22 | CLAUDE.md onprem 조정 |
| ch2  | 2.5 클러스터 | ⚠️ | 2026-04-22 | 사전 부트스트랩, 상태 확인만 |
| ch2  | 2.6 빌드/배포 | ✅ | 2026-04-22 | 로컬 registry + aarch64 build |
| ch2  | 2.7 첫 커밋 | ⬜ | | |
| ch3  | 3.2 GitOps 도구 | ⬜ | | |
| ch3  | 3.3 기능 추가 | ⬜ | | |
| ch3  | 3.4 CI | ⬜ | | |
| ch3  | 3.5 CI-CD 연결 | ⬜ | | |
| ch4  | 4.2 메트릭 모니터링 | ⬜ | | |
| ch4  | 4.3 로그 수집 | ⬜ | | |
| ch4  | 4.4 알림 | ⬜ | | |
| ch5  | 5.2 트래픽 관리 | ⬜ | | |
| ch5  | 5.3 무중단 배포 | ⬜ | | |
| ch6  | 6.1 캐시 | ⬜ | | |
| ch6  | 6.2 시크릿 관리 | ⬜ | | |
| ch6  | 6.3 Canary 전환 | ⬜ | | |
| ch7  | 7.2 멀티 노드풀 | ⬜ | | |
| ch7  | 7.3 App of Apps | ⬜ | | |
| ch7  | 7.4 멀티테넌시 | ⬜ | | |
| ch8  | 8.1 메시징 | ⬜ | | |
| ch8  | 8.2 트레이싱 | ⬜ | | |
| ch8  | 8.3 CronJob | ⬜ | | |
| ch9  | 9.1 저장소 분석 | ⬜ | | |
| ch9  | 9.2 회고 | ⬜ | | |
| ch9  | 9.3 온보딩 문서 | ⬜ | | |
| ch9  | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9  | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| 컨테이너 레지스트리 | 로컬 `registry:2` Pod + NodePort 30500 | Harbor, 호스트 docker registry, ctr import | 가장 단순, 클러스터 안에 있어 CI 흐름 일관성 |
| CNI | Calico v3.31.2 | Cilium, Flannel | 저자 B.001/U 원본 검증 스택 그대로 |
| LoadBalancer | MetalLB v0.15.3 (L2) | kube-vip | B.001/U 원본 |
| Gateway API | NGINX Gateway Fabric v2.3.0 | Cilium Gateway, Envoy Gateway | B.001/U 원본 |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25-alpine | 2026-04-22 ch2.6 |
| Notiflex 이미지 | v0.1.0 (arm64) | 2026-04-22 ch2.6 |
| Kubernetes | v1.35.0 | kubeadm |
| containerd | 2.2.2 | |

## 현재 리소스

| 노드 | 머신 사양 | 역할 | 주요 워크로드 |
|------|----------|------|-------------|
| cp-k8s | 2 vCPU / 4 GB | control-plane | etcd, apiserver, scheduler, controller-manager |
| w1-k8s | 2 vCPU / 4 GB | worker | (빈 상태) |
| w2-k8s | 2 vCPU / 4 GB | worker | notiflex-api |
| w3-k8s | 2 vCPU / 4 GB | worker | (빈 상태) |
| w4-k8s | 2 vCPU / 8 GB | worker | notiflex-api, registry |
| w5-k8s | 2 vCPU / 2 GB | worker | (빈 상태) |

## 엣지 케이스 / 트러블슈팅 이력

- **kubeadm `--cluster-name` CLI 플래그 미지원** — ClusterConfiguration YAML 없이는 cluster name 변경 불가. 기본 `kubernetes` 유지.
- **kubeconfig merge 시 이름 충돌** — 기존 kubeconfig의 `kubernetes` cluster / `kubernetes-admin` user와 충돌 → `onprem-notiflex-cluster`/`onprem-notiflex-admin`으로 유니크화 후 병합.
- **Colima 기본 x86_64 vs VM arm64 아키텍처 불일치** — `exec format error`. Colima를 aarch64로 재생성.
- **`docker build` 시 캐시된 amd64 이미지 우선** — `--platform=linux/arm64` 명시 필요.
- **containerd `IfNotPresent` 정책 + digest 변경** — 태그만으로 재pull 안 됨. 각 노드에서 `crictl rmi`로 수동 정리.
