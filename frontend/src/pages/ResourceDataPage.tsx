import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Input, Select, Button, Space, Typography, Tag, Drawer, Descriptions, Spin,
  Modal, Transfer, message,
} from 'antd'
import {
  ArrowLeftOutlined, SearchOutlined, ReloadOutlined, EyeOutlined,
  DownloadOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Title } = Typography

export default function ResourceDataPage() {
  const { resource } = useParams<{ resource: string }>()
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchBy, setSearchBy] = useState('')
  const [sortBy, setSortBy] = useState('')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('asc')
  const [selectedRow, setSelectedRow] = useState<any>(null)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [exportModalOpen, setExportModalOpen] = useState(false)
  const [selectedColumns, setSelectedColumns] = useState<string[]>([])
  const queryClient = useQueryClient()

  const { data, isLoading, isFetching, refetch } = useQuery({
    queryKey: ['resource-data', resource, page, search, searchBy, sortBy, sortDir],
    queryFn: () =>
      api.get(`/api/v1/resources/${resource}`, {
        params: { page, limit: 25, search, search_by: searchBy, sort_by: sortBy, sort_dir: sortDir },
      }).then((r) => r.data),
  })

  const rows = data?.data || []
  const total = data?.meta?.total || 0
  const allColumns = rows.length > 0 ? Object.keys(rows[0]) : []

  const exportMutation = useMutation({
    mutationFn: (body: any) => api.post('/api/v1/exports', body),
    onSuccess: (r) => {
      message.success('Export job submitted! Cek halaman Export.')
      setExportModalOpen(false)
      queryClient.invalidateQueries({ queryKey: ['exports'] })
    },
    onError: (err: any) => {
      message.error('Gagal submit export: ' + (err.response?.data?.message || err.message))
    },
  })

  const columns = allColumns.slice(0, 10).map((key) => ({
    title: key,
    dataIndex: key,
    key,
    ellipsis: true,
    width: key === 'id' || key === 'recid' ? 80 : undefined,
    sorter: true,
  })).concat([{
    title: 'Aksi',
    key: 'action',
    width: 80,
    render: (_: any, record: any) => (
      <Button
        size="small"
        icon={<EyeOutlined />}
        onClick={() => { setSelectedRow(record); setDrawerOpen(true) }}
      />
    ),
  }])

  const handleTableChange = (pagination: any, _filters: any, sorter: any) => {
    setPage(pagination.current)
    if (sorter.field) {
      setSortBy(sorter.field)
      setSortDir(sorter.order === 'ascend' ? 'asc' : 'desc')
    }
  }

  const handleExport = () => {
    exportMutation.mutate({
      resource_name: resource,
      columns: selectedColumns.length > 0 ? selectedColumns : allColumns,
      format: 'csv',
    })
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/data')}>Kembali</Button>
        <Title level={4} style={{ margin: 0 }}>{resource?.toUpperCase()}</Title>
        <Tag color="blue">{total} records</Tag>
        <Button icon={<DownloadOutlined />} onClick={() => {
          setSelectedColumns(allColumns)
          setExportModalOpen(true)
        }}>
          Export CSV
        </Button>
      </Space>

      <Card
        extra={
          <Space>
            <Select
              allowClear
              placeholder="Cari di kolom"
              style={{ width: 150 }}
              value={searchBy || undefined}
              onChange={(v) => setSearchBy(v || '')}
              options={allColumns.map((c) => ({ value: c, label: c }))}
            />
            <Input
              placeholder="Cari..."
              prefix={<SearchOutlined />}
              style={{ width: 200 }}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              onPressEnter={() => refetch()}
            />
            <Button icon={<ReloadOutlined />} onClick={() => refetch()} />
          </Space>
        }
      >
        <Table
          dataSource={rows}
          columns={columns}
          rowKey={(record) => record.id || record.recid || Math.random()}
          loading={isLoading || isFetching}
          onChange={handleTableChange}
          pagination={{
            current: page,
            pageSize: 25,
            total,
            showSizeChanger: false,
            showTotal: (t) => `Total ${t} records`,
          }}
          scroll={{ x: 'max-content' }}
          size="small"
        />
      </Card>

      <Drawer
        title="Detail Record"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={600}
      >
        {selectedRow ? (
          <Descriptions column={1} bordered size="small">
            {Object.entries(selectedRow).map(([key, value]) => (
              <Descriptions.Item key={key} label={key}>
                {value === null ? '-' : String(value)}
              </Descriptions.Item>
            ))}
          </Descriptions>
        ) : <Spin />}
      </Drawer>

      <Modal
        title="Export CSV — Pilih Kolom"
        open={exportModalOpen}
        onCancel={() => setExportModalOpen(false)}
        onOk={handleExport}
        confirmLoading={exportMutation.isPending}
        width={600}
      >
        <Transfer
          dataSource={allColumns.map((c) => ({ key: c, title: c }))}
          titles={['Available', 'Selected']}
          targetKeys={selectedColumns}
          onChange={(keys) => setSelectedColumns(keys as string[])}
          render={(item) => item.title}
          listStyle={{ width: 250, height: 400 }}
        />
      </Modal>
    </div>
  )
}
