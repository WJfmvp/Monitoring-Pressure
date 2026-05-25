<template>
  <div v-loading="loading">
    <h2 class="page-title" style="margin-bottom: 16px">班级整体压力分布</h2>

    <!-- 顶部指标卡片 -->
    <el-row :gutter="16">
      <el-col :span="5">
        <div class="metric-card">
          <span class="label">班级学生总数</span>
          <span class="value">{{ summary.total_students || 0 }}</span>
        </div>
      </el-col>
      <el-col :span="5">
        <div class="metric-card">
          <span class="label">已评估</span>
          <span class="value">
            {{ summary.assessed_count || 0 }}
            <span class="sub">/ {{ summary.total_students || 0 }}</span>
          </span>
        </div>
      </el-col>
      <el-col :span="5">
        <div class="metric-card warn">
          <span class="label">预警人数（中+高）</span>
          <span class="value">{{ summary.warned_count || 0 }}</span>
        </div>
      </el-col>
      <el-col :span="5">
        <div class="metric-card danger">
          <span class="label">高压力预警</span>
          <span class="value">{{ summary.high_risk_count || 0 }}</span>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="metric-card">
          <span class="label">班级综合均分</span>
          <span class="value">{{ summary.avg_total || '-' }}</span>
        </div>
      </el-col>
    </el-row>

    <!-- 中部图表区 -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="10">
        <div class="page-card">
          <h3 class="card-title">压力等级分布</h3>
          <v-chart v-if="hasItems" class="chart" :option="pieOption" autoresize />
          <el-empty v-else description="暂无评估数据" />
        </div>
      </el-col>
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">班级维度平均分</h3>
          <v-chart v-if="hasItems" class="chart" :option="barOption" autoresize />
          <el-empty v-else description="暂无评估数据" />
        </div>
      </el-col>
    </el-row>

    <!-- 学生压力排行 -->
    <div class="page-card" style="margin-top: 16px">
      <div class="card-header">
        <h3 class="card-title" style="margin: 0">学生压力排行</h3>
        <div>
          <el-radio-group v-model="filterLevel" size="small">
            <el-radio-button :value="''">全部</el-radio-button>
            <el-radio-button :value="2">高压力</el-radio-button>
            <el-radio-button :value="1">中等</el-radio-button>
            <el-radio-button :value="0">正常</el-radio-button>
            <el-radio-button :value="-1">未评估</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <el-table
        :data="filteredItems"
        stripe
        :row-class-name="rowClass"
        @row-click="onRowClick"
      >
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="username" label="学生" min-width="120">
          <template #default="{ row }">
            <span>{{ row.username }}</span>
            <el-tag v-if="row.sex === 1" size="small" effect="plain" style="margin-left: 6px">
              男
            </el-tag>
            <el-tag
              v-else-if="row.sex === 2"
              size="small"
              effect="plain"
              type="warning"
              style="margin-left: 6px"
            >
              女
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_id" label="学号" width="120" />
        <el-table-column label="预警等级" width="120">
          <template #default="{ row }">
            <el-tag v-if="!row.has_result" type="info">未评估</el-tag>
            <el-tag v-else :type="LEVEL_TAG[row.warning_status]">
              {{ LEVEL_LABEL[row.warning_status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="综合" width="100">
          <template #default="{ row }">
            <span v-if="row.has_result" :style="{ color: scoreColor(row.total_score) }">
              {{ row.total_score?.toFixed?.(1) }}
            </span>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="行为" width="80">
          <template #default="{ row }">
            <span v-if="row.has_result">{{ row.behavior_score?.toFixed?.(1) }}</span>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="学业" width="80">
          <template #default="{ row }">
            <span v-if="row.has_result">{{ row.academic_score?.toFixed?.(1) }}</span>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="心理" width="80">
          <template #default="{ row }">
            <span v-if="row.has_result">{{ row.psychological_score?.toFixed?.(1) }}</span>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="assessed_at" label="最近评估" min-width="160">
          <template #default="{ row }">
            {{ row.assessed_at || '—' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.has_result"
              :type="row.warning_status >= 1 ? 'danger' : 'primary'"
              size="small"
              link
              @click.stop="openDiagnosis(row)"
            >
              <el-icon style="margin-right: 4px"><DocumentCopy /></el-icon>
              压力诊断书
            </el-button>
            <span v-else class="muted">—</span>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import { PieChart as EChartsPie, BarChart as EChartsBar } from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'
import VChart from 'vue-echarts'
import { getClassOverview } from '@/api/teacher'

use([CanvasRenderer, EChartsPie, EChartsBar, TooltipComponent, LegendComponent, GridComponent])

const router = useRouter()
const loading = ref(false)
const summary = ref({})
const items = ref([])
const filterLevel = ref('')

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']
const LEVEL_COLOR = { 0: '#67c23a', 1: '#e6a23c', 2: '#f56c6c', '-1': '#909399' }

const hasItems = computed(() => (summary.value?.assessed_count || 0) > 0)

const filteredItems = computed(() => {
  if (filterLevel.value === '' || filterLevel.value === undefined) {
    return items.value
  }
  if (filterLevel.value === -1) {
    return items.value.filter(it => !it.has_result)
  }
  return items.value.filter(it => it.has_result && it.warning_status === filterLevel.value)
})

const pieOption = computed(() => {
  const dist = summary.value.level_distribution || []
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: ['45%', '70%'],
        label: { formatter: '{b}\n{c} 人' },
        data: dist.map(d => ({
          name: d.label,
          value: d.count,
          itemStyle: { color: LEVEL_COLOR[d.level] },
        })),
      },
    ],
  }
})

const barOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 50, right: 30, top: 30, bottom: 30 },
  xAxis: { type: 'category', data: ['行为', '学业', '心理', '综合'] },
  yAxis: { type: 'value', max: 100 },
  series: [
    {
      type: 'bar',
      barWidth: 40,
      label: { show: true, position: 'top' },
      itemStyle: {
        color: ({ dataIndex }) =>
          ['#6e8efb', '#67c23a', '#a777e3', '#f56c6c'][dataIndex],
        borderRadius: [6, 6, 0, 0],
      },
      data: [
        summary.value.avg_behavior || 0,
        summary.value.avg_academic || 0,
        summary.value.avg_psychological || 0,
        summary.value.avg_total || 0,
      ],
    },
  ],
}))

