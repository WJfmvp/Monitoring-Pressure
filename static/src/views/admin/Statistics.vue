<template>
  <div v-loading="loading">
    <div class="page-card">
      <h2 class="page-title">数据统计</h2>
      <p class="hint">基于现有数据的跨表分析视图，前端实时聚合。</p>
    </div>

    <!-- 顶部 4 个指标 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">系统用户数</span>
          <span class="value">{{ stats.userCount || 0 }}</span>
          <span class="sub">学生 {{ stats.studentCount }} · 教师 {{ stats.teacherCount }} · 管理员 {{ stats.adminCount }}</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card warn">
          <span class="label">累计评估次数</span>
          <span class="value">{{ stats.stressCount || 0 }}</span>
          <span class="sub">预警占比 {{ stats.warningRate }}%</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">已采集成绩</span>
          <span class="value">{{ stats.academicCount || 0 }}</span>
          <span class="sub">来自 {{ stats.importCount }} 次导入</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">干预建议模板</span>
          <span class="value">{{ stats.suggestionCount || 0 }}</span>
          <span class="sub">覆盖 {{ stats.suggLevels }} 个等级</span>
        </div>
      </el-col>
    </el-row>

    <!-- 第一行：角色饼图 + 评估时间分布 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">用户角色分布</h3>
          <v-chart v-if="userList.length" class="chart" :option="rolePieOption" autoresize />
          <el-empty v-else description="暂无数据" />
        </div>
      </el-col>
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">近 30 天评估时间分布</h3>
          <v-chart v-if="stressList.length" class="chart" :option="timelineOption" autoresize />
          <el-empty v-else description="暂无评估记录" />
        </div>
      </el-col>
    </el-row>

    <!-- 第二行：等级分布堆叠 + 干预建议分类 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">预警等级分布（按时间）</h3>
          <v-chart v-if="stressList.length" class="chart" :option="stackedOption" autoresize />
          <el-empty v-else description="暂无评估记录" />
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">干预建议库分类</h3>
          <v-chart v-if="suggList.length" class="chart" :option="suggCategoryOption" autoresize />
          <el-empty v-else description="暂无干预建议" />
        </div>
      </el-col>
    </el-row>

    <!-- 第三行：高压力 Top10 -->
    <div class="page-card" style="margin-top: 16px">
      <h3 class="card-title">高压力 Top 10 学生</h3>
      <el-empty v-if="topStudents.length === 0" description="暂无高压力学生" />
      <el-table v-else :data="topStudents" stripe size="small">
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="user_id" label="学生 ID" width="120" />
        <el-table-column label="综合得分" width="120">
          <template #default="{ row }">
            <span :style="{ color: scoreColor(row.total_score) }">
              {{ row.total_score?.toFixed?.(1) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="预警等级" width="120">
          <template #default="{ row }">
            <el-tag :type="LEVEL_TAG[row.warning_status]">
              {{ LEVEL_LABEL[row.warning_status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="behavior_score" label="行为" width="100" />
        <el-table-column prop="academic_score" label="学业" width="100" />
        <el-table-column prop="psychological_score" label="心理" width="100" />
        <el-table-column prop="created_at" label="评估时间">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import {
  PieChart as EChartsPie,
  LineChart as EChartsLine,
  BarChart as EChartsBar,
} from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'
import {
  getUserList,
  getStressResultList,
  getAcademicList,
  getAcademicImportList,
  getInterventionSuggestionList,
} from '@/api/admin'

use([
  CanvasRenderer,
  EChartsPie,
  EChartsLine,
  EChartsBar,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent,
])

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']

const loading = ref(false)
const userList = ref([])
const stressList = ref([])
const academicList = ref([])
const suggList = ref([])

const stats = reactive({
  userCount: 0,
  studentCount: 0,
  teacherCount: 0,
  adminCount: 0,
  stressCount: 0,
  warningRate: 0,
  academicCount: 0,
  importCount: 0,
  suggestionCount: 0,
  suggLevels: 0,
})

const rolePieOption = computed(() => {
  const counts = { student: 0, teacher: 0, admin: 0 }
  userList.value.forEach(u => {
    counts[u.role] = (counts[u.role] || 0) + 1
  })
  return {
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: counts.student, name: '学生', itemStyle: { color: '#6e8efb' } },
          { value: counts.teacher, name: '教师', itemStyle: { color: '#67c23a' } },
          { value: counts.admin, name: '管理员', itemStyle: { color: '#f56c6c' } },
        ],
        label: { formatter: '{b}\n{c} 人' },
      },
    ],
  }
})

const timelineOption = computed(() => {
  // 近 30 天每天的评估数量
  const buckets = {}
  const now = new Date()
  for (let i = 29; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    const key = d.toISOString().slice(0, 10)
    buckets[key] = 0
  }
  stressList.value.forEach(r => {
    const date = (r.created_at || '').slice(0, 10)
    if (buckets[date] !== undefined) buckets[date]++
  })
  const keys = Object.keys(buckets)
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 30, bottom: 40 },
    xAxis: { type: 'category', data: keys.map(k => k.slice(5)) },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'bar',
        barWidth: '60%',
        itemStyle: { color: '#409eff', borderRadius: [3, 3, 0, 0] },
        data: keys.map(k => buckets[k]),
      },
    ],
  }
})

