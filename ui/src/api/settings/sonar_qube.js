import request from '@/utils/request'

export function listSonarQube() {
  return request({
    url: '/settings/sonar_qube',
    method: 'get',
  })
}

export function createSonarQube(data) {
  return request({
    url: '/settings/sonar_qube',
    method: 'post',
    data,
  })
}

export function updateSonarQube(id, data) {
  return request({
    url: `/settings/sonar_qube/${id}`,
    method: 'put',
    data,
  })
}

export function deleteSonarQube(id) {
  return request({
    url: `/settings/sonar_qube/${id}`,
    method: 'delete',
  })
}