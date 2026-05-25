<template>
  <div>
    <div class="page-card profile-hero" :class="roleClass">
      <div class="avatar">
        <span>{{ userInitial }}</span>
      </div>
      <div class="info">
        <h2 class="username">{{ info?.username || '未命名用户' }}</h2>
        <div class="info-row">
          <el-tag :type="ROLE_TAG[info?.role]">{{ roleLabel }}</el-tag>
          <span class="meta">学号/ID {{ info?.user_id }}</span>
          <span class="meta">·</span>
          <span class="meta">{{ info?.telephone || '-' }}</span>
        </div>
        <p class="bio">{{ welcomeTip }}</p>
      </div>
      <div class="actions">
        <el-button type="primary" @click="$router.push('/profile/complete')">
          <el-icon style="margin-right: 4px"><Edit /></el-icon>
          完善资料
        </el-button>
        <el-button @click="onRefresh" :loading="loading">
          <el-icon style="margin-right: 4px"><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户 ID">{{ info?.user_id }}</el-descriptions-item>
            <el-descriptions-item label="角色">{{ roleLabel }}</el-descriptions-item>
            <el-descriptions-item label="用户名">{{ info?.username || '-' }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ info?.telephone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="性别">{{ sexLabel }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ info?.email || '-' }}</el-descriptions-item>
          </el-descriptions>

          <div class="completion">
            <span>资料完整度</span>
            <el-progress :percentage="completionRate" :stroke-width="10" />
          </div>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">常用功能</h3>
          <div class="quick-list">
            <div class="quick-item" v-for="q in quickItems" :key="q.label" @click="$router.push(q.path)">
              <el-icon :size="20" :color="q.color"><component :is="q.icon" /></el-icon>
              <div class="quick-text">
                <h4>{{ q.label }}</h4>
                <p>{{ q.desc }}</p>
              </div>
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">账号安全</h3>
          <div class="security-row">
            <div class="security-info">
              <h4>登录密码</h4>
              <p>建议每 3 个月更新一次密码</p>
            </div>
            <el-button @click="$router.push('/reset-password')">修改密码</el-button>
          </div>
          <el-divider />
          <div class="security-row">
            <div class="security-info">
              <h4>偏好设置</h4>
              <p>主题色、字号、动画等本地偏好</p>
            </div>
            <el-button @click="$router.push('/settings')">前往设置</el-button>
          </div>
          <el-divider />
          <div class="security-row">
            <div class="security-info">
              <h4>退出登录</h4>
              <p>清除当前 Token 和本地缓存的用户信息</p>
            </div>
            <el-button type="danger" plain @click="onLogout">退出</el-button>
          </div>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card help-card">
          <h3 class="card-title">需要帮助？</h3>
          <p class="hint">查看常见问题或了解更多平台信息。</p>
          <el-button @click="$router.push('/help')" style="margin-right: 8px">
            <el-icon style="margin-right: 4px"><QuestionFilled /></el-icon>
            帮助中心
          </el-button>
          <el-button @click="$router.push('/about')">
            <el-icon style="margin-right: 4px"><InfoFilled /></el-icon>
            关于平台
          </el-button>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

const info = computed(() => userStore.userInfo)

const ROLE_LABEL = { student: '学生', teacher: '教师', admin: '管理员' }
const ROLE_TAG = { student: '', teacher: 'success', admin: 'danger' }
const roleLabel = computed(() => ROLE_LABEL[info.value?.role] || info.value?.role || '-')
const roleClass = computed(() => `role-${info.value?.role || 'guest'}`)

const sexLabel = computed(() => {
  const s = info.value?.sex
  if (s === 1) return '男'
  if (s === 2) return '女'
  return '未填写'
})

const userInitial = computed(() => {
  const name = info.value?.username || ''
  return name ? name.charAt(0).toUpperCase() : 'U'
})

