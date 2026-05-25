<template>
  <div>
    <!-- 欢迎横幅 -->
    <div class="page-card hero">
      <div class="hero-text">
        <h2>{{ greeting }}，{{ userStore.userInfo?.username || '同学' }} 👋</h2>
        <p>{{ todayHint }}</p>
      </div>
      <div class="hero-date">
        <div class="date-day">{{ todayDay }}</div>
        <div class="date-meta">{{ todayMeta }}</div>
      </div>
    </div>

    <!-- 4 个指标 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="6">
        <div class="metric-card" :style="{ borderLeftColor: levelColor }">
          <span class="label">最新压力等级</span>
          <span class="value" :style="{ color: levelColor }">{{ levelLabel }}</span>
          <span class="sub">{{ latest?.created_at?.slice(0, 10) || '—' }}</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">综合得分</span>
          <span class="value">{{ latest?.total_score?.toFixed?.(1) ?? '-' }}</span>
          <span class="sub">{{ classCompareText }}</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">心理得分</span>
          <span class="value">{{ latest?.psychological_score?.toFixed?.(1) ?? '-' }}</span>
          <span class="sub">{{ psychTrendText }}</span>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card">
          <span class="label">学业得分</span>
          <span class="value">{{ latest?.academic_score?.toFixed?.(1) ?? '-' }}</span>
          <span class="sub">{{ academicTrendText }}</span>
        </div>
      </el-col>
    </el-row>

    <!-- 任务区 + 趋势 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="8">
        <div class="page-card task-card">
          <h3 class="card-title">
            <el-icon><Bell /></el-icon>
            待办事项
          </h3>
          <ul class="task-list">
            <li v-for="t in tasks" :key="t.key" :class="t.done ? 'done' : ''">
              <el-icon :color="t.done ? '#67c23a' : '#e6a23c'">
                <component :is="t.done ? 'CircleCheck' : 'Warning'" />
              </el-icon>
              <div class="task-body">
                <h4>{{ t.title }}</h4>
                <p>{{ t.desc }}</p>
              </div>
              <el-button v-if="!t.done" link type="primary" @click="$router.push(t.path)">
                前往
              </el-button>
            </li>
          </ul>
        </div>
      </el-col>
      <el-col :span="16">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><DataLine /></el-icon>
            心理测评分数变化
          </h3>
          <v-chart v-if="psychList.length" class="chart" :option="lineOption" autoresize />
          <el-empty v-else description="还没有历史问卷，先去填一份吧" />
        </div>
      </el-col>
    </el-row>

    <!-- 维度对比 + 干预记录 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><PieChart /></el-icon>
            维度得分对比
          </h3>
          <v-chart v-if="latest" class="chart" :option="barOption" autoresize />
          <el-empty v-else description="尚未生成压力评估" />
        </div>
      </el-col>
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><ChatLineSquare /></el-icon>
            最近干预建议
          </h3>
          <el-empty v-if="interventions.length === 0" description="暂无建议推送" />
          <el-table v-else :data="interventions.slice(0, 4)" size="small">
            <el-table-column label="建议">
              <template #default="{ row }">
                <div class="intervention-title">{{ suggestionOf(row).title || `建议 #${row.suggestion_id}` }}</div>
                <div class="intervention-time">{{ formatTime(row.push_time) }}</div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag size="small" :type="FB_TAG[row.feedback_status] || 'info'">
                  {{ FB_LABEL[row.feedback_status] || '-' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-button
            link
            type="primary"
            @click="$router.push('/student/intervention')"
            style="margin-top: 8px"
          >
            查看全部 →
          </el-button>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import { BarChart as EChartsBar, LineChart as EChartsLine } from 'echarts/charts'
import {
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'
import {
  getLatestStressResult,
  getMyInterventions,
  getPsychologicalList,
} from '@/api/student'
import { useUserStore } from '@/stores/user'

use([CanvasRenderer, EChartsBar, EChartsLine, TooltipComponent, GridComponent, LegendComponent])

const userStore = useUserStore()

const latest = ref(null)
const interventions = ref([])
const psychList = ref([])

const LEVEL_LABELS = ['正常', '中等', '高压力']
const LEVEL_COLORS = ['#67c23a', '#e6a23c', '#f56c6c']
const FB_LABEL = ['未反馈', '有效', '无效']
const FB_TAG = ['info', 'success', 'danger']

const levelLabel = computed(() => {
  const lvl = latest.value?.warning_status
  return lvl == null ? '—' : LEVEL_LABELS[lvl] || '-'
})
const levelColor = computed(() => {
  const lvl = latest.value?.warning_status
  return lvl == null ? '#909399' : LEVEL_COLORS[lvl] || '#303133'
})

// 时间 / 问候
const now = new Date()
const greeting = computed(() => {
  const h = now.getHours()
  if (h < 6) return '夜深了'
  if (h < 12) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})
const todayDay = computed(() => now.getDate())
const todayMeta = computed(() => {
  const m = now.getMonth() + 1
  const weeks = ['日', '一', '二', '三', '四', '五', '六']
  return `${m} 月 · 周${weeks[now.getDay()]}`
})
const todayHint = computed(() => {
  if (!latest.value) return '快去完成一次心理测评，了解自己当下的状态。'
  if (latest.value.warning_status === 2) return '系统监测到你近期压力较大，记得多关注自己的情绪。'
  if (latest.value.warning_status === 1) return '注意劳逸结合，压力管理是长期的事情。'
  return '当前状态良好，继续保持节奏 ✨'
})

// 任务区
const tasks = computed(() => {
  const items = []

  // 1. 是否需要心理测评
  const recent = psychList.value[0]
  const daysSince = recent ? daysBetween(recent.assess_date, now) : 999
  items.push({
    key: 'psy',
    done: daysSince <= 7,
    title: daysSince <= 7 ? '本周已完成心理测评' : '提交一次心理测评',
    desc: recent ? `最近一次：${formatDate(recent.assess_date)}` : '尚未有任何记录',
    path: '/student/psychological',
  })

  // 2. 是否完善了资料
  items.push({
    key: 'profile',
    done: !!(userStore.userInfo?.email && userStore.userInfo?.sex),
    title: '完善个人资料',
    desc: '完善邮箱和性别后，老师能更好地了解你',
    path: '/profile/complete',
  })

  // 3. 试试减压工具
  items.push({
    key: 'wellness',
    done: false,
    title: '试试放松小工具',
    desc: '呼吸训练、番茄钟，让自己慢下来',
    path: '/wellness',
  })

  return items
})

function daysBetween(date, ref) {
  if (!date) return 999
  const a = new Date(date)
  if (Number.isNaN(a.valueOf())) return 999
  return Math.floor((ref - a) / (1000 * 60 * 60 * 24))
}

// 图表
const barOption = computed(() => ({
  tooltip: {},
  grid: { left: 40, right: 20, top: 30, bottom: 30 },
  xAxis: { type: 'category', data: ['行为', '学业', '心理', '综合'] },
  yAxis: { type: 'value', max: 100 },
  series: [
    {
      type: 'bar',
      barWidth: 36,
      itemStyle: {
        color: ({ dataIndex }) => ['#6e8efb', '#67c23a', '#a777e3', '#f56c6c'][dataIndex],
        borderRadius: [6, 6, 0, 0],
      },
      data: latest.value
        ? [
            latest.value.behavior_score,
            latest.value.academic_score,
            latest.value.psychological_score,
            latest.value.total_score,
          ]
        : [],
    },
  ],
}))

const lineOption = computed(() => {
  const sorted = [...psychList.value].sort(
    (a, b) => new Date(a.assess_date) - new Date(b.assess_date)
  )
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 30, bottom: 30 },
    xAxis: {
      type: 'category',
      data: sorted.map(it => formatDate(it.assess_date)),
    },
    yAxis: { type: 'value', max: 100 },
    series: [
      {
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: { width: 3, color: '#a777e3' },
        itemStyle: { color: '#a777e3' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(167,119,227,0.3)' },
              { offset: 1, color: 'rgba(167,119,227,0.02)' },
            ],
          },
        },
        data: sorted.map(it => it.questionnaire_score || 0),
      },
    ],
  }
})

