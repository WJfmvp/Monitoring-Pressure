<template>
  <div class="page-card">
    <h2 class="page-title">成绩导入记录</h2>

    <div class="page-toolbar">
      <el-input v-model="filters.file_name" placeholder="文件名" clearable style="width: 200px" />
      <el-input
        v-model="filters.operator_id"
        placeholder="操作人 ID"
        clearable
        style="width: 160px"
      />
      <el-select v-model="filters.status" clearable placeholder="状态" style="width: 140px">
        <el-option label="处理中/失败" :value="0" />
        <el-option label="成功" :value="1" />
      </el-select>
      <el-button type="primary" @click="reload">查询</el-button>
    </div>

    <el-table :data="rows" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="file_name" label="文件名" />
      <el-table-column prop="operator_id" label="操作人" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? '成功' : '处理中/失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="row_count" label="行数" width="100" />
      <el-table-column prop="created_at" label="导入时间" width="180">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
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
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getAcademicImportList } from '@/api/admin'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filters = ref({ file_name: '', operator_id: '', status: '' })

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value }
  if (filters.value.file_name) p.file_name = filters.value.file_name
  if (filters.value.operator_id) p.operator_id = filters.value.operator_id
  if (filters.value.status !== '' && filters.value.status != null) p.status = filters.value.status
  return p
}

async function reload() {
  loading.value = true
  try {
    const data = await getAcademicImportList(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 19) : ''
}

onMounted(reload)
</script>