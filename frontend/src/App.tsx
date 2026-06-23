import { Routes, Route, Navigate } from 'react-router-dom'
import { App as AntApp } from 'antd'
import MainLayout from './layouts/MainLayout'
import Dashboard from './pages/Dashboard'
import Login from './pages/Login'
import NotFound from './pages/NotFound'
import Forbidden from './pages/Forbidden'
import ErrorBoundary from './components/ErrorBoundary'
import ProtectedRoute from './components/ProtectedRoute'

function App() {
  return (
    <ErrorBoundary>
      <AntApp>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/403" element={<Forbidden />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="data" element={<div style={{ padding: 24 }}>Data Produksi</div>} />
            <Route path="data/:resource" element={<div style={{ padding: 24 }}>Resource Table</div>} />
            <Route path="analytics/ai-chat" element={<div style={{ padding: 24 }}>AI Chat</div>} />
            <Route path="exports" element={<div style={{ padding: 24 }}>Export</div>} />
            <Route
              path="admin/users"
              element={
                <ProtectedRoute roles={['admin', 'eng']}>
                  <div style={{ padding: 24 }}>User Management</div>
                </ProtectedRoute>
              }
            />
            <Route
              path="admin/roles"
              element={
                <ProtectedRoute roles={['admin', 'eng']}>
                  <div style={{ padding: 24 }}>Roles & Permissions</div>
                </ProtectedRoute>
              }
            />
            <Route
              path="admin/db-connections"
              element={
                <ProtectedRoute roles={['admin', 'eng']}>
                  <div style={{ padding: 24 }}>DB Connections</div>
                </ProtectedRoute>
              }
            />
            <Route
              path="admin/resource-db-configs"
              element={
                <ProtectedRoute roles={['admin', 'eng']}>
                  <div style={{ padding: 24 }}>Resource DB Configs</div>
                </ProtectedRoute>
              }
            />
            <Route path="settings" element={<div style={{ padding: 24 }}>Settings</div>} />
          </Route>
          <Route path="*" element={<NotFound />} />
        </Routes>
      </AntApp>
    </ErrorBoundary>
  )
}

export default App
