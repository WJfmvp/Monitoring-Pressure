<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">学生压力评估查询</h2>

      <div class="page-toolbar">
        <el-input
          v-model="filters.student_id"
          placeholder="学生 ID（必填）"
          clearable
          style="width: 200px"
        />
        <el-date-picker
          v-model="filters.start_date"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="起始日期"
        />
        <el-date-picker
          v-model="filters.end_date"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="截止日期"
        />
        <el-select
          v-model="filters.warning_status"
          placeholder="预警状态"
          clearable
          style="width: 140px"
        >
          <el-option label="正常" :value="0" />
          <el-option label="中等" :value="1" />
          <el-option label="高压力" :value="2" />
        </el-select>
        <el-button type="primary" @click="reload">查询</el-button>
      </div>

      <el-empty v-if="!filters.student_id" description="请先输入学生 ID" />
      <template v-else>
        <el-table :data="rows" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="created_at" label="评估时间" width="180">
            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column prop="behavior_score" label="行为" width="80" />
          <el-table-column prop="academic_score" label="学业" width="80" />
          <el-table-column prop="psychological_score" label="心理" width="80" />
          <el-table-column prop="total_score" label="综合" width="80" />
          <el-table-column prop="model_predict_level" label="预测等级" width="110">
            <template #default="{ row }">
              <el-tag :type="LEVEL_TAG[row.model_predict_level]">
                {{ LEVEL_LABEL[row.model_predict_level] || '-' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="warning_status" label="预警等级" width="110">
            <template #default="{ row }">
              <el-tag :type="LEVEL_TAG[row.warning_status]">
                {{ LEVEL_LABEL[row.warning_status] || '-' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>

        <div class="page-toolbar" style="justify-content: flex-end; margin-top: 16px">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            background
            @current-change="reload"
            @size-change="reload"
          />
        </div>
      </template>
    </div>

    <div v-if="rows.length" class="page-card" style="margin-top: 16px">
      <h3 class="card-title">综合得分趋势</h3>
      <v-chart class="chart" :option="lineOption" autoresize />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import { LineChart as EChartsLine } from 'echarts/charts'
import { TooltipComponent, GridComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { getStudentStressList } from '@/api/teacher'

use([CanvasRenderer, EChartsLine, TooltipComponent, GridComponent])

const route = useRoute()

const LEVEL_LABEL = ['正常', '中等', '高压力']
const LEVEL_TAG = ['success', 'warning', 'danger']

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({ student_id: '', start_date: '', end_date: '', warning_status: '' })

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value, student_id: filters.value.student_id }
  if (filters.value.start_date) p.start_date = filters.value.start_date
  if (filters.value.end_date) p.end_date = filters.value.end_date
  if (filters.value.warning_status !== '' && filters.value.warning_status != null)
    p.warning_status = filters.value.warning_status
  return p
}

async function reload() {
  if (!filters.value.student_id) {
    ElMessage.warning('请先输入学生 ID')
    return
  }
  loading.value = true
  try {
    const data = await getStudentStressList(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

const lineOption = computed(() => {
  const sorted = [...rows.value].sort(
    (a, b) => new Date(a.created_at) - new Date(b.created_at)
  )
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 30, bottom: 40 },
    xAxis: { type: 'category', data: sorted.map(r => formatTime(r.created_at)) },
    yAxis: { type: 'value', max: 100 },
    series: [
      {
        type: 'line',
        smooth: true,
        symbol: 'circle',
        itemStyle: { color: '#409eff' },
        data: sorted.map(r => r.total_score),
      },
    ],
  }
})

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 16) : ''
}

onMounted(() => {
  if (route.query.student_id) {
    filters.value.student_id = String(route.query.student_id)
    reload()
  }
})
</script>

<style scoped>
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
.chart {
  height: 280px;
  width: 100%;
}
</style>