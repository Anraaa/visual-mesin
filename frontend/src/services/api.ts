import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '../stores/authStore'
import type { APIResponse } from '../types'

class APIClient {
  private instance: AxiosInstance

  constructor(baseURL?: string) {
    this.instance = axios.create({
      baseURL: baseURL || import.meta.env.VITE_API_URL || 'http://localhost:8080',
      timeout: 50000,
      headers: { 'Content-Type': 'application/json' },
    })

    this.instance.interceptors.request.use((config) => {
      const token = useAuthStore.getState().token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    this.instance.interceptors.response.use(
      (res) => res,
      (error) => {
        if (error.response?.status === 401) {
          useAuthStore.getState().logout()
          window.location.href = '/login'
        }
        return Promise.reject(error)
      },
    )
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
    const res = await this.instance.get<APIResponse<T>>(url, config)
    return res.data
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
    const res = await this.instance.post<APIResponse<T>>(url, data, config)
    return res.data
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
    const res = await this.instance.put<APIResponse<T>>(url, data, config)
    return res.data
  }

  async delete<T = void>(url: string, config?: AxiosRequestConfig): Promise<APIResponse<T>> {
    const res = await this.instance.delete<APIResponse<T>>(url, config)
    return res.data
  }

  async download(url: string): Promise<Blob> {
    const res = await this.instance.get(url, { responseType: 'blob' })
    return res.data
  }
}

const api = new APIClient()
export default api
