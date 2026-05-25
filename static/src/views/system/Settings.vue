<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">偏好设置</h2>
      <p class="hint">本地保存，不会同步到服务器。</p>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">外观</h3>

          <div class="setting-row">
            <div>
              <h4>主题色</h4>
              <p class="desc">选择系统的主色调，影响按钮、链接等。</p>
            </div>
            <el-radio-group v-model="prefs.themeColor" @change="applyTheme">
              <el-radio-button
                v-for="t in THEMES"
                :key="t.value"
                :value="t.value"
              >
                <span class="theme-dot" :style="{ background: t.color }"></span>
                {{ t.label }}
              </el-radio-button>
            </el-radio-group>
          </div>

          <el-divider />

          <div class="setting-row">
            <div>
              <h4>字号</h4>
              <p class="desc">仅影响正文区域，菜单和图表保持原样。</p>
            </div>
            <el-radio-group v-model="prefs.fontSize" @change="applyFontSize">
              <el-radio-button value="small">小</el-radio-button>
              <el-radio-button value="medium">默认</el-radio-button>
              <el-radio-button value="large">大</el-radio-button>
            </el-radio-group>
          </div>

          <el-divider />

          <div class="setting-row">
            <div>
              <h4>页面切换动画</h4>
              <p class="desc">关闭后，页面切换更迅速但少了过渡效果。</p>
            </div>
            <el-switch v-model="prefs.transitions" @change="savePrefs" />
          </div>

          <el-divider />

          <div class="setting-row">
            <div>
              <h4>表格密度</h4>
              <p class="desc">紧凑模式可以一屏显示更多数据。</p>
            </div>
            <el-radio-group v-model="prefs.tableDensity" @change="savePrefs">
              <el-radio-button value="loose">宽松</el-radio-button>
              <el-radio-button value="default">默认</el-radio-button>
              <el-radio-button value="compact">紧凑</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </el-col>

      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">账号与安全</h3>

          <div class="account-row">
            <div>
              <h4>当前账号</h4>
              <p class="desc">{{ userStore.userInfo?.username || '-' }} ·
                {{ ROLE_LABEL[userStore.role] || userStore.role }}</p>
              <p class="desc">{{ userStore.userInfo?.telephone || '-' }}</p>
            </div>
          </div>

          <el-divider />

          <el-button @click="$router.push('/reset-password')" style="width: 100%">
            <el-icon style="margin-right: 4px"><Lock /></el-icon>
            修改密码
          </el-button>

          <el-button
            @click="$router.push('/profile/complete')"
            style="width: 100%; margin: 12px 0 0 0"
          >
            <el-icon style="margin-right: 4px"><User /></el-icon>
            完善资料
          </el-button>

          <el-button @click="onLogout" type="danger" plain style="width: 100%; margin-top: 12px">
            <el-icon style="margin-right: 4px"><SwitchButton /></el-icon>
            退出登录
          </el-button>
        </div>

        <div class="page-card" style="margin-top: 16px">
          <h3 class="card-title">本地数据</h3>
          <p class="desc">心情打卡、番茄钟等本地数据存储概况。</p>
          <ul class="storage-list">
            <li>心情打卡：{{ storageStats.mood }} 条</li>
            <li>番茄钟记录：{{ storageStats.pomo }} 天</li>
            <li>偏好设置：{{ storageStats.prefs ? '已保存' : '默认' }}</li>
          </ul>
          <el-button @click="resetPrefs" plain style="margin-right: 8px">重置默认偏好</el-button>
          <el-popconfirm title="确定清空所有本地数据？此操作不可恢复" @confirm="clearAll">
            <template #reference>
              <el-button type="danger" plain>清空本地数据</el-button>
            </template>
          </el-popconfirm>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const PREFS_KEY = 'mp_prefs'
const ROLE_LABEL = { student: '学生', teacher: '教师', admin: '管理员' }

const THEMES = [
  { value: 'blue', label: '默认蓝', color: '#409eff' },
  { value: 'purple', label: '紫色', color: '#a777e3' },
  { value: 'green', label: '青绿', color: '#42b883' },
  { value: 'orange', label: '橙色', color: '#fa8c16' },
]

const defaultPrefs = () => ({
  themeColor: 'blue',
  fontSize: 'medium',
  transitions: true,
  tableDensity: 'default',
})

const prefs = reactive(defaultPrefs())

const storageStats = computed(() => {
  let mood = 0
  let pomo = 0
  try {
    mood = Object.keys(JSON.parse(localStorage.getItem('mp_mood_log') || '{}')).length
    pomo = Object.keys(JSON.parse(localStorage.getItem('mp_pomo_log') || '{}')).length
  } catch {}
  return {
    mood,
    pomo,
    prefs: !!localStorage.getItem(PREFS_KEY),
  }
})

function loadPrefs() {
  try {
    const raw = JSON.parse(localStorage.getItem(PREFS_KEY) || '{}')
    Object.assign(prefs, { ...defaultPrefs(), ...raw })
  } catch {
    Object.assign(prefs, defaultPrefs())
  }
  applyAll()
}

function savePrefs() {
  localStorage.setItem(PREFS_KEY, JSON.stringify(prefs))
  applyAll()
}

const THEME_MAP = {
  blue: '#409eff',
  purple: '#a777e3',
  green: '#42b883',
  orange: '#fa8c16',
}

function applyTheme() {
  const color = THEME_MAP[prefs.themeColor] || '#409eff'
  document.documentElement.style.setProperty('--el-color-primary', color)
  savePrefs()
}

function applyFontSize() {
  const map = { small: '13px', medium: '14px', large: '16px' }
  document.documentElement.style.setProperty('--mp-base-font', map[prefs.fontSize])
  document.body.style.fontSize = map[prefs.fontSize]
  savePrefs()
}

function applyAll() {
  applyTheme()
  applyFontSize()
}

function resetPrefs() {
  Object.assign(prefs, defaultPrefs())
  savePrefs()
  ElMessage.success('已恢复默认偏好')
}

function clearAll() {
  ['mp_prefs', 'mp_mood_log', 'mp_pomo_log'].forEach(k => localStorage.removeItem(k))
  Object.assign(prefs, defaultPrefs())
  applyAll()
  ElMessage.success('已清空本地数据')
}

function onLogout() {
  userStore.logout()
  router.push('/login')
}

onMounted(loadPrefs)
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 6px 0 0 0;
}
.card-title {
  margin: 0 0 12px 0;
}
.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 8px 0;
}
.setting-row h4 {
  margin: 0 0 4px 0;
}
.setting-row .desc {
  color: #909399;
  font-size: 12px;
  margin: 0;
}
.theme-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 4px;
  vertical-align: middle;
}
.account-row h4 {
  margin: 0 0 4px 0;
}
.account-row .desc {
  color: #909399;
  font-size: 13px;
  margin: 2px 0;
}
.storage-list {
  color: #606266;
  padding-left: 18px;
  margin: 8px 0 12px 0;
}
</style>