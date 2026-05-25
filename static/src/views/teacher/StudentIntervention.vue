<template>
  <div class="page-card">
    <h2 class="page-title">学生干预记录</h2>

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
        v-model="filters.feedback_status"
        clearable
        placeholder="反馈状态"
        style="width: 140px"
      >
        <el-option label="未反馈" :value="0" />
        <el-option label="有效" :value="1" />
        <el-option label="无效" :value="2" />
      </el-select>
      <el-button type="primary" @click="reload">查询</el-button>
    </div>

    <el-empty v-if="!filters.student_id" description="请先输入学生 ID" />
    <template v-else>
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="建议" min-width="220">
          <template #default="{ row }">
            <div class="suggestion-title">{{ row.suggestion?.title || `建议 #${row.suggestion_id}` }}</div>
            <div class="suggestion-meta">
              {{ row.suggestion?.category || '-' }} · {{ LEVEL_LABEL[row.suggestion?.level] || '-' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="result_id" label="评估 ID" width="100" />
        <el-table-column prop="push_time" label="推送时间" width="180">
          <template #default="{ row }">{{ formatTime(row.push_time) }}</template>
        </el-table-column>
        <el-table-column label="建议内容" min-width="300">
          <template #default="{ row }">
            <el-tooltip
              v-if="row.suggestion?.content"
              :content="row.suggestion.content"
              placement="top"
            >
              <span class="content-preview">{{ row.suggestion.content }}</span>
            </el-tooltip>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="feedback" label="反馈内容" min-width="160" />
        <el-table-column prop="feedback_status" label="反馈状态" width="120">
          <template #default="{ row }">
            <el-tag :type="STATUS_TAG[row.feedback_status] || 'info'">
              {{ STATUS_LABEL[row.feedback_status] || '-' }}
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
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getStudentInterventionList } from '@/api/teacher'

const route = useRoute()

const STATUS_LABEL = ['未反馈', '有效', '无效']
const STATUS_TAG = ['info', 'success', 'danger']
const LEVEL_LABEL = ['正常', '中等', '高压']

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({ student_id: '', start_date: '', end_date: '', feedback_status: '' })

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value, student_id: filters.value.student_id }
  if (filters.value.start_date) p.start_date = filters.value.start_date
  if (filters.value.end_date) p.end_date = filters.value.end_date
  if (filters.value.feedback_status !== '' && filters.value.feedback_status != null)
    p.feedback_status = filters.value.feedback_status
  return p
}

async function reload() {
  if (!filters.value.student_id) {
    ElMessage.warning('请先输入学生 ID')
    return
  }
  loading.value = true
  try {
    const data = await getStudentInterventionList(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 19) : ''
}

onMounted(() => {
  if (route.query.student_id) {
    filters.value.student_id = String(route.query.student_id)
    reload()
  }
})
</script>

<style scoped>
.suggestion-title {
  font-weight: 600;
  color: #303133;
}
.suggestion-meta {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
}
.content-preview {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  white-space: nowrap;
}
</style>
