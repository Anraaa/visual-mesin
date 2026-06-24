import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, InputNumber, Select, Space, Tag, Popconfirm,
  Typography, message, Tooltip,
} from 'antd'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, CheckCircleOutlined, DatabaseOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Text } = Typography

export default function DBConnectionsPage() {
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<any>(null)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['db-connections'],
    queryFn: () => api.get('/api/v1/db-connections'),
  })

  const createMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/db-connections', values),
    onSuccess: () => {
      message.success('Koneksi database berhasil ditambahkan')
      closeModal()
      queryClient.invalidateQueries({ queryKey: ['db-connections'] })
    },
  })

  const updateMutation = useMutation({
    mutationFn: (values: any) => api.put(`/api/v1/db-connections/${editing.id}`, values),
    onSuccess: () => {
      message.success('Koneksi database berhasil diupdate')
      closeModal()
      queryClient.invalidateQueries({ queryKey: ['db-connections'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/db-connections/${id}`),
    onSuccess: () => {
      message.success('Koneksi database berhasil dihapus')
      queryClient.invalidateQueries({ queryKey: ['db-connections'] })
    },
  })

  const testMutation = useMutation({
    mutationFn: (id: number) => api.post(`/api/v1/db-connections/${id}/test`),
    onSuccess: () => {
      message.success('Test koneksi berhasil')
      queryClient.invalidateQueries({ queryKey: ['db-connections'] })
    },
    onError: (err: any) => {
      message.error('Test koneksi gagal: ' + (err.response?.data?.message || err.message))
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

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Nama', dataIndex: 'name', key: 'name' },
    { title: 'Host', dataIndex: 'host', key: 'host' },
    { title: 'Port', dataIndex: 'port', key: 'port', width: 80 },
    { title: 'Database', dataIndex: 'database_name', key: 'database_name' },
    { title: 'Driver', dataIndex: 'driver', key: 'driver', render: (v: string) => <Tag style={{ fontWeight: 600 }}>{v}</Tag> },
    {
      title: 'Status', dataIndex: 'is_active', key: 'is_active',
      render: (v: boolean) => <Tag color={v ? 'green' : 'red'} style={{ fontWeight: 600 }}>{v ? 'Aktif' : 'Nonaktif'}</Tag>,
    },
    {
      title: 'Test', key: 'test',
      render: (_: any, record: any) => (
        <Tooltip title={record.last_test_message || ''}>
          <Button
            size="small"
            icon={<CheckCircleOutlined />}
            loading={testMutation.isPending}
            onClick={() => testMutation.mutate(record.id)}
            style={{ borderRadius: 8, fontSize: 12 }}
          >
            Test
          </Button>
        </Tooltip>
      ),
    },
    {
      title: 'Aksi', key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} style={{ borderRadius: 8 }} />
          <Popconfirm title="Hapus koneksi?" onConfirm={() => deleteMutation.mutate(record.id)}>
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
          <h4>Database Connections</h4>
          <p className="page-subtitle">Kelola koneksi ke database resource produksi</p>
        </div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)} className="btn-glow" style={{ borderRadius: 10 }}>
          Tambah Koneksi
        </Button>
      </div>
      <Card className="modern-card" bodyStyle={{ padding: 0 }}>
        <Table
          dataSource={data?.data || []}
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
            <span style={{ fontWeight: 600 }}>{editing ? 'Edit Koneksi Database' : 'Tambah Koneksi Database'}</span>
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
          <Form.Item name="name" label="Nama Koneksi" rules={[{ required: true }]}>
            <Input placeholder="My Database" />
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
              <Input placeholder="localhost" />
            </Form.Item>
            <Form.Item name="port" label="Port" rules={[{ required: true }]} initialValue={3306}>
              <InputNumber min={1} max={65535} style={{ width: 120 }} />
            </Form.Item>
          </Space>
          <Form.Item name="database_name" label="Nama Database" rules={[{ required: true }]}>
            <Input placeholder="visual_mesin" />
          </Form.Item>
          <Space style={{ width: '100%' }} size="large">
            <Form.Item name="username" label="Username" rules={[{ required: true }]} style={{ flex: 1 }}>
              <Input placeholder="root" />
            </Form.Item>
            <Form.Item name="password" label="Password" style={{ flex: 1 }}>
              <Input.Password placeholder="Password" />
            </Form.Item>
          </Space>
        </Form>
      </Modal>
    </div>
  )
}
