import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Card, Col, Row, Typography, Spin, Skeleton, Input, Empty, Space } from 'antd'
import {
  BuildOutlined, ExperimentOutlined, FireOutlined, ScissorOutlined,
  ToolOutlined, DashboardOutlined, AlertOutlined, FileTextOutlined,
  ContainerOutlined, ControlOutlined, BookOutlined, ThunderboltOutlined,
  BarcodeOutlined, RightOutlined, SearchOutlined,
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
  const [search, setSearch] = useState('')

  const { data: groupsData, isLoading } = useQuery({
    queryKey: ['resource-groups'],
    queryFn: () => api.get('/api/v1/data-produksi-config/groups'),
  })

  const groups = groupsData?.data || []

  const filteredGroups = useMemo(() => {
    if (!search) return groups
    const q = search.toLowerCase()
    return groups
      .map((g: any) => ({
        ...g,
        items: (g.items || []).filter(
          (i: any) =>
            (i.label || i.resource_name).toLowerCase().includes(q) ||
            i.resource_name.toLowerCase().includes(q)
        ),
      }))
      .filter((g: any) => g.items.length > 0)
  }, [groups, search])

  if (isLoading) {
    return (
      <div className="page-enter">
        <div className="page-header">
          <div>
            <h4>Data Produksi</h4>
            <p className="page-subtitle">Memuat data resource...</p>
          </div>
        </div>
        <Row gutter={[16, 16]}>
          {[1, 2, 3, 4].map((i) => (
            <Col xs={24} sm={12} lg={8} xl={6} key={i}>
              <Card style={{ borderRadius: 14, height: 180 }}>
                <Skeleton active paragraph={{ rows: 2 }} />
              </Card>
            </Col>
          ))}
        </Row>
      </div>
    )
  }

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Data Produksi</h4>
          <p className="page-subtitle">Pilih tabel resource untuk melihat data produksi</p>
        </div>
        <Space>
          <Input
            placeholder="Cari resource..."
            prefix={<SearchOutlined style={{ color: 'var(--text-tertiary)' }} />}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ width: 240, borderRadius: 10 }}
            allowClear
          />
        </Space>
      </div>

      {filteredGroups.length === 0 ? (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '40vh' }}>
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={search ? 'Tidak ada resource yang cocok' : 'Belum ada resource terdaftar'}
          />
        </div>
      ) : (
        <Row gutter={[16, 16]}>
          {filteredGroups.map((group: any) => {
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
                      width: 44, height: 44, borderRadius: 12,
                      background: `${group.color}18`,
                      display: 'flex', alignItems: 'center', justifyContent: 'center',
                      fontSize: 22, color: group.color,
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

                  <Text strong style={{ fontSize: 15, color: 'var(--text-primary)', display: 'block', marginBottom: 12 }}>
                    {group.name}
                  </Text>

                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                    {(group.items || []).map((item: any) => (
                      <div
                        key={item.resource_name}
                        onClick={() => navigate(`/data/${item.resource_name}`)}
                        style={{
                          display: 'inline-flex', alignItems: 'center', gap: 4,
                          padding: '4px 12px', borderRadius: 8,
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
      )}
    </div>
  )
}
