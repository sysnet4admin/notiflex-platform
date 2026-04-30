# ADR-014: OpenTelemetry와 Tempo를 사용한 분산 트레이싱 도입

## 상태

**채택 (Accepted)**

## 맥락 (Context)

-   Notiflex 플랫폼이 마이크로서비스 아키텍처로 발전함에 따라, 사용자 요청이 여러 서비스를 거치며 처리됩니다.
-   서비스 간 호출 관계가 복잡해지면서 전체 요청의 흐름을 추적하고, 병목 구간을 식별하거나 장애 발생 지점을 특정하기가 점점 어려워지고 있습니다.
-   성능 최적화와 안정적인 서비스 운영을 위해 분산 트레이싱(Distributed Tracing) 시스템 도입이 필수적입니다.

## 결정 (Decision)

**OpenTelemetry (OTel) SDK를 애플리케이션에 통합하고, 수집된 트레이스(Trace)를 Grafana Tempo를 통해 저장 및 조회합니다.**

-   애플리케이션 코드(Go)에 OpenTelemetry SDK를 사용하여 트레이스 데이터를 생성하고 전송(export)합니다.
-   수집된 트레이스 데이터는 Grafana Tempo에 저장합니다. Tempo는 대규모 트레이스를 위해 설계된 비용 효율적인 분산 트레이싱 백엔드입니다.
-   Tempo는 `ops-pool`에 배포하여 다른 운영 도구(Prometheus, Loki)와 함께 관리합니다.
-   Grafana 대시보드를 통해 메트릭(Metrics), 로그(Logs), 트레이스(Traces)를 통합적으로 조회(MLT)하여 Observability를 극대화합니다.

## 이유 (Rationale)

-   **표준화된 계측**: OpenTelemetry는 CNCF의 표준 프로젝트로, 코드 계측 방법을 표준화하여 특정 벤더에 대한 종속성을 없애줍니다. OTel SDK를 한번 적용하면 다양한 백엔드 시스템으로 손쉽게 전환할 수 있습니다.
-   **Grafana 스택과의 완벽한 통합**: Tempo는 Grafana Labs에서 개발하여 Prometheus, Loki와 자연스럽게 연동됩니다. Grafana UI 내에서 로그 라인을 클릭하여 관련된 트레이스를 바로 확인하거나, 트레이스에서 특정 메트릭으로 드릴다운하는 등 유기적인 분석이 가능합니다.
-   **비용 효율적인 대규모 저장**: Tempo는 트레이스 전체를 인덱싱하는 대신, 오브젝트 스토리지(GCS, S3 등)에 압축하여 저장하므로 기존 방식(e.g., Jaeger, Zipkin)에 비해 대규모 데이터를 훨씬 저렴한 비용으로 장기간 보관할 수 있습니다.
-   **운영 단순성**: Tempo는 단일 바이너리로 실행 가능하며, 설치 및 운영이 비교적 간단합니다. `ops-pool`에 다른 관측 가능성 도구와 함께 배치하여 관리 부담을 줄일 수 있습니다.
