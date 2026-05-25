<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">最新压力检测结果</h2>

      <el-empty v-if="!latest" description="暂无评估结果" />

      <template v-else>
        <el-row :gutter="16">
          <el-col :span="8">
            <div class="metric-card">
              <span class="label">综合得分</span>
              <span class="value">{{ latest.total_score?.toFixed?.(2) ?? '-' }}</span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="metric-card">
              <span class="label">模型预测等级</span>
              <span class="value" :style="{ color: LEVEL_COLORS[latest.model_predict_level] }">
                {{ LEVEL_LABELS[latest.model_predict_level] || '-' }}
              </span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="metric-card">
              <span class="label">综合预警等级</span>
              <span class="value" :style="{ color: LEVEL_COLORS[latest.warning_status] }">
                {{ LEVEL_LABELS[latest.warning_status] || '-' }}
              </span>
            </div>
          </el-col>
        </el-row>

        <el-descriptions :column="2" border style="margin-top: 16px">
          <el-descriptions-item label="行为得分">{{ latest.behavior_score }}</el-descriptions-item>
          <el-descriptions-item label="学业得分">{{ latest.academic_score }}</el-descriptions-item>
          <el-descriptions-item label="心理得分">
            {{ latest.psychological_score }}
          </el-descriptions-item>
          <el-descriptions-item label="预测概率">
            {{ ((latest.predict_probability || 0) * 100).toFixed(2) }}%
          </el-descriptions-item>
          <el-descriptions-item label="预警时间">{{ latest.warning_time || '-' }}</el-descriptions-item>
          <el-descriptions-item label="生成时间">{{ latest.created_at }}</el-descriptions-item>
        </el-descriptions>

        <h3 class="card-title" style="margin-top: 24px">维度对比</h3>
        <v-chart class="chart" :option="barOption" autoresize />
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
// 别名：避免与 Element Plus 同名图标冲突
import { BarChart as EChartsBar } from 'echarts/charts'
import { TooltipComponent, GridComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { getLatestStressResult } from '@/api/student'

use([CanvasRenderer, EChartsBar, TooltipComponent, GridComponent])

const latest = ref(null)

const LEVEL_LABELS = ['正常', '中等', '高压力']
const LEVEL_COLORS = ['#67c23a', '#e6a23c', '#f56c6c']

const barOption = computed(() => ({
  tooltip: {},
  grid: { left: 40, right: 20, top: 30, bottom: 30 },
  xAxis: { type: 'category', data: ['行为', '学业', '心理', '综合'] },
  yAxis: { type: 'value', max: 100 },
  series: [
    {
      type: 'bar',
      barWidth: 32,
      itemStyle: { color: '#6e8efb', borderRadius: [4, 4, 0, 0] },
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

onMounted(async () => {
  try {
    latest.value = await getLatestStressResult()
  } catch (_) {}
})
</script>

<style scoped>
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
}
.chart {
  height: 320px;
  width: 100%;
}
</style>