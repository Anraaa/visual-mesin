import type { User } from '../types'
import { useAuthStore } from '../stores/authStore'

export function can(permission: string, user?: User): boolean {
  const u = user || useAuthStore.getState().user
  if (!u) return false
  // admin tetap harus punya permission yang sesuai
  if (!u.permissions || u.permissions.length === 0) return false
  if (u.permissions.includes(permission)) return true
  const [action, ...rest] = permission.split('-')
  const module = rest.join('-')
  const anyPerm = `view-any-${module}`
  return u.permissions.includes(anyPerm)
}

export function canAny(permissions: string[], user?: User): boolean {
  return permissions.some((p) => can(p, user))
}

export function canModule(module: string, user?: User): boolean {
  return canAny([`view-${module}`, `view-any-${module}`], user)
}
