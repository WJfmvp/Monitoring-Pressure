<template>
  <div class="page-card">
    <h2 class="page-title">我的成绩</h2>

    <div class="page-toolbar">
      <el-input
        v-model="filters.exam_name"
        placeholder="考试名称"
        clearable
        style="width: 180px"
      />
      <el-input v-model="filters.term" placeholder="学期" clearable style="width: 140px" />
      <el-input
        v-model.number="filters.min_score"
        placeholder="最低分"
        clearable
        style="width: 120px"
      />
      <el-input
        v-model.number="filters.max_score"
        placeholder="最高分"
        clearable
        style="width: 120px"
      />
      <el-button type="primary" @click="reload">查询</el-button>
      <el-button @click="onReset">重置</el-button>
    </div>

    <el-table :data="rows" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="exam_name" label="考试名称" />
      <el-table-column prop="term" label="学期" width="120" />
      <el-table-column prop="exam_date" label="考试日期" width="130">
        <template #default="{ row }">{{ formatDate(row.exam_date) }}</template>
      </el-table-column>
      <el-table-column prop="exam_score" label="成绩" width="90" />
      <el-table-column prop="average_score" label="平均分" width="90" />
      <el-table-column prop="score_fluctuation" label="波动" width="90" />
      <el-table-column prop="class_rank" label="班级排名" width="100" />
      <el-table-column prop="grade_rank" label="年级排名" width="100" />
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
import { getMyAcademicRecords } from '@/api/student'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({ exam_name: '', term: '', min_score: '', max_score: '' })

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value }
  if (filters.value.exam_name) p.exam_name = filters.value.exam_name
  if (filters.value.term) p.term = filters.value.term
  if (filters.value.min_score !== '' && filters.value.min_score != null)
    p.min_score = filters.value.min_score
  if (filters.value.max_score !== '' && filters.value.max_score != null)
    p.max_score = filters.value.max_score
  return p
}

async function reload() {
  loading.value = true
  try {
    const data = await getMyAcademicRecords(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function onReset() {
  filters.value = { exam_name: '', term: '', min_score: '', max_score: '' }
  page.value = 1
  reload()
}

function formatDate(t) {
  if (!t) return ''
  return String(t).slice(0, 10)
}

onMounted(reload)
</script>