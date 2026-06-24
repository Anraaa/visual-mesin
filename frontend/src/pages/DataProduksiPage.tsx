import { useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Card, Col, Row, Typography, Spin } from 'antd'
import {
  BuildOutlined, ExperimentOutlined, FireOutlined, ScissorOutlined,
  ToolOutlined, DashboardOutlined, AlertOutlined, FileTextOutlined,
  ContainerOutlined, ControlOutlined, BookOutlined, ThunderboltOutlined,
  BarcodeOutlined, RightOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Text } = Typography

const iconMap: Record<string, any> = {
  BuildOutlined: <BuildOutlined />,
  ExperimentOutlined: <ExperimentOutlined />,
  FireOutlined: <FireOutlined />,
  ScissorOutlined: <ScissorOutlined />,
  ToolOutlined: <ToolOutlined />,
  DashboardOutlined: <DashboardOutlined />,
  AlertOutlined: <AlertOutlined />,
  FileTextOutlined: <FileTextOutlined />,
  ContainerOutlined: <ContainerOutlined />,
  ControlOutlined: <ControlOutlined />,
  BookOutlined: <BookOutlined />,
  ThunderboltOutlined: <ThunderboltOutlined />,
  BarcodeOutlined: <BarcodeOutlined />,
}

export default function DataProduksiPage() {
  const navigate = useNavigate()

  const { data: groupsData, isLoading } = useQuery({
    queryKey: ['resource-groups'],
    queryFn: () => api.get('/api/v1/data-produksi-config/groups'),
  })

  const groups = groupsData?.data || []

  if (isLoading) {
    return <div style={{ textAlign: 'center', padding: 80 }}><Spin size="large" /></div>
  }

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Data Produksi</h4>
          <p className="page-subtitle">Pilih tabel resource untuk melihat data produksi</p>
        </div>
      </div>

      <Row gutter={[16, 16]}>
        {groups.map((group: any) => {
          const iconEl = iconMap[group.icon] || <BuildOutlined />
          return (
            <Col xs={24} sm={12} lg={8} xl={6} key={group.id}>
              <Card
                className="resource-card"
                style={{
                  borderLeft: `3px solid ${group.color}`,
                  background: `linear-gradient(135deg, ${group.color}15, ${group.color}05)`,
                  height: '100%',
                }}
                bodyStyle={{ padding: 16 }}
              >
                <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: 14 }}>
                  <div style={{
                    width: 40, height: 40, borderRadius: 10,
                    background: `${group.color}18`,
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    fontSize: 20, color: group.color,
                  }}>
                    {iconEl}
                  </div>
                  <div style={{
                    fontSize: 11, color: 'var(--text-tertiary)',
                    background: 'var(--bg-hover)', padding: '2px 10px',
                    borderRadius: 20, fontWeight: 500,
                  }}>
                    {group.items?.length || 0} tables
                  </div>
                </div>

                <Text strong style={{ fontSize: 14, color: 'var(--text-primary)', display: 'block', marginBottom: 10 }}>
                  {group.name}
                </Text>

                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                  {(group.items || []).map((item: any) => (
                    <div
                      key={item.resource_name}
                      onClick={() => navigate(`/data/${item.resource_name}`)}
                      style={{
                        display: 'inline-flex', alignItems: 'center', gap: 4,
                        padding: '3px 10px', borderRadius: 6,
                        background: `${group.color}12`, color: group.color,
                        fontSize: 12, fontWeight: 600, cursor: 'pointer',
                        transition: 'all 0.2s',
                        border: `1px solid ${group.color}20`,
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.background = `${group.color}25`
                        e.currentTarget.style.transform = 'translateY(-1px)'
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.background = `${group.color}12`
                        e.currentTarget.style.transform = 'translateY(0)'
                      }}
                    >
                      {item.label || item.resource_name}
                      <RightOutlined style={{ fontSize: 9, opacity: 0.6 }} />
                    </div>
                  ))}
                </div>
              </Card>
            </Col>
          )
        })}
      </Row>
    </div>
  )
}
