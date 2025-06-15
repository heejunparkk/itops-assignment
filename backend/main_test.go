package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// 테스트용 이슈 데이터 생성 함수
func createTestIssue() Issue {
	return Issue{
		Title:       "테스트 이슈",
		Description: "이것은 테스트 이슈입니다.",
		Status:      StatusPending,
		User:        nil,
	}
}

// 테스트용 이슈 데이터 생성 함수 (사용자 포함)
func createTestIssueWithUser() Issue {
	return Issue{
		Title:       "테스트 이슈 (사용자 있음)",
		Description: "이것은 사용자가 있는 테스트 이슈입니다.",
		Status:      StatusInProgress,
		User:        &users[0],
	}
}

// 이슈 목록 조회 테스트
func TestGetIssues(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 테스트 서버 생성
	req, err := http.NewRequest("GET", "/issues", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getIssuesHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusOK)
	}

	// 응답 본문 확인
	var response []Issue
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 이슈 목록이 비어있지 않은지 확인
	if len(response) == 0 {
		t.Errorf("이슈 목록이 비어있습니다")
	}
}

// 상태별 이슈 필터링 테스트
func TestGetIssuesWithStatusFilter(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 테스트할 상태
	testStatus := StatusPending

	// 테스트 서버 생성
	req, err := http.NewRequest("GET", "/issues?status="+testStatus, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getIssuesHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusOK)
	}

	// 응답 본문 확인
	var response []Issue
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 모든 이슈가 요청한 상태인지 확인
	for _, issue := range response {
		if issue.Status != testStatus {
			t.Errorf("필터링된 이슈 중 잘못된 상태가 있습니다: %v, 기대값: %v", issue.Status, testStatus)
		}
	}
}

// 이슈 생성 테스트
func TestCreateIssue(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 테스트 이슈 생성
	testIssue := createTestIssue()
	issueJSON, _ := json.Marshal(testIssue)

	// 테스트 서버 생성
	req, err := http.NewRequest("POST", "/issues", bytes.NewBuffer(issueJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createIssueHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusCreated)
	}

	// 응답 본문 확인
	var response Issue
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 생성된 이슈 확인
	if response.Title != testIssue.Title {
		t.Errorf("생성된 이슈의 제목이 다릅니다: %v, 기대값: %v", response.Title, testIssue.Title)
	}
	if response.Description != testIssue.Description {
		t.Errorf("생성된 이슈의 설명이 다릅니다: %v, 기대값: %v", response.Description, testIssue.Description)
	}
	if response.Status != testIssue.Status {
		t.Errorf("생성된 이슈의 상태가 다릅니다: %v, 기대값: %v", response.Status, testIssue.Status)
	}
	if response.ID <= 0 {
		t.Errorf("생성된 이슈의 ID가 유효하지 않습니다: %v", response.ID)
	}
}

// 이슈 상세 조회 테스트
func TestGetIssue(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 테스트할 이슈 ID
	testIssueID := 1

	// 테스트 서버 생성
	req, err := http.NewRequest("GET", "/issues/"+strconv.Itoa(testIssueID), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getIssueHandler(w, r)
	})
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusOK)
	}

	// 응답 본문 확인
	var response Issue
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 조회된 이슈 확인
	if response.ID != testIssueID {
		t.Errorf("조회된 이슈의 ID가 다릅니다: %v, 기대값: %v", response.ID, testIssueID)
	}
}

// 이슈 업데이트 테스트
func TestUpdateIssue(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 테스트할 이슈 ID
	testIssueID := 1

	// 업데이트할 이슈 데이터
	updatedIssue := Issue{
		Title:       "업데이트된 이슈",
		Description: "이것은 업데이트된 이슈입니다.",
		Status:      StatusPending,
		User:        nil,
	}
	issueJSON, _ := json.Marshal(updatedIssue)

	// 테스트 서버 생성
	req, err := http.NewRequest("PATCH", "/issues/"+strconv.Itoa(testIssueID), bytes.NewBuffer(issueJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateIssueHandler(w, r)
	})
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusOK)
	}

	// 응답 본문 확인
	var response Issue
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 업데이트된 이슈 확인
	if response.ID != testIssueID {
		t.Errorf("업데이트된 이슈의 ID가 다릅니다: %v, 기대값: %v", response.ID, testIssueID)
	}
	if response.Title != updatedIssue.Title {
		t.Errorf("업데이트된 이슈의 제목이 다릅니다: %v, 기대값: %v", response.Title, updatedIssue.Title)
	}
	if response.Description != updatedIssue.Description {
		t.Errorf("업데이트된 이슈의 설명이 다릅니다: %v, 기대값: %v", response.Description, updatedIssue.Description)
	}
}

// 비즈니스 규칙 테스트: 담당자가 없는 경우 PENDING 상태만 가능
func TestBusinessRuleUserNullStatusPending(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 담당자가 없고 상태가 PENDING이 아닌 이슈 생성
	testIssue := Issue{
		Title:       "테스트 이슈",
		Description: "이것은 테스트 이슈입니다.",
		Status:      StatusInProgress, // PENDING이 아닌 상태
		User:        nil,              // 담당자 없음
	}
	issueJSON, _ := json.Marshal(testIssue)

	// 테스트 서버 생성
	req, err := http.NewRequest("POST", "/issues", bytes.NewBuffer(issueJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createIssueHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인 - 비즈니스 규칙 위반으로 400 Bad Request 예상
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusBadRequest)
	}
}

// 비즈니스 규칙 테스트: COMPLETED 또는 CANCELLED 상태에서는 담당자 변경 불가
func TestBusinessRuleCompletedCannotChangeUser(t *testing.T) {
	// 테스트 전 원본 이슈 데이터 백업
	originalIssues := issues
	defer func() {
		issues = originalIssues
	}()

	// 완료된 이슈 생성 및 저장
	completedIssue := Issue{
		ID:          getNextID(),
		Title:       "완료된 이슈",
		Description: "이것은 완료된 이슈입니다.",
		Status:      StatusCompleted,
		User:        &users[0],
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	issues = append(issues, completedIssue)

	// 담당자 변경 시도
	updatedIssue := completedIssue
	updatedIssue.User = &users[1] // 다른 담당자로 변경
	issueJSON, _ := json.Marshal(updatedIssue)

	// 테스트 서버 생성
	req, err := http.NewRequest("PATCH", "/issues/"+strconv.Itoa(completedIssue.ID), bytes.NewBuffer(issueJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateIssueHandler(w, r)
	})
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인 - 비즈니스 규칙 위반으로 400 Bad Request 예상
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusBadRequest)
	}
}

// 사용자 목록 조회 테스트
func TestGetUsers(t *testing.T) {
	// 테스트 서버 생성
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUsersHandler)
	handler.ServeHTTP(rr, req)

	// 상태 코드 확인
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("핸들러가 잘못된 상태 코드를 반환했습니다: %v, 기대값: %v", status, http.StatusOK)
	}

	// 응답 본문 확인
	var response []User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("응답을 JSON으로 파싱할 수 없습니다: %v", err)
	}

	// 사용자 목록이 비어있지 않은지 확인
	if len(response) == 0 {
		t.Errorf("사용자 목록이 비어있습니다")
	}

	// 사용자 수가 예상과 일치하는지 확인
	if len(response) != len(users) {
		t.Errorf("사용자 수가 일치하지 않습니다: %v, 기대값: %v", len(response), len(users))
	}
}
