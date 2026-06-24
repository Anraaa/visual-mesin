import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, Form, Input, Button, Typography, message } from 'antd'
import { UserOutlined, LockOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import api from '../services/api'

const { Title } = Typography

export default function Login() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()
  const { darkMode } = useThemeStore()

  const onFinish = async (values: { email: string; password: string }) => {
    setLoading(true)
    try {
      const res = await api.post<{ token: string; user: any }>('/api/v1/auth/login', values)
      const { token, user } = res.data
      setAuth(token, user)
      message.success('Login berhasil')
      navigate('/dashboard')
    } catch {
      message.error('Login gagal. Periksa email/NIP dan password.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: darkMode
        ? 'linear-gradient(135deg, #0a0a0f 0%, #12121a 30%, #1a1a2e 60%, #0f0f1a 100%)'
        : 'linear-gradient(135deg, #667eea 0%, #764ba2 50%, #8b5cf6 100%)',
      position: 'relative',
      overflow: 'hidden',
    }}>
      {/* Animated background orbs */}
      <div style={{
        position: 'absolute',
        width: 500,
        height: 500,
        borderRadius: '50%',
        background: darkMode
          ? 'radial-gradient(circle, rgba(22,119,255,0.08) 0%, transparent 70%)'
          : 'radial-gradient(circle, rgba(255,255,255,0.12) 0%, transparent 70%)',
        top: -150,
        right: -100,
        animation: 'login-float 12s ease-in-out infinite',
        pointerEvents: 'none',
      }} />
      <div style={{
        position: 'absolute',
        width: 400,
        height: 400,
        borderRadius: '50%',
        background: darkMode
          ? 'radial-gradient(circle, rgba(114,46,209,0.08) 0%, transparent 70%)'
          : 'radial-gradient(circle, rgba(255,255,255,0.08) 0%, transparent 70%)',
        bottom: -100,
        left: -80,
        animation: 'login-float2 15s ease-in-out infinite',
        pointerEvents: 'none',
      }} />
      <div style={{
        position: 'absolute',
        width: 300,
        height: 300,
        borderRadius: '50%',
        background: darkMode
          ? 'radial-gradient(circle, rgba(19,194,194,0.05) 0%, transparent 70%)'
          : 'radial-gradient(circle, rgba(255,255,255,0.06) 0%, transparent 70%)',
        top: '40%',
        left: '60%',
        animation: 'login-float 18s ease-in-out infinite reverse',
        pointerEvents: 'none',
      }} />

      {/* Subtle grid pattern overlay */}
      <div style={{
        position: 'absolute',
        inset: 0,
        opacity: darkMode ? 0.03 : 0.04,
        backgroundImage: 'radial-gradient(circle, #fff 1px, transparent 1px)',
        backgroundSize: '40px 40px',
        pointerEvents: 'none',
      }} />

      <div style={{
        animation: 'scaleIn 0.5s cubic-bezier(0.16, 1, 0.3, 1)',
      }}>
        <Card
          className="glass-card"
          style={{
            width: '100%',
            maxWidth: 420,
            padding: 0,
            position: 'relative',
            overflow: 'hidden',
          }}
        >
          {/* Top accent bar */}
          <div style={{
            height: 4,
            background: 'linear-gradient(90deg, var(--primary-color), var(--primary-hover), #8b5cf6, var(--primary-color))',
            backgroundSize: '300% 100%',
            animation: 'gradient-shift 4s ease infinite',
          }} />

          <div style={{ padding: '32px 24px 28px' }}>
            <div style={{ textAlign: 'center', marginBottom: 28 }}>
              <div style={{
                width: 56,
                height: 56,
                borderRadius: 14,
                background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
                marginBottom: 16,
                fontSize: 24,
                fontWeight: 800,
                color: '#fff',
                boxShadow: `0 8px 24px ${darkMode ? 'rgba(22,119,255,0.3)' : 'rgba(22,119,255,0.25)'}`,
              }}>
                VM
              </div>
              <Title level={4} style={{ margin: 0, marginBottom: 4, fontWeight: 700 }}>
                Visual Mesin
              </Title>
              <Typography.Text type="secondary" style={{ fontSize: 13, color: 'var(--text-tertiary)' }}>
                Silakan login untuk melanjutkan
              </Typography.Text>
            </div>

            <Form layout="vertical" onFinish={onFinish} autoComplete="off" size="large" style={{ gap: 0 }}>
              <Form.Item
                name="email"
                rules={[{ required: true, message: 'Masukkan email atau NIP' }]}
                style={{ marginBottom: 16 }}
              >
                <Input
                  prefix={<UserOutlined style={{ color: 'var(--text-tertiary)' }} />}
                  placeholder="Email / NIP"
                  variant="filled"
                  style={{ height: 44, borderRadius: 10, border: 'none', background: darkMode ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)' }}
                />
              </Form.Item>
              <Form.Item
                name="password"
                rules={[{ required: true, message: 'Masukkan password' }]}
                style={{ marginBottom: 20 }}
              >
                <Input.Password
                  prefix={<LockOutlined style={{ color: 'var(--text-tertiary)' }} />}
                  placeholder="Password"
                  iconRender={(visible) => visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />}
                  variant="filled"
                  style={{ height: 44, borderRadius: 10, border: 'none', background: darkMode ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)' }}
                />
              </Form.Item>
              <Form.Item style={{ marginBottom: 0 }}>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={loading}
                  block
                  size="large"
                  className="btn-glow"
                  style={{
                    height: 44,
                    borderRadius: 10,
                    fontSize: 15,
                    fontWeight: 600,
                    letterSpacing: 0.5,
                  }}
                >
                  Login
                </Button>
              </Form.Item>
            </Form>
          </div>

          {/* Bottom credential info */}
          <div style={{
            padding: '12px 24px',
            borderTop: '1px solid var(--border-color)',
            textAlign: 'center',
          }}>
            <Typography.Text style={{ color: 'var(--text-tertiary)', fontSize: 11, display: 'block', marginBottom: 4 }}>
              Demo Credentials
            </Typography.Text>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4, fontSize: 11 }}>
              <div>
                <span style={{ color: 'var(--text-tertiary)' }}>Admin: </span>
                <span style={{ color: 'var(--text-secondary)', fontWeight: 500 }}>admin@admin.com / password</span>
              </div>
              <div>
                <span style={{ color: 'var(--text-tertiary)' }}>User: </span>
                <span style={{ color: 'var(--text-secondary)', fontWeight: 500 }}>user@visualmesin.com / user123</span>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  )
}