const stackedOption = computed(() => {
  // 按周分组，每周三种等级数量
  const weekly = {}
  stressList.value.forEach(r => {
    const d = new Date(r.created_at)
    if (Number.isNaN(d.valueOf())) return
    // 取该日期所在周的周一
    const monday = new Date(d)
    monday.setDate(d.getDate() - ((d.getDay() + 6) % 7))
    const key = monday.toISOString().slice(0, 10)
    if (!weekly[key]) weekly[key] = [0, 0, 0]
    const lvl = r.warning_status ?? 0
    if (lvl >= 0 && lvl <= 2) weekly[key][lvl]++
  })
  const keys = Object.keys(weekly).sort()
  return {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { bottom: 0, data: ['正常', '中等', '高压力'] },
    grid: { left: 40, right: 20, top: 30, bottom: 50 },
    xAxis: { type: 'category', data: keys.map(k => k.slice(5)) },
    yAxis: { type: 'value' },
    series: [
      {
        name: '正常',
        type: 'bar',
        stack: 'total',
        itemStyle: { color: '#67c23a' },
        data: keys.map(k => weekly[k][0]),
      },
      {
        name: '中等',
        type: 'bar',
        stack: 'total',
        itemStyle: { color: '#e6a23c' },
        data: keys.map(k => weekly[k][1]),
      },
      {
        name: '高压力',
        type: 'bar',
        stack: 'total',
        itemStyle: { color: '#f56c6c' },
        data: keys.map(k => weekly[k][2]),
      },
    ],
  }
})

const suggCategoryOption = computed(() => {
  const cats = {}
  suggList.value.forEach(s => {
    const c = s.category || '未分类'
    cats[c] = (cats[c] || 0) + 1
  })
  return {
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: '60%',
        data: Object.entries(cats).map(([name, value]) => ({ name, value })),
        label: { formatter: '{b}\n{c} 条' },
      },
    ],
  }
})

const topStudents = computed(() => {
  return [...stressList.value]
    .filter(r => r.warning_status >= 1)
    .sort((a, b) => (b.total_score || 0) - (a.total_score || 0))
    .slice(0, 10)
})

function scoreColor(s) {
  if (s == null) return ''
  if (s >= 75) return '#f56c6c'
  if (s >= 55) return '#e6a23c'
  return '#67c23a'
}

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 19) : ''
}

async function loadAll() {
  loading.value = true
  try {
    const [users, stress, academic, imports, suggs] = await Promise.all([
      getUserList().catch(() => ({ list: [] })),
      getStressResultList({ page: 1, page_size: 500 }).catch(() => ({ list: [], total: 0 })),
      getAcademicList({ page: 1, page_size: 1 }).catch(() => ({ list: [], total: 0 })),
      getAcademicImportList({ page: 1, page_size: 1 }).catch(() => ({ list: [], total: 0 })),
      getInterventionSuggestionList({ page: 1, page_size: 200 }).catch(() => ({ list: [], total: 0 })),
    ])

    userList.value = users?.list || []
    stressList.value = stress?.list || []
    suggList.value = suggs?.list || []

    stats.userCount = userList.value.length
    stats.studentCount = userList.value.filter(u => u.role === 'student').length
    stats.teacherCount = userList.value.filter(u => u.role === 'teacher').length
    stats.adminCount = userList.value.filter(u => u.role === 'admin').length

    stats.stressCount = stress?.total || stressList.value.length
    const warned = stressList.value.filter(r => r.warning_status >= 1).length
    stats.warningRate = stress?.total
      ? Math.round((warned / stressList.value.length) * 1000) / 10
      : 0

    stats.academicCount = academic?.total || 0
    stats.importCount = imports?.total || 0

    stats.suggestionCount = suggs?.total || suggList.value.length
    stats.suggLevels = new Set(suggList.value.map(s => s.level)).size
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 6px 0 0 0;
}
.metric-card {
  background: #fff;
  border-radius: 10px;
  padding: 18px 22px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.metric-card .label {
  color: #909399;
  font-size: 13px;
}
.metric-card .value {
  font-size: 26px;
  font-weight: 600;
}
.metric-card .sub {
  color: #909399;
  font-size: 12px;
}
.metric-card.warn {
  border-left: 4px solid #e6a23c;
}
.metric-card.warn .value {
  color: #e6a23c;
}
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
.chart {
  height: 300px;
  width: 100%;
}
</style>