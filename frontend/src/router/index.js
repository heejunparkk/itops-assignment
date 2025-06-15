import { createRouter, createWebHistory } from 'vue-router'
import IssueList from '../pages/IssueList.vue'
import IssueForm from '../pages/IssueForm.vue'

const routes = [
  {
    path: '/',
    redirect: '/issues',
  },
  {
    path: '/issues',
    name: 'IssueList',
    component: IssueList,
  },
  {
    path: '/issues/new',
    component: IssueForm,
  },
  {
    path: '/issues/:id',
    component: IssueForm,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
