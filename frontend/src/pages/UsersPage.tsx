import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, Select, Space, Tag, Popconfirm, message, Typography,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, TeamOutlined } from '@ant-design/icons'
import api from '../services/api'

const { Title } = Typography

export default function UsersPage() {
  const [modalOpen, setModalOpen] = useState(false)
  const [roleModalOpen, setRoleModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<any>(null)
  const [selectedUserId, setSelectedUserId] = useState<number>(0)
  const [form] = Form.useForm()
  const [roleForm] = Form.useForm()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => api.get('/api/v1/users').then((r) => r.data),
  })

  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => api.get('/api/v1/roles').then((r) => r.data),
  })

  const updateMutation = useMutation({
    mutationFn: (values: any) =>
      editingUser
        ? api.put(`/api/v1/users/${editingUser.id}`, values)
        : Promise.resolve(),
    onSuccess: () => {
      message.success('User berhasil diupdate')
      setModalOpen(false)
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/users/${id}`),
    onSuccess: () => {
      message.success('User berhasil dihapus')
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })

  const assignRoleMutation = useMutation({
    mutationFn: (values: { role_id: number }) =>
      api.post(`/api/v1/users/${selectedUserId}/assign-role`, values),
    onSuccess: () => {
      message.success('Role berhasil ditetapkan')
      setRoleModalOpen(false)
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })

  const openEdit = (record: any) => {
    setEditingUser(record)
    form.setFieldsValue(record)
    setModalOpen(true)
  }

  const openAssignRole = (record: any) => {
    setSelectedUserId(record.id)
    roleForm.resetFields()
    setRoleModalOpen(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Nama', dataIndex: 'user_name', key: 'user_name' },
    { title: 'Email', dataIndex: 'email', key: 'email' },
    { title: 'Level', dataIndex: 'user_level', key: 'user_level', render: (v: string) => <Tag>{v}</Tag> },
    {
      title: 'Roles', dataIndex: 'roles', key: 'roles',
      render: (roles: any[]) => roles?.map((r: any) => <Tag key={r.id} color="blue">{r.name}</Tag>),
    },
    {
      title: 'Aksi', key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<TeamOutlined />} onClick={() => openAssignRole(record)}>Role</Button>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)} />
          <Popconfirm title="Hapus user?" onConfirm={() => deleteMutation.mutate(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={4}>Manajemen Users</Title>
      <Card>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 25, total: data?.meta?.total }}
        />
      </Card>

      <Modal
        title="Edit User"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={updateMutation.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(values) => updateMutation.mutate(values)}>
          <Form.Item name="user_name" label="Nama" rules={[{ required: true }]}>
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
          <Form.Item name="department" label="Departemen">
            <Input />
          </Form.Item>
          <Form.Item name="jabatan" label="Jabatan">
            <Input />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="Assign Role"
        open={roleModalOpen}
        onCancel={() => setRoleModalOpen(false)}
        onOk={() => roleForm.submit()}
      >
        <Form form={roleForm} layout="vertical" onFinish={(values) => assignRoleMutation.mutate(values)}>
          <Form.Item name="role_id" label="Role" rules={[{ required: true }]}>
            <Select
              options={(rolesData?.data || []).map((r: any) => ({ value: r.id, label: r.name }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
