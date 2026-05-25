<template>
  <div class="full-center">
    <el-card class="login-card">
      <h2 class="title">学业压力监测平台</h2>
      <p class="subtitle">请登录以继续</p>

      <el-tabs v-model="tab" class="tabs">
        <el-tab-pane label="账号密码" name="password">
          <el-form
            ref="pwdRef"
            :model="pwdForm"
            :rules="pwdRules"
            label-position="top"
            @keyup.enter="onPasswordLogin"
          >
            <el-form-item label="手机号" prop="telephone">
              <el-input v-model="pwdForm.telephone" placeholder="11 位手机号" maxlength="11" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="pwdForm.password" type="password" show-password />
            </el-form-item>
            <el-button type="primary" :loading="loading" class="submit" @click="onPasswordLogin">
              登录
            </el-button>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="验证码" name="code">
          <el-form
            ref="codeRef"
            :model="codeForm"
            :rules="codeRules"
            label-position="top"
            @keyup.enter="onCodeLogin"
          >
            <el-form-item label="手机号" prop="telephone">
              <el-input v-model="codeForm.telephone" maxlength="11" />
            </el-form-item>
            <el-form-item label="验证码" prop="token">
              <div class="code-row">
                <el-input v-model="codeForm.token" placeholder="6 位验证码" />
                <el-button :disabled="codeCountdown > 0" @click="onSendCode">
                  {{ codeCountdown > 0 ? `${codeCountdown}s` : '获取验证码' }}
                </el-button>
              </div>
            </el-form-item>
            <el-button type="primary" :loading="loading" class="submit" @click="onCodeLogin">
              登录
            </el-button>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="footer-links">
        <router-link to="/register">注册新账号</router-link>
        <router-link to="/reset-password">忘记密码？</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { loginByPassword, loginByCode, sendVerifyCode } from '@/api/account'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const tab = ref('password')
const loading = ref(false)

const phoneRule = {
  pattern: /^1\d{10}$/,
  message: '请输入正确的 11 位手机号',
  trigger: 'blur',
}

const pwdRef = ref(null)
const pwdForm = ref({ telephone: '', password: '' })
const pwdRules = {
  telephone: [{ required: true, message: '请输入手机号', trigger: 'blur' }, phoneRule],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const codeRef = ref(null)
const codeForm = ref({ telephone: '', token: '' })
const codeRules = {
  telephone: [{ required: true, message: '请输入手机号', trigger: 'blur' }, phoneRule],
  token: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

const codeCountdown = ref(0)
let timer = null

function startCountdown() {
  codeCountdown.value = 60
  timer = setInterval(() => {
    codeCountdown.value -= 1
    if (codeCountdown.value <= 0) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

onBeforeUnmount(() => timer && clearInterval(timer))

async function onSendCode() {
  if (!/^1\d{10}$/.test(codeForm.value.telephone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  await sendVerifyCode(codeForm.value.telephone)
  ElMessage.success('验证码已发送')
  startCountdown()
}

function postLoginRedirect() {
  const redirect = route.query.redirect
  if (redirect && typeof redirect === 'string' && redirect !== '/login') {
    router.replace(redirect)
  } else {
    router.replace('/dispatch')
  }
}

async function onPasswordLogin() {
  await pwdRef.value.validate()
  loading.value = true
  try {
    const res = await loginByPassword(pwdForm.value)
    userStore.setAuth(res.token, res.user_info)
    ElMessage.success('登录成功')
    postLoginRedirect()
  } finally {
    loading.value = false
  }
}

async function onCodeLogin() {
  await codeRef.value.validate()
  loading.value = true
  try {
    const res = await loginByCode(codeForm.value)
    userStore.setAuth(res.token, res.user_info)
    ElMessage.success('登录成功')
    postLoginRedirect()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-card {
  width: 420px;
  border-radius: 12px;
}
.title {
  margin: 0 0 4px 0;
  text-align: center;
  font-size: 22px;
}
.subtitle {
  text-align: center;
  color: #909399;
  margin: 0 0 16px 0;
}
.tabs {
  margin-top: 8px;
}
.submit {
  width: 100%;
}
.code-row {
  display: flex;
  gap: 8px;
}
.code-row .el-input {
  flex: 1;
}
.footer-links {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}
.footer-links a {
  color: #409eff;
}
</style>