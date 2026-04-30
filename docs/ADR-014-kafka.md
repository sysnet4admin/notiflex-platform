# ADR-013: Strimzi를 이용한 Kafka 클러스터 도입

## 상태

**채택 (Accepted)**

## 맥락 (Context)

-   실시간 알림 처리를 위한 비동기 메시징 시스템이 필요합니다.
-   대용량 트래픽을 안정적으로 처리하고, MSA(Microservice Architecture) 환경에서 각 컴포넌트 간의 결합도를 낮추어야 합니다.
-   메시지 큐를 통해 특정 서비스의 장애가 다른 서비스로 전파되는 것을 막고, 시스템 전체의 회복탄력성을 높여야 합니다.

## 결정 (Decision)

**Strimzi Operator를 사용하여 Kubernetes 네이티브 방식으로 Kafka 클러스터를 `worker-pool`에 배포합니다.**

-   Apache Kafka를 메시징 시스템으로 선택합니다.
-   Strimzi 프로젝트가 제공하는 Operator를 활용하여 Kubernetes 클러스터 내에서 Kafka를 선언적으로 관리합니다.
-   최신 버전의 Kafka에서 지원하는 KRaft(Kafka Raft) 모드를 사용하여 ZooKeeper 없이 클러스터를 구성, 운영 부담을 줄입니다.
-   I/O가 빈번한 Kafka 브로커는 `worker-pool` 노드에 배치하여 다른 워크로드와의 리소스 경합을 방지합니다.

## 이유 (Rationale)

-   **업계 표준 및 생태계**: Kafka는 고성능 실시간 데이터 스트리밍의 사실상 표준이며, 풍부한 커넥터와 라이브러리 생태계를 갖추고 있습니다.
-   **Strimzi Operator의 편리성**: Strimzi는 Kubernetes 환경에서 Kafka의 설치, 설정, 업그레이드, 보안(TLS, SASL), 모니터링 등 복잡한 운영 작업을 자동화하여 개발자와 운영자가 비즈니스 로직에 더 집중할 수 있게 해줍니다.
-   **KRaft 모드의 이점**: ZooKeeper 의존성을 제거함으로써 아키텍처를 단순화하고, 관리 포인트를 줄이며, 더 빠른 클러스터 부팅과 장애 복구를 가능하게 합니다. 리소스 사용량 또한 감소합니다.
-   **워크로드 격리**: I/O 집약적인 Kafka 워크로드를 `worker-pool`에 격리함으로써, CPU 집약적인 API 서버나 다른 운영 도구의 성능에 영향을 주지 않고 안정적인 메시징 환경을 보장합니다.
