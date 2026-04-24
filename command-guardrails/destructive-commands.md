# 파괴적 명령 단계 제어

## kubectl delete

```
1. 무엇을 삭제하는가? (리소스 종류·이름·namespace)
2. 영향 범위는? (PV는 PVC 삭제 시 삭제? Service는 Pod 삭제 시 어떻게?)
3. 백업했는가? (etcd snapshot 또는 manifest export)
4. 사용자 확인 받기
5. 실행
6. 결과 확인 + JOURNEY.md 기록
```

## helm uninstall

```
1. release 이름 + namespace 확인
2. 의존하는 리소스 확인 (CRD instance가 있으면 cleanup 필요)
3. 사용자 확인 받기
4. 실행
5. 잔존 PVC/PV 확인
```

## gcloud container clusters delete

```
1. 클러스터 안의 데이터 백업 (etcd state, PV → 외부 저장)
2. 의존 리소스 확인 (LB, PV, secret manager 참조)
3. async로 실행 + 다른 정리와 병렬
4. 완료 후 고아 리소스 (LB, PV, secret) 별도 정리
```
