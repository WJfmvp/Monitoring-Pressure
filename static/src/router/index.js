import { createRouter, createWebHashHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const Login = () => import('@/views/account/Login.vue')
const Register = () => import('@/views/account/Register.vue')
const ResetPassword = () => import('@/views/account/ResetPassword.vue')
const BasicLayout = () => import('@/layouts/BasicLayout.vue')

const Profile = () => import('@/views/profile/Profile.vue')
const CompleteInfo = () => import('@/views/profile/CompleteInfo.vue')

// shared
const Wellness = () => import('@/views/wellness/Wellness.vue')
const Help = () => import('@/views/system/Help.vue')
const About = () => import('@/views/system/About.vue')
const Settings = () => import('@/views/system/Settings.vue')

// student
const StudentHome = () => import('@/views/student/Home.vue')
const StudentPsychological = () => import('@/views/student/Psychological.vue')
const StudentAcademic = () => import('@/views/student/Academic.vue')
const StudentStressResult = () => import('@/views/student/StressResult.vue')
const StudentIntervention = () => import('@/views/student/Intervention.vue')

// teacher
const TeacherHome = () => import('@/views/teacher/Home.vue')
const TeacherAcademic = () => import('@/views/teacher/StudentAcademic.vue')
const TeacherStress = () => import('@/views/teacher/StudentStress.vue')
const TeacherIntervention = () => import('@/views/teacher/StudentIntervention.vue')
const TeacherClassOverview = () => import('@/views/teacher/ClassOverview.vue')
const TeacherStudentDiagnosis = () => import('@/views/teacher/StudentDiagnosis.vue')

// admin
const AdminHome = () => import('@/views/admin/Home.vue')
const AdminUserList = () => import('@/views/admin/UserList.vue')
const AdminAcademicImport = () => import('@/views/admin/AcademicImport.vue')
const AdminAcademicImportList = () => import('@/views/admin/AcademicImportList.vue')
const AdminAcademicList = () => import('@/views/admin/AcademicList.vue')
const AdminStressList = () => import('@/views/admin/StressResultList.vue')
const AdminStressPredict = () => import('@/views/admin/StressPredict.vue')
const AdminIntervention = () => import('@/views/admin/Intervention.vue')
const AdminStatistics = () => import('@/views/admin/Statistics.vue')

// 角色常量
export const ROLE = {
  STUDENT: 'student',
  TEACHER: 'teacher',
  ADMIN: 'admin',
}

const routes = [
  { path: '/login', component: Login, meta: { public: true, title: '登录' } },
  { path: '/register', component: Register, meta: { public: true, title: '注册' } },
  { path: '/reset-password', component: ResetPassword, meta: { public: true, title: '重置密码' } },

  {
    path: '/',
    component: BasicLayout,
    redirect: '/dispatch',
    children: [
      // 角色分发占位（守卫里会重定向）
      { path: 'dispatch', component: { template: '<div></div>' }, meta: { title: '加载中' } },

      { path: 'profile', component: Profile, meta: { title: '个人中心' } },
      { path: 'profile/complete', component: CompleteInfo, meta: { title: '完善资料' } },

      // 跨角色共用
      { path: 'wellness', component: Wellness, meta: { title: '放松一下' } },
      { path: 'help', component: Help, meta: { title: '帮助中心' } },
      { path: 'about', component: About, meta: { title: '关于平台' } },
      { path: 'settings', component: Settings, meta: { title: '偏好设置' } },

      // 学生
      {
        path: 'student/home',
        component: StudentHome,
        meta: { title: '学生主页', roles: [ROLE.STUDENT, ROLE.ADMIN] },
      },
      {
        path: 'student/psychological',
        component: StudentPsychological,
        meta: { title: '心理测评', roles: [ROLE.STUDENT, ROLE.ADMIN] },
      },
      {
        path: 'student/academic',
        component: StudentAcademic,
        meta: { title: '我的成绩', roles: [ROLE.STUDENT, ROLE.ADMIN] },
      },
      {
        path: 'student/stress',
        component: StudentStressResult,
        meta: { title: '压力检测结果', roles: [ROLE.STUDENT, ROLE.ADMIN] },
      },
      {
        path: 'student/intervention',
        component: StudentIntervention,
        meta: { title: '干预建议', roles: [ROLE.STUDENT, ROLE.ADMIN] },
      },

      // 教师
      {
        path: 'teacher/home',
        component: TeacherHome,
        meta: { title: '教师主页', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },
      {
        path: 'teacher/class/overview',
        component: TeacherClassOverview,
        meta: { title: '班级压力分布', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },
      {
        path: 'teacher/student/diagnosis',
        component: TeacherStudentDiagnosis,
        meta: { title: '压力诊断书', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },
      {
        path: 'teacher/student/academic',
        component: TeacherAcademic,
        meta: { title: '学生成绩', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },
      {
        path: 'teacher/student/stress',
        component: TeacherStress,
        meta: { title: '学生压力', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },
      {
        path: 'teacher/student/intervention',
        component: TeacherIntervention,
        meta: { title: '学生干预记录', roles: [ROLE.TEACHER, ROLE.ADMIN] },
      },

      // 管理员
      {
        path: 'admin/home',
        component: AdminHome,
        meta: { title: '管理员主页', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/users',
        component: AdminUserList,
        meta: { title: '用户管理', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/academic-import',
        component: AdminAcademicImport,
        meta: { title: 'Excel 导入', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/academic-imports',
        component: AdminAcademicImportList,
        meta: { title: '导入记录', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/academic',
        component: AdminAcademicList,
        meta: { title: '成绩数据', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/stress',
        component: AdminStressList,
        meta: { title: '压力评估列表', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/stress-predict',
        component: AdminStressPredict,
        meta: { title: '压力预测', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/intervention',
        component: AdminIntervention,
        meta: { title: '干预建议管理', roles: [ROLE.ADMIN] },
      },
      {
        path: 'admin/statistics',
        component: AdminStatistics,
        meta: { title: '数据统计', roles: [ROLE.ADMIN] },
      },
    ],
  },

  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

const ROLE_HOME = {
  student: '/student/home',
  teacher: '/teacher/home',
  admin: '/admin/home',
}

router.beforeEach(to => {
  const userStore = useUserStore()
  const meta = to.meta || {}

  if (meta.title) document.title = `${meta.title} - 学业压力监测平台`

  if (meta.public) {
    // 已登录访问登录/注册页 → 直接跳到角色主页
    if (userStore.isLogin && (to.path === '/login' || to.path === '/register')) {
      return ROLE_HOME[userStore.role] || '/profile'
    }
    return true
  }

  if (!userStore.isLogin) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  // /dispatch → 按角色重定向到对应主页
  if (to.path === '/dispatch') {
    return ROLE_HOME[userStore.role] || '/profile'
  }

  // 角色权限校验
  if (Array.isArray(meta.roles) && meta.roles.length > 0) {
    if (!meta.roles.includes(userStore.role)) {
      ElMessage.warning('当前角色无访问权限')
      return ROLE_HOME[userStore.role] || '/profile'
    }
  }

  return true
})

export default router
