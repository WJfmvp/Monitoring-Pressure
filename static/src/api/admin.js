import request from './request'

export function getAdminHome() {
  return request({ url: '/admin/home', method: 'get' })
}

// 用户管理
export function getUserList() {
  return request({ url: '/admin/user/list', method: 'get' })
}

export function updateUserRole(payload) {
  return request({
    url: '/admin/user/updateRole',
    method: 'post',
    data: payload, // { telephone, role }
  })
}

// Excel 成绩导入
export function importAcademicExcel(file) {
  const form = new FormData()
  form.append('file', file)
  return request({
    url: '/admin/academic/import',
    method: 'post',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getAcademicImportList(params) {
  return request({ url: '/admin/academic/import/list', method: 'get', params })
}

export function getAcademicList(params) {
  return request({ url: '/admin/academic/list', method: 'get', params })
}

// 压力评估列表
export function getStressResultList(params) {
  return request({ url: '/admin/stress/result/list', method: 'get', params })
}

// 干预建议
export function createInterventionSuggestion(payload) {
  return request({
    url: '/admin/intervention/create',
    method: 'post',
    data: payload,
  })
}

export function getInterventionSuggestionList(params) {
  return request({ url: '/admin/intervention/list', method: 'get', params })
}

// 压力预测：调用 Python 随机森林服务
export function predictStress(payload) {
  return request({
    url: '/admin/stress/predict',
    method: 'post',
    data: payload,
  })
}