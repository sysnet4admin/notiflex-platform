# Notiflex Platform — Journey

## 진행 현황 (run-53)

| 장 | 절 | 상태 | 완료일 |
|---|---|:---:|---|
| 2 | 2.2 install-check | ✅ | 2026-04-24 |
| 2 | 2.3 gcloud | ✅ | 2026-04-24 |
| 2 | 2.4 github-repo | ✅ | 2026-04-24 |
| 2 | 2.5 gke-cluster | ✅ | 2026-04-24 |
| 2 | 2.6 build-deploy | ✅ | 2026-04-24 |
| 2 | 2.7 first-commit | ✅ | 2026-04-24 |
| 2 | [별도] /update-docs 스킬 | ✅ | 2026-04-24 |
| 3 | 3.2 argocd | ✅ | 2026-04-24 |
| 3 | 3.3 rolling | ✅ | 2026-04-24 |
| 3 | 3.4 github-actions | ✅ | 2026-04-24 |
| 3 | 3.5 ci-argocd | ✅ | 2026-04-24 |
| 3 | [별도] CLAUDE.md 규칙 | ✅ | 2026-04-24 |

## 도구 선택 기록

- **GitOps**: ArgoCD (vs Flux/Jenkins X/Spinnaker) — UI + App of Apps + ArgoCD 생태계
- **CI**: GitHub Actions (vs Cloud Build/GitLab CI/Jenkins) — 저장소 일치 + 학습 단순

## 현재 버전

- notiflex-api: v0.1.1 (또는 sha-XXX, CI 자동)
- ArgoCD: 7.8.2 (Helm chart)
- Cluster: GKE notiflex-cluster

## 현재 리소스

- 노드풀: default-pool (e2-medium × 2, Spot)
- Namespace: notiflex (Deployment + Service), argocd (ArgoCD)
- ArgoCD Application: notiflex-smb (Synced + Healthy)

## 트러블슈팅 이력

(없음)
