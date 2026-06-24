import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Card, Table, Button, Tag, Typography, message, Tooltip, Space,
} from 'antd'
import {
  DownloadOutlined, ReloadOutlined, CheckCircleOutlined,
  ClockCircleOutlined, CloseCircleOutlined, LoadingOutlined,
  ExportOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Text } = Typography

const statusConfig: Record<string, { color: string; icon: React.ReactNode }> = {
  queued: { color: 'default', icon: <ClockCircleOutlined /> },
  processing: { color: 'processing', icon: <LoadingOutlined /> },
  completed: { color: 'success', icon: <CheckCircleOutlined /> },
  failed: { color: 'error', icon: <CloseCircleOutlined /> },
}

export default function ExportPage() {
  const [page, setPage] = useState(1)
  const [polling, setPolling] = useState(false)

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['exports', page],
    queryFn: () => api.get('/api/v1/exports', { params: { page, limit: 25 } }),
    refetchInterval: polling ? 3000 : false,
  })

  const handleDownload = async (record: any) => {
    try {
      const blob = await api.download(`/api/v1/exports/${record.id}/download`)
      const url = window.URL.createObjectURL(new Blob([blob]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `${record.resource_name}_${record.id}.${record.format}`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
      message.success('Download dimulai')
    } catch {
      message.error('Gagal download file')
    }
  }

  const hasProcessing = (data?.data || []).some(
    (r: any) => r.status === 'queued' || r.status === 'processing',
  )
  if (hasProcessing !== polling) setPolling(hasProcessing)

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Resource', dataIndex: 'resource_name', key: 'resource_name' },
    {
      title: 'Status', dataIndex: 'status', key: 'status',
      render: (v: string) => {
        const cfg = statusConfig[v] || { color: 'default', icon: null }
        return <Tag color={cfg.color} icon={cfg.icon} style={{ fontWeight: 600 }}>{v}</Tag>
      },
    },
    { title: 'Format', dataIndex: 'format', key: 'format', render: (v: string) => <Tag style={{ fontWeight: 600 }}>{v?.toUpperCase()}</Tag> },
    {
      title: 'Progress', key: 'progress',
      render: (_: any, r: any) =>
        r.total_rows > 0 ? (
          <Text style={{ fontSize: 13 }}>{r.processed_rows} / {r.total_rows}</Text>
        ) : '-',
    },
    {
      title: 'File', dataIndex: 'file_url', key: 'file_url',
      render: (v: string, r: any) =>
        v ? (
          <Button type="link" icon={<DownloadOutlined />} onClick={() => handleDownload(r)} style={{ padding: 0 }}>
            Download
          </Button>
        ) : '-',
    },
    {
      title: 'Error', dataIndex: 'error_message', key: 'error_message',
      ellipsis: true,
      render: (v: string) => v ? <Tooltip title={v}><Tag color="red">Error</Tag></Tooltip> : '-',
    },
    {
      title: 'Dibuat', dataIndex: 'created_at', key: 'created_at',
      render: (v: string) => new Date(v).toLocaleString('id-ID'),
      width: 170,
    },
  ]

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Export Data</h4>
          <p className="page-subtitle">Kelola dan download file export</p>
        </div>
        <Button icon={<ReloadOutlined />} onClick={() => refetch()} style={{ borderRadius: 10 }}>
          Refresh
        </Button>
      </div>
      <Card className="modern-card" bodyStyle={{ padding: 0 }}>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          onChange={(p) => setPage(p.current || 1)}
          pagination={{
            current: page,
            pageSize: 25,
            total: data?.meta?.total,
            showTotal: (t) => `Total ${t} export`,
          }}
        />
      </Card>
    </div>
  )
}
