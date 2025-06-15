<template>
  <div class="issue-list">
    <div class="issue-list-header">
      <h2>이슈 목록</h2>
      <router-link to="/issues/new" class="btn btn-primary">새 이슈 생성</router-link>
    </div>

    <div class="filter-section">
      <div class="filter-group">
        <label>상태 필터링:</label>
        <select v-model="statusFilter" class="form-control" @change="changeFilter(statusFilter)">
          <option value="ALL">전체</option>
          <option value="PENDING">PENDING</option>
          <option value="IN_PROGRESS">IN_PROGRESS</option>
          <option value="COMPLETED">COMPLETED</option>
          <option value="CANCELLED">CANCELLED</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>이슈 목록을 불러오는 중입니다... ({{ statusFilter }})</p>
    </div>
    <div v-else-if="filteredIssues.length === 0" class="no-issues">
      <div class="empty-state">
        <i class="empty-icon">📋</i>
        <h3>이슈가 없습니다</h3>
        <p v-if="statusFilter !== 'ALL'">선택하신 '{{ statusFilter }}' 상태의 이슈가 없습니다.</p>
        <p v-else>아직 등록된 이슈가 없습니다. 새 이슈를 생성해보세요!</p>
        <router-link to="/issues/new" class="btn btn-outline-primary">새 이슈 생성하기</router-link>
      </div>
    </div>
    <div v-else class="issue-cards">
      <div
        v-for="issue in filteredIssues"
        :key="issue.id"
        class="card issue-card"
        @click="goToIssueDetail(issue.id)"
      >
        <div class="issue-card-header">
          <h3>{{ issue.title }}</h3>
          <div
            class="status-badge"
            :class="{
              'status-pending': issue.status === 'PENDING',
              'status-in-progress': issue.status === 'IN_PROGRESS',
              'status-completed': issue.status === 'COMPLETED',
              'status-cancelled': issue.status === 'CANCELLED',
            }"
          >
            {{ issue.status }}
          </div>
        </div>
        <div class="issue-card-body">
          <p class="issue-description">{{ issue.description }}</p>
        </div>
        <div class="issue-card-footer">
          <div class="issue-meta">
            <span class="issue-assignee">
              담당자: {{ issue.user ? issue.user.name : '없음' }}
            </span>
            <span class="issue-date"> 생성일: {{ formatDate(issue.createdAt) }} </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'IssueList',
  setup() {
    const router = useRouter()
    const statusFilter = ref('ALL')
    const issues = ref([])
    const loading = ref(false)
    const error = ref(null)

    // 이슈 목록 가져오기
    const fetchIssues = async () => {
      loading.value = true
      error.value = null
      console.log('필터 적용:', statusFilter.value)
      try {
        let url = 'http://localhost:8080/issues'
        if (statusFilter.value !== 'ALL') {
          url += `?status=${statusFilter.value}`
        }
        console.log('요청 URL:', url)
        const response = await fetch(url)
        if (!response.ok) {
          throw new Error('이슈 목록을 가져오는 중 오류가 발생했습니다.')
        }
        const data = await response.json()
        console.log('받은 데이터:', data)
        // 백엔드에서 null을 반환할 경우 빈 배열로 처리
        issues.value = data === null ? [] : data
      } catch (err) {
        console.error('이슈 목록 가져오기 오류:', err)
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    // 상태 필터 변경 시 이슈 다시 가져오기
    const changeFilter = (status) => {
      statusFilter.value = status
      fetchIssues()
    }

    // 상태별 필터링
    const filteredIssues = computed(() => {
      return issues.value
    })

    // 이슈 상세 페이지로 이동
    const goToIssueDetail = (issueId) => {
      router.push(`/issues/${issueId}`)
    }

    // 날짜 포맷팅 함수
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('ko-KR', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    }

    onMounted(() => {
      fetchIssues()
    })

    return {
      statusFilter,
      filteredIssues,
      changeFilter,
      goToIssueDetail,
      formatDate,
      loading,
      error,
    }
  },
}
</script>

<style scoped>
.issue-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.issue-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-section {
  margin-bottom: 20px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-group select {
  width: auto;
}

.issue-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.issue-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
}

.issue-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.issue-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.issue-card-header h3 {
  margin: 0;
  font-size: 18px;
}

.status-badge {
  padding: 5px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.status-pending {
  background-color: #f8d7da;
  color: #721c24;
}

.status-in-progress {
  background-color: #cce5ff;
  color: #004085;
}

.status-completed {
  background-color: #d4edda;
  color: #155724;
}

.status-cancelled {
  background-color: #e2e3e5;
  color: #383d41;
}

.issue-card-body {
  margin-top: 10px;
}

.issue-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #6c757d;
}

.issue-assignee {
  text-align: left;
}

.issue-date {
  text-align: right;
}

.issue-card-footer {
  margin-top: 15px;
  width: 100%;
}

.loading {
  text-align: center;
  padding: 40px 20px;
  font-size: 18px;
  color: #6c757d;
}

.spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top-color: #007bff;
  animation: spin 1s ease-in-out infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.no-issues {
  text-align: center;
  padding: 40px 20px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.empty-state h3 {
  margin: 0 0 10px 0;
  color: #343a40;
}

.empty-state p {
  margin: 0 0 20px 0;
  color: #6c757d;
}

.btn-outline-primary {
  color: #007bff;
  border-color: #007bff;
  background-color: transparent;
  padding: 8px 16px;
  border-radius: 4px;
  text-decoration: none;
  transition: all 0.2s;
}

.btn-outline-primary:hover {
  color: #fff;
  background-color: #007bff;
}
</style>
