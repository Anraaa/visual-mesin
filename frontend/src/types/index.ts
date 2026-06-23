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
