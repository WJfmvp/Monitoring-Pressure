<template>
  <div v-loading="loading">
    <!-- 顶栏：返回 + 标题 -->
    <div class="page-card top-bar">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回班级总览
      </el-button>
      <h2 class="page-title" style="margin: 0">压力诊断书</h2>
      <div></div>
    </div>

    <template v-if="data">
      <!-- 学生信息 + 当前等级 -->
      <div class="page-card hero-card" :class="heroClass">
        <div class="hero-left">
          <div class="hero-name">
            <span class="name">{{ data.student?.username }}</span>
            <el-tag v-if="data.student?.sex === 1" size="small">男</el-tag>
            <el-tag v-else-if="data.student?.sex === 2" size="small" type="warning">女</el-tag>
            <span class="sub">学号 {{ data.student?.user_id }}</span>
          </div>
          <div class="meta">
            手机 {{ data.student?.telephone || '-' }} ·
            邮箱 {{ data.student?.email || '-' }} ·
            最近评估 {{ data.latest?.assessed_at || '—' }}
          </div>
        </div>

        <div class="hero-right" v-if="latest">
          <div class="big-tag">
            <span class="tag-label">当前预警</span>
            <span class="tag-level" :style="{ color: LEVEL_COLOR[latest.warning_status] }">
              {{ LEVEL_LABEL[latest.warning_status] }}
            </span>
          </div>
          <div class="big-tag">
            <span class="tag-label">综合得分</span>
            <span class="tag-level" :style="{ color: LEVEL_COLOR[latest.warning_status] }">
              {{ latest.total_score?.toFixed?.(1) }}
            </span>
          </div>
          <div class="big-tag">
            <span class="tag-label">模型置信度</span>
            <span class="tag-level">
              {{ ((latest.predict_probability || 0) * 100).toFixed(1) }}%
            </span>
          </div>
        </div>
        <el-empty v-else description="该学生暂无评估数据" />
      </div>

      <el-row :gutter="16" style="margin-top: 16px" v-if="latest">
        <el-col :span="10">
          <div class="page-card">
            <h3 class="card-title">压力维度雷达</h3>
            <v-chart class="chart" :option="radarOption" autoresize />
          </div>
        </el-col>
        <el-col :span="14">
          <div class="page-card risk-card">
            <h3 class="card-title">
              <el-icon><Warning /></el-icon>
              风险预警说明
            </h3>
            <ul class="risk-list">
              <li v-for="(t, i) in data.explanations" :key="i" :class="riskItemClass(t)">
                {{ t }}
              </li>
            </ul>
          </div>
        </el-col>
      </el-row>

      <!-- 趋势图 -->
      <div class="page-card" style="margin-top: 16px" v-if="data.history?.length">
        <h3 class="card-title">压力趋势（近 {{ data.history.length }} 次评估）</h3>
        <v-chart class="chart trend-chart" :option="trendOption" autoresize />
      </div>

      <!-- 干预建议 -->
      <div class="page-card" style="margin-top: 16px" v-if="latest">
        <h3 class="card-title">建议干预方向</h3>
        <el-empty
          v-if="!data.suggestions?.length"
          description="暂无对应等级的干预建议模板，可前往「干预建议管理」配置"
        />
        <el-row :gutter="12" v-else>
          <el-col :span="12" v-for="s in data.suggestions" :key="s.id">
            <div class="sugg-card" :class="`level-${s.level}`">
              <div class="sugg-header">
                <el-tag :type="LEVEL_TAG[s.level]" size="small">
                  {{ LEVEL_LABEL[s.level] }}
                </el-tag>
                <el-tag v-if="s.category" size="small" effect="plain">{{ s.category }}</el-tag>
                <span class="sugg-title">{{ s.title }}</span>
              </div>
              <div class="sugg-content">{{ s.content }}</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import { RadarChart as EChartsRadar, LineChart as EChartsLine } from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
  RadarComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'
import { getStudentDiagnosis } from '@/api/teacher'

use([
  CanvasRenderer,
  EChartsRadar,
  EChartsLine,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  RadarComponent,
])

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const data = ref(null)

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']
const LEVEL_COLOR = ['#67c23a', '#e6a23c', '#f56c6c']

const latest = computed(() => data.value?.latest)

const heroClass = computed(() => {
  if (!latest.value) return ''
  if (latest.value.warning_status === 2) return 'hero-danger'
  if (latest.value.warning_status === 1) return 'hero-warn'
  return 'hero-ok'
})

