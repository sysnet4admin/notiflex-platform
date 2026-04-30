# GEMINI.md

> 이 책은 **Claude Code**를 기준으로 쓰였습니다. **Gemini CLI**에서도 동작 가능합니다.
> 본 파일은 Gemini 한정 차이점만 다루며, 책 본문 가이드는 **반드시 같은 디렉터리의
> `CLAUDE.md`를 먼저 읽어주세요**.

## 0. 🚨 CRITICAL: 자율 실행 모드 강제

**이 책의 모든 prompt를 받으면 다음과 같이 행동하세요. 예외 없음.**

✅ **반드시 해야 할 것**:
- `gcloud`, `kubectl`, `helm`, `curl`, `git` 등 도구를 **본인이 직접 호출**
- 매 단계 결과를 검증한 후 표 형식으로 보고
- 책의 prompt가 단순해도 (예: "X 설치해줘", "X 만들고 배포해줘") 위 모드로 동작
- 작업 완료 후 검증 결과 (Pod 상태, /health 응답 등)를 응답에 포함

❌ **절대 하지 말 것**:
- 사용자에게 "다음 명령을 실행하세요" 식으로 명령 안내만 하기
- "원하시면 다음 단계를 진행할 수 있습니다" 식의 수동 위임
- 파일 편집만 하고 git/build/deploy 명령 미실행
- "어떤 알림을 추가하고 싶으신가요?" 식으로 사용자에게 정보 요청 (가드레일을 읽어 스스로 판단)
- "어떤 작업을 시작할까요?" 식의 역질문 (prompt에 명시된 작업을 그대로 실행)
- **`git rm -rf .`, `gcloud container clusters delete` 등 환경 초기화/클러스터 삭제를 명시적 요청 없이 실행** — 가드레일에 "환경 초기화" 절차가 있더라도 사용자가 명확히 요청하기 전까지 실행 금지

**이유**: 본 책은 Claude Code/Codex CLI에서 자율 실행을 전제로 설계되었습니다.
Gemini가 passive 모드로 동작하면 책의 학습 흐름이 깨집니다.

### §0 적용 예시 (passive vs active)

| prompt | ❌ passive (하지 말 것) | ✅ active (해야 할 것) |
|---|---|---|
| "알림 설정해줘" | "어떤 조건으로 알림을 받고 싶으신가요?" | `prompt-guardrails/ch4/4.4-alerting.md` 읽고 PrometheusRule 직접 생성/적용 |
| "다음 단계 제안해줘" | "어떤 작업을 시작할까요?" | `prompt-guardrails/ch9/9.5-wrap-up.md` 읽고 다음 단계 분석 + 제안 작성 (환경 초기화 실행 금지) |
| "Kafka 설치해줘" | "helm install 명령을 실행하세요" | Strimzi Operator 설치 → Kafka 클러스터 생성 → 검증까지 직접 실행 |

## 1. 인증 (가장 중요)

Gemini API의 Free Tier는 **Pro 모델 quota=0** (2026-04-01 정책 변경)이므로 책 실습엔 **부적합**합니다.

| 옵션 | 가능한 모델 | 비용 | 추천도 |
|---|---|---|---|
| **Vertex AI** (GCP project + ADC) | gemini-2.5-pro, gemini-3.x preview | GCP 빌링 (신규 $300 크레딧) | ⭐ **권장** |
| Google AI Studio Pay-as-you-go | 모두 | $1.25/1M input, $10/1M output | 신용카드 등록 필요 |
| Free Tier (API key) | Flash만 (5~25 RPD) | 무료 | ❌ Agentic 작업 불가 |

### Vertex AI 셋업 (권장 경로)

