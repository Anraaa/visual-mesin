import { Card, Col, Row, Typography, Space } from 'antd'
import {
  CheckCircleOutlined, CloseCircleOutlined, DashboardOutlined,
  RiseOutlined, SettingOutlined, ToolOutlined, SwapRightOutlined,
} from '@ant-design/icons'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'

const { Title, Text } = Typography

const productionData = [
  { name: 'RTBA1', OK: 180, NG: 12 },
  { name: 'RTBA2', OK: 165, NG: 8 },
  { name: 'RTBA3', OK: 195, NG: 15 },
  { name: 'Extruder', OK: 220, NG: 5 },
  { name: 'Curing', OK: 190, NG: 10 },
  { name: 'Trimming', OK: 210, NG: 7 },
]

const pieData = [
  { name: 'OK', value: 85, color: '#52c41a' },
  { name: 'NG', value: 15, color: '#ff4d4f' },
]

const statCards = [
  { title: 'Total Produksi', value: '1.250', icon: <DashboardOutlined />, gradient: 'gradient-blue', suffix: 'pcs' },
  { title: 'OK', value: '1.180', icon: <CheckCircleOutlined />, gradient: 'gradient-green', suffix: 'pcs' },
  { title: 'NG', value: '70', icon: <CloseCircleOutlined />, gradient: 'gradient-red', suffix: 'pcs' },
  { title: 'Yield Rate', value: '94.4', icon: <RiseOutlined />, gradient: 'gradient-purple', suffix: '%' },
]

const machineStats = [
  { name: 'Mesin Aktif', value: 24, icon: <SettingOutlined />, gradient: 'gradient-cyan' },
  { name: 'Dalam Perawatan', value: 3, icon: <ToolOutlined />, gradient: 'gradient-orange' },
]

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div style={{
        background: 'var(--bg-elevated)',
        padding: '14px 18px',
        borderRadius: 12,
        boxShadow: 'var(--shadow-lg)',
        border: '1px solid var(--border-color)',
      }}>
        <Text strong style={{ fontSize: 13, display: 'block', marginBottom: 6 }}>{label}</Text>
        {payload.map((p: any, i: number) => (
          <div key={i} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 13 }}>
            <span style={{ width: 8, height: 8, borderRadius: '50%', background: p.color }} />
            <span style={{ color: 'var(--text-secondary)' }}>{p.name}:</span>
            <span style={{ fontWeight: 700 }}>{p.value}</span>
          </div>
        ))}
      </div>
    )
  }
  return null
}

export default function Dashboard() {
  return (
    <div className="page-enter">
      {/* Page Header */}
      <div className="page-header">
        <div>
          <h4>Dashboard Produksi</h4>
          <p className="page-subtitle">Ringkasan produksi hari ini</p>
        </div>
        <Space>
          <span style={{ fontSize: 13, color: 'var(--text-tertiary)' }}>
            <SwapRightOutlined style={{ marginRight: 4 }} />
            Update terbaru
          </span>
        </Space>
      </div>

      {/* Stat Cards */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {statCards.map((card, i) => (
          <Col xs={24} sm={12} lg={6} key={i}>
            <Card className={`stat-card ${card.gradient}`} bordered={false}>
              <div className="stat-bg-icon">{card.icon}</div>
              <div style={{ position: 'relative', zIndex: 1 }}>
                <div className="stat-value">
                  {card.value}
                  {card.suffix && <span className="stat-suffix">{card.suffix}</span>}
                </div>
                <div className="stat-label">{card.title}</div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      {/* Machine Stats */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {machineStats.map((card, i) => (
          <Col xs={24} sm={12} lg={6} key={i}>
            <Card className={`stat-card ${card.gradient}`} bordered={false}>
              <div className="stat-bg-icon">{card.icon}</div>
              <div style={{ position: 'relative', zIndex: 1 }}>
                <div className="stat-value">{card.value}</div>
                <div className="stat-label">{card.name}</div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      {/* Charts */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card
            className="modern-card"
            title={
              <Space>
                <div style={{
                  width: 8, height: 8, borderRadius: '50%',
                  background: 'var(--primary-color)',
                }} />
                <span style={{ fontWeight: 600, fontSize: 15 }}>Produksi per Mesin</span>
              </Space>
            }
            styles={{ body: { padding: '20px 16px 8px' } }}
          >
            <ResponsiveContainer width="100%" height={340}>
              <BarChart data={productionData} margin={{ top: 8, right: 16, bottom: 0, left: -16 }}>
                <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                <XAxis
                  dataKey="name"
                  tick={{ fill: 'var(--text-secondary)', fontSize: 12 }}
                  axisLine={false}
                  tickLine={false}
                />
                <YAxis
                  tick={{ fill: 'var(--text-secondary)', fontSize: 12 }}
                  axisLine={false}
                  tickLine={false}
                />
                <Tooltip content={<CustomTooltip />} cursor={{ fill: 'var(--bg-hover)' }} />
                <Bar dataKey="OK" fill="#52c41a" radius={[6, 6, 0, 0]} maxBarSize={36} />
                <Bar dataKey="NG" fill="#ff4d4f" radius={[6, 6, 0, 0]} maxBarSize={36} />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card
            className="modern-card"
            title={
              <Space>
                <div style={{
                  width: 8, height: 8, borderRadius: '50%',
                  background: 'var(--primary-color)',
                }} />
                <span style={{ fontWeight: 600, fontSize: 15 }}>Rasio OK / NG</span>
              </Space>
            }
            styles={{ body: { padding: '20px' } }}
          >
            <ResponsiveContainer width="100%" height={340}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={75}
                  outerRadius={120}
                  dataKey="value"
                  paddingAngle={3}
                  startAngle={90}
                  endAngle={-270}
                >
                  {pieData.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={entry.color}
                      stroke="none"
                      style={{ filter: `drop-shadow(0 2px 8px ${entry.color}40)` }}
                    />
                  ))}
                </Pie>
                <Tooltip content={<CustomTooltip />} />
                <text
                  x="50%" y="48%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{ fontSize: 32, fontWeight: 800, fill: 'var(--text-primary)' }}
                >
                  85%
                </text>
                <text
                  x="50%" y="62%"
                  textAnchor="middle"
                  dominantBaseline="middle"
                  style={{ fontSize: 13, fill: 'var(--text-tertiary)', fontWeight: 500 }}
                >
                  Yield Rate
                </text>
              </PieChart>
            </ResponsiveContainer>
            <div style={{ display: 'flex', justifyContent: 'center', gap: 32, marginTop: 8 }}>
              {pieData.map((item) => (
                <div key={item.name} style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <span style={{
                    width: 10, height: 10, borderRadius: '50%',
                    background: item.color,
                    boxShadow: `0 2px 6px ${item.color}50`,
                  }} />
                  <span style={{ fontSize: 13, color: 'var(--text-secondary)' }}>{item.name}</span>
                  <span style={{ fontSize: 13, fontWeight: 700 }}>{item.value}%</span>
                </div>
              ))}
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
