<template>
  <div v-loading="loading">
    <!-- 欢迎区 -->
    <div class="page-card hero">
      <div>
        <h2>{{ greeting }}，管理员 👋</h2>
        <p>{{ todayHint }}</p>
      </div>
      <div class="hero-stats">
        <div class="hero-stat">
          <span class="num">{{ stats.userCount }}</span>
          <span class="lbl">总用户</span>
        </div>
        <div class="hero-stat">
          <span class="num">{{ stats.stressTotal }}</span>
          <span class="lbl">评估记录</span>
        </div>
        <div class="hero-stat">
          <span class="num">{{ stats.warningCount }}</span>
          <span class="lbl">预警人次</span>
        </div>
      </div>
    </div>

    <!-- 4 个 metric -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6">
        <div class="metric-card primary">
          <span class="label">学生 / 教师 / 管理员</span>
          <span class="value">{{ stats.studentCount }} · {{ stats.teacherCount }} · {{ stats.adminCount }}</span>
          <span class="sub">点击进入用户管理</span>
          <span class="metric-link" @click="$router.push('/admin/users')">前往 →</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card warn">
          <span class="label">本周高压力学生</span>
          <span class="value">{{ stats.highRiskWeek }}</span>
          <span class="sub">需立即关注</span>
          <span class="metric-link" @click="$router.push('/admin/stress')">查看列表 →</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">导入批次</span>
          <span class="value">{{ stats.importCount }}</span>
          <span class="sub">最近上传</span>
          <span class="metric-link" @click="$router.push('/admin/academic-imports')">导入记录 →</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">干预建议模板</span>
          <span class="value">{{ stats.suggestionCount }}</span>
          <span class="sub">在用方案</span>
          <span class="metric-link" @click="$router.push('/admin/intervention')">维护 →</span>
        </div>
      </el-col>
    </el-row>

    <!-- 中部图表 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><PieChart /></el-icon>
            最近评估等级分布
          </h3>
          <v-chart v-if="hasStress" class="chart" :option="pieOption" autoresize />
          <el-empty v-else description="暂无评估数据" />
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card system-card">
          <h3 class="card-title">
            <el-icon><Monitor /></el-icon>
            系统状态
          </h3>
          <ul class="status-list">
            <li>
              <span class="dot ok"></span>
              <span class="key">Go 后端 :8080</span>
              <el-tag size="small" type="success">在线</el-tag>
            </li>
            <li>
              <span class="dot" :class="mlStatus.ok ? 'ok' : 'fail'"></span>
              <span class="key">Python ML :8000</span>
              <el-tag size="small" :type="mlStatus.ok ? 'success' : 'danger'">
                {{ mlStatus.ok ? '在线' : '离线' }}
              </el-tag>
            </li>
            <li>
              <span class="dot ok"></span>
              <span class="key">MySQL</span>
              <el-tag size="small" type="success">已连接</el-tag>
            </li>
            <li>
              <span class="dot info"></span>
              <span class="key">验证码存储</span>
              <el-tag size="small" type="info">内存兜底</el-tag>
            </li>
          </ul>
          <p class="status-tip">如运行预测出错，请确认 Python 服务在线</p>
        </div>
      </el-col>
    </el-row>

    <!-- 快捷操作 + 最近活动 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><Clock /></el-icon>
            最近评估活动
          </h3>
          <el-empty v-if="recentStress.length === 0" description="暂无数据" />
          <el-table v-else :data="recentStress" size="small">
            <el-table-column prop="user_id" label="学生 ID" width="100" />
            <el-table-column label="预警" width="100">
              <template #default="{ row }">
                <el-tag :type="LEVEL_TAG[row.warning_status]" size="small">
                  {{ LEVEL_LABEL[row.warning_status] }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total_score" label="综合" width="80" />
            <el-table-column prop="created_at" label="时间">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><Lightning /></el-icon>
            快捷操作
          </h3>
          <div class="quick-grid">
            <div class="quick-item" @click="$router.push('/admin/academic-import')">
              <el-icon :size="22" color="#409eff"><Upload /></el-icon>
              <span>Excel 导入</span>
            </div>
            <div class="quick-item" @click="$router.push('/admin/stress-predict')">
              <el-icon :size="22" color="#67c23a"><MagicStick /></el-icon>
              <span>运行预测</span>
            </div>
            <div class="quick-item" @click="$router.push('/admin/users')">
              <el-icon :size="22" color="#f56c6c"><UserFilled /></el-icon>
              <span>用户管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/admin/intervention')">
              <el-icon :size="22" color="#e6a23c"><ChatLineSquare /></el-icon>
              <span>干预建议</span>
            </div>
            <div class="quick-item" @click="$router.push('/admin/statistics')">
              <el-icon :size="22" color="#a777e3"><TrendCharts /></el-icon>
              <span>数据统计</span>
            </div>
            <div class="quick-item" @click="$router.push('/help')">
              <el-icon :size="22" color="#909399"><QuestionFilled /></el-icon>
              <span>帮助中心</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 这里必须别名，否则会遮蔽全局注册的 Element Plus 图标 <PieChart />
import { PieChart as EChartsPie } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import {
  getUserList,
  getStressResultList,
  getAcademicImportList,
  getInterventionSuggestionList,
} from '@/api/admin'

use([CanvasRenderer, EChartsPie, TooltipComponent, LegendComponent])

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']

const loading = ref(false)
const stressList = ref([])
const stats = reactive({
  userCount: 0,
  studentCount: 0,
  teacherCount: 0,
  adminCount: 0,
  stressTotal: 0,
  warningCount: 0,
  highRiskWeek: 0,
  importCount: 0,
  suggestionCount: 0,
})
const mlStatus = reactive({ ok: true })

const now = new Date()
const greeting = computed(() => {
  const h = now.getHours()
  if (h < 6) return '夜深了'
  if (h < 12) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})
const todayHint = computed(() => {
  if (stats.warningCount > 0) return `当前 ${stats.warningCount} 人次处于预警状态，建议优先处理。`
  return '今天没有新的预警，系统运行良好。'
})

const hasStress = computed(() => stressList.value.length > 0)

const pieOption = computed(() => {
  const buckets = [0, 0, 0]
  stressList.value.forEach(r => {
    const lvl = r.warning_status ?? 0
    if (lvl >= 0 && lvl <= 2) buckets[lvl]++
  })
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: ['45%', '70%'],
        label: { formatter: '{b}\n{c}' },
        data: [
          { value: buckets[0], name: '正常', itemStyle: { color: '#67c23a' } },
          { value: buckets[1], name: '中等', itemStyle: { color: '#e6a23c' } },
          { value: buckets[2], name: '高压力', itemStyle: { color: '#f56c6c' } },
        ],
      },
    ],
  }
})

