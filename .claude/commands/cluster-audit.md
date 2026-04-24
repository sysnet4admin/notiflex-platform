---
description: 클러스터 전반 상태 점검 (Pod, 리소스, 관측 가능성, ArgoCD)
---

# /cluster-audit

`gke-sysnet4admin_book_gitaiops` 컨텍스트의 클러스터를 점검:

## 1. Pod 상태 (전체 namespace)

```bash
kubectl --context gke-sysnet4admin_book_gitaiops get pods -A
```
- 모든 Pod이 Running/Completed인지
- CrashLoopBackOff, Error, Pending 있으면 즉시 보고

## 2. 노드 리소스

```bash
kubectl --context gke-sysnet4admin_book_gitaiops top nodes
kubectl --context gke-sysnet4admin_book_gitaiops top pods -A --sort-by=cpu | head -10
```
- CPU/Memory 사용률 80% 초과 시 보고

## 3. Prometheus 동작 확인

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n monitoring port-forward svc/kube-prometheus-kube-prome-prometheus 19090:9090 &
PF=$!; sleep 3
curl -s "http://localhost:19090/api/v1/targets" | jq '.data.activeTargets | length'  # 18+ 기대
curl -s "http://localhost:19090/api/v1/alerts" | jq '.data.alerts | length'
kill $PF
```

## 4. Loki ready

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n logging port-forward svc/loki 13100:3100 &
PF=$!; sleep 3
curl -s http://localhost:13100/ready  # "ready" 기대
kill $PF
```

## 5. ArgoCD Application 상태

```bash
kubectl --context gke-sysnet4admin_book_gitaiops -n argocd get app -o custom-columns=NAME:.metadata.name,SYNC:.status.sync.status,HEALTH:.status.health.status
```
- 모든 app이 Synced + Healthy인지

## 6. 결과 요약

위 5개 점검 결과를 표 형식으로 보고. 이슈가 있으면 ⚠️ 표시 + 추천 조치.
