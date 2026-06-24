import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'

interface AuthState {
  token: string | null
  user: User | null
  setAuth: (token: string, user: User) => void
  updateUser: (user: Partial<User>) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      setAuth: (token, user) => set({ token, user }),
      updateUser: (partial) => set((s) => ({ user: s.user ? { ...s.user, ...partial } : null })),
      logout: () => set({ token: null, user: null }),
      isAuthenticated: () => !!get().token,
    }),
    { name: 'auth-storage' },
  ),
)

export const useUserToken = () => useAuthStore((s) => s.token)
export const useUserInfo = () => useAuthStore((s) => s.user)
export const useUserLevel = () => useAuthStore((s) => s.user?.user_level)
