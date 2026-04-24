# Notiflex 아키텍처 (run-53 ch6 시점)

## 컴포넌트 토폴로지

```
[외부 사용자]
     ↓
[Gateway API (gke-l7)]  ← ch5.2
     ↓ /
[HTTPRoute notiflex-route]  ← ch5.2
     ↓
[Service notiflex-api (ClusterIP)]
     ↓
[Rollout notiflex-api (Canary 20/50/80)]  ← ch5.3 → ch6.3
     ↓
[Pod notiflex-api × 2] (image: notiflex/api:v0.3.0)
     ↓ INCR (TCP 6379)
[Service valkey-primary]
     ↓
[StatefulSet valkey-primary × 1]
```

## 관측 가능성 스택

- Prometheus (kube-prometheus-stack v66.2.1) — 메트릭, 18 targets
- Grafana — 대시보드
- Alertmanager — 알림 라우팅
- PrometheusRule notiflex-alerts (NotiflexPodRestartTooMany, NotiflexPodDown)
- Loki + Fluent Bit — 로그 수집

## GitOps

- ArgoCD v7.8.2 (notiflex-smb Application, k8s/smb/ path)
- GitHub Actions CI (.github/workflows/ci.yaml) — image build + sed deployment.yaml + commit
- Self-heal + auto-prune

## Secret 관리

- Secret Manager (`valkey-password`)
- CSI Driver + GCP Provider
- Workload Identity (notiflex-sa@ ↔ notiflex namespace)

## 노드 구성

- default-pool: e2-medium × 2, Spot
- Future: api-pool / worker-pool / ops-pool (ch7.2)
