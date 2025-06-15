# ItOps 개발자 채용 과제 - 이슈 관리 시스템

## 과제 개요

- **제출 형태**: 개인 GitHub 레포지토리에 코드를 업로드하고, 해당 레포지토리 주소를 메일로 제출
- **제출 내용**: GitHub 레포지토리 링크 (예: https://github.com/username/itops-assignment)
- **구현 범위**: 백엔드 API + 프론트엔드 대시보드 + SQL 문제 답안

## 과제 프로젝트 구조

```
itops-assignment/
├── backend/              # Go 백엔드 API 서버 (보일러플레이트 미제공)
├── frontend/             # Vue.js 프론트엔드 (보일러플레이트 제공)
│   ├── src/
│   └── ...
├── SQL_ASSIGNMENT.md     # SQL 문제 설명
└── README.md             # 프로젝트 요구사항 및 실행 가이드
```

### 제출시 GitHub 레포지토리 구성 예시

```
itops-assignment/
├── backend/              # Go 백엔드 API 서버 답안
├── frontend/             # Vue.js 프론트엔드 답안
└── sql/                  # SQL 문제 답안
    ├── problem1.sql      # 문제 1 답안
    └── problem2.sql      # 문제 2 답안
```

## 기술 스택

### 백엔드

- **언어**: Go

### 프론트엔드

- **프레임워크**: Vue 3 (Composition API, Options API 모두 사용 가능)
- **상태 관리**: 필요시 Vuex 사용 가능
- **스타일링**: CSS, SCSS, CSS-in-JS 등 자유롭게 선택
- **포트**: 5173번 사용 (Vite 기본 포트)

### SQL

- **데이터베이스**: Oracle

## 실행 방법

### 1. 백엔드 서버 구현 및 실행

- **백엔드는 보일러 플레이트가 따로 없고 처음부터 구현해야 합니다.**

### 2. 프론트엔드 개발 서버 실행

- **프론트엔드는 보일러플레이트가 제공됩니다.**

```bash
cd frontend
npm install
npm run dev
```

- 개발 서버가 `http://localhost:5173`에서 실행됩니다.

## 데이터 모델

### User 구조체

```go
type User struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}
```

### Issue 구조체

```go
type Issue struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    User        *User     `json:"user,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

## 비즈니스 규칙

### 이슈 상태 (Status)

- 유효한 상태값: `PENDING`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED`

### 사용자 관리

- 시스템에 미리 정의된 사용자:

```json
[
  { "id": 1, "name": "김개발" },
  { "id": 2, "name": "이디자인" },
  { "id": 3, "name": "박기획" }
]
```

### 상태 변경 규칙

- 제목, 설명, 상태, 담당자 변경 가능
- 요청 data에 명시되지 않은 필드는 업데이트하지 않음
- 업데이트 시, 담당자가 없는 상태에서는 `IN_PROGRESS`, `COMPLETED` 상태로 변경 불가
- 상태가 `PENDING`일 때 담당자 할당 시 상태를 변경
  - 따로 상태를 지정하지 않는 경우 `IN_PROGRESS`로 변경되어야 함
- 담당자 제거(userId -> null) 시 상태는 `PENDING` 으로 변경
- `COMPLETED` 또는 `CANCELLED` 상태에서는 해당 issue 업데이트 불가

## 백엔드 구현 요구사항

### 필수 구현 기능

#### 1. API 엔드포인트 구현 (필수)

- [ ] POST /issue - 이슈 생성
- [ ] GET /issues - 이슈 목록 조회 (상태별 필터링 포함)
- [ ] GET /issue/:id - 이슈 상세 조회
- [ ] PATCH /issue/:id - 이슈 수정

#### 2.비즈니스 로직 구현 (필수)

- [ ] 이슈 상태 유효성 검증 (PENDING, IN_PROGRESS, COMPLETED, CANCELLED)
- [ ] 담당자 할당/해제 로직
- [ ] 상태 변경 규칙 적용
- [ ] 완료/취소된 이슈 수정 제한
- [ ] 부분 업데이트 처리 (명시되지 않은 필드는 업데이트하지 않음)

## 백엔드 API 명세

### API 공통 사항

- 필수 파라미터 누락, 유효하지 않은 데이터 요청 시 에러 응답 반환

### 1. 이슈 생성 [POST] /issue

- 담당자(userId)가 있는 경우: 상태를 `IN_PROGRESS`로 설정
  - 존재하지 않는 사용자를 담당자로 지정할 수 없습니다.
- 담당자(userId)가 없는 경우: 상태를 `PENDING`으로 설정

**요청 예시:**

```json
{
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "userId": 1
}
```

**응답 예시 (201 Created):**

```json
{
  "id": 1,
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "status": "IN_PROGRESS",
  "user": { "id": 1, "name": "김개발" },
  "createdAt": "2025-06-02T10:00:00Z",
  "updatedAt": "2025-06-02T10:00:00Z"
}
```

### 2. 이슈 목록 조회 [GET] /issues

- 쿼리 파라미터: `status` (선택사항)
- status 파라미터가 없는 경우 전체 이슈 조회
- 예시: `/issues?status=PENDING`

**응답 예시 (200 OK):**

```json
{
  "issues": [
    {
      "id": 1,
      "title": "버그 수정 필요",
      "description": "로그인 페이지에서 오류 발생",
      "status": "PENDING",
      "createdAt": "2025-06-02T10:00:00Z",
      "updatedAt": "2025-06-02T10:05:00Z"
    }
  ]
}
```

### 3. 이슈 상세 조회 [GET] /issue/:id

**응답 예시 (200 OK):**

```json
{
  "id": 1,
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "status": "PENDING",
  "createdAt": "2025-06-02T10:00:00Z",
  "updatedAt": "2025-06-02T10:05:00Z"
}
```

### 4. 이슈 수정 [PATCH] /issue/:id

**요청 예시:**

```json
{
  "title": "로그인 버그 수정",
  "status": "IN_PROGRESS",
  "userId": 2
}
```

**응답 예시 (200 OK):**

```json
{
  "id": 1,
  "title": "로그인 버그 수정",
  "description": "로그인 페이지에서 오류 발생",
  "status": "IN_PROGRESS",
  "user": { "id": 2, "name": "이디자인" },
  "createdAt": "2025-06-02T10:00:00Z",
  "updatedAt": "2025-06-02T10:10:00Z"
}
```

### 5. 에러 응답

- API는 적절한 HTTP 상태 코드, 메시지와 함께 에러 응답을 반환
- 에러 응답은 아래 형식을 따름

```json
{
  "error": "에러 메시지",
  "code": 400
}
```

- 응답에 포함되는 HTTP 상태 코드는 REST API 설계 원칙에 따라 적절하게 선택

## 프론트엔드 구현 요구사항

### 필수 구현 화면

#### 1. 이슈 목록 페이지 (필수)

- [ ] 이슈 목록 표시 (제목, 상태, 담당자, 생성일)
- [ ] 상태별 필터링 (전체, PENDING, IN_PROGRESS, COMPLETED, CANCELLED)
- [ ] 새 이슈 생성 버튼
- [ ] 각 이슈 클릭 시 상세 페이지로 이동

#### 2. 이슈 생성/상세/수정 페이지 (필수)

- [ ] 이슈 정보 표시 및 편집 (제목, 설명, 상태)
- [ ] 담당자 지정 기능 (드롭다운: 김개발, 이디자인, 박기획 중 선택)
- [ ] 상태 변경 기능 (PENDING, IN_PROGRESS, COMPLETED, CANCELLED)
- [ ] 저장 기능
- [ ] 목록으로 돌아가기

### 프론트엔드 처리 로직 (UI 제약사항)

- **상태 변경 제약**: 담당자가 지정되어 있지 않다면 이슈의 상태를 변경할 수 없음
- **완료/취소된 이슈**: `COMPLETED` 또는 `CANCELLED` 상태에서는 담당자/상태 변경 불가

### UI 구현 가이드

#### 상태 드롭다운 제어

- 담당자가 `null`인 경우: `PENDING` 상태만 선택 가능, 다른 상태는 비활성화
- 담당자가 있는 경우: 모든 상태 선택 가능

#### 담당자 드롭다운 제어

- `COMPLETED` 또는 `CANCELLED` 상태인 경우: 담당자 드롭다운 비활성화
- 다른 상태인 경우: 담당자 변경 가능

## SQL 문제 요구사항

- SQL_ASSIGNMENT.md 파일을 참고하여 문제를 풀어주세요.

## 추가 고려사항

- API 테스트 코드
- 백엔드 서버 구조(아키텍쳐) 기술
- 프론트엔드 컴포넌트 구조 기술
- 쿼리 최적화
