<template>
  <div class="page-card">
    <h2 class="page-title">用户管理</h2>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="user_id" label="用户 ID" width="120" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="telephone" label="手机号" width="140" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column prop="sex" label="性别" width="80">
        <template #default="{ row }">
          {{ row.sex === 1 ? '男' : row.sex === 2 ? '女' : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="role" label="角色" width="120">
        <template #default="{ row }">
          <el-tag :type="ROLE_TAG[row.role] || ''">{{ ROLE_LABEL[row.role] || row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openRoleDialog(row)">修改角色</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="修改用户角色" width="420px">
      <el-form label-width="80px">
        <el-form-item label="手机号">
          <el-input v-model="dialogForm.telephone" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="dialogForm.role" style="width: 100%">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onConfirm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getUserList, updateUserRole } from '@/api/admin'

const ROLE_LABEL = { student: '学生', teacher: '教师', admin: '管理员' }
const ROLE_TAG = { student: '', teacher: 'success', admin: 'danger' }

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const saving = ref(false)
const dialogForm = reactive({ telephone: '', role: 'student' })

async function load() {
  loading.value = true
  try {
    const data = await getUserList()
    list.value = data?.list || []
  } finally {
    loading.value = false
  }
}

function openRoleDialog(row) {
  dialogForm.telephone = row.telephone
  dialogForm.role = row.role || 'student'
  dialogVisible.value = true
}

async function onConfirm() {
  saving.value = true
  try {
    await updateUserRole({ telephone: dialogForm.telephone, role: dialogForm.role })
    ElMessage.success('修改成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>