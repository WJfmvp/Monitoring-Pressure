<template>
  <div class="page-card">
    <h2 class="page-title">学生成绩查询</h2>

    <div class="page-toolbar">
      <el-input
        v-model="filters.student_id"
        placeholder="学生 ID（必填）"
        clearable
        style="width: 200px"
      />
      <el-input v-model="filters.exam_name" placeholder="考试名称" clearable style="width: 180px" />
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
      <el-button type="primary" @click="reload">查询</el-button>
    </div>

    <el-empty v-if="!filters.student_id" description="请先输入学生 ID" />
    <template v-else>
      <el-table :data="rows" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="exam_name" label="考试名称" />
        <el-table-column prop="term" label="学期" width="120" />
        <el-table-column prop="exam_date" label="考试日期" width="130">
          <template #default="{ row }">{{ formatDate(row.exam_date) }}</template>
        </el-table-column>
        <el-table-column prop="exam_score" label="成绩" width="90" />
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
    </template>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getStudentAcademicList } from '@/api/teacher'

const route = useRoute()
const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({ student_id: '', exam_name: '', start_date: '', end_date: '' })

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value, student_id: filters.value.student_id }
  if (filters.value.exam_name) p.exam_name = filters.value.exam_name
  if (filters.value.start_date) p.start_date = filters.value.start_date
  if (filters.value.end_date) p.end_date = filters.value.end_date
  return p
}

async function reload() {
  if (!filters.value.student_id) {
    ElMessage.warning('请先输入学生 ID')
    return
  }
  loading.value = true
  try {
    const data = await getStudentAcademicList(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function formatDate(t) {
  return t ? String(t).slice(0, 10) : ''
}

onMounted(() => {
  if (route.query.student_id) {
    filters.value.student_id = String(route.query.student_id)
    reload()
  }
})
</script>