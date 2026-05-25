import request from './request'

export function getTeacherHome() {
  return request({ url: '/teacher/home', method: 'get' })
}

export function getStudentAcademicList(params) {
  return request({ url: '/teacher/student/academic/list', method: 'get', params })
}

export function getStudentStressList(params) {
  return request({ url: '/teacher/student/stress/list', method: 'get', params })
}

export function getStudentInterventionList(params) {
  return request({ url: '/teacher/student/intervention/list', method: 'get', params })
}

// 班级整体压力分布
export function getClassOverview() {
  return request({ url: '/teacher/class/overview', method: 'get' })
}

// 单个学生压力诊断书
export function getStudentDiagnosis(studentId) {
  return request({
    url: '/teacher/student/diagnosis',
    method: 'get',
    params: { student_id: studentId },
  })
}