# Architecture Decision Records

## ADR-001: Gateway API for external traffic (ch5.2)

**Date**: 2026-04-24
**Status**: Accepted

### Context
notiflex-api를 외부에 노출해야 함. Ingress, Service LoadBalancer, Gateway API, Istio 후보.

### Decision
**Gateway API (GKE gke-l7-regional-external-managed)** 채택.

### Rationale
- K8s SIG 차세대 표준 (Ingress 후속), 향후 5년 안정적
- 역할 분리 (인프라팀: Gateway / 앱팀: HTTPRoute)
- GKE 네이티브 (`--gateway-api=standard`로 ch2.5에서 활성화)
- HealthCheckPolicy로 /health 명시적 지정 가능

### Consequences
- proxy-only-subnet 1회 생성 필요 (ch5.2 트러블슈팅 포인트)
- 미래 ch6.3 Canary와 자연스러운 통합

---

## ADR-002: Argo Rollouts Blue/Green (ch5.3)

**Date**: 2026-04-24
**Status**: Accepted

### Context
무중단 + 안전한 배포 필요. K8s Rolling Update, Argo Rollouts, Flagger, Spinnaker 후보.

### Decision
**Argo Rollouts Blue/Green (autoPromotionEnabled: false)** 채택.

### Rationale
- ArgoCD 같은 Argo 생태계 (UI/CLI 일관)
- preview Service로 green 단독 검증 후 수동 promote
- Rolling Update 대비 명시적 cutover, 즉시 롤백
- ch6.3 Canary로 자연스러운 진화 (전략만 변경)

### Consequences
- replicas 2배 일시적 (BG 진행 중)
- preview Service 추가 (k8s/smb/service-preview.yaml)
