<template>
  <div class="issue-form">
    <div class="card">
      <div class="card-header">
        <h2>{{ isNewIssue ? '새 이슈 생성' : '이슈 상세/수정' }}</h2>
      </div>
      <div class="card-body">
        <form @submit.prevent="saveIssue">
          <div class="form-group">
            <label for="title">제목</label>
            <input
              type="text"
              id="title"
              v-model="issue.title"
              class="form-control"
              :disabled="isIssueReadOnly"
              required
            />
          </div>

          <div class="form-group">
            <label for="description">설명</label>
            <textarea
              id="description"
              v-model="issue.description"
              class="form-control"
              rows="5"
              :disabled="isIssueReadOnly"
              required
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="status">상태</label>
              <select
                id="status"
                v-model="issue.status"
                class="form-control"
                :disabled="!issue.user || isIssueReadOnly"
              >
                <option value="PENDING">PENDING</option>
                <option value="IN_PROGRESS" :disabled="!issue.user">IN_PROGRESS</option>
                <option value="COMPLETED" :disabled="!issue.user">COMPLETED</option>
                <option value="CANCELLED" :disabled="!issue.user">CANCELLED</option>
              </select>
              <small v-if="!issue.user" class="form-text text-muted">
                담당자가 지정되어야 상태를 변경할 수 있습니다.
              </small>
            </div>

            <div class="form-group">
              <label for="user">담당자</label>
              <select
                id="user"
                v-model="selectedUserId"
                class="form-control"
                :disabled="isUserDisabled"
              >
                <option :value="null">담당자 없음</option>
                <option v-for="user in users" :key="user.id" :value="user.id">
                  {{ user.name }}
                </option>
              </select>
            </div>
          </div>

          <div class="form-meta" v-if="!isNewIssue">
            <div class="meta-item">
              <strong>생성일:</strong> {{ formatDate(issue.createdAt) }}
            </div>
            <div class="meta-item">
              <strong>수정일:</strong> {{ formatDate(issue.updatedAt) }}
            </div>
          </div>

          <div class="form-actions">
            <button
              type="button"
              class="btn btn-secondary"
              @click="goToList"
            >
              목록으로
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              :disabled="isIssueReadOnly"
            >
              저장
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

export default {
  name: 'IssueForm',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const issueId = computed(() => route.params.id)
    const isNewIssue = computed(() => !issueId.value)
    const loading = ref(false)
    const selectedUserId = ref(null)

    // 기본 이슈 객체
    const defaultIssue = {
      id: null,
      title: '',
      description: '',
      status: 'PENDING',
      user: null,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }
    
    // 이슈의 초기 상태를 저장하기 위한 변수
    const initialStatus = ref('')

    const issue = ref({ ...defaultIssue })

    // 이슈 데이터 불러오기
    const fetchIssue = async () => {
      if (isNewIssue.value) {
        issue.value = { ...defaultIssue }
        return
      }

      loading.value = true
      try {
        const response = await fetch(`http://localhost:8080/issue/${issueId.value}`)
        if (!response.ok) {
          throw new Error('이슈를 불러오는 중 오류가 발생했습니다.')
        }
        const data = await response.json()
        issue.value = { ...data }
        // 초기 상태 저장
        initialStatus.value = data.status
        // 담당자가 있으면 선택된 사용자 ID 설정
        if (data.user) {
          selectedUserId.value = data.user.id
        } else {
          selectedUserId.value = null
        }
      } catch (err) {
        console.error('이슈 로딩 오류:', err)
        alert('이슈를 불러오는 중 오류가 발생했습니다.')
        router.push('/issues')
      } finally {
        loading.value = false
      }
    }

    // 사용자 목록
    const users = ref([])
    
    // 사용자 목록 불러오기
    const loadUsers = async () => {
      try {
        const response = await fetch('http://localhost:8080/users')
        if (!response.ok) {
          throw new Error('사용자 목록을 불러오는 중 오류가 발생했습니다.')
        }
        users.value = await response.json()
      } catch (err) {
        console.error('사용자 목록 로딩 오류:', err)
      }
    }

    // 이슈 저장
    const saveIssue = async () => {
      loading.value = true
      try {
        const issueData = {
          ...issue.value,
          updatedAt: new Date().toISOString(),
        }

        let url = 'http://localhost:8080/issue'
        let method = 'POST'

        if (!isNewIssue.value) {
          url = `http://localhost:8080/issue/${issueId.value}`
          method = 'PATCH'
        }

        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(issueData),
        })

        if (!response.ok) {
          throw new Error('이슈를 저장하는 중 오류가 발생했습니다.')
        }

        router.push('/issues')
      } catch (err) {
        console.error('이슈 저장 오류:', err)
        alert('이슈를 저장하는 중 오류가 발생했습니다.')
      } finally {
        loading.value = false
      }
    }

    // 목록으로 돌아가기
    const goToList = () => {
      router.push('/issues')
    }

    // 날짜 포맷팅 함수
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('ko-KR', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // 이슈가 읽기 전용인지 여부 (초기 상태가 COMPLETED 또는 CANCELLED인 경우)
    const isIssueReadOnly = computed(() => {
      // 새 이슈 생성 시에는 항상 편집 가능
      if (isNewIssue.value) return false;
      // 기존 이슈 수정 시에는 초기 상태가 완료/취소인 경우에만 읽기 전용
      return initialStatus.value === 'COMPLETED' || initialStatus.value === 'CANCELLED'
    })

    // 상태 변경 가능 여부
    const isStatusDisabled = computed(() => {
      // 새 이슈 생성 시에는 담당자가 없어도 상태 변경 가능
      if (isNewIssue.value) return !issue.value.user;
      // 기존 이슈 수정 시에는 담당자가 없거나 완료/취소 상태면 비활성화
      return !issue.value.user || isIssueReadOnly.value
    })

    // 담당자 변경 가능 여부
    const isUserDisabled = computed(() => {
      // 새 이슈 생성 시에는 항상 담당자 변경 가능
      if (isNewIssue.value) return false;
      // 기존 이슈 수정 시에만 완료/취소 상태인 경우 담당자 변경 불가
      return isIssueReadOnly.value
    })

    // 담당자 변경 감지
    watch(selectedUserId, (newValue) => {
      if (newValue === null) {
        // 담당자가 제거되면 상태는 PENDING으로 변경
        issue.value.user = null
        issue.value.status = 'PENDING'
      } else {
        // 담당자가 선택되면 즉시 user 객체 업데이트
        const selectedUser = users.value.find(u => u.id === newValue)
        issue.value.user = selectedUser
        
        // 새 이슈 생성 시에만 자동으로 상태 변경 적용
        // 기존 이슈 수정 시에는 사용자가 직접 설정한 상태 유지
        if (isNewIssue.value && issue.value.status === 'PENDING') {
          issue.value.status = 'IN_PROGRESS'
        }
      }
    })

    onMounted(() => {
      loadUsers()
      fetchIssue()
    })

    return {
      issue,
      isNewIssue,
      loading,
      users,
      selectedUserId,
      saveIssue,
      goToList,
      formatDate,
      isIssueReadOnly,
      isStatusDisabled,
      isUserDisabled,
      initialStatus
    }
  }
}
</script>

<style scoped>
.issue-form {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  background-color: #f8f9fa;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
}

.card-body {
  padding: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-meta {
  display: flex;
  gap: 2rem;
  margin: 1.5rem 0;
  color: #666;
  font-size: 0.9rem;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
}

textarea.form-control {
  resize: vertical;
}
</style>
