import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type ColorPreset = 'blue' | 'cyan' | 'purple' | 'orange' | 'red'
export type SidebarStyle = 'vertical' | 'mini'

interface ThemeState {
  darkMode: boolean
  collapsed: boolean
  colorPreset: ColorPreset
  sidebarStyle: SidebarStyle
  breadcrumb: boolean
  toggleTheme: () => void
  toggleCollapsed: () => void
  setCollapsed: (v: boolean) => void
  setColorPreset: (preset: ColorPreset) => void
  setSidebarStyle: (style: SidebarStyle) => void
  toggleBreadcrumb: () => void
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      darkMode: false,
      collapsed: false,
      colorPreset: 'blue',
      sidebarStyle: 'vertical',
      breadcrumb: true,
      toggleTheme: () => set((state) => ({ darkMode: !state.darkMode })),
      toggleCollapsed: () => set((state) => ({ collapsed: !state.collapsed })),
      setCollapsed: (v) => set({ collapsed: v }),
      setColorPreset: (preset) => set({ colorPreset: preset }),
      setSidebarStyle: (style) => set({ sidebarStyle: style }),
      toggleBreadcrumb: () => set((state) => ({ breadcrumb: !state.breadcrumb })),
    }),
    { name: 'theme-storage' },
  ),
)

export const useThemeMode = () => useThemeStore((s) => s.darkMode)
export const useSidebarCollapsed = () => useThemeStore((s) => s.collapsed)
export const useColorPreset = () => useThemeStore((s) => s.colorPreset)
export const useSidebarStyle = () => useThemeStore((s) => s.sidebarStyle)
export const useBreadcrumbEnabled = () => useThemeStore((s) => s.breadcrumb)
