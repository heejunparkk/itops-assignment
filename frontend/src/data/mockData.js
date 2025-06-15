// 미리 정의된 사용자 목록
export const users = [
  { id: 1, name: '김개발' },
  { id: 2, name: '이디자인' },
  { id: 3, name: '박기획' },
]

// 목 이슈 데이터
export const issues = [
  {
    id: 1,
    title: '로그인 페이지 버그 수정',
    description: '사용자가 로그인 시 간헐적으로 오류가 발생하는 문제를 해결해야 합니다.',
    status: 'IN_PROGRESS',
    user: { id: 1, name: '김개발' },
    createdAt: '2025-05-28T09:30:00Z',
    updatedAt: '2025-05-28T14:20:00Z',
  },
  {
    id: 2,
    title: '메인 페이지 UI 개선',
    description: '메인 페이지의 레이아웃을 더 직관적으로 개선하고 사용자 경험을 향상시켜야 합니다.',
    status: 'PENDING',
    userId: null,
    user: null,
    createdAt: '2025-05-29T11:15:00Z',
    updatedAt: '2025-05-29T11:15:00Z',
  },
  {
    id: 3,
    title: '데이터베이스 성능 최적화',
    description: '쿼리 성능이 느려서 페이지 로딩 시간이 길어지는 문제를 해결했습니다.',
    status: 'COMPLETED',
    user: { id: 3, name: '박기획' },
    createdAt: '2025-05-30T08:00:00Z',
    updatedAt: '2025-05-30T16:45:00Z',
  },
  {
    id: 4,
    title: '모바일 반응형 이슈',
    description:
      '모바일 환경에서 레이아웃이 깨지는 문제가 발견되었으나 우선순위가 낮아 취소되었습니다.',
    status: 'CANCELLED',
    user: { id: 2, name: '이디자인' },
    createdAt: '2025-05-31T13:20:00Z',
    updatedAt: '2025-05-31T10:30:00Z',
  },
]
