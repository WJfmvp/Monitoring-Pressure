<template>
  <div class="page-card">
    <h2 class="page-title">成绩数据</h2>

    <div class="page-toolbar">
      <el-input v-model="filters.user_id" placeholder="学生 ID" clearable style="width: 160px" />
      <el-input v-model="filters.import_id" placeholder="导入 ID" clearable style="width: 140px" />
      <el-input v-model="filters.exam_name" placeholder="考试名称" clearable style="width: 180px" />
      <el-input v-model="filters.term" placeholder="学期" clearable style="width: 140px" />
      <el-select v-model="filters.source_type" clearable placeholder="来源" style="width: 140px">
        <el-option label="手动录入" :value="1" />
        <el-option label="Excel 导入" :value="2" />
      </el-select>
      <el-input
        v-model.number="filters.min_score"
        placeholder="最低分"
        clearable
        style="width: 110px"
      />
      <el-input
        v-model.number="filters.max_score"
        placeholder="最高分"
        clearable
        style="width: 110px"
      />
      <el-button type="primary" @click="reload">查询</el-button>
    </div>

    <el-table :data="rows" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_id" label="学生 ID" width="100" />
      <el-table-column prop="exam_name" label="考试名称" />
      <el-table-column prop="term" label="学期" width="100" />
      <el-table-column prop="exam_score" label="成绩" width="80" />
      <el-table-column prop="class_rank" label="班级排名" width="100" />
      <el-table-column prop="source_type" label="来源" width="100">
        <template #default="{ row }">
          {{ row.source_type === 1 ? '手动录入' : 'Excel 导入' }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="入库时间" width="180">
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
import { getAcademicList } from '@/api/admin'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({
  user_id: '',
  import_id: '',
  exam_name: '',
  term: '',
  source_type: '',
  min_score: '',
  max_score: '',
})

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value }
  for (const k of ['user_id', 'import_id', 'exam_name', 'term', 'source_type', 'min_score', 'max_score']) {
    const v = filters.value[k]
    if (v !== '' && v != null) p[k] = v
  }
  return p
}

async function reload() {
  loading.value = true
  try {
    const data = await getAcademicList(buildParams())
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