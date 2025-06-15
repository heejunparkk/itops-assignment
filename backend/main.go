package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 이슈 상태 상수
const (
	StatusPending    = "PENDING"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
)

// User 모델
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Issue 모델
type Issue struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	User        *User     `json:"user"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 전역 변수로 사용자와 이슈 데이터 저장 (실제로는 DB 사용)
var users = []User{
	{ID: 1, Name: "김개발"},
	{ID: 2, Name: "이디자인"},
	{ID: 3, Name: "박기획"},
}

var issues = []Issue{
	{
		ID:          1,
		Title:       "첫 번째 이슈",
		Description: "이것은 첫 번째 이슈입니다.",
		Status:      StatusPending,
		User:        nil,
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		UpdatedAt:   time.Now().Add(-48 * time.Hour),
	},
	{
		ID:          2,
		Title:       "두 번째 이슈",
		Description: "이것은 두 번째 이슈입니다.",
		Status:      StatusInProgress,
		User:        &users[0],
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	},
	{
		ID:          3,
		Title:       "세 번째 이슈",
		Description: "이것은 세 번째 이슈입니다.",
		Status:      StatusCompleted,
		User:        &users[1],
		CreatedAt:   time.Now().Add(-12 * time.Hour),
		UpdatedAt:   time.Now().Add(-6 * time.Hour),
	},
}

// 다음 ID 생성 함수
func getNextID() int {
	maxID := 0
	for _, issue := range issues {
		if issue.ID > maxID {
			maxID = issue.ID
		}
	}
	return maxID + 1
}

// 사용자 ID로 사용자 찾기
func findUserByID(userID int) *User {
	for i, user := range users {
		if user.ID == userID {
			return &users[i]
		}
	}
	return nil
}

// CORS 미들웨어
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 이슈 생성 핸들러
func createIssueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var issue Issue
	err := json.NewDecoder(r.Body).Decode(&issue)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 필수 필드 검증
	if issue.Title == "" || issue.Description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	// 상태 검증
	if !isValidStatus(issue.Status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// 담당자가 없는데 상태가 PENDING이 아닌 경우
	if issue.User == nil && issue.Status != StatusPending {
		http.Error(w, "Cannot set status other than PENDING without a user", http.StatusBadRequest)
		return
	}

	// 새 이슈 생성
	now := time.Now()
	newIssue := Issue{
		ID:          getNextID(),
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 담당자 설정
	if issue.User != nil {
		user := findUserByID(issue.User.ID)
		if user == nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		newIssue.User = user
	}

	// 상태를 사용자가 직접 지정한 값으로 유지
	// 사용자가 직접 PENDING으로 설정한 경우 담당자가 있어도 PENDING 상태 유지

	// 이슈 저장
	issues = append(issues, newIssue)

	// 응답
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newIssue)
}

// 이슈 목록 조회 핸들러
func getIssuesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 상태 필터링
	status := r.URL.Query().Get("status")
	var filteredIssues []Issue

	if status != "" && isValidStatus(status) {
		for _, issue := range issues {
			if issue.Status == status {
				filteredIssues = append(filteredIssues, issue)
			}
		}
	} else {
		// 모든 이슈 복사
		filteredIssues = make([]Issue, len(issues))
		copy(filteredIssues, issues)
	}

	// ID 내림차순으로 정렬 (최신순)
	sort.Slice(filteredIssues, func(i, j int) bool {
		return filteredIssues[i].ID > filteredIssues[j].ID
	})

	// 응답
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredIssues)
}

// 이슈 상세 조회 핸들러
func getIssueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 ID 추출
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path[2])
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	// 이슈 찾기
	for _, issue := range issues {
		if issue.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(issue)
			return
		}
	}

	http.Error(w, "Issue not found", http.StatusNotFound)
}

// 이슈 수정 핸들러
func updateIssueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 ID 추출
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path[2])
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	// 요청 본문 파싱
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 이슈 찾기
	var targetIssue *Issue
	var index int
	for i, issue := range issues {
		if issue.ID == id {
			targetIssue = &issues[i]
			index = i
			break
		}
	}

	if targetIssue == nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	// 완료/취소된 이슈는 수정 불가
	if targetIssue.Status == StatusCompleted || targetIssue.Status == StatusCancelled {
		http.Error(w, "Cannot update completed or cancelled issues", http.StatusBadRequest)
		return
	}

	// 필드 업데이트
	if title, ok := updates["title"].(string); ok && title != "" {
		targetIssue.Title = title
	}

	if description, ok := updates["description"].(string); ok {
		targetIssue.Description = description
	}

	// 담당자 업데이트
	var userUpdated bool
	if userMap, ok := updates["user"].(map[string]interface{}); ok {
		if userID, ok := userMap["id"].(float64); ok {
			user := findUserByID(int(userID))
			if user == nil {
				http.Error(w, "User not found", http.StatusBadRequest)
				return
			}
			targetIssue.User = user
			userUpdated = true
		}
	} else if updates["user"] == nil {
		targetIssue.User = nil
		userUpdated = true
	}

	// 상태 업데이트
	if status, ok := updates["status"].(string); ok && isValidStatus(status) {
		// 담당자가 없는데 PENDING이 아닌 상태로 변경 시도
		if targetIssue.User == nil && status != StatusPending {
			http.Error(w, "Cannot set status other than PENDING without a user", http.StatusBadRequest)
			return
		}
		targetIssue.Status = status
	}

	// 샀용자가 직접 상태를 설정한 경우에는 상태를 변경하지 않음
	// 담당자 업데이트로 인한 자동 상태 변경은 사용자가 직접 상태를 설정한 경우에는 적용하지 않음

	// 담당자가 제거되었으면 상태를 PENDING으로 변경
	if userUpdated && targetIssue.User == nil {
		targetIssue.Status = StatusPending
	}

	// 업데이트 시간 갱신
	targetIssue.UpdatedAt = time.Now()

	// 이슈 업데이트
	issues[index] = *targetIssue

	// 응답
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(targetIssue)
}

// 사용자 목록 조회 핸들러
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 상태 유효성 검사
func isValidStatus(status string) bool {
	return status == StatusPending ||
		status == StatusInProgress ||
		status == StatusCompleted ||
		status == StatusCancelled
}

// 라우터 설정
func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// API 엔드포인트
	mux.HandleFunc("/issue", createIssueHandler)
	mux.HandleFunc("/issues", getIssuesHandler)
	mux.HandleFunc("/users", getUsersHandler)
	mux.HandleFunc("/issue/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getIssueHandler(w, r)
		} else if r.Method == "PATCH" {
			updateIssueHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// CORS 미들웨어 적용
	return enableCORS(mux)
}

func main() {
	handler := setupRoutes()
	port := 8080

	fmt.Printf("서버가 http://localhost:%d 에서 실행 중입니다...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
