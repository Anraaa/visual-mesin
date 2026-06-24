import { useCallback, useMemo, useState, useEffect } from 'react'
import { Outlet, useNavigate, useLocation, Navigate } from 'react-router-dom'
import {
  Layout, Menu, Button, theme, Dropdown, Avatar, Badge, List, Typography, Popover, Space, Breadcrumb,
  Tooltip, Divider, Switch, Drawer, Result,
} from 'antd'
import { MenuOutlined } from '@ant-design/icons'
import {
  DashboardOutlined,
  MenuFoldOutlined, MenuUnfoldOutlined,
  UserOutlined, LogoutOutlined,
  MoonOutlined, SunOutlined,
  DatabaseOutlined, TeamOutlined, SafetyOutlined, LinkOutlined,
  ExportOutlined, BellOutlined, HistoryOutlined,
  RobotOutlined, BgColorsOutlined, SettingOutlined,
} from '@ant-design/icons'
import { useAuthStore, useUserLevel } from '../stores/authStore'
import { useThemeStore, useBreadcrumbEnabled, type ColorPreset } from '../stores/themeStore'
import { useWebSocket } from '../hooks/useWebSocket'
import { useNotificationStore } from '../stores/notificationStore'
import { canModule } from '../utils/permissions'
import api from '../services/api'

const { Header, Sider, Content } = Layout
const { Text } = Typography

const labelMap: Record<string, string> = {
  dashboard: 'Dashboard', data: 'Data Produksi', 'ai-chat': 'AI Chat',
  exports: 'Export', admin: 'Administrasi', users: 'User Management', roles: 'Role Management',
  'role-permission': 'Role Permission',   'resource-connection': 'Resource Connection', 'data-produksi-config': 'Data Produksi Config', 'activity-logs': 'Activity Log',
  analytics: 'Analytics',
}

const colorPresetOptions: { value: ColorPreset; label: string; color: string }[] = [
  { value: 'blue', label: 'Biru', color: '#1677ff' },
  { value: 'cyan', label: 'Cyan', color: '#13c2c2' },
  { value: 'purple', label: 'Ungu', color: '#722ed1' },
  { value: 'orange', label: 'Oranye', color: '#fa8c16' },
  { value: 'red', label: 'Merah', color: '#f5222d' },
]

