import { useState, useMemo } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Button, theme, Dropdown, Avatar, Typography, Space, Tag } from 'antd'
import type { MenuProps } from 'antd'
import {
  DashboardOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
  MoonOutlined,
  SunOutlined,
  DatabaseOutlined,
  RobotOutlined,
  ExportOutlined,
  SafetyCertificateOutlined,
  ApiOutlined,
  AppstoreOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import BreadcrumbNav from '../components/BreadcrumbNav'
import type { MenuItem } from '../types'

const { Header, Sider, Content } = Layout

type AntMenuItem = Required<MenuProps>['items'][number]

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { darkMode, toggleTheme } = useThemeStore()
  const { token: themeToken } = theme.useToken()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenu: MenuProps = {
    items: [
      {
        key: 'profile',
        icon: <UserOutlined />,
        label: (
          <Space direction="vertical" size={0}>
            <span>{user?.user_name || 'User'}</span>
            <Tag color="blue" style={{ fontSize: 11, margin: 0 }}>
              {user?.user_level?.toUpperCase() || 'PROD'}
            </Tag>
          </Space>
        ),
      },
      { type: 'divider' },
      { key: 'settings', icon: <SettingOutlined />, label: 'Settings' },
      { type: 'divider' },
      { key: 'logout', icon: <LogoutOutlined />, label: 'Logout', danger: true },
    ],
    onClick: ({ key }) => {
      if (key === 'logout') handleLogout()
      if (key === 'settings') navigate('/settings')
    },
  }

  const allMenuItems: MenuItem[] = [
    { key: '/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' },
    { key: '/data', icon: <DatabaseOutlined />, label: 'Data Produksi' },
    { key: '/analytics/ai-chat', icon: <RobotOutlined />, label: 'AI Chat' },
    { key: '/exports', icon: <ExportOutlined />, label: 'Export' },
    {
      key: '/admin',
      icon: <SafetyCertificateOutlined />,
      label: 'Admin',
      roles: ['admin', 'eng'],
      children: [
        { key: '/admin/users', icon: <UserOutlined />, label: 'Users' },
        { key: '/admin/roles', icon: <SafetyCertificateOutlined />, label: 'Roles & Permissions' },
        { key: '/admin/db-connections', icon: <ApiOutlined />, label: 'DB Connections' },
        { key: '/admin/resource-db-configs', icon: <AppstoreOutlined />, label: 'Resource DB' },
      ],
    },
    { key: '/settings', icon: <SettingOutlined />, label: 'Settings' },
  ]

  const visibleMenuItems = useMemo(() => {
    const filterByRole = (items: MenuItem[]): AntMenuItem[] => {
      return items
        .filter((item) => !item.roles || (user && item.roles.includes(user.user_level)))
        .map((item) => ({
          key: item.key,
          icon: item.icon,
          label: item.label,
          children: item.children ? filterByRole(item.children) : undefined,
        }))
    }
    return filterByRole(allMenuItems)
  }, [user])

  const selectedKeys = useMemo(() => {
    const path = location.pathname
    const parts = path.split('/').filter(Boolean)
    const keys: string[] = [path]
    if (parts.length > 1) {
      keys.push('/' + parts[0])
    }
    return keys
  }, [location.pathname])

  const openKeys = useMemo(() => {
    const path = location.pathname
    const parts = path.split('/').filter(Boolean)
    if (parts.length > 1) {
      return ['/' + parts[0]]
    }
    return []
  }, [location.pathname])

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme={darkMode ? 'dark' : 'light'}
        style={{ borderRight: `1px solid ${themeToken.colorBorderSecondary}` }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontWeight: 700,
            fontSize: collapsed ? 14 : 18,
            letterSpacing: 1,
            borderBottom: `1px solid ${themeToken.colorBorderSecondary}`,
          }}
        >
          {collapsed ? 'VM' : 'Visual Mesin'}
        </div>
        <Menu
          theme={darkMode ? 'dark' : 'light'}
          mode="inline"
          selectedKeys={selectedKeys}
          defaultOpenKeys={openKeys}
          items={visibleMenuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 24px',
            background: darkMode ? '#141414' : '#fff',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: `1px solid ${themeToken.colorBorderSecondary}`,
            height: 64,
          }}
        >
          <Space>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
            />
          </Space>
          <Space size="middle">
            <Button
              type="text"
              icon={darkMode ? <SunOutlined /> : <MoonOutlined />}
              onClick={toggleTheme}
            />
            <Dropdown menu={userMenu} placement="bottomRight" trigger={['click']}>
              <Space style={{ cursor: 'pointer' }}>
                <Avatar
                  icon={<UserOutlined />}
                  style={{ backgroundColor: themeToken.colorPrimary }}
                />
                <Typography.Text
                  ellipsis
                  style={{ maxWidth: 120, verticalAlign: 'middle' }}
                >
                  {user?.user_name || 'User'}
                </Typography.Text>
              </Space>
            </Dropdown>
          </Space>
        </Header>
        <Content style={{ margin: 24 }}>
          <BreadcrumbNav />
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
