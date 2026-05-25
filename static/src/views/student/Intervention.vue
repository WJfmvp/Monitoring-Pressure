<template>
  <div class="page-card">
    <h2 class="page-title">干预建议记录</h2>
    <el-empty v-if="!loading && list.length === 0" description="暂无建议推送" />
    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="建议" min-width="220">
        <template #default="{ row }">
          <div class="suggestion-title">{{ suggestionOf(row).title || `建议 #${row.suggestion_id}` }}</div>
          <div class="suggestion-meta">
            {{ suggestionOf(row).category || '-' }} · {{ LEVEL_LABEL[suggestionOf(row).level] || '-' }}
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
            v-if="suggestionOf(row).content"
            :content="suggestionOf(row).content"
            placement="top"
          >
            <span class="content-preview">{{ suggestionOf(row).content }}</span>
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
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getMyInterventions } from '@/api/student'

const list = ref([])
const loading = ref(false)

const STATUS_LABEL = ['未反馈', '有效', '无效']
const STATUS_TAG = ['info', 'success', 'danger']
const LEVEL_LABEL = ['正常', '中等', '高压']

function suggestionOf(row) {
  return row?.suggestion || row?.Suggestion || {
    title: row?.suggestion_title || row?.title || '',
    content: row?.suggestion_content || row?.content || '',
    category: row?.suggestion_category || row?.category || '',
    level: row?.suggestion_level ?? row?.level,
  }
}

function formatTime(t) {
  if (!t) return ''
  return String(t).replace('T', ' ').slice(0, 19)
}

async function load() {
  loading.value = true
  try {
    const data = await getMyInterventions()
    list.value = Array.isArray(data) ? data : []
  } finally {
    loading.value = false
  }
}

onMounted(load)
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
