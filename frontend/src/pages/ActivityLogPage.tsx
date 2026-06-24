import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Card, Table, Tag, Typography, Button, Space, Segmented,
} from 'antd'
import { ReloadOutlined, HistoryOutlined } from '@ant-design/icons'
import api from '../services/api'
import { useAuthStore } from '../stores/authStore'

const { Text } = Typography

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
  const [pageSize, setPageSize] = useState(50)
  const [scope, setScope] = useState<string>('all')
  const user = useAuthStore((s) => s.user)

  const endpoint = scope === 'me' ? '/api/v1/activity-logs/me' : '/api/v1/activity-logs'

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['activity-logs', scope, page, pageSize],
    queryFn: () =>
      api.get(endpoint, { params: { page, limit: pageSize } }),
  })

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: 'Action', dataIndex: 'action', key: 'action',
      render: (v: string) => (
        <Tag color={actionColors[v] || 'default'} style={{ fontWeight: 600 }}>{v}</Tag>
      ),
    },
    {
      title: 'Deskripsi', dataIndex: 'description', key: 'description',
      ellipsis: true,
      render: (v: string) => <Text style={{ fontSize: 13 }}>{v}</Text>,
    },
    {
      title: 'Entity', dataIndex: 'entity_type', key: 'entity_type',
      render: (v: string) => v ? <Tag style={{ fontWeight: 500 }}>{v}</Tag> : '-',
    },
    { title: 'Entity ID', dataIndex: 'entity_id', key: 'entity_id', width: 80 },
    {
      title: 'User', key: 'user', width: 120,
      render: (_: any, r: any) => r.user?.user_name || '-',
    },
    {
      title: 'Waktu', dataIndex: 'created_at', key: 'created_at', width: 180,
      render: (v: string) => new Date(v).toLocaleString('id-ID'),
    },
  ]

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Activity Log</h4>
          <p className="page-subtitle">Riwayat aktivitas pengguna</p>
        </div>
        <Space>
          <Segmented
            options={[
              { label: 'Semua', value: 'all' },
              { label: `Saya (${user?.user_name || ''})`, value: 'me' },
            ]}
            value={scope}
            onChange={(v) => { setScope(v as string); setPage(1) }}
          />
          <Button icon={<ReloadOutlined />} onClick={() => refetch()} style={{ borderRadius: 10 }}>
            Refresh
          </Button>
        </Space>
      </div>
      <Card className="modern-card" bodyStyle={{ padding: 0 }}>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          onChange={(p) => { setPage(p.current || 1); if (p.pageSize) setPageSize(p.pageSize) }}
          pagination={{
            current: page,
            pageSize,
            total: data?.meta?.total,
            showSizeChanger: true,
            showQuickJumper: true,
            pageSizeOptions: ['10', '25', '50', '100'],
            size: 'small',
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total}`,
          }}
          scroll={{ x: 'max-content' }}
        />
      </Card>
    </div>
  )
}