function rowClass({ row }) {
  if (!row.has_result) return 'row-muted'
  if (row.warning_status === 2) return 'row-danger'
  if (row.warning_status === 1) return 'row-warn'
  return ''
}

function scoreColor(s) {
  if (s == null) return ''
  if (s >= 75) return '#f56c6c'
  if (s >= 55) return '#e6a23c'
  return '#67c23a'
}

function openDiagnosis(row) {
  if (!row.has_result) return
  router.push({
    path: '/teacher/student/diagnosis',
    query: { student_id: row.user_id },
  })
}

function onRowClick(row) {
  if (row.has_result) openDiagnosis(row)
}

async function load() {
  loading.value = true
  try {
    const data = await getClassOverview()
    summary.value = data?.summary || {}
    items.value = data?.items || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.metric-card {
  background: #fff;
  border-radius: 10px;
  padding: 18px 22px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
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
  font-size: 14px;
  color: #909399;
  font-weight: 400;
}
.metric-card.warn {
  border-left: 4px solid #e6a23c;
}
.metric-card.warn .value {
  color: #e6a23c;
}
.metric-card.danger {
  border-left: 4px solid #f56c6c;
}
.metric-card.danger .value {
  color: #f56c6c;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.chart {
  height: 300px;
  width: 100%;
}
.muted {
  color: #c0c4cc;
}

:deep(.row-danger) {
  background-color: rgba(245, 108, 108, 0.08) !important;
}
:deep(.row-warn) {
  background-color: rgba(230, 162, 60, 0.06) !important;
}
:deep(.row-muted td) {
  color: #909399;
}
:deep(.el-table__row) {
  cursor: pointer;
}
</style>
