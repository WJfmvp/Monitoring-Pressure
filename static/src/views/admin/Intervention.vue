<template>
  <div class="page-card">
    <h2 class="page-title">干预建议管理</h2>

    <div class="page-toolbar">
      <el-select v-model="filters.level" clearable placeholder="等级" style="width: 140px">
        <el-option label="低" :value="0" />
        <el-option label="中" :value="1" />
        <el-option label="高" :value="2" />
      </el-select>
      <el-select v-model="filters.category" clearable placeholder="分类" style="width: 140px">
        <el-option label="学习" value="学习" />
        <el-option label="心理" value="心理" />
        <el-option label="生活" value="生活" />
      </el-select>
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
      <el-button type="success" @click="dialogVisible = true">新建建议</el-button>
    </div>

    <el-table :data="rows" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column prop="level" label="等级" width="100">
        <template #default="{ row }">
          <el-tag :type="LEVEL_TAG[row.level]">{{ LEVEL_LABEL[row.level] }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="分数区间" width="160">
        <template #default="{ row }">
          {{ row.min_score }} - {{ row.max_score }}
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容" />
      <el-table-column prop="created_at" label="创建时间" width="180">
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

    <el-dialog v-model="dialogVisible" title="新建干预建议" width="560px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="100" />
        </el-form-item>
        <el-form-item label="等级" prop="level">
          <el-radio-group v-model="form.level">
            <el-radio :value="0">低</el-radio>
            <el-radio :value="1">中</el-radio>
            <el-radio :value="2">高</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category" clearable>
            <el-option label="学习" value="学习" />
            <el-option label="心理" value="心理" />
            <el-option label="生活" value="生活" />
          </el-select>
        </el-form-item>
        <el-form-item label="最低分">
          <el-input-number v-model="form.min_score" :min="0" :max="100" :precision="2" />
        </el-form-item>
        <el-form-item label="最高分">
          <el-input-number v-model="form.max_score" :min="0" :max="100" :precision="2" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onCreate">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createInterventionSuggestion, getInterventionSuggestionList } from '@/api/admin'

const LEVEL_LABEL = ['低', '中', '高']
const LEVEL_TAG = ['success', 'warning', 'danger']

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filters = ref({ level: '', category: '', min_score: '', max_score: '' })

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref(null)
const form = ref({ title: '', level: 0, category: '', min_score: 0, max_score: 100, content: '' })

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  level: [{ required: true, message: '请选择等级', trigger: 'change' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
}

function buildParams() {
  const p = { page: page.value, page_size: pageSize.value }
  if (filters.value.level !== '' && filters.value.level != null) p.level = filters.value.level
  if (filters.value.category) p.category = filters.value.category
  if (filters.value.min_score !== '' && filters.value.min_score != null)
    p.min_score = filters.value.min_score
  if (filters.value.max_score !== '' && filters.value.max_score != null)
    p.max_score = filters.value.max_score
  return p
}

async function reload() {
  loading.value = true
  try {
    const data = await getInterventionSuggestionList(buildParams())
    rows.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

async function onCreate() {
  await formRef.value.validate()
  saving.value = true
  try {
    await createInterventionSuggestion(form.value)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    form.value = { title: '', level: 0, category: '', min_score: 0, max_score: 100, content: '' }
    reload()
  } finally {
    saving.value = false
  }
}

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 19) : ''
}

onMounted(reload)
</script>