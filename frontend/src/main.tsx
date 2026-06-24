import { StrictMode, lazy, Suspense, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ConfigProvider, theme as antTheme, Spin } from 'antd'
import { useThemeStore, useColorPreset } from './stores/themeStore'
import MainLayout from './layouts/MainLayout'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, refetchOnWindowFocus: false, staleTime: 30_000 },
  },
})

const Login = lazy(() => import('./pages/Login'))
const Dashboard = lazy(() => import('./pages/Dashboard'))
const UsersPage = lazy(() => import('./pages/UsersPage'))
const RolesPage = lazy(() => import('./pages/RolesPage'))
const RolePermissionPage = lazy(() => import('./pages/RolePermissionPage'))
const ResourceConnectionPage = lazy(() => import('./pages/ResourceConnectionPage'))
const DataProduksiConfigPage = lazy(() => import('./pages/DataProduksiConfigPage'))
const DataProduksiPage = lazy(() => import('./pages/DataProduksiPage'))
const ResourceDataPage = lazy(() => import('./pages/ResourceDataPage'))
const ExportPage = lazy(() => import('./pages/ExportPage'))
const ActivityLogPage = lazy(() => import('./pages/ActivityLogPage'))
const AIChat = lazy(() => import('./pages/AIChat'))
const Forbidden = lazy(() => import('./pages/Forbidden'))
const NotFound = lazy(() => import('./pages/NotFound'))

function Lazy({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: 300,
        flexDirection: 'column',
        gap: 16,
      }}>
        <Spin size="large" />
        <span style={{ color: 'var(--text-tertiary)', fontSize: 13 }}>Memuat...</span>
      </div>
    }>
      {children}
    </Suspense>
  )
}

const antdThemeTokens = {
  borderRadius: 10,
  borderRadiusSM: 6,
  controlHeight: 36,
  fontSize: 14,
  colorLink: '#1677ff',
}

function getAntdAlgorithm(darkMode: boolean, colorPreset: string) {
  const base = darkMode ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm

  const presetColors: Record<string, string> = {
    blue: '#1677ff', cyan: '#13c2c2', purple: '#722ed1', orange: '#fa8c16', red: '#f5222d',
  }

  const presetPrimary = presetColors[colorPreset] || '#1677ff'

  return (seedToken: any) => {
    const derived = base(seedToken)
    return {
      ...derived,
      colorPrimary: presetPrimary,
      colorLink: presetPrimary,
      colorLinkHover: presetPrimary,
    }
  }
}

function ThemeSync() {
  const { darkMode } = useThemeStore()
  const colorPreset = useColorPreset()

  useEffect(() => {
    document.documentElement.setAttribute('data-theme-mode', darkMode ? 'dark' : 'light')
    document.documentElement.setAttribute('data-color-preset', colorPreset)
  }, [darkMode, colorPreset])

  return null
}

function Root() {
  const { darkMode, colorPreset } = useThemeStore()
  const [ready, setReady] = useState(false)

  useEffect(() => { setReady(true) }, [])

  if (!ready) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        background: '#f4f5f7',
        flexDirection: 'column',
        gap: 16,
      }}>
        <div style={{
          width: 48, height: 48, borderRadius: 14,
          background: 'linear-gradient(135deg, #1677ff, #4096ff)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: 22, fontWeight: 800, color: '#fff',
        }}>
          VM
        </div>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <ConfigProvider
      theme={{
        algorithm: getAntdAlgorithm(darkMode, colorPreset),
        token: antdThemeTokens,
      }}
    >
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <ThemeSync />
          <Routes>
            <Route path="/login" element={<Lazy><Login /></Lazy>} />
            <Route path="/" element={<MainLayout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Lazy><Dashboard /></Lazy>} />
              <Route path="admin/users" element={<Lazy><UsersPage /></Lazy>} />
              <Route path="admin/roles" element={<Lazy><RolesPage /></Lazy>} />
              <Route path="admin/role-permission" element={<Lazy><RolePermissionPage /></Lazy>} />
              <Route path="admin/resource-connection" element={<Lazy><ResourceConnectionPage /></Lazy>} />
              <Route path="admin/data-produksi-config" element={<Lazy><DataProduksiConfigPage /></Lazy>} />
              <Route path="admin/activity-logs" element={<Lazy><ActivityLogPage /></Lazy>} />
              <Route path="data" element={<Lazy><DataProduksiPage /></Lazy>} />
              <Route path="data/:resource" element={<Lazy><ResourceDataPage /></Lazy>} />
              <Route path="exports" element={<Lazy><ExportPage /></Lazy>} />
              <Route path="analytics/ai-chat" element={<Lazy><AIChat /></Lazy>} />
              <Route path="403" element={<Lazy><Forbidden /></Lazy>} />
              <Route path="*" element={<Lazy><NotFound /></Lazy>} />
            </Route>
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ConfigProvider>
  )
}

createRoot(document.getElementById('root')!).render(
  <StrictMode><Root /></StrictMode>,
)