```bash
# 1. Vertex AI API 활성화
gcloud services enable aiplatform.googleapis.com

# 2. Application Default Credentials 설정
gcloud auth application-default login --no-launch-browser
gcloud auth application-default set-quota-project <PROJECT_ID>

# 3. settings.json에서 authType 변경 (환경변수만으로는 부족)
#    ~/.gemini/settings.json → "selectedAuthType": "vertex-ai" 로 변경
python3 -c "
import json
with open('$(echo ~)/.gemini/settings.json') as f:
    s = json.load(f)
s['selectedAuthType'] = 'vertex-ai'
s.setdefault('security', {}).setdefault('auth', {})['selectedType'] = 'vertex-ai'
with open('$(echo ~)/.gemini/settings.json', 'w') as f:
    json.dump(s, f, indent=2)
"

# 4. 실행 시 환경변수 함께 전달
export GOOGLE_CLOUD_PROJECT=<PROJECT_ID>
export GOOGLE_CLOUD_LOCATION=us-central1   # gemini Pro 가용 region
```

> ⚠️ `~/.gemini/.env`에 GEMINI_API_KEY가 있으면 Vertex AI 환경변수가 무시됨.
> settings.json의 `selectedAuthType`을 `"vertex-ai"`로 바꾸는 것이 핵심.

## 2. 실행 명령

| Claude Code | Gemini CLI |
|---|---|
| `claude --dangerously-skip-permissions` | `gemini --yolo --model gemini-2.5-pro --include-directories <work_dir>` |

`--yolo`: 모든 도구 호출 자동 승인 (= claude의 dangerously-skip-permissions 대응)
`--include-directories <work_dir>`: 작업 저장소(notiflex-platform)를 추가 워크스페이스로 포함
`-p "<prompt>"`: 비대화형 (headless) 모드, 단일 prompt 실행
`-r/--resume`: 이전 세션 이어가기 (codex와 달리 메모리 보존 가능 ⭐)

## 3. 알려진 동작 차이 (vs Codex/Claude)

### 보수적 도구 사용 (가장 큰 차이)
- `--yolo` 모드여도 **shell 명령(git/build/kubectl) 실행은 보수적**
- 파일 read/edit는 적극적이나 **명령 실행을 사용자에게 위임**하는 경향
- 책의 "X 만들고 배포해줘"를 "사용자가 진행하도록 가이드 작성"으로 해석 가능
- **권고**: 명시적으로 "**도구를 직접 호출하여 실행하라**" prompt에 포함

### Hang 패턴 (v0.38.2에서 발생, v0.39.1에서 해소)
run-01(v0.38.2)에서 ch3.5c/ch5.3c/ch6.2c/ch8.1c 4회 hang이 발생했으나, run-02(v0.39.1 + 20분 타임아웃 규칙)에서 0회로 해소됨.

**대응 규칙** (버전 업그레이드 후에도 유지):
- 작업 시작 후 **20분** 내에 진행 로그가 없으면 해당 작업을 중단
- 더 구체적인 focused prompt(예: "Strimzi Operator만 먼저 설치해줘")로 단계 분리하여 재시도
- 재시도 시 `gemini --resume <session>`으로 이전 컨텍스트 이어가기 가능

### 컨텍스트 파일 자동 로드
- Gemini CLI는 `settings.json`의 `contextFileName`에 정의된 파일을 자동 로드
- 기본값에 `AGENTS.md` 포함됨 → 본 GEMINI.md와 함께 양쪽 다 로드 가능
- 책 본문(CLAUDE.md)도 명시적으로 읽도록 prompt에 안내 권장

### 경로 혼동 패턴 ⚠️ (run-02에서 반복 발생)
`--cd _Book_GitAIOps`로 실행 시 `notiflex-platform`을 형제 디렉터리로 인식하지 못하고 `_Book_GitAIOps` 하위에 파일을 생성하는 문제.

**영향 단계**: update-docs 스킬 생성, settings.local.json, 멀티테넌시 매니페스트 등
**대응**:
- 파일 생성 요청 시 **절대 경로**를 명시: `"/Users/.../notiflex-platform/k8s/enterprise/"`
- 생성 완료 후 파일 존재 여부 반드시 확인
- `find` 명령으로 실제 위치 탐색 후 올바른 경로로 이동

### /update-docs 세션 메모리 없음 한계 (run-02 발견)
Gemini `-p` 비대화형 모드는 매 호출이 fresh session이라 "이번 장에서 완료한 항목"을 자동 파악 불가.
`/update-docs` 실행 시 prompt에 완료 항목을 명시해야 함:

