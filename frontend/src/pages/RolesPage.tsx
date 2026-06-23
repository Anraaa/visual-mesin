import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, Space, Tag, Popconfirm, message, Typography,
} from 'antd'
import { PlusOutlined, DeleteOutlined, SafetyOutlined } from '@ant-design/icons'
import api from '../services/api'

const { Title } = Typography

export default function RolesPage() {
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['roles'],
    queryFn: () => api.get('/api/v1/roles').then((r) => r.data),
  })

  const createMutation = useMutation({
    mutationFn: (values: { name: string }) => api.post('/api/v1/roles', values),
    onSuccess: () => {
      message.success('Role berhasil dibuat')
      setModalOpen(false)
      form.resetFields()
      queryClient.invalidateQueries({ queryKey: ['roles'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/roles/${id}`),
    onSuccess: () => {
      message.success('Role berhasil dihapus')
      queryClient.invalidateQueries({ queryKey: ['roles'] })
    },
  })

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'Nama', dataIndex: 'name', key: 'name', render: (v: string) => <Tag color="blue">{v}</Tag> },
    {
      title: 'Permissions', dataIndex: 'permissions', key: 'permissions',
      render: (perms: any[]) => perms?.map((p: any) => <Tag key={p.id}>{p.name}</Tag>) || '-',
    },
    {
      title: 'Aksi', key: 'action',
      render: (_: any, record: any) => (
        <Popconfirm title="Hapus role?" onConfirm={() => deleteMutation.mutate(record.id)}>
          <Button size="small" danger icon={<DeleteOutlined />} />
        </Popconfirm>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Manajemen Roles</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          Tambah Role
        </Button>
      </div>
      <Card>
        <Table
          dataSource={data?.data || []}
          columns={columns}
          rowKey="id"
          loading={isLoading}
        />
      </Card>

      <Modal
        title="Tambah Role"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={createMutation.isPending}
      >
        <Form form={form} layout="vertical" onFinish={(values) => createMutation.mutate(values)}>
          <Form.Item name="name" label="Nama Role" rules={[{ required: true }]}>
            <Input placeholder="admin, eng, tech, prod" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
