<template>
  <div>
    <div class="page-card">
      <h2 class="page-title">Excel 成绩导入</h2>
      <p class="hint">支持 .xlsx / .xls 格式，操作人取自当前管理员的 JWT。</p>
    </div>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="14">
        <div class="page-card">
          <h3 class="card-title">
            <el-icon><Upload /></el-icon>
            上传文件
          </h3>

          <el-upload
            drag
            :auto-upload="false"
            :on-change="onChange"
            :show-file-list="true"
            :limit="1"
            :file-list="fileList"
            accept=".xlsx,.xls"
            class="upload-box"
          >
            <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
            <div class="el-upload__text">
              点击或拖拽 Excel 文件到此处<br />
              <span class="hint">只接受 .xlsx / .xls，单次最多 1 个文件</span>
            </div>
          </el-upload>

          <div class="actions">
            <el-button
              type="primary"
              :loading="uploading"
              :disabled="!currentFile"
              @click="onUpload"
            >
              <el-icon style="margin-right: 4px"><Upload /></el-icon>
              开始上传
            </el-button>
            <el-button @click="reset">清空</el-button>
          </div>

          <el-alert v-if="result" :title="result" type="success" :closable="false" show-icon />
        </div>
      </el-col>

      <el-col :span="10">
        <div class="page-card guide-card">
          <h3 class="card-title">
            <el-icon><Document /></el-icon>
            导入说明
          </h3>
          <ol class="guide-list">
            <li>
              <strong>文件格式</strong>
              <p>仅支持 .xlsx 和 .xls 两种格式，单文件不超过 20MB。</p>
            </li>
            <li>
              <strong>必填列</strong>
              <p>user_id、exam_name、term、exam_date、raw_score 至少包含，其余字段会在数据清洗时自动计算（如平均分、波动）。</p>
            </li>
            <li>
              <strong>导入失败</strong>
              <p>导入记录会在「导入记录」页面显示，状态为「失败」时可查看 error_message。</p>
            </li>
            <li>
              <strong>重复导入</strong>
              <p>同一文件可重复上传，每次会生成新的导入批次（import_id 不同）。</p>
            </li>
          </ol>
        </div>
      </el-col>
    </el-row>

    <div class="page-card" style="margin-top: 16px">
      <div class="card-header">
        <h3 class="card-title" style="margin: 0">
          <el-icon><Clock /></el-icon>
          最近导入记录
        </h3>
        <el-button link type="primary" @click="$router.push('/admin/academic-imports')">
          查看全部 →
        </el-button>
      </div>

      <el-table :data="recentImports" v-loading="historyLoading" stripe size="small">
        <el-table-column prop="id" label="批次 ID" width="100" />
        <el-table-column prop="file_name" label="文件名" />
        <el-table-column prop="operator_id" label="操作人" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === 1 ? 'success' : 'warning'">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="import_count" label="条数" width="100" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { importAcademicExcel, getAcademicImportList } from '@/api/admin'

const currentFile = ref(null)
const fileList = ref([])
const uploading = ref(false)
const result = ref('')

const recentImports = ref([])
const historyLoading = ref(false)

function onChange(file) {
  currentFile.value = file.raw
  fileList.value = [file]
}

function reset() {
  currentFile.value = null
  fileList.value = []
  result.value = ''
}

async function onUpload() {
  if (!currentFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  uploading.value = true
  try {
    await importAcademicExcel(currentFile.value)
    result.value = `已成功导入文件：${currentFile.value.name}`
    ElMessage.success('导入成功')
    loadRecent()
  } finally {
    uploading.value = false
  }
}

async function loadRecent() {
  historyLoading.value = true
  try {
    const data = await getAcademicImportList({ page: 1, page_size: 5 })
    recentImports.value = data?.list || []
  } finally {
    historyLoading.value = false
  }
}

function formatTime(t) {
  return t ? String(t).replace('T', ' ').slice(0, 19) : ''
}

onMounted(loadRecent)
</script>

<style scoped>
.hint {
  color: #909399;
  margin: 6px 0 0 0;
  font-size: 13px;
}
.card-title {
  margin: 0 0 12px 0;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.actions {
  margin: 16px 0;
  display: flex;
  gap: 8px;
}
.upload-box {
  width: 100%;
}

.guide-list {
  padding-left: 18px;
  color: #606266;
  margin: 0;
}
.guide-list li {
  margin-bottom: 14px;
}
.guide-list li strong {
  display: block;
  color: #303133;
  margin-bottom: 4px;
}
.guide-list li p {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>