import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, Select, Space, Tag, Popconfirm, message, Typography,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons'
import api from '../services/api'

const { Text } = Typography

export default function UsersPage() {
  const [createOpen, setCreateOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<any>(null)
  const [createForm] = Form.useForm()
  const [editForm] = Form.useForm()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => api.get('/api/v1/users'),
  })

  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => api.get('/api/v1/roles'),
  })

  const roleOptions = (rolesData?.data || []).map((r: any) => ({ value: r.id, label: r.name }))

  const createMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/users', values),
    onSuccess: () => {
      message.success('User berhasil dibuat')
      setCreateOpen(false)
      createForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || 'Gagal membuat user')
    },
  })

  const updateMutation = useMutation({
    mutationFn: async (values: any) => {
      await api.put(`/api/v1/users/${editingUser.id}`, values)
    },
    onSuccess: () => {
      message.success('User berhasil diupdate')
      setEditOpen(false)
      setEditingUser(null)
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || 'Gagal mengupdate user')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/users/${id}`),
    onSuccess: () => {
      message.success('User berhasil dihapus')
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })

  const openEdit = (record: any) => {
    setEditingUser(record)
    editForm.setFieldsValue({
      ...record,
      role_ids: (record.roles || []).map((r: any) => r.id),
    })
    setEditOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Nama', dataIndex: 'user_name', key: 'user_name' },
    { title: 'Email', dataIndex: 'email', key: 'email' },
    {
      title: 'Level', dataIndex: 'user_level', key: 'user_level',
      render: (v: string) => (
        <Tag color={v === 'admin' ? 'red' : v === 'eng' ? 'blue' : v === 'tech' ? 'orange' : 'green'}>
          {v}
        </Tag>
      ),
    },
    {
      title: 'Roles', dataIndex: 'roles', key: 'roles',
      render: (roles: any[]) => roles?.map((r: any) => <Tag key={r.id} color="blue" style={{ marginBottom: 2 }}>{r.name}</Tag>),
    },
    {
      title: 'Aksi', key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} style={{ borderRadius: 8 }} />
          <Popconfirm title="Hapus user?" onConfirm={() => deleteMutation.mutate(record.id)}>
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
          <h4>Manajemen Users</h4>
          <p className="page-subtitle">Kelola semua pengguna sistem</p>
        </div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { createForm.resetFields(); setCreateOpen(true) }}>
          Tambah User
        </Button>
      </div>
      <Card className="modern-card" bodyStyle={{ padding: 0 }}>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 25, total: data?.meta?.total, showTotal: (t) => `Total ${t} user` }}
        />
      </Card>

      <Modal
        title={<Space><UserOutlined style={{ color: 'var(--primary-color)' }} /><span style={{ fontWeight: 600 }}>Tambah User</span></Space>}
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={() => createForm.submit()}
        confirmLoading={createMutation.isPending}
        okText="Simpan"
        cancelText="Batal"
        width={520}
      >
        <Form form={createForm} layout="vertical" onFinish={(values) => createMutation.mutate(values)}>
          <Form.Item name="user_name" label="Nama" rules={[{ required: true, message: 'Nama wajib diisi' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email', message: 'Email wajib diisi' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="Password" rules={[{ required: true, min: 6, message: 'Password minimal 6 karakter' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="user_level" label="Level" initialValue="prod">
            <Select
              options={[
                { value: 'admin', label: 'Admin' },
                { value: 'eng', label: 'Engineer' },
                { value: 'tech', label: 'Teknisi' },
                { value: 'prod', label: 'Produksi' },
              ]}
            />
          </Form.Item>
          <Form.Item name="nip" label="NIP">
            <Input />
          </Form.Item>
          <Form.Item name="role_ids" label="Roles">
            <Select mode="multiple" options={roleOptions} placeholder="Pilih role" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={<Space><EditOutlined style={{ color: 'var(--primary-color)' }} /><span style={{ fontWeight: 600 }}>Edit User</span></Space>}
        open={editOpen}
        onCancel={() => { setEditOpen(false); setEditingUser(null) }}
        onOk={() => editForm.submit()}
        confirmLoading={updateMutation.isPending}
        okText="Simpan"
        cancelText="Batal"
        width={520}
      >
        <Form form={editForm} layout="vertical" onFinish={(values) => updateMutation.mutate(values)}>
          <Form.Item name="user_name" label="Nama" rules={[{ required: true, message: 'Nama wajib diisi' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="Email">
            <Input />
          </Form.Item>
          <Form.Item name="user_level" label="Level">
            <Select
              options={[
                { value: 'admin', label: 'Admin' },
                { value: 'eng', label: 'Engineer' },
                { value: 'tech', label: 'Teknisi' },
                { value: 'prod', label: 'Produksi' },
              ]}
            />
          </Form.Item>
          <Form.Item name="nip" label="NIP">
            <Input />
          </Form.Item>
          <Form.Item name="department" label="Departemen">
            <Input />
          </Form.Item>
          <Form.Item name="jabatan" label="Jabatan">
            <Input />
          </Form.Item>
          <Form.Item name="role_ids" label="Roles">
            <Select mode="multiple" options={roleOptions} placeholder="Pilih role" />
          </Form.Item>
          <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
            Kosongkan field password jika tidak ingin mengganti password.
          </Text>
          <Form.Item name="old_password" label="Password Lama">
            <Input.Password placeholder="Wajib jika mengganti password" />
          </Form.Item>
          <Form.Item name="password" label="Password Baru">
            <Input.Password placeholder="Kosongkan jika tidak diganti" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
