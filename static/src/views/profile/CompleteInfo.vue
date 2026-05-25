<template>
  <div class="page-card">
    <h2 class="page-title">完善资料</h2>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 520px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="性别" prop="sex">
        <el-radio-group v-model="form.sex">
          <el-radio :value="1">男</el-radio>
          <el-radio :value="2">女</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="loading" @click="onSubmit">保存</el-button>
        <el-button @click="$router.push('/profile')">返回</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { completeInformation } from '@/api/account'

const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = ref({ username: '', sex: 1, email: '' })

const rules = {
  email: [
    {
      type: 'email',
      message: '请输入正确的邮箱',
      trigger: 'blur',
    },
  ],
}

onMounted(async () => {
  if (!userStore.userInfo) {
    try {
      await userStore.fetchProfile()
    } catch (_) {}
  }
  const info = userStore.userInfo
  if (info) {
    form.value.username = info.username || ''
    form.value.sex = info.sex || 1
    form.value.email = info.email || ''
  }
})

async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await completeInformation(form.value)
    await userStore.fetchProfile()
    ElMessage.success('保存成功')
  } finally {
    loading.value = false
  }
}
</script>