<template>
  <div class="full-center">
    <el-card class="form-card">
      <h2 class="title">重置密码</h2>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="手机号" prop="telephone">
          <el-input v-model="form.telephone" maxlength="11" />
        </el-form-item>
        <el-form-item label="验证码" prop="token">
          <div class="code-row">
            <el-input v-model="form.token" />
            <el-button :disabled="countdown > 0" @click="onSendCode">
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm">
          <el-input v-model="form.confirm" type="password" show-password />
        </el-form-item>

        <el-button type="primary" class="submit" :loading="loading" @click="onSubmit">
          重置密码
        </el-button>
      </el-form>

      <div class="footer-links">
        <router-link to="/login">返回登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { resetPassword, sendVerifyCode } from '@/api/account'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const countdown = ref(0)
let timer = null

const form = ref({ telephone: '', token: '', password: '', confirm: '' })

const rules = {
  telephone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1\d{10}$/, message: '手机号格式不正确', trigger: 'blur' },
  ],
  token: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' },
  ],
  confirm: [
    {
      validator: (_, v, cb) => (v === form.value.password ? cb() : cb(new Error('两次密码不一致'))),
      trigger: 'blur',
    },
  ],
}

onBeforeUnmount(() => timer && clearInterval(timer))

async function onSendCode() {
  if (!/^1\d{10}$/.test(form.value.telephone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  await sendVerifyCode(form.value.telephone)
  ElMessage.success('验证码已发送')
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function onSubmit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await resetPassword({
      telephone: form.value.telephone,
      password: form.value.password,
      token: form.value.token,
    })
    ElMessage.success('密码重置成功')
    router.push('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-card {
  width: 420px;
  border-radius: 12px;
}
.title {
  text-align: center;
  margin: 0 0 16px 0;
}
.code-row {
  display: flex;
  gap: 8px;
}
.code-row .el-input {
  flex: 1;
}
.submit {
  width: 100%;
}
.footer-links {
  text-align: center;
  margin-top: 14px;
  font-size: 13px;
}
.footer-links a {
  color: #409eff;
}
</style>