```
/update-docs. 이번 6장 완료: ch6.1 Valkey INCR 연동, ch6.2 GKE Secret Manager CSI+WI,
ch6.3 Canary 전략 전환, ch6.4 architecture.md 생성. JOURNEY.md ch6 ✅ 갱신해줘
```

### Loki 설치 주의 (run-02 발견)
최신 `grafana/loki` 차트는 object storage 요건으로 기본 설정에서 Pending 발생.
**대응**: `grafana/loki-stack` 차트 또는 `grafana/loki` SingleBinary + chunksCache.enabled=false 명시.

### web_fetch fallback 오류
- Gemini가 웹 검색 시 `gemini-3-flash-preview` 모델로 fallback하는데 access 없어 오류 발생
- **대응**: 웹 검색보다 `decision-guides/`, `prompt-guardrails/` 파일을 직접 읽도록 유도
- 오류가 발생해도 실제 작업에는 영향 없음 — 계속 진행

### 모델 선택
- Free Tier: Flash만 가능 (Pro 차단)
- Vertex AI: gemini-2.5-pro 등 Pro 사용 가능
- 책 실습엔 **gemini-2.5-pro 이상** 권장 (agentic capability)

### Skill / Subagent 부재
- Claude Code의 `/skill`, `Agent` subagent는 Gemini에 없음
- `/update-docs` 실행 시 Gemini가 `gemini skills` 형식으로 만들려 할 수 있음 → **반드시 `.claude/commands/update-docs.md` 파일로 만들어야 함** (절대 경로 명시 권장)
- 스킬 파일 생성 prompt 예시:
  ```
  /Users/.../notiflex-platform/.claude/commands/update-docs.md 파일을
  prompt-guardrails/ch2/update-docs-skill.md 가드레일의 마크다운 블록 내용으로 생성해줘
  ```

## 4. 책 본문 진행

이하 모든 가이드는 `CLAUDE.md`와 동일하게 따릅니다:
- `decision-guides/` — 도구 선택 근거
- `prompt-guardrails/` — 단계별 실행 지침
- `result-templates/` — 검증 체크리스트

세 디렉터리는 도구 명칭(Edit/Write 등)과 무관하므로 어느 에이전트에서나 동일하게 동작합니다.

## 5. Prompt 형식 (책 본문 그대로 사용)

위 §0 (자율 실행 모드 강제)가 본 GEMINI.md를 통해 전달되면, **책 prompt를 그대로 사용**해도 자율 실행됩니다:

| 책 본문 prompt | Gemini 동작 (§0 적용 시) |
|---|---|
| "Notiflex 앱 만들고 배포해줘" | ✅ git/build/deploy/검증 자율 실행 |
| "ArgoCD 설치해줘" | ✅ kubectl/helm 자율 실행 + 검증 |

**§0 미적용 시 우회**: 만약 GEMINI.md가 자동 로드되지 않는 환경(다른 디렉터리에서 실행 등)에서는 prompt에 "**도구를 직접 호출하여 실행해라**"를 명시적으로 추가하세요.

## 6. 3-prompt 패턴 호환 (codex와 동일 한계)

책의 일부 서브챕터는 **탐색 → 비교 → 실행** 3단계 prompt 패턴을 사용합니다.
Gemini CLI의 `-p` 비대화형 모드는 매 호출이 fresh session(메모리 없음)이므로,
**비교 단계의 "다른 도구도 있다고 했는데"** 같은 모호한 prompt는 컨텍스트 손실 발생.

**권장 우회**: 비교 prompt에 명시적 주제어 포함
- 책 원본: "다른 도구도 있다고 했는데, 비교하면 어때?"
- Gemini용: "**ArgoCD 외 다른 GitOps 도구**도 있다고 했는데, 비교하면 어때?"

대안: `gemini --resume <session>`으로 세션 메모리 보존 (codex 대비 강점).

## 7. 다른 에이전트 AI

- Codex CLI: `AGENTS.md` (별도 파일) — 본 책 저장소 루트에 있음
- Claude Code: `CLAUDE.md` (책 본문 가이드)
