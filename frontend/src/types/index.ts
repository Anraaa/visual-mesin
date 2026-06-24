export interface User {
  id: number
  nip?: string
  user_name: string
  user_level: 'admin' | 'eng' | 'tech' | 'prod'
  email?: string
  avatar_url?: string
  department?: string
  jabatan?: string
  roles?: Role[]
  permissions?: string[]
}

export interface Role {
  id: number
  name: string
  guard_name: string
}

export interface Permission {
  id: number
  name: string
  guard_name: string
}

export interface MenuItem {
  key: string
  icon: React.ReactNode
  label: string
  children?: MenuItem[]
  roles?: string[]
}

export interface NavItem {
  key: string
  icon: React.ReactNode
  label: string
  children?: NavItem[]
  roles?: string[]
}

export interface APIResponse<T = unknown> {
  success: boolean
  message: string
  data: T
  meta?: {
    current_page: number
    per_page: number
    total: number
    last_page: number
  }
}

export interface UserToken {
  accessToken: string
  refreshToken?: string
}
