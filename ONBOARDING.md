# Notiflex Platform — Onboarding (30분)

## 1. 접근 정보

- **kubectl context**: `gke-sysnet4admin_book_gitaiops`
- **cluster**: GKE notiflex-cluster v1.35.1-gke.1396002 (asia-northeast3-a)
- **이미지 레지스트리**: `asia-northeast3-docker.pkg.dev/project-75fce205-dfa5-4975-a56/notiflex/api:vX.Y.Z`
- **ArgoCD UI**: `kubectl port-forward svc/argocd-server -n argocd 8443:443` → https://localhost:8443 (admin / `kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d`)
- **Grafana UI**: `kubectl port-forward svc/kube-prometheus-grafana -n monitoring 3000:80` → http://localhost:3000 (admin / prom-operator)
- **Gateway IP**: `kubectl -n notiflex get gateway notiflex-gateway -o jsonpath='{.status.addresses[0].value}'` → curl http://<IP>/health

## 2. 배포 플로우

```
git push (app/** 변경)
  → GitHub Actions (.github/workflows/ci.yaml)
  → gcloud builds submit + push (image: sha-XXX 태그)
  → ci.yaml이 deployment.yaml의 image tag 자동 sed 치환 + commit (`[skip ci]`)
  → ArgoCD polling (3분) 또는 webhook → notiflex-smb sync
  → Argo Rollouts Canary (20→50→80→100, 30s pause)
  → Gateway API로 외부 서비스
```

## 3. Q&A

1. **Canary abort**: `kubectl argo rollouts abort notiflex-api -n notiflex`
2. **Loki 로그 조회**: `{namespace_name="notiflex"} |= "error"` (Grafana Explore)
3. **Tempo trace**: 앱 OTel SDK 통합 후 `OTEL_EXPORTER_OTLP_ENDPOINT=tempo.monitoring:4317`
4. **Kafka topic 추가**: `KafkaTopic` CRD apply (k8s/kafka/)
5. **새 테넌트**: `argocd/notiflex-<tenant>.yaml` + `k8s/<tenant>/` 추가 → root-app이 자동 발견
6. **알림 확인**: Alertmanager UI 또는 `kubectl get prometheusrule -A`

## 4. 클러스터 복구

- 노드 재시작: `gcloud compute instances reset <NODE_NAME> --zone=asia-northeast3-a`
- 클러스터 재생성: `gcloud container clusters delete && create` (ch2.5 가드레일 참조)
- 전체 환경 정리: `prompt-guardrails/ch9/restart-to-ready.md`
