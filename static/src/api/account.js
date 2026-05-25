import request from './request'

// 验证码
export function sendVerifyCode(telephone) {
  return request({
    url: '/account/sendVerifyCode',
    method: 'get',
    params: { telephone },
  })
}

// 账号密码登录
export function loginByPassword(payload) {
  return request({
    url: '/account/login',
    method: 'post',
    data: payload, // { telephone, password }
  })
}

// 验证码登录
export function loginByCode(payload) {
  return request({
    url: '/account/login/code',
    method: 'post',
    data: payload, // { telephone, token: <verify_code> }
  })
}

export function register(payload) {
  return request({
    url: '/account/register',
    method: 'post',
    data: payload, // { username, password, telephone, token }
  })
}

export function resetPassword(payload) {
  return request({
    url: '/account/resetPassword',
    method: 'post',
    data: payload, // { telephone, password, token }
  })
}

// 已登录
export function getProfile() {
  return request({ url: '/account/profile', method: 'get' })
}

export function completeInformation(payload) {
  return request({
    url: '/account/completeInformation',
    method: 'post',
    data: payload, // { username, sex, email }
  })
}
