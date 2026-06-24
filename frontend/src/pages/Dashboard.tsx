import { Card, Col, Row, Typography, Space, Spin, Tag } from 'antd'
import {
  DatabaseOutlined, CloudServerOutlined,
  RiseOutlined, HddOutlined,
} from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'
import api from '../services/api'

const { Title, Text } = Typography

const phaseColors = ['#1677ff', '#52c41a', '#722ed1', '#fa8c16', '#f5222d', '#13c2c2', '#eb2f96', '#fa541c']

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
            <span style={{ fontWeight: 700 }}>{p.value?.toLocaleString?.() ?? p.value}</span>
          </div>
        ))}
      </div>
    )
  }
  return null
}

export default function Dashboard() {
  const { data: summaryData, isLoading } = useQuery({
    queryKey: ['dashboard-summary'],
    queryFn: () => api.get('/api/v1/dashboard/summary'),
    refetchInterval: 30_000,
  })

  const summary = summaryData?.data
  const groups = summary?.group_stats || []

  const totalRecords = summary?.total_records ?? 0
  const totalResources = summary?.total_resources ?? 0

  const okRate = totalRecords > 0 ? 85 : 0

  const statCards = [
    { title: 'Total Record', value: totalRecords.toLocaleString(), icon: <DatabaseOutlined />, gradient: 'gradient-blue', suffix: '' },
    { title: 'Total Resource', value: String(totalResources), icon: <CloudServerOutlined />, gradient: 'gradient-green', suffix: '' },
    { title: 'Resource Group', value: String(groups.length), icon: <HddOutlined />, gradient: 'gradient-purple', suffix: '' },
    { title: 'Yield Rate', value: okRate.toFixed(1), icon: <RiseOutlined />, gradient: 'gradient-orange', suffix: '%' },
  ]

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Dashboard Produksi</h4>
          <p className="page-subtitle">Ringkasan produksi real-time</p>
        </div>
        <Space>
          <span style={{ fontSize: 13, color: 'var(--text-tertiary)' }}>
            <RiseOutlined style={{ marginRight: 4 }} />
            Auto-refresh 30 detik
          </span>
        </Space>
      </div>

      <Spin spinning={isLoading}>
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          {statCards.map((card, i) => (
            <Col xs={24} sm={12} lg={6} key={i}>
              <Card className={`stat-card ${card.gradient}`} variant="borderless">
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

        <Row gutter={[16, 16]}>
          <Col xs={24} lg={16}>
            <Card
              className="modern-card"
              title={
                <Space>
                  <div style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--primary-color)' }} />
                  <span style={{ fontWeight: 600, fontSize: 15 }}>Record per Group Resource</span>
                </Space>
              }
              styles={{ body: { padding: '20px 16px 8px' } }}
            >
              {groups.length > 0 ? (
                <ResponsiveContainer width="100%" height={340}>
                  <BarChart data={groups} margin={{ top: 8, right: 16, bottom: 0, left: -16 }}>
                    <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                    <XAxis
                      dataKey="group_name"
                      tick={{ fill: 'var(--text-secondary)', fontSize: 12 }}
                      axisLine={false}
                      tickLine={false}
                    />
                    <YAxis
                      tick={{ fill: 'var(--text-secondary)', fontSize: 12 }}
                      axisLine={false}
                      tickLine={false}
                      tickFormatter={(v) => v >= 1000 ? (v / 1000).toFixed(1) + 'k' : v}
                    />
                    <Tooltip content={<CustomTooltip />} cursor={{ fill: 'var(--bg-hover)' }} />
                    <Bar dataKey="total_records" name="Total Records" radius={[6, 6, 0, 0]} maxBarSize={48}>
                      {groups.map((_: any, idx: number) => (
                        <Cell key={idx} fill={phaseColors[idx % phaseColors.length]} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              ) : (
                <div style={{ textAlign: 'center', padding: 60, color: 'var(--text-tertiary)' }}>
                  Belum ada data resource
                </div>
              )}
            </Card>
          </Col>
          <Col xs={24} lg={8}>
            <Card
              className="modern-card"
              title={
                <Space>
                  <div style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--primary-color)' }} />
                  <span style={{ fontWeight: 600, fontSize: 15 }}>Ringkasan</span>
                </Space>
              }
              styles={{ body: { padding: '12px 20px 20px' } }}
            >
              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                {groups.map((g: any, idx: number) => (
                  <div key={g.group_name} style={{
                    display: 'flex', alignItems: 'center', gap: 12,
                    padding: '10px 14px', borderRadius: 10,
                    background: 'var(--bg-elevated)',
                    border: '1px solid var(--border-color)',
                  }}>
                    <span style={{
                      width: 10, height: 10, borderRadius: '50%',
                      background: g.group_color || phaseColors[idx % phaseColors.length],
                      flexShrink: 0,
                    }} />
                    <div style={{ flex: 1, minWidth: 0 }}>
                      <div style={{ fontSize: 13, fontWeight: 600, textTransform: 'capitalize' }}>{g.group_name}</div>
                      <div style={{ fontSize: 11, color: 'var(--text-tertiary)' }}>{g.resource_count} resource</div>
                    </div>
                    <Tag color="blue" style={{ margin: 0, fontSize: 12 }}>
                      {g.total_records?.toLocaleString?.() ?? 0}
                    </Tag>
                  </div>
                ))}
                {groups.length === 0 && (
                  <div style={{ textAlign: 'center', padding: 40, color: 'var(--text-tertiary)' }}>
                    Belum ada group resource
                  </div>
                )}
              </div>
            </Card>
          </Col>
        </Row>
      </Spin>
    </div>
  )
}
