import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, InputNumber, Select, Space, Tag, Popconfirm,
  Typography, message, Tooltip, Switch,
} from 'antd'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, CheckCircleOutlined, CloseCircleOutlined,
  DatabaseOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Text } = Typography

const allResources = [
  { label: 'Material', group: true },
  { value: 'item_measurement', label: 'item_measurement' },
  { value: 'materials', label: 'materials' },
  { value: 'monitoringtl1', label: 'monitoringtl1' },
  { value: 'rscpc1', label: 'rscpc1' },
  { value: 'rtltl1', label: 'rtltl1' },
  { label: 'Curing', group: true },
  { value: 'rtci1', label: 'rtci1' },
  { value: 'rtctr1', label: 'rtctr1' },
  { label: 'Production Control', group: true },
  { value: 'order_report', label: 'order_report' },
  { label: 'Recipe', group: true },
  { value: 'recipe1queue', label: 'recipe1queue' },
  { value: 'recipe1', label: 'recipe1' },
  { value: 'recipe_history', label: 'recipe_history' },
  { label: 'Building', group: true },
  { value: 'rtba1', label: 'rtba1' },
  { value: 'rtba2', label: 'rtba2' },
  { value: 'rtba3', label: 'rtba3' },
  { value: 'rtbc1', label: 'rtbc1' },
  { value: 'rtbc2', label: 'rtbc2' },
  { value: 'rtbc3', label: 'rtbc3' },
  { value: 'rtbc4', label: 'rtbc4' },
  { value: 'rtbe1', label: 'rtbe1' },
  { value: 'rtbe2', label: 'rtbe2' },
  { label: 'Trimming', group: true },
  { value: 'trimmings', label: 'trimmings' },
  { label: 'rteex1', group: true },
  { value: 'rteex1', label: 'rteex1' },
  { label: 'rteex2', group: true },
  { value: 'rteex2', label: 'rteex2' },
  { value: 'recorddatapcs', label: 'recorddatapcs' },
  { value: 'recorddatacyclic', label: 'recorddatacyclic' },
  { label: 'rteex3', group: true },
  { value: 'alarm_history', label: 'alarm_history' },
  { value: 'batch_report', label: 'batch_report' },
  { value: 'datalog', label: 'datalog' },
  { value: 'rteex3head', label: 'rteex3head' },
]

