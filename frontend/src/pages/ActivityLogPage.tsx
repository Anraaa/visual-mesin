import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Card, Table, Tag, Typography, Button, Space, Segmented,
} from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import api from '../services/api'
import { useAuthStore } from '../stores/authStore'

const { Title } = Typography

const actionColors: Record<string, string> = {
  login: 'blue',
  logout: 'default',
  create: 'green',
  update: 'orange',
  delete: 'red',
  export: 'purple',
  export_completed: 'cyan',
}

export default function ActivityLogPage() {
  const [page, setPage] = useState(1)
  const [scope, setScope] = useState<string>('all')
  const user = useAuthStore((s) => s.user)

  const endpoint = scope === 'me' ? '/api/v1/activity-logs/me' : '/api/v1/activity-logs'

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['activity-logs', scope, page],
    queryFn: () =>
      api.get(endpoint, { params: { page, limit: 50 } }).then((r) => r.data),
  })

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: 'Action', dataIndex: 'action', key: 'action',
      render: (v: string) => (
        <Tag color={actionColors[v] || 'default'}>{v}</Tag>
      ),
    },
    { title: 'Deskripsi', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: 'Entity', dataIndex: 'entity_type', key: 'entity_type', render: (v: string) => v ? <Tag>{v}</Tag> : '-' },
    { title: 'Entity ID', dataIndex: 'entity_id', key: 'entity_id', width: 80 },
    {
      title: 'User', key: 'user', width: 120,
      render: (_: any, r: any) => r.user?.name || r.user?.user_name || '-',
    },
    {
      title: 'Waktu', dataIndex: 'created_at', key: 'created_at', width: 180,
      render: (v: string) => new Date(v).toLocaleString('id-ID'),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Activity Log</Title>
        <Space>
          <Segmented
            options={[
              { label: 'Semua', value: 'all' },
              { label: `Saya (${user?.name || ''})`, value: 'me' },
            ]}
            value={scope}
            onChange={(v) => { setScope(v as string); setPage(1) }}
          />
          <Button icon={<ReloadOutlined />} onClick={() => refetch()}>Refresh</Button>
        </Space>
      </div>
      <Card>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          onChange={(p) => setPage(p.current || 1)}
          pagination={{
            current: page,
            pageSize: 50,
            total: data?.meta?.total,
            showTotal: (t) => `Total ${t} aktivitas`,
          }}
          scroll={{ x: 900 }}
        />
      </Card>
    </div>
  )
}
