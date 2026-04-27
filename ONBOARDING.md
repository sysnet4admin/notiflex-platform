# Notiflex Platform 온보딩 가이드 (30분)

## 사전 준비

- gcloud CLI 설치 + GCP 인증
- kubectl + helm 설치
- GitHub 계정 (sysnet4admin/notiflex-platform 접근)

## 클러스터 접속

```bash
gcloud container clusters get-credentials notiflex-cluster \
  --zone=asia-northeast3-a --project=project-75fce205-dfa5-4975-a56
kubectl config rename-context ... gke-sysnet4admin_book_gitaiops
```

## 주요 접근 방법

**API (외부)**:
```bash
curl http://35.216.122.221/health   # {"status":"ok"}
curl http://35.216.122.221/version  # {"version":"v0.1.1"}
```

**ArgoCD UI**:
```bash
kubectl port-forward svc/argocd-server -n argocd 8443:443
# https://localhost:8443 (admin/초기비밀번호)
```

**Grafana**:
```bash
kubectl port-forward svc/kube-prometheus-grafana -n monitoring 3000:80
# http://localhost:3000
```

## 아키텍처 요약

```
외부 → Gateway API (35.216.122.221) → HTTPRoute → notiflex-api Service
  → Rollout (Canary 20→50→80%) → Pod (api-pool 노드)
  → Valkey (공유 카운터)

모니터링: Prometheus + Grafana + Loki + Tempo (monitoring ns)
메시징: Kafka 4.1.0 KRaft (kafka ns)
GitOps: ArgoCD App of Apps (root-app → notiflex-smb + notiflex-enterprise)
```

## 배포 플로우

1. `app/` 수정 → git push
2. GitHub Actions CI (WIF) → 이미지 빌드 → sha-XXXXXXX 태그
3. deployment.yaml 이미지 태그 자동 업데이트 + push [skip ci]
4. ArgoCD 감지 → Canary 배포 (20→50→80%)

## 자주 묻는 질문

**Q: Canary 배포를 중단하려면?**
```bash
kubectl argo rollouts abort notiflex-api -n notiflex
```

**Q: 로그는 어떻게 검색해?**
Grafana → Explore → Loki → `{namespace="notiflex"} |= "error"`

**Q: 트레이스는?**
Grafana → Explore → Tempo → TraceID로 조회

**Q: 새 테넌트 추가는?**
k8s/enterprise/ 디렉터리를 복사해서 새 네임스페이스로 수정 + argocd/ Application YAML 추가

**Q: Kafka 토픽 추가는?**
`k8s/kafka/` 에 KafkaTopic YAML 추가 → git push → ArgoCD 자동 배포

**Q: 알림 규칙 추가는?**
`k8s/monitoring/notiflex-alerts.yaml` 에 PrometheusRule 추가 → git push