const welcomeTip = computed(() => {
  if (info.value?.role === 'admin') return '管理员账号，拥有平台全部功能权限。'
  if (info.value?.role === 'teacher') return '教师账号，可查看班级整体压力分布和学生诊断书。'
  return '欢迎使用学业压力监测平台，记得定期完成心理测评。'
})

const completionRate = computed(() => {
  if (!info.value) return 0
  let filled = 0
  if (info.value.username) filled++
  if (info.value.telephone) filled++
  if (info.value.email) filled++
  if (info.value.sex) filled++
  return Math.round((filled / 4) * 100)
})

const quickItems = computed(() => {
  if (info.value?.role === 'admin') {
    return [
      { label: '管理员主页', desc: '系统概览与快捷入口', path: '/admin/home', icon: 'House', color: '#409eff' },
      { label: '用户管理', desc: '查看和修改用户角色', path: '/admin/users', icon: 'UserFilled', color: '#67c23a' },
      { label: '数据统计', desc: '跨表分析与图表', path: '/admin/statistics', icon: 'TrendCharts', color: '#a777e3' },
    ]
  }
  if (info.value?.role === 'teacher') {
    return [
      { label: '班级压力分布', desc: '一眼看完全班状态', path: '/teacher/class/overview', icon: 'TrendCharts', color: '#409eff' },
      { label: '学生数据查询', desc: '按学号查具体数据', path: '/teacher/home', icon: 'User', color: '#67c23a' },
      { label: '放松一下', desc: '番茄钟 / 呼吸训练', path: '/wellness', icon: 'Sunny', color: '#e6a23c' },
    ]
  }
  return [
    { label: '心理测评', desc: '记录今天的状态', path: '/student/psychological', icon: 'EditPen', color: '#a777e3' },
    { label: '我的成绩', desc: '查看历史成绩', path: '/student/academic', icon: 'Document', color: '#409eff' },
    { label: '压力检测结果', desc: '看看最新评估', path: '/student/stress', icon: 'DataAnalysis', color: '#67c23a' },
    { label: '放松一下', desc: '番茄钟 / 呼吸训练', path: '/wellness', icon: 'Sunny', color: '#e6a23c' },
  ]
})

async function onRefresh() {
  loading.value = true
  try {
    await userStore.fetchProfile()
  } finally {
    loading.value = false
  }
}

function onLogout() {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  if (!info.value) onRefresh()
})
</script>

<style scoped>
.profile-hero {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px 32px;
  background: linear-gradient(120deg, #6e8efb 0%, #a777e3 100%);
  color: #fff;
}
.profile-hero.role-teacher {
  background: linear-gradient(120deg, #43c6ac 0%, #6ec96e 100%);
}
.profile-hero.role-admin {
  background: linear-gradient(120deg, #1f2d3d 0%, #485563 100%);
}
.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30px;
  font-weight: 600;
  border: 3px solid rgba(255, 255, 255, 0.3);
}
.info {
  flex: 1;
}
.username {
  margin: 0 0 6px 0;
  color: #fff;
  font-size: 24px;
}
.info-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.info-row .meta {
  color: rgba(255, 255, 255, 0.9);
  font-size: 13px;
}
.bio {
  margin: 8px 0 0 0;
  color: rgba(255, 255, 255, 0.85);
  font-size: 13px;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}

.completion {
  margin-top: 18px;
  color: #606266;
  font-size: 13px;
}
.completion span {
  display: block;
  margin-bottom: 4px;
}

.quick-list .quick-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  background: #f8f9fc;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.15s;
}
.quick-list .quick-item:hover {
  background: #ecf5ff;
  transform: translateX(2px);
}
.quick-text {
  flex: 1;
}
.quick-text h4 {
  margin: 0 0 2px 0;
  font-size: 14px;
}
.quick-text p {
  margin: 0;
  color: #909399;
  font-size: 12px;
}

.security-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.security-info h4 {
  margin: 0 0 4px 0;
}
.security-info p {
  margin: 0;
  color: #909399;
  font-size: 12px;
}
.help-card .hint {
  color: #909399;
  margin: 0 0 12px 0;
}
</style>