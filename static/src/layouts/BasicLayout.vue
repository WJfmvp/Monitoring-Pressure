<template>
  <el-container class="layout-root">
    <el-aside :width="collapsed ? '64px' : '220px'" class="layout-aside">
      <div class="brand">
        <el-icon><DataAnalysis /></el-icon>
        <span v-show="!collapsed">学业压力监测</span>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="collapsed"
        router
        background-color="#1f2d3d"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <template v-for="item in menus" :key="item.path">
          <el-sub-menu v-if="item.children?.length" :index="item.path">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.path"
              :index="child.path"
            >
              {{ child.title }}
            </el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb>
            <el-breadcrumb-item>{{ roleLabel }}</el-breadcrumb-item>
            <el-breadcrumb-item>{{ route.meta.title || '' }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <el-dropdown trigger="click" @command="handleCommand">
          <span class="user-trigger">
            <el-avatar :size="28">{{ userInitial }}</el-avatar>
            <span class="user-name">{{ userStore.userInfo?.username || '未登录' }}</span>
            <el-icon><CaretBottom /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                个人中心
              </el-dropdown-item>
              <el-dropdown-item command="complete">
                <el-icon><Edit /></el-icon>
                完善资料
              </el-dropdown-item>
              <el-dropdown-item command="settings">
                <el-icon><Setting /></el-icon>
                偏好设置
              </el-dropdown-item>
              <el-dropdown-item command="help" divided>
                <el-icon><QuestionFilled /></el-icon>
                帮助中心
              </el-dropdown-item>
              <el-dropdown-item command="about">
                <el-icon><InfoFilled /></el-icon>
                关于平台
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>

      <el-main class="layout-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const collapsed = ref(false)

const ROLE_LABEL = { student: '学生', teacher: '教师', admin: '管理员' }
const roleLabel = computed(() => ROLE_LABEL[userStore.role] || '用户')

const userInitial = computed(() => {
  const name = userStore.userInfo?.username || ''
  return name ? name.charAt(0).toUpperCase() : 'U'
})

// 跨角色公共菜单（拼接到每个角色菜单尾部）
const SHARED_TAIL = [
  { path: '/wellness', title: '放松一下', icon: 'Sunny' },
  { path: '/help', title: '帮助中心', icon: 'QuestionFilled' },
]

// 菜单按角色分发
const STUDENT_MENU_BASE = [
  { path: '/student/home', title: '主页', icon: 'House' },
  {
    path: '/student/group',
    title: '我的数据',
    icon: 'DataLine',
    children: [
      { path: '/student/academic', title: '我的成绩' },
      { path: '/student/psychological', title: '心理测评' },
      { path: '/student/stress', title: '压力检测结果' },
      { path: '/student/intervention', title: '干预建议' },
    ],
  },
]

const TEACHER_MENU_BASE = [
  { path: '/teacher/home', title: '主页', icon: 'House' },
  { path: '/teacher/class/overview', title: '班级压力分布', icon: 'TrendCharts' },
  {
    path: '/teacher/group',
    title: '学生数据',
    icon: 'User',
    children: [
      { path: '/teacher/student/academic', title: '学生成绩' },
      { path: '/teacher/student/stress', title: '学生压力' },
      { path: '/teacher/student/intervention', title: '干预记录' },
    ],
  },
]

const ADMIN_MENU_BASE = [
  { path: '/admin/home', title: '主页', icon: 'House' },
  { path: '/admin/users', title: '用户管理', icon: 'UserFilled' },
  {
    path: '/admin/academic-group',
    title: '成绩数据',
    icon: 'Document',
    children: [
      { path: '/admin/academic-import', title: 'Excel 导入' },
      { path: '/admin/academic-imports', title: '导入记录' },
      { path: '/admin/academic', title: '成绩列表' },
    ],
  },
  {
    path: '/admin/stress-group',
    title: '压力检测',
    icon: 'TrendCharts',
    children: [
      { path: '/admin/stress', title: '评估列表' },
      { path: '/admin/stress-predict', title: '运行预测' },
    ],
  },
  { path: '/admin/intervention', title: '干预建议', icon: 'ChatLineSquare' },
  { path: '/admin/statistics', title: '数据统计', icon: 'DataAnalysis' },
]

const STUDENT_MENU = [...STUDENT_MENU_BASE, ...SHARED_TAIL]
const TEACHER_MENU = [...TEACHER_MENU_BASE, ...SHARED_TAIL]
const ADMIN_MENU = [...ADMIN_MENU_BASE, ...SHARED_TAIL]

const menus = computed(() => {
  if (userStore.role === 'student') return STUDENT_MENU
  if (userStore.role === 'teacher') return TEACHER_MENU
  if (userStore.role === 'admin') return ADMIN_MENU
  return []
})

function handleCommand(cmd) {
  const map = {
    profile: '/profile',
    complete: '/profile/complete',
    settings: '/settings',
    help: '/help',
    about: '/about',
  }
  if (map[cmd]) {
    router.push(map[cmd])
    return
  }
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}

onMounted(async () => {
  // 已有 token 但本地缓存空时，刷新一下用户信息
  if (userStore.isLogin && !userStore.userInfo) {
    try {
      await userStore.fetchProfile()
    } catch (_) {
      /* 拦截器已处理 */
    }
  }
})
</script>

<style scoped>
.layout-root {
  height: 100vh;
}

.layout-aside {
  background: #1f2d3d;
  transition: width 0.2s;
  overflow-x: hidden;
}

.layout-aside .el-menu {
  border-right: none;
}

.brand {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 18px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  background: #15212e;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #303133;
}

.user-name {
  font-size: 14px;
}

.layout-main {
  background: #f3f5f8;
  padding: 18px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>