const classCompareText = computed(() => {
  if (!latest.value) return '—'
  return latest.value.total_score >= 55 ? '高于轻度阈值' : '处于健康区间'
})

const psychTrendText = computed(() => {
  if (psychList.value.length < 2) return '需要更多数据'
  const sorted = [...psychList.value].sort(
    (a, b) => new Date(a.assess_date) - new Date(b.assess_date)
  )
  const delta = sorted[sorted.length - 1].questionnaire_score - sorted[sorted.length - 2].questionnaire_score
  if (Math.abs(delta) < 1) return '与上次持平'
  return delta > 0 ? `较上次 +${delta.toFixed(1)}` : `较上次 ${delta.toFixed(1)}`
})

const academicTrendText = computed(() => {
  if (!latest.value) return '—'
  if (latest.value.academic_score >= 75) return '需关注'
  if (latest.value.academic_score >= 55) return '中等'
  return '良好'
})

function suggestionOf(row) {
  return row?.suggestion || row?.Suggestion || {
    title: row?.suggestion_title || row?.title || '',
    content: row?.suggestion_content || row?.content || '',
    category: row?.suggestion_category || row?.category || '',
    level: row?.suggestion_level ?? row?.level,
  }
}

function formatDate(t) {
  return t ? String(t).slice(0, 10) : ''
}
function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 16) : ''
}

onMounted(async () => {
  try {
    latest.value = await getLatestStressResult()
  } catch (_) {}
  try {
    const data = await getMyInterventions()
    interventions.value = Array.isArray(data) ? data : []
  } catch (_) {}
  try {
    const data = await getPsychologicalList()
    psychList.value = Array.isArray(data) ? data : []
  } catch (_) {}
})
</script>

<style scoped>
.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(120deg, #6e8efb 0%, #a777e3 100%);
  color: #fff;
  padding: 22px 28px;
}
.hero-text h2 {
  margin: 0 0 6px 0;
  font-size: 22px;
}
.hero-text p {
  margin: 0;
  opacity: 0.9;
}
.hero-date {
  text-align: center;
  color: #fff;
}
.hero-date .date-day {
  font-size: 42px;
  font-weight: 600;
  line-height: 1.1;
}
.hero-date .date-meta {
  font-size: 13px;
  opacity: 0.85;
}

.metric-card {
  background: #fff;
  border-radius: 10px;
  padding: 18px 22px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 4px;
  border-left: 4px solid #dcdfe6;
}
.metric-card .label {
  color: #909399;
  font-size: 13px;
}
.metric-card .value {
  font-size: 28px;
  font-weight: 600;
}
.metric-card .sub {
  color: #909399;
  font-size: 12px;
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

.task-card {
  min-height: 320px;
}
.task-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.task-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}
.task-list li:last-child {
  border-bottom: none;
}
.task-list li.done .task-body h4 {
  text-decoration: line-through;
  color: #909399;
}
.task-body {
  flex: 1;
}
.task-body h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
}
.task-body p {
  margin: 0;
  color: #909399;
  font-size: 12px;
}
.intervention-title {
  color: #303133;
  font-size: 13px;
  font-weight: 600;
}
.intervention-time {
  color: #909399;
  font-size: 12px;
  margin-top: 3px;
}
</style>