function getMenuItems(userLevel?: string, permissions?: string[]) {
  const hasPerm = (module: string) => userLevel === 'admin' || canModule(module, { user_level: userLevel, permissions } as any)

  const items: any[] = []

  if (hasPerm('dashboard')) {
    items.push({ key: '/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' })
  }
  if (hasPerm('data-produksi')) {
    items.push({ key: '/data', icon: <DatabaseOutlined />, label: 'Data Produksi' })
  }
  if (hasPerm('export')) {
    items.push({ key: '/exports', icon: <ExportOutlined />, label: 'Export' })
  }
  if (hasPerm('ai-chat')) {
    items.push({
      key: 'analytics', icon: <RobotOutlined />, label: 'AI Chat',
      children: [
        { key: '/analytics/ai-chat', icon: <RobotOutlined />, label: 'Chat Assistant' },
      ],
    })
  }

  const adminChildren: { key: string; icon: any; label: string; perm: string }[] = [
    { key: '/admin/users', icon: <UserOutlined />, label: 'User Management', perm: 'user' },
    { key: '/admin/roles', icon: <SafetyOutlined />, label: 'Role Management', perm: 'role' },
        { key: '/admin/role-permission', icon: <SafetyOutlined />, label: 'Role Permission', perm: 'role' },
        { key: '/admin/resource-connection', icon: <DatabaseOutlined />, label: 'Resource Connection', perm: 'resource-connection' },
        { key: '/admin/data-produksi-config', icon: <DatabaseOutlined />, label: 'Data Produksi Config', perm: 'data-produksi-config' },
    { key: '/admin/activity-logs', icon: <HistoryOutlined />, label: 'Activity Log', perm: 'activity-log' },
  ]

  const visibleAdmin = adminChildren.filter((c) => hasPerm(c.perm))
  if (visibleAdmin.length > 0) {
    items.push({
      key: 'admin', icon: <SafetyOutlined />, label: 'Administrasi',
      children: visibleAdmin.map((c) => ({ key: c.key, icon: c.icon, label: c.label })),
    })
  }

  return items
}

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { token, user, logout } = useAuthStore()
  const userLevel = useUserLevel()
  const {
    darkMode, collapsed, colorPreset,
    toggleTheme, toggleCollapsed, setColorPreset, toggleBreadcrumb,
  } = useThemeStore()
  const breadcrumbEnabled = useBreadcrumbEnabled()
  const { token: antdToken } = theme.useToken()
  const { notifications, addNotification, markAllRead, unreadCount } = useNotificationStore()
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768)

  useEffect(() => {
    const handleResize = () => setIsMobile(window.innerWidth < 768)
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  useEffect(() => {
    const t = setTimeout(() => {
      api.get<any>('/api/v1/auth/me').then((res) => {
        if (res?.data) useAuthStore.getState().updateUser(res.data)
      }).catch(() => {})
    }, 100)
    return () => clearTimeout(t)
  }, [location.pathname])

  if (!token) return <Navigate to="/login" replace />

  const wsHandlers: Record<string, (payload: any) => void> = {
    notification: useCallback((payload: any) => {
      addNotification({ type: payload.type, title: payload.title, message: payload.message })
    }, [addNotification]),
    pong: () => {},
  }

  useWebSocket(wsHandlers)

  const handleLogout = () => { logout(); navigate('/login') }

  const userMenu = {
    items: [
      {
        key: 'profile', icon: <UserOutlined />, label: (
          <div>
            <div style={{ fontWeight: 600 }}>{user?.user_name || 'User'}</div>
            <div style={{ fontSize: 12, color: 'var(--text-tertiary)' }}>{user?.email || ''}</div>
          </div>
        ),
      },
      { type: 'divider' as const },
      { key: 'logout', icon: <LogoutOutlined />, label: 'Logout', danger: true },
    ],
    onClick: ({ key }: { key: string }) => { if (key === 'logout') handleLogout() },
  }

  const selectedKey = '/' + location.pathname.split('/').slice(1, 3).join('/')
  const menuItems = useMemo(() => getMenuItems(userLevel, user?.permissions), [userLevel, user?.permissions])

  const hasAnyPerm = userLevel === 'admin' || (user?.permissions && user.permissions.length > 0)
  const routePermMap: Record<string, string> = {
    dashboard: 'dashboard',
    data: 'data-produksi',
    exports: 'export',
    analytics: 'ai-chat',
    users: 'user',
    roles: 'role',
    'role-permission': 'role',
    'resource-connection': 'resource-connection',
    'data-produksi-config': 'data-produksi-config',
    'activity-logs': 'activity-log',
  }
  const pathSegment = location.pathname.split('/').filter(Boolean)[0] || ''
  const routeModule = routePermMap[pathSegment === 'admin' ? location.pathname.split('/')[2] || '' : pathSegment]

  const breadcrumbItems = useMemo(() => {
    const parts = location.pathname.split('/').filter(Boolean)
    const crumbs: { key: string; title: React.ReactNode }[] = [
      { key: '/', title: <span onClick={() => navigate('/dashboard')} style={{ cursor: 'pointer' }}>Home</span> },
    ]
    let path = ''
    for (const part of parts) {
      path += `/${part}`
      const label = labelMap[part] || part.charAt(0).toUpperCase() + part.slice(1).replace(/-/g, ' ')
      crumbs.push({ key: path, title: <span onClick={() => navigate(path)} style={{ cursor: 'pointer' }}>{label}</span> })
    }
    return crumbs
  }, [location.pathname, navigate])

  const notificationContent = (
    <div style={{ width: 380, maxHeight: 420, overflow: 'auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '12px 16px', borderBottom: `1px solid ${antdToken.colorBorderSecondary}` }}>
        <Text strong style={{ fontSize: 15 }}>Notifikasi</Text>
        <Button type="link" size="small" onClick={markAllRead} style={{ fontSize: 12 }}>Tandai Dibaca</Button>
      </div>
      {notifications.length === 0 ? (
        <div style={{ padding: 32, textAlign: 'center', color: 'var(--text-tertiary)', fontSize: 13 }}>Tidak ada notifikasi</div>
      ) : (
        <List
          dataSource={notifications.slice(0, 20)}
          renderItem={(item) => (
            <List.Item
              style={{
                padding: '10px 16px',
                background: item.read ? 'transparent' : 'var(--primary-bg)',
                cursor: 'pointer',
                transition: 'background 0.2s',
                borderLeft: item.read ? 'none' : `3px solid var(--primary-color)`,
              }}
              onClick={() => useNotificationStore.getState().markRead(item.id)}
            >
              <List.Item.Meta
                title={<Text strong={!item.read} style={{ fontSize: 13 }}>{item.title}</Text>}
                description={<Text type="secondary" style={{ fontSize: 12 }}>{item.message}</Text>}
              />
            </List.Item>
          )}
        />
      )}
    </div>
  )

  const siderWidth = collapsed ? 68 : 260

  const sidebarContent = (
    <>
      <div style={{
        height: 60,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontWeight: 700,
        fontSize: collapsed ? 18 : 22,
        color: 'var(--primary-color)',
        borderBottom: '1px solid var(--sidebar-border)',
        letterSpacing: collapsed ? 0 : 1.5,
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        transition: 'all 0.3s',
        padding: collapsed ? '0 10px' : '0 20px',
      }}>
        <div style={{
          width: 34,
          height: 34,
          borderRadius: 10,
          background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: '#fff',
          fontSize: 16,
          fontWeight: 800,
          marginRight: collapsed && !isMobile ? 0 : 12,
          flexShrink: 0,
          boxShadow: `0 4px 12px ${darkMode ? 'rgba(22,119,255,0.3)' : 'rgba(22,119,255,0.25)'}`,
        }}>
          VM
        </div>
        {(collapsed && !isMobile) ? null : (
          <span style={{
            background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            backgroundClip: 'text',
          }}>
            Visual Mesin
          </span>
        )}
      </div>
      <div style={{ padding: '8px 0', overflow: 'auto', height: 'calc(100vh - 60px)', '::-webkit-scrollbar': { width: 0 } } as React.CSSProperties}>
        <Menu
          theme={darkMode ? 'dark' : 'light'}
          mode="inline"
          selectedKeys={[selectedKey]}
          defaultOpenKeys={['admin', 'analytics']}
          items={menuItems}
          onClick={({ key }) => {
            navigate(key)
            if (isMobile) setMobileMenuOpen(false)
          }}
          style={{
            background: 'transparent',
            borderInlineEnd: 'none',
          }}
        />
      </div>
    </>
  )

  return (
    <Layout style={{ minHeight: '100vh', background: 'var(--bg-layout)' }}>
      {!isMobile && (
        <Sider
          trigger={null}
          collapsible
          collapsed={collapsed}
          width={260}
          theme={darkMode ? 'dark' : 'light'}
          style={{
            background: 'var(--sidebar-bg)',
            borderRight: '1px solid var(--sidebar-border)',
            position: 'fixed',
            left: 0,
            top: 0,
            bottom: 0,
            zIndex: 100,
            transition: 'all 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
            overflow: 'hidden',
          }}
        >
          {sidebarContent}
        </Sider>
      )}

      <Drawer
        placement="left"
        width={280}
        open={isMobile && mobileMenuOpen}
        onClose={() => setMobileMenuOpen(false)}
        styles={{ body: { padding: 0, background: 'var(--sidebar-bg)' } }}
        closable={false}
      >
        <div style={{ background: 'var(--sidebar-bg)', height: '100%' }}>
          <div style={{
            height: 60,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '0 20px',
            borderBottom: '1px solid var(--sidebar-border)',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <div style={{
                width: 34, height: 34, borderRadius: 10,
                background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                color: '#fff', fontSize: 16, fontWeight: 800,
                boxShadow: `0 4px 12px ${darkMode ? 'rgba(22,119,255,0.3)' : 'rgba(22,119,255,0.25)'}`,
              }}>
                VM
              </div>
              <span style={{
                fontWeight: 700, fontSize: 20,
                background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                backgroundClip: 'text',
              }}>
                Visual Mesin
              </span>
            </div>
            <Button type="text" icon={<MenuOutlined />} onClick={() => setMobileMenuOpen(false)} />
          </div>
          <Menu
            theme={darkMode ? 'dark' : 'light'}
            mode="inline"
            selectedKeys={[selectedKey]}
            defaultOpenKeys={['admin', 'analytics']}
            items={menuItems}
            onClick={({ key }) => {
              navigate(key)
              setMobileMenuOpen(false)
            }}
            style={{ background: 'transparent', borderInlineEnd: 'none' }}
          />
        </div>
      </Drawer>

      <Layout style={{
        marginLeft: isMobile ? 0 : siderWidth,
        transition: 'margin-left 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
        background: 'var(--bg-layout)',
        minHeight: '100vh',
      }}>
        <Header style={{
          padding: isMobile ? '0 12px' : '0 24px',
          background: 'var(--header-bg)',
          backdropFilter: 'blur(16px)',
          WebkitBackdropFilter: 'blur(16px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: '1px solid var(--header-border)',
          height: 60,
          position: 'sticky',
          top: 0,
          zIndex: 99,
          transition: 'all 0.3s ease',
        }}>
          <Space>
            {isMobile ? (
              <Button
                type="text"
                icon={<MenuOutlined />}
                onClick={() => setMobileMenuOpen(true)}
                style={{
                  fontSize: 20,
                  color: 'var(--text-secondary)',
                  width: 36,
                  height: 36,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              />
            ) : (
              <Button
                type="text"
                icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                onClick={toggleCollapsed}
                style={{
                  fontSize: 18,
                  color: 'var(--text-secondary)',
                  width: 36,
                  height: 36,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              />
            )}
            {breadcrumbEnabled && !isMobile && (
              <Breadcrumb
                items={breadcrumbItems}
                style={{ marginLeft: 8 }}
              />
            )}
          </Space>
          <Space size={isMobile ? 0 : 4}>
            <Tooltip title="Activity Log">
              <Button
                type="text"
                icon={<HistoryOutlined />}
                onClick={() => navigate('/admin/activity-logs')}
                style={{
                  color: 'var(--text-secondary)',
                  width: 36,
                  height: 36,
                  display: isMobile ? 'none' : 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              />
            </Tooltip>
            <Popover content={notificationContent} trigger="click" placement="bottomRight">
              <Badge count={unreadCount()} size="small" offset={[-2, 2]}>
                <Button
                  type="text"
                  icon={<BellOutlined />}
                  style={{
                    color: 'var(--text-secondary)',
                    width: 36,
                    height: 36,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                />
              </Badge>
            </Popover>
            {!isMobile && (
              <Tooltip title={darkMode ? 'Mode Terang' : 'Mode Gelap'}>
                <Button
                  type="text"
                  icon={darkMode ? <SunOutlined /> : <MoonOutlined />}
                  onClick={toggleTheme}
                  style={{
                    color: 'var(--text-secondary)',
                    fontSize: 18,
                    width: 36,
                    height: 36,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                />
              </Tooltip>
            )}
            {!isMobile && (
              <Tooltip title="Pengaturan Tampilan">
                <Button
                  type="text"
                  icon={<BgColorsOutlined />}
                  onClick={() => setSettingsOpen(true)}
                  style={{
                    color: 'var(--text-secondary)',
                    fontSize: 18,
                    width: 36,
                    height: 36,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                />
              </Tooltip>
            )}
            <Dropdown menu={userMenu} placement="bottomRight" overlayStyle={{ minWidth: 200 }}>
              <div style={{
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: 10,
                padding: '6px 12px 6px 8px',
                borderRadius: 10,
                marginLeft: 8,
                transition: 'all 0.2s',
              }}
                className="user-dropdown-trigger"
              >
                <Avatar
                  icon={<UserOutlined />}
                  style={{
                    background: `linear-gradient(135deg, var(--primary-color), ${darkMode ? '#4096ff' : '#0958d9'})`,
                    boxShadow: '0 2px 8px rgba(22,119,255,0.3)',
                  }}
                  size={isMobile ? 28 : 30}
                />
                {!isMobile && (
                  <span style={{ color: 'var(--text-primary)', fontSize: 14, fontWeight: 500, lineHeight: 1 }}>
                    {user?.user_name || 'User'}
                  </span>
                )}
              </div>
            </Dropdown>
          </Space>
        </Header>
        <Content style={{ margin: isMobile ? 12 : 24, minHeight: 280 }}>
          <div className="page-enter">
            {!hasAnyPerm ? (
              <div style={{
                display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
                minHeight: 400, textAlign: 'center',
              }}>
                <div style={{
                  width: 80, height: 80, borderRadius: 20,
                  background: 'var(--primary-bg)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontSize: 36, color: 'var(--primary-color)', marginBottom: 24,
                }}>
                  <SafetyOutlined />
                </div>
                <h3 style={{ margin: '0 0 8px' }}>Akses Tidak Tersedia</h3>
                <p style={{ color: 'var(--text-secondary)', maxWidth: 400, margin: '0 0 24px', lineHeight: 1.6 }}>
                  Anda tidak memiliki permission untuk mengakses sistem ini.
                  Silakan hubungi Administrator untuk mendapatkan akses.
                </p>
                <Button type="primary" danger icon={<LogoutOutlined />} onClick={handleLogout}>
                  Logout
                </Button>
              </div>
            ) : routeModule && userLevel !== 'admin' && !canModule(routeModule, user) ? (
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
                <Result
                  status="403"
                  title={<span style={{ fontSize: 64, fontWeight: 800, color: 'var(--primary-color)', lineHeight: 1 }}>403</span>}
                  subTitle="Anda tidak memiliki akses ke halaman ini"
                  extra={
                    <Button type="primary" onClick={() => navigate('/dashboard')} className="btn-glow" style={{ borderRadius: 10, height: 44, paddingInline: 28 }}>
                      Kembali ke Dashboard
                    </Button>
                  }
                />
              </div>
            ) : (
              <Outlet />
            )}
          </div>
        </Content>
      </Layout>

      <Drawer
        title={
          <Space>
            <SettingOutlined style={{ fontSize: 18, color: 'var(--primary-color)' }} />
            <span style={{ fontWeight: 600 }}>Pengaturan Tampilan</span>
          </Space>
        }
        placement="right"
        onClose={() => setSettingsOpen(false)}
        open={settingsOpen}
        width={340}
        styles={{ body: { padding: 24 } }}
      >
        <Divider orientation="left" plain style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)' }}>
          Mode Tampilan
        </Divider>
        <div style={{ display: 'flex', gap: 12, marginBottom: 24 }}>
          <div
            onClick={() => { if (darkMode) toggleTheme() }}
            style={{
              flex: 1, padding: 16, borderRadius: 10, cursor: 'pointer', textAlign: 'center',
              border: !darkMode ? '2px solid var(--primary-color)' : '2px solid var(--border-color-strong)',
              background: !darkMode ? 'var(--primary-bg)' : 'transparent',
              transition: 'all 0.2s',
            }}
          >
            <SunOutlined style={{ fontSize: 24, color: !darkMode ? 'var(--primary-color)' : 'var(--text-secondary)' }} />
            <div style={{ marginTop: 8, fontWeight: 500, fontSize: 13 }}>Terang</div>
          </div>
          <div
            onClick={() => { if (!darkMode) toggleTheme() }}
            style={{
              flex: 1, padding: 16, borderRadius: 10, cursor: 'pointer', textAlign: 'center',
              border: darkMode ? '2px solid var(--primary-color)' : '2px solid var(--border-color-strong)',
              background: darkMode ? 'var(--primary-bg)' : 'transparent',
              transition: 'all 0.2s',
            }}
          >
            <MoonOutlined style={{ fontSize: 24, color: darkMode ? 'var(--primary-color)' : 'var(--text-secondary)' }} />
            <div style={{ marginTop: 8, fontWeight: 500, fontSize: 13 }}>Gelap</div>
          </div>
        </div>

        <Divider orientation="left" plain style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)' }}>
          Warna Tema
        </Divider>
        <div style={{ display: 'flex', gap: 14, marginBottom: 24, flexWrap: 'wrap' }}>
          {colorPresetOptions.map((opt) => (
            <Tooltip key={opt.value} title={opt.label}>
              <div
                onClick={() => setColorPreset(opt.value)}
                style={{
                  width: 40, height: 40, borderRadius: '50%', cursor: 'pointer',
                  background: `linear-gradient(135deg, ${opt.color}, ${opt.color}dd)`,
                  border: colorPreset === opt.value ? '3px solid var(--text-primary)' : '3px solid transparent',
                  transition: 'all 0.2s',
                  boxShadow: colorPreset === opt.value
                    ? `0 0 0 3px ${opt.color}30, 0 4px 12px ${opt.color}40`
                    : '0 2px 8px rgba(0,0,0,0.06)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                {colorPreset === opt.value && (
                  <div style={{
                    width: 10, height: 10, borderRadius: '50%',
                    background: '#fff',
                  }} />
                )}
              </div>
            </Tooltip>
          ))}
        </div>

        <Divider orientation="left" plain style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-secondary)' }}>
          Pengaturan Lain
        </Divider>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: 14, color: 'var(--text-primary)' }}>Tampilkan Breadcrumb</span>
            <Switch checked={breadcrumbEnabled} onChange={toggleBreadcrumb} />
          </div>
        </Space>
      </Drawer>
    </Layout>
  )
}