const recentStress = computed(() =>
  [...stressList.value]
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    .slice(0, 6)
)

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 16) : ''
}

async function loadAll() {
  loading.value = true
  try {
    const [users, stress, imports, suggs] = await Promise.all([
      getUserList().catch(() => ({ list: [] })),
      getStressResultList({ page: 1, page_size: 100 }).catch(() => ({ list: [], total: 0 })),
      getAcademicImportList({ page: 1, page_size: 1 }).catch(() => ({ list: [], total: 0 })),
      getInterventionSuggestionList({ page: 1, page_size: 1 }).catch(() => ({ list: [], total: 0 })),
    ])
    const userList = users?.list || []
    stats.userCount = userList.length
    stats.studentCount = userList.filter(u => u.role === 'student').length
    stats.teacherCount = userList.filter(u => u.role === 'teacher').length
    stats.adminCount = userList.filter(u => u.role === 'admin').length

    stressList.value = stress?.list || []
    stats.stressTotal = stress?.total || stressList.value.length
    stats.warningCount = stressList.value.filter(r => r.warning_status >= 1).length

    // 本周（近 7 天）高风险
    const weekAgo = new Date()
    weekAgo.setDate(weekAgo.getDate() - 7)
    stats.highRiskWeek = stressList.value.filter(
      r => r.warning_status === 2 && new Date(r.created_at) >= weekAgo
    ).length

    stats.importCount = imports?.total || 0
    stats.suggestionCount = suggs?.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(120deg, #1f2d3d 0%, #3a4f63 100%);
  color: #fff;
  padding: 24px 32px;
}
.hero h2 {
  margin: 0 0 6px 0;
  font-size: 22px;
}
.hero p {
  margin: 0;
  opacity: 0.85;
}
.hero-stats {
  display: flex;
  gap: 32px;
}
.hero-stat {
  text-align: center;
}
.hero-stat .num {
  display: block;
  font-size: 26px;
  font-weight: 600;
}
.hero-stat .lbl {
  display: block;
  font-size: 12px;
  opacity: 0.7;
}

.metric-card {
  background: #fff;
  border-radius: 10px;
  padding: 18px 22px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: relative;
}
.metric-card .label {
  color: #909399;
  font-size: 13px;
}
.metric-card .value {
  font-size: 24px;
  font-weight: 600;
}
.metric-card .sub {
  color: #909399;
  font-size: 12px;
}
.metric-link {
  position: absolute;
  right: 16px;
  bottom: 12px;
  color: #409eff;
  font-size: 12px;
  cursor: pointer;
}
.metric-card.primary {
  border-top: 3px solid #409eff;
}
.metric-card.warn {
  border-top: 3px solid #e6a23c;
}
.metric-card.warn .value {
  color: #e6a23c;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.chart {
  height: 280px;
  width: 100%;
}

.system-card .status-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.status-list li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px dashed #ebeef5;
}
.status-list li:last-child {
  border-bottom: none;
}
.status-list .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #909399;
}
.status-list .dot.ok {
  background: #67c23a;
}
.status-list .dot.fail {
  background: #f56c6c;
}
.status-list .dot.info {
  background: #409eff;
}
.status-list .key {
  flex: 1;
  color: #303133;
}
.status-tip {
  margin: 12px 0 0 0;
  color: #909399;
  font-size: 12px;
}

.quick-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
}
.quick-item {
  background: #f8f9fc;
  border-radius: 8px;
  padding: 18px 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.15s;
}
.quick-item:hover {
  background: #ecf5ff;
  transform: translateY(-2px);
}
.quick-item span {
  display: block;
  margin-top: 6px;
  font-size: 13px;
}
</style>