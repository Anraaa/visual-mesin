import { useLocation, Link } from 'react-router-dom'
import { Breadcrumb } from 'antd'
import { HomeOutlined } from '@ant-design/icons'
import { useMemo } from 'react'

const labelMap: Record<string, string> = {
  dashboard: 'Dashboard',
  data: 'Data Produksi',
  'ai-chat': 'AI Chat',
  exports: 'Export',
  admin: 'Admin',
  users: 'Users',
  roles: 'Roles',
  'db-connections': 'DB Connections',
  settings: 'Settings',
}

export default function BreadcrumbNav() {
  const location = useLocation()

  const items = useMemo(() => {
    const parts = location.pathname.split('/').filter(Boolean)
    const crumbs = [
      {
        key: '/',
        title: (
          <Link to="/">
            <HomeOutlined style={{ marginRight: 4 }} />
            Home
          </Link>
        ),
      },
    ]

    let path = ''
    for (const part of parts) {
      path += `/${part}`
      const label = labelMap[part] || part.charAt(0).toUpperCase() + part.slice(1)
      crumbs.push({
        key: path,
        title: <Link to={path}>{label}</Link>,
      })
    }

    return crumbs
  }, [location.pathname])

  return <Breadcrumb items={items} style={{ marginBottom: 16 }} />
}