export default function ResourceConnectionPage() {
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<any>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['resource-db-configs'],
    queryFn: () => api.get('/api/v1/resource-db-configs'),
  })

  const rows = data?.data || []

  const configuredResources = rows.map((r: any) => r.resource_name)

  const createMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/resource-db-configs', values),
    onSuccess: () => {
      message.success('Koneksi resource berhasil ditambahkan')
      closeModal()
      queryClient.invalidateQueries({ queryKey: ['resource-db-configs'] })
    },
    onError: (err: any) => {
      message.error(err.response?.data?.message || err.message)
    },
  })

  const updateMutation = useMutation({
    mutationFn: (values: any) => api.put(`/api/v1/resource-db-configs/${editing.id}`, values),
    onSuccess: () => {
      message.success('Koneksi resource berhasil diupdate')
      closeModal()
      queryClient.invalidateQueries({ queryKey: ['resource-db-configs'] })
    },
    onError: (err: any) => {
      message.error(err.response?.data?.message || err.message)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/resource-db-configs/${id}`),
    onSuccess: () => {
      message.success('Koneksi resource berhasil dihapus')
      queryClient.invalidateQueries({ queryKey: ['resource-db-configs'] })
    },
  })

  const testMutation = useMutation({
    mutationFn: (id: number) => api.post(`/api/v1/resource-db-configs/${id}/test`),
    onSuccess: () => {
      message.success('Test koneksi berhasil')
      queryClient.invalidateQueries({ queryKey: ['resource-db-configs'] })
    },
    onError: (err: any) => {
      message.error('Test koneksi gagal: ' + (err.response?.data?.message || err.message))
    },
  })

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, is_active }: { id: number; is_active: boolean }) =>
      api.put(`/api/v1/resource-db-configs/${id}`, { is_active }),
    onSuccess: () => {
      message.success('Status koneksi berhasil diubah')
      queryClient.invalidateQueries({ queryKey: ['resource-db-configs'] })
    },
  })

  const closeModal = () => {
    setModalOpen(false)
    setEditing(null)
    form.resetFields()
  }

  const openEdit = (record: any) => {
    setEditing(record)
    form.setFieldsValue(record)
    setModalOpen(true)
  }

  const resourceOptions = allResources
    .filter((r) => !r.group)
    .map((r) => ({
      ...r,
      disabled: !editing && configuredResources.includes(r.value),
    }))

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: 'Resource', dataIndex: 'resource_name', key: 'resource_name',
      render: (v: string) => <Tag color="blue" style={{ fontWeight: 600 }}>{v}</Tag>,
    },
    {
      title: 'Label', dataIndex: 'label', key: 'label',
      render: (v: string) => v || '-',
    },
    { title: 'Host', dataIndex: 'host', key: 'host' },
    { title: 'Port', dataIndex: 'port', key: 'port', width: 80 },
    { title: 'Database', dataIndex: 'database_name', key: 'database_name' },
    {
      title: 'Aktif', dataIndex: 'is_active', key: 'is_active', width: 80,
      render: (_: any, record: any) => (
        <Switch
          checked={record.is_active}
          size="small"
          loading={toggleActiveMutation.isPending}
          onChange={(checked) =>
            toggleActiveMutation.mutate({ id: record.id, is_active: checked })
          }
        />
      ),
    },
    {
      title: 'Test', key: 'test', width: 100,
      render: (_: any, record: any) => (
        <Tooltip title={record.last_test_message || ''}>
          <Button
            size="small"
            icon={record.is_last_test_success === true ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
            loading={testMutation.isPending}
            onClick={() => testMutation.mutate(record.id)}
            style={{ borderRadius: 8, fontSize: 12 }}
            danger={record.is_last_test_success === false}
          >
            Test
          </Button>
        </Tooltip>
      ),
    },
    {
      title: 'Aksi', key: 'action', width: 120,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} style={{ borderRadius: 8 }} />
          <Popconfirm title="Hapus koneksi resource ini?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} style={{ borderRadius: 8 }} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div className="page-enter">
      <div className="page-header">
        <div>
          <h4>Resource Connection</h4>
          <p className="page-subtitle">Atur koneksi database per resource produksi</p>
        </div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)} className="btn-glow" style={{ borderRadius: 10 }}>
          Tambah Koneksi
        </Button>
      </div>
      <Card className="modern-card" bodyStyle={{ padding: 0 }}>
        <Table
          dataSource={rows}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 25, total: data?.meta?.total }}
        />
      </Card>

      <Modal
        title={
          <Space>
            <DatabaseOutlined style={{ color: 'var(--primary-color)' }} />
            <span style={{ fontWeight: 600 }}>{editing ? 'Edit Koneksi Resource' : 'Tambah Koneksi Resource'}</span>
          </Space>
        }
        open={modalOpen}
        onCancel={closeModal}
        onOk={() => form.submit()}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
        okText={editing ? 'Simpan' : 'Tambah'}
        cancelText="Batal"
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={(values) => {
            const mut = editing ? updateMutation : createMutation
            mut.mutate(values)
          }}
        >
          <Form.Item name="resource_name" label="Pilih Resource" rules={[{ required: true }]}>
            <Select
              placeholder="Pilih resource produksi..."
              disabled={!!editing}
              showSearch
              filterOption={(v, o) => (o?.label as string || '').toLowerCase().includes(v.toLowerCase())}
              options={resourceOptions}
            />
          </Form.Item>
          <Form.Item name="label" label="Label (opsional)">
            <Input placeholder="Misal: Server Produksi 1" />
          </Form.Item>
          <Form.Item name="driver" label="Driver" rules={[{ required: true }]} initialValue="mariadb">
            <Select options={[
              { value: 'mysql', label: 'MySQL' },
              { value: 'mariadb', label: 'MariaDB' },
              { value: 'postgresql', label: 'PostgreSQL' },
            ]} />
          </Form.Item>
          <Space style={{ width: '100%' }} size="large">
            <Form.Item name="host" label="Host" rules={[{ required: true }]} style={{ flex: 1 }}>
              <Input placeholder="192.168.1.100" />
            </Form.Item>
            <Form.Item name="port" label="Port" rules={[{ required: true }]} initialValue={3306}>
              <InputNumber min={1} max={65535} style={{ width: 120 }} />
            </Form.Item>
          </Space>
          <Form.Item name="database_name" label="Nama Database" rules={[{ required: true }]}>
            <Input placeholder="db_produksi" />
          </Form.Item>
          <Space style={{ width: '100%' }} size="large">
            <Form.Item name="username" label="Username" rules={[{ required: true }]} style={{ flex: 1 }}>
              <Input placeholder="root" />
            </Form.Item>
            <Form.Item name="password" label="Password" style={{ flex: 1 }}>
              <Input.Password placeholder="Password" />
            </Form.Item>
          </Space>
          {editing && (
            <Form.Item name="is_active" label="Aktif" valuePropName="checked">
              <Switch />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </div>
  )
}
