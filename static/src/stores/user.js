import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getProfile } from '@/api/account'
import { getToken, setToken, getUser, setUser, clearAuth } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(getToken())
  const userInfo = ref(getUser())

  const isLogin = computed(() => Boolean(token.value))
  const role = computed(() => userInfo.value?.role || '')

  function setAuth(t, info) {
    token.value = t || ''
    userInfo.value = info || null
    setToken(token.value)
    setUser(userInfo.value)
  }

  async function fetchProfile() {
    const data = await getProfile()
    // GetProfile 返回 { user_info: {...} }
    const info = data?.user_info || data
    userInfo.value = info
    setUser(info)
    return info
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    clearAuth()
  }

  return { token, userInfo, isLogin, role, setAuth, fetchProfile, logout }
})