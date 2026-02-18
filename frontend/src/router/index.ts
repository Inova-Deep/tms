import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'
import MyTrainingView from '@/views/MyTrainingView.vue'
import EmployeesView from '@/views/EmployeesView.vue'
import EmployeeDetailView from '@/views/EmployeeDetailView.vue'
import ProfilesView from '@/views/ProfilesView.vue'
import SessionsView from '@/views/SessionsView.vue'
import CertificationsView from '@/views/CertificationsView.vue'
import CoursesView from '@/views/CoursesView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', component: DashboardView },
    { path: '/employees', component: EmployeesView },
    { path: '/employees/:id', component: EmployeeDetailView },
    { path: '/profiles', component: ProfilesView },
    { path: '/sessions', component: SessionsView },
    { path: '/certifications', component: CertificationsView },
    { path: '/courses', component: CoursesView },
    { path: '/my-training', component: MyTrainingView },
    { path: '/help', name: 'Help', component: () => import('@/views/HelpView.vue') },
    { path: '/:pathMatch(.*)*', component: NotFoundView },
  ],
})

export default router
