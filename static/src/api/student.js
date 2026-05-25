import request from './request'

export function getStudentHome() {
  return request({ url: '/student/home', method: 'get' })
}

export function getMyAcademicRecords(params) {
  return request({ url: '/student/academic', method: 'get', params })
}

export function submitPsychological(payload) {
  return request({
    url: '/student/psychological/submit',
    method: 'post',
    data: payload,
  })
}

export function getPsychologicalList() {
  return request({ url: '/student/psychological/list', method: 'get' })
}

export function getLatestStressResult() {
  return request({ url: '/student/stress/result/latest', method: 'get' })
}

export function getMyInterventions() {
  return request({ url: '/student/intervention/list', method: 'get' })
}