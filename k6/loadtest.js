import http from 'k6/http'
import { check, sleep, group } from 'k6'
import { Rate, Trend } from 'k6/metrics'

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

const loginFailureRate = new Rate('login_failures')
const apiResponseTime = new Trend('api_response_time')

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 200 },
    { duration: '30s', target: 500 },
    { duration: '1m', target: 500 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    login_failures: ['rate<0.05'],
  },
}

export default function () {
  group('Auth Flow', () => {
    const loginPayload = JSON.stringify({
      email: 'admin@admin.com',
      password: 'password',
    })

    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, loginPayload, {
      headers: { 'Content-Type': 'application/json' },
    })

    apiResponseTime.add(loginRes.timings.duration)
    loginFailureRate.add(loginRes.status !== 200)

    check(loginRes, {
      'login status is 200': (r) => r.status === 200,
      'login returns token': (r) => {
        try { return JSON.parse(r.body).data?.token != null }
        catch { return false }
      },
    })

    if (loginRes.status === 200) {
      let token
      try { token = JSON.parse(loginRes.body).data.token } catch { return }

      group('Authenticated APIs', () => {
        const authHeaders = {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        }

        const profileRes = http.get(`${BASE_URL}/api/v1/auth/me`, {
          headers: authHeaders,
        })
        apiResponseTime.add(profileRes.timings.duration)
        check(profileRes, {
          'profile status is 200': (r) => r.status === 200,
          'profile has user data': (r) => {
            try { return JSON.parse(r.body).data?.user_name != null }
            catch { return false }
          },
        })

        const rolesRes = http.get(`${BASE_URL}/api/v1/roles`, {
          headers: authHeaders,
        })
        apiResponseTime.add(rolesRes.timings.duration)
        check(rolesRes, {
          'roles status is 200': (r) => r.status === 200,
        })

        const usersRes = http.get(`${BASE_URL}/api/v1/users`, {
          headers: authHeaders,
        })
        apiResponseTime.add(usersRes.timings.duration)
        check(usersRes, {
          'users status is 200': (r) => r.status === 200,
        })
      })
    }
  })

  sleep(1)
}
