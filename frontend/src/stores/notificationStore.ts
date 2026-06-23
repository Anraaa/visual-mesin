import { create } from 'zustand'

export interface Notification {
  id: string
  type: string
  title: string
  message: string
  timestamp: number
  read: boolean
}

interface NotificationState {
  notifications: Notification[]
  addNotification: (notif: Omit<Notification, 'id' | 'timestamp' | 'read'>) => void
  markRead: (id: string) => void
  markAllRead: () => void
  clearAll: () => void
  unreadCount: () => number
}

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],
  addNotification: (notif) =>
    set((state) => ({
      notifications: [
        {
          ...notif,
          id: Date.now().toString(),
          timestamp: Date.now(),
          read: false,
        },
        ...state.notifications,
      ],
    })),
  markRead: (id) =>
    set((state) => ({
      notifications: state.notifications.map((n) =>
        n.id === id ? { ...n, read: true } : n,
      ),
    })),
  markAllRead: () =>
    set((state) => ({
      notifications: state.notifications.map((n) => ({ ...n, read: true })),
    })),
  clearAll: () => set({ notifications: [] }),
  unreadCount: () => get().notifications.filter((n) => !n.read).length,
}))
