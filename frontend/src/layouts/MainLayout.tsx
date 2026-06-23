import { useState, useCallback } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  Layout, Menu, Button, theme, Dropdown, Avatar, Badge, List, Typography, Popover, Space,
} from 'antd'
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
  TeamOutlined,
  SafetyOutlined,
  LinkOutlined,
  ExportOutlined,
  BellOutlined,
  HistoryOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { useWebSocket } from '../hooks/useWebSocket'
import { useNotificationStore } from '../stores/notificationStore'

const { Header, Sider, Content } = Layout
const { Text } = Typography

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' },
  { key: '/data', icon: <DatabaseOutlined />, label: 'Data Produksi' },
  { key: '/exports', icon: <ExportOutlined />, label: 'Export' },
  {
    key: 'admin',
    icon: <TeamOutlined />,
    label: 'Administrasi',
    children: [
      { key: '/admin/users', icon: <UserOutlined />, label: 'Users' },
      { key: '/admin/roles', icon: <SafetyOutlined />, label: 'Roles' },
      { key: '/admin/db-connections', icon: <LinkOutlined />, label: 'DB Connections' },
    ],
  },
]

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { darkMode, toggleTheme } = useThemeStore()
  const { token: themeToken } = theme.useToken()
  const { notifications, addNotification, markAllRead, unreadCount } = useNotificationStore()

  const wsHandlers: Record<string, (payload: any) => void> = {
    notification: useCallback((payload: any) => {
      addNotification({
        type: payload.type,
        title: payload.title,
        message: payload.message,
      })
    }, [addNotification]),
    pong: () => {},
  }

  useWebSocket(wsHandlers)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenu = {
    items: [
      { key: 'profile', icon: <UserOutlined />, label: 'Profile' },
      { key: 'settings', icon: <SettingOutlined />, label: 'Settings' },
      { type: 'divider' as const },
      { key: 'logout', icon: <LogoutOutlined />, label: 'Logout', danger: true },
    ],
    onClick: ({ key }: { key: string }) => {
      if (key === 'logout') handleLogout()
    },
  }

  const selectedKey = '/' + location.pathname.split('/').slice(1, 3).join('/')

  const notificationContent = (
    <div style={{ width: 360, maxHeight: 400, overflow: 'auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 12px', borderBottom: '1px solid #f0f0f0' }}>
        <Text strong>Notifikasi</Text>
        <Button type="link" size="small" onClick={markAllRead}>Tandai Dibaca</Button>
      </div>
      {notifications.length === 0 ? (
        <div style={{ padding: 24, textAlign: 'center', color: '#999' }}>Tidak ada notifikasi</div>
      ) : (
        <List
          dataSource={notifications.slice(0, 20)}
          renderItem={(item) => (
            <List.Item
              style={{
                padding: '8px 12px',
                background: item.read ? 'transparent' : '#e6f4ff',
                cursor: 'pointer',
              }}
              onClick={() => useNotificationStore.getState().markRead(item.id)}
            >
              <List.Item.Meta
                title={<Text strong={!item.read}>{item.title}</Text>}
                description={<Text type="secondary" style={{ fontSize: 12 }}>{item.message}</Text>}
              />
            </List.Item>
          )}
        />
      )}
    </div>
  )

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
            fontWeight: 'bold',
            fontSize: collapsed ? 14 : 18,
          }}
        >
          {collapsed ? 'VM' : 'Visual Mesin'}
        </div>
        <Menu
          theme={darkMode ? 'dark' : 'light'}
          mode="inline"
          selectedKeys={[selectedKey]}
          defaultOpenKeys={['admin']}
          items={menuItems}
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
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <Space size="middle">
            <Button
              type="text"
              icon={<HistoryOutlined />}
              onClick={() => navigate('/admin/activity-logs')}
            />
            <Popover
              content={notificationContent}
              trigger="click"
              placement="bottomRight"
            >
              <Badge count={unreadCount()} size="small">
                <Button type="text" icon={<BellOutlined />} />
              </Badge>
            </Popover>
            <Button
              type="text"
              icon={darkMode ? <SunOutlined /> : <MoonOutlined />}
              onClick={toggleTheme}
            />
            <Dropdown menu={userMenu} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
                <Avatar icon={<UserOutlined />} />
                <span>{user?.name || 'User'}</span>
              </div>
            </Dropdown>
          </Space>
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
