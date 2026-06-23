import { Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import Dashboard from './pages/Dashboard'
import Login from './pages/Login'
import UsersPage from './pages/UsersPage'
import RolesPage from './pages/RolesPage'
import DBConnectionsPage from './pages/DBConnectionsPage'
import DataProduksiPage from './pages/DataProduksiPage'
import ResourceDataPage from './pages/ResourceDataPage'
import ExportPage from './pages/ExportPage'
import AIChat from './pages/AIChat'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="admin/users" element={<UsersPage />} />
        <Route path="admin/roles" element={<RolesPage />} />
        <Route path="admin/db-connections" element={<DBConnectionsPage />} />
        <Route path="data" element={<DataProduksiPage />} />
        <Route path="data/:resource" element={<ResourceDataPage />} />
        <Route path="exports" element={<ExportPage />} />
        <Route path="analytics/ai-chat" element={<AIChat />} />
      </Route>
    </Routes>
  )
}

export default App