const radarOption = computed(() => ({
  tooltip: {},
  radar: {
    indicator: [
      { name: '行为压力', max: 100 },
      { name: '学业压力', max: 100 },
      { name: '心理压力', max: 100 },
      { name: '综合风险', max: 100 },
    ],
    radius: '65%',
    splitNumber: 4,
    axisName: { color: '#606266' },
  },
  series: [
    {
      type: 'radar',
      symbolSize: 8,
      areaStyle: { opacity: 0.25 },
      lineStyle: { width: 2 },
      itemStyle: { color: LEVEL_COLOR[latest.value?.warning_status || 0] },
      data: [
        {
          name: '当前评估',
          value: [
            latest.value?.behavior_score || 0,
            latest.value?.academic_score || 0,
            latest.value?.psychological_score || 0,
            latest.value?.total_score || 0,
          ],
        },
      ],
    },
  ],
}))

const trendOption = computed(() => {
  const history = data.value?.history || []
  const labels = history.map(h => h.assessed_at?.slice(5, 10))
  return {
    tooltip: { trigger: 'axis' },
    legend: { bottom: 0, data: ['综合', '行为', '学业', '心理'] },
    grid: { left: 40, right: 20, top: 30, bottom: 50 },
    xAxis: { type: 'category', data: labels },
    yAxis: { type: 'value', max: 100 },
    series: [
      {
        name: '综合',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        lineStyle: { width: 3 },
        itemStyle: { color: '#409eff' },
        data: history.map(h => h.total_score),
      },
      {
        name: '行为',
        type: 'line',
        smooth: true,
        itemStyle: { color: '#6e8efb' },
        data: history.map(h => h.behavior_score),
      },
      {
        name: '学业',
        type: 'line',
        smooth: true,
        itemStyle: { color: '#67c23a' },
        data: history.map(h => h.academic_score),
      },
      {
        name: '心理',
        type: 'line',
        smooth: true,
        itemStyle: { color: '#a777e3' },
        data: history.map(h => h.psychological_score),
      },
    ],
  }
})

function riskItemClass(t) {
  if (t.startsWith('⚠')) return 'risk-warn'
  if (t.startsWith('✓')) return 'risk-ok'
  return ''
}

function goBack() {
  if (window.history.length > 1) router.back()
  else router.push('/teacher/class/overview')
}

async function load() {
  const sid = route.query.student_id
  if (!sid) {
    ElMessage.warning('未指定 student_id')
    return
  }
  loading.value = true
  try {
    data.value = await getStudentDiagnosis(sid)
  } finally {
    loading.value = false
  }
}

watch(() => route.query.student_id, () => load())

onMounted(load)
</script>

<style scoped>
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
}

.hero-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-left: 4px solid #909399;
  flex-wrap: wrap;
  gap: 16px;
}
.hero-ok {
  border-left-color: #67c23a;
}
.hero-warn {
  border-left-color: #e6a23c;
  background: linear-gradient(90deg, rgba(230, 162, 60, 0.05), #fff);
}
.hero-danger {
  border-left-color: #f56c6c;
  background: linear-gradient(90deg, rgba(245, 108, 108, 0.08), #fff);
}
.hero-left .hero-name {
  display: flex;
  align-items: center;
  gap: 10px;
}
.hero-left .name {
  font-size: 22px;
  font-weight: 600;
}
.hero-left .sub {
  color: #909399;
  font-size: 13px;
}
.hero-left .meta {
  margin-top: 8px;
  color: #606266;
  font-size: 13px;
}
.hero-right {
  display: flex;
  gap: 28px;
}
.big-tag {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}
.big-tag .tag-label {
  color: #909399;
  font-size: 12px;
}
.big-tag .tag-level {
  font-size: 22px;
  font-weight: 600;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.chart {
  height: 320px;
  width: 100%;
}
.trend-chart {
  height: 280px;
}

.risk-card .risk-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.risk-list li {
  padding: 8px 12px;
  border-radius: 6px;
  margin-bottom: 8px;
  background: #f5f7fa;
  line-height: 1.6;
}
.risk-list li.risk-warn {
  background: rgba(245, 108, 108, 0.08);
  color: #f56c6c;
  font-weight: 500;
}
.risk-list li.risk-ok {
  background: rgba(103, 194, 58, 0.08);
  color: #67c23a;
}

.sugg-card {
  background: #f8f9fc;
  border-radius: 8px;
  padding: 14px 16px;
  margin-bottom: 12px;
  border-left: 3px solid #909399;
}
.sugg-card.level-0 {
  border-left-color: #67c23a;
}
.sugg-card.level-1 {
  border-left-color: #e6a23c;
}
.sugg-card.level-2 {
  border-left-color: #f56c6c;
}
.sugg-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.sugg-title {
  font-weight: 600;
  margin-left: 4px;
}
.sugg-content {
  color: #606266;
  font-size: 13px;
  line-height: 1.6;
}
</style>
