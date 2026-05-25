import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth } from '@/utils/auth'

// 后端有几种返回结构，需要在拦截器内做兼容：
//   1. util.ResponseSuccess  → { code: 0,   msg, data }
//   2. util.ResponseError    → { code: !=0, msg }
//   3. 部分手写            → { code: 200, msg, data }
//   4. 部分手写            → { message, data }
//   5. 部分错误            → { error: "..." }
const AUTH_ERROR_CODES = new Set([1007, 1008, 1009, 1010])

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api',
  timeout: 15000,
})

service.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

function handleAuthFailure(msg) {
  clearAuth()
  ElMessage.error(msg || '登录已失效，请重新登录')
  // 通过整页跳转规避循环重定向
  if (location.hash !== '#/login' && location.pathname !== '/login') {
    location.href = '/login'
  }
}

service.interceptors.response.use(
  res => {
    const data = res.data

    // 形式 1/2/3：带 code 字段
    if (data && typeof data.code === 'number') {
      if (data.code === 0 || data.code === 200) {
        // data 可能放在 data.data，也可能本身就是 payload
        return data.data !== undefined ? data.data : data
      }
      const msg = data.msg || data.message || '请求失败'
      if (AUTH_ERROR_CODES.has(data.code)) {
        handleAuthFailure(msg)
      } else {
        ElMessage.error(msg)
      }
      return Promise.reject(new Error(msg))
    }

    // 形式 5：error 字段
    if (data && data.error) {
      ElMessage.error(data.error)
      return Promise.reject(new Error(data.error))
    }

    // 形式 4：{message, data} 或纯对象
    if (data && data.data !== undefined) return data.data
    return data
  },
  err => {
    if (err.response) {
      const { status, data } = err.response
      if (status === 401 || status === 403) {
        handleAuthFailure(data?.msg || data?.error || '未授权')
        return Promise.reject(err)
      }
      const msg = data?.msg || data?.message || data?.error || `请求失败(${status})`
      ElMessage.error(msg)
    } else {
      ElMessage.error(err.message || '网络异常')
    }
    return Promise.reject(err)
  }
)

export default service