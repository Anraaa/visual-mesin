import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Button, Modal, Form, Input, InputNumber, Select, Space, Tag, Popconfirm,
  Typography, message, Switch, Collapse, ColorPicker, Tooltip,
} from 'antd'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, DatabaseOutlined,
  BuildOutlined, ExperimentOutlined, FireOutlined, ScissorOutlined,
  ToolOutlined, DashboardOutlined, AlertOutlined, FileTextOutlined,
  ContainerOutlined, ControlOutlined, BookOutlined, ThunderboltOutlined,
  BarcodeOutlined,
} from '@ant-design/icons'
import api from '../services/api'

const { Title, Text } = Typography

const iconOptions: { value: string; icon: any }[] = [
  { value: 'BuildOutlined', icon: <BuildOutlined /> },
  { value: 'ExperimentOutlined', icon: <ExperimentOutlined /> },
  { value: 'FireOutlined', icon: <FireOutlined /> },
  { value: 'ScissorOutlined', icon: <ScissorOutlined /> },
  { value: 'ToolOutlined', icon: <ToolOutlined /> },
  { value: 'DashboardOutlined', icon: <DashboardOutlined /> },
  { value: 'AlertOutlined', icon: <AlertOutlined /> },
  { value: 'FileTextOutlined', icon: <FileTextOutlined /> },
  { value: 'ContainerOutlined', icon: <ContainerOutlined /> },
  { value: 'ControlOutlined', icon: <ControlOutlined /> },
  { value: 'BookOutlined', icon: <BookOutlined /> },
  { value: 'ThunderboltOutlined', icon: <ThunderboltOutlined /> },
  { value: 'BarcodeOutlined', icon: <BarcodeOutlined /> },
]

const iconMap: Record<string, any> = {}
for (const opt of iconOptions) {
  iconMap[opt.value] = opt.icon
}

const dataTypeOptions = [
  { value: 'varchar', label: 'VARCHAR' },
  { value: 'int', label: 'INT' },
  { value: 'bigint', label: 'BIGINT' },
  { value: 'decimal', label: 'DECIMAL' },
  { value: 'datetime', label: 'DATETIME' },
  { value: 'date', label: 'DATE' },
  { value: 'text', label: 'TEXT' },
  { value: 'longtext', label: 'LONGTEXT' },
  { value: 'boolean', label: 'BOOLEAN (TINYINT)' },
  { value: 'enum', label: 'ENUM' },
]

interface ColumnInput {
  key: string
  column_name: string
  data_type: string
  length?: number
  decimal_places?: number
  enum_values?: string
  is_nullable: boolean
  default_value?: string
  is_primary: boolean
  is_auto_increment: boolean
  sort_order: number
}

let colKeyCounter = 0
function newColKey() { return `col_${++colKeyCounter}` }

function emptyColumn(): ColumnInput {
  return {
    key: newColKey(),
    column_name: '',
    data_type: 'varchar',
    length: 255,
    decimal_places: undefined,
    enum_values: '',
    is_nullable: true,
    default_value: '',
    is_primary: false,
    is_auto_increment: false,
    sort_order: 0,
  }
}

export default function DataProduksiConfigPage() {
  const [groupModalOpen, setGroupModalOpen] = useState(false)
  const [editingGroup, setEditingGroup] = useState<any>(null)
  const [groupForm] = Form.useForm()

  const [itemModalOpen, setItemModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [itemForm] = Form.useForm()

  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [createResourceName, setCreateResourceName] = useState('')
  const [createLabel, setCreateLabel] = useState('')
  const [createGroupId, setCreateGroupId] = useState<number | null>(null)
  const [columns, setColumns] = useState<ColumnInput[]>([emptyColumn()])

  const queryClient = useQueryClient()

  const { data: groupsData, isLoading } = useQuery({
    queryKey: ['resource-groups'],
    queryFn: () => api.get('/api/v1/data-produksi-config/groups'),
  })

  const groups = groupsData?.data || []

  // Group mutations
  const createGroupMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/data-produksi-config/groups', values),
    onSuccess: () => { message.success('Grup berhasil dibuat'); closeGroupModal(); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  const updateGroupMutation = useMutation({
    mutationFn: (values: any) => api.put(`/api/v1/data-produksi-config/groups/${editingGroup.id}`, values),
    onSuccess: () => { message.success('Grup berhasil diupdate'); closeGroupModal(); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  const deleteGroupMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/data-produksi-config/groups/${id}`),
    onSuccess: () => { message.success('Grup berhasil dihapus'); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  // Item mutations
  const createItemMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/data-produksi-config/items', values),
    onSuccess: () => { message.success('Resource berhasil ditambahkan'); closeItemModal(); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  const updateItemMutation = useMutation({
    mutationFn: (values: any) => api.put(`/api/v1/data-produksi-config/items/${editingItem.id}`, values),
    onSuccess: () => { message.success('Resource berhasil diupdate'); closeItemModal(); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  const deleteItemMutation = useMutation({
    mutationFn: (id: number) => api.delete(`/api/v1/data-produksi-config/items/${id}`),
    onSuccess: () => { message.success('Resource berhasil dihapus'); queryClient.invalidateQueries({ queryKey: ['resource-groups'] }) },
    onError: (err: any) => { message.error(err.response?.data?.message || err.message) },
  })

  // Create resource with table
  const createTableMutation = useMutation({
    mutationFn: (values: any) => api.post('/api/v1/data-produksi-config/create-table', values),
    onSuccess: () => {
      message.success('Resource beserta tabel berhasil dibuat')
      closeCreateModal()
      queryClient.invalidateQueries({ queryKey: ['resource-groups'] })
    },
    onError: (err: any) => {
      message.error(err.response?.data?.message || err.message)
    },
  })

  const closeGroupModal = () => { setGroupModalOpen(false); setEditingGroup(null); groupForm.resetFields() }
  const closeItemModal = () => { setItemModalOpen(false); setEditingItem(null); itemForm.resetFields() }
  const closeCreateModal = () => {
    setCreateModalOpen(false)
    setCreateResourceName('')
    setCreateLabel('')
    setCreateGroupId(null)
    setColumns([emptyColumn()])
  }

  const openEditGroup = (group: any) => {
    setEditingGroup(group)
    groupForm.setFieldsValue({ name: group.name, color: group.color, icon: group.icon, sort_order: group.sort_order })
    setGroupModalOpen(true)
  }

  const openEditItem = (item: any) => {
    setEditingItem(item)
    itemForm.setFieldsValue({ label: item.label, sort_order: item.sort_order, is_active: item.is_active })
    setItemModalOpen(true)
  }

  const openCreateResource = (groupId: number) => {
    setCreateGroupId(groupId)
    setCreateModalOpen(true)
  }

  const addColumn = () => setColumns((prev) => [...prev, emptyColumn()])

  const removeColumn = (key: string) => {
    setColumns((prev) => prev.filter((c) => c.key !== key))
  }

  const updateColumn = (key: string, field: string, value: any) => {
    setColumns((prev) => prev.map((c) => (c.key === key ? { ...c, [field]: value } : c)))
  }

  const handleCreateResource = () => {
    if (!createResourceName || !createGroupId) {
      message.warning('Nama resource harus diisi')
      return
    }
    const validCols = columns.filter((c) => c.column_name)
    if (validCols.length === 0) {
      message.warning('Minimal satu kolom harus diisi')
      return
    }
    createTableMutation.mutate({
      group_id: createGroupId,
      resource_name: createResourceName,
      label: createLabel,
      columns: validCols.map((c) => ({
        column_name: c.column_name,
        data_type: c.data_type,
        length: c.length || null,
        decimal_places: c.decimal_places || null,
        enum_values: c.enum_values || '',
        is_nullable: c.is_nullable,
        default_value: c.default_value || '',
        is_primary: c.is_primary,
        is_auto_increment: c.is_auto_increment,
        sort_order: c.sort_order,
      })),
    })
  }

  const itemColumns = [
    { title: 'Resource', dataIndex: 'resource_name', key: 'resource_name', render: (v: string) => <Tag color="blue">{v}</Tag> },
    { title: 'Label', dataIndex: 'label', key: 'label', render: (v: string) => v || '-' },
    { title: 'Sort', dataIndex: 'sort_order', key: 'sort_order', width: 60 },
    {
      title: 'Aktif', dataIndex: 'is_active', key: 'is_active', width: 80,
      render: (_: any, record: any) => <Switch checked={record.is_active} size="small" disabled />,
    },
    {
      title: 'Aksi', key: 'action', width: 160,
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEditItem(record)} />
          <Popconfirm title="Hapus resource ini?" onConfirm={() => deleteItemMutation.mutate(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <div>
          <h4>Data Produksi Config</h4>
          <p className="page-subtitle">Kelola grup, resource, dan buat tabel baru</p>
        </div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingGroup(null); groupForm.resetFields(); setGroupModalOpen(true) }}>
          Tambah Grup
        </Button>
      </div>

      {isLoading ? <Card loading /> : groups.map((group: any) => (
        <Card
          key={group.id}
          style={{ marginBottom: 16, borderLeft: `3px solid ${group.color}` }}
          title={
            <Space>
              <span style={{ color: group.color, fontSize: 18 }}>{iconMap[group.icon] || <BuildOutlined />}</span>
              <Text strong style={{ fontSize: 15 }}>{group.name}</Text>
              <Tag color={group.color}>{group.items?.length || 0} resources</Tag>
            </Space>
          }
          extra={
            <Space>
              <Button size="small" onClick={() => openCreateResource(group.id)} icon={<PlusOutlined />}>Tambah Resource</Button>
              <Button size="small" icon={<EditOutlined />} onClick={() => openEditGroup(group)} />
              <Popconfirm title="Hapus grup?" onConfirm={() => deleteGroupMutation.mutate(group.id)}>
                <Button size="small" danger icon={<DeleteOutlined />} />
              </Popconfirm>
            </Space>
          }
        >
          <Table
            dataSource={group.items || []}
            columns={itemColumns}
            rowKey="id"
            pagination={false}
            size="small"
          />
        </Card>
      ))}

      {/* Group Modal */}
      <Modal
        title={editingGroup ? 'Edit Grup' : 'Tambah Grup'}
        open={groupModalOpen}
        onCancel={closeGroupModal}
        onOk={() => groupForm.submit()}
        confirmLoading={createGroupMutation.isPending || updateGroupMutation.isPending}
      >
        <Form
          form={groupForm}
          layout="vertical"
          onFinish={(values) => {
            const mut = editingGroup ? updateGroupMutation : createGroupMutation
            const color = typeof values.color === 'string' ? values.color : values.color?.toHexString?.() || '#1677ff'
            mut.mutate({ ...values, color })
          }}
        >
          <Form.Item name="name" label="Nama Grup" rules={[{ required: true }]}>
            <Input placeholder="Building Sensitive" />
          </Form.Item>
          <Space style={{ width: '100%' }} size="large">
            <Form.Item name="color" label="Warna" initialValue="#1677ff">
              <ColorPicker format="hex" showText />
            </Form.Item>
            <Form.Item name="icon" label="Ikon" initialValue="BuildOutlined">
              <Select style={{ width: 160 }}>
                {iconOptions.map((opt) => (
                  <Select.Option key={opt.value} value={opt.value}>
                    <Space>{opt.icon}<span>{opt.value}</span></Space>
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          </Space>
          <Form.Item name="sort_order" label="Urutan" initialValue={0}>
            <InputNumber min={0} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Item Modal */}
      <Modal
        title="Edit Resource"
        open={itemModalOpen}
        onCancel={closeItemModal}
        onOk={() => itemForm.submit()}
        confirmLoading={updateItemMutation.isPending}
      >
        <Form
          form={itemForm}
          layout="vertical"
          onFinish={(values) => updateItemMutation.mutate(values)}
        >
          <Form.Item name="label" label="Label">
            <Input placeholder="Display label" />
          </Form.Item>
          <Space style={{ width: '100%' }} size="large">
            <Form.Item name="sort_order" label="Urutan" initialValue={0}>
              <InputNumber min={0} />
            </Form.Item>
            <Form.Item name="is_active" label="Aktif" valuePropName="checked" initialValue={true}>
              <Switch />
            </Form.Item>
          </Space>
        </Form>
      </Modal>

      {/* Create Resource Modal */}
      <Modal
        title="Tambah Resource Baru + Tabel"
        open={createModalOpen}
        onCancel={closeCreateModal}
        onOk={handleCreateResource}
        confirmLoading={createTableMutation.isPending}
        width={800}
        okText="Buat Resource"
      >
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Space style={{ width: '100%' }}>
            <Form.Item label="Nama Resource" style={{ marginBottom: 0 }} required>
              <Input
                placeholder="rtbs1"
                value={createResourceName}
                onChange={(e) => setCreateResourceName(e.target.value)}
                style={{ width: 200 }}
              />
            </Form.Item>
            <Form.Item label="Label" style={{ marginBottom: 0 }}>
              <Input
                placeholder="Building Sensitive 1"
                value={createLabel}
                onChange={(e) => setCreateLabel(e.target.value)}
                style={{ width: 250 }}
              />
            </Form.Item>
          </Space>

          <div style={{ borderTop: '1px solid var(--border-color)', paddingTop: 16 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
              <Text strong>Definisi Kolom</Text>
              <Button size="small" icon={<PlusOutlined />} onClick={addColumn}>Tambah Kolom</Button>
            </div>
            {columns.map((col, idx) => (
              <Card
                key={col.key}
                size="small"
                style={{ marginBottom: 8 }}
                extra={
                  columns.length > 1 ? (
                    <Button size="small" danger icon={<DeleteOutlined />} onClick={() => removeColumn(col.key)} />
                  ) : null
                }
              >
                <Space style={{ width: '100%' }} wrap>
                  <Input
                    placeholder="column_name"
                    value={col.column_name}
                    onChange={(e) => updateColumn(col.key, 'column_name', e.target.value)}
                    style={{ width: 150 }}
                    addonBefore={`${idx + 1}`}
                  />
                  <Select
                    value={col.data_type}
                    onChange={(v) => updateColumn(col.key, 'data_type', v)}
                    style={{ width: 140 }}
                    options={dataTypeOptions}
                  />
                  {(col.data_type === 'varchar') && (
                    <InputNumber
                      placeholder="Length"
                      value={col.length}
                      onChange={(v) => updateColumn(col.key, 'length', v)}
                      min={1} max={65535}
                      style={{ width: 90 }}
                      addonBefore="Len"
                    />
                  )}
                  {(col.data_type === 'decimal') && (
                    <>
                      <InputNumber
                        placeholder="Precision"
                        value={col.length}
                        onChange={(v) => updateColumn(col.key, 'length', v)}
                        min={1} max={65}
                        style={{ width: 90 }}
                        addonBefore="Prec"
                      />
                      <InputNumber
                        placeholder="Scale"
                        value={col.decimal_places}
                        onChange={(v) => updateColumn(col.key, 'decimal_places', v)}
                        min={0} max={30}
                        style={{ width: 90 }}
                        addonBefore="Scale"
                      />
                    </>
                  )}
                  {(col.data_type === 'enum') && (
                    <Input
                      placeholder="'val1','val2'"
                      value={col.enum_values}
                      onChange={(e) => updateColumn(col.key, 'enum_values', e.target.value)}
                      style={{ width: 200 }}
                    />
                  )}
                  <Tooltip title="Nullable">
                    <Select
                      value={col.is_nullable ? 'yes' : 'no'}
                      onChange={(v) => updateColumn(col.key, 'is_nullable', v === 'yes')}
                      style={{ width: 70 }}
                      options={[
                        { value: 'yes', label: 'NULL' },
                        { value: 'no', label: 'NN' },
                      ]}
                    />
                  </Tooltip>
                  <Tooltip title="Primary Key">
                    <Select
                      value={col.is_primary ? 'yes' : 'no'}
                      onChange={(v) => updateColumn(col.key, 'is_primary', v === 'yes')}
                      style={{ width: 40 }}
                      options={[
                        { value: 'no', label: ' ' },
                        { value: 'yes', label: 'PK' },
                      ]}
                    />
                  </Tooltip>
                  <Tooltip title="Auto Increment">
                    <Select
                      value={col.is_auto_increment ? 'yes' : 'no'}
                      onChange={(v) => updateColumn(col.key, 'is_auto_increment', v === 'yes')}
                      style={{ width: 40 }}
                      options={[
                        { value: 'no', label: ' ' },
                        { value: 'yes', label: 'AI' },
                      ]}
                    />
                  </Tooltip>
                  <Input
                    placeholder="Default"
                    value={col.default_value}
                    onChange={(e) => updateColumn(col.key, 'default_value', e.target.value)}
                    style={{ width: 120 }}
                  />
                </Space>
              </Card>
            ))}
          </div>
        </Space>
      </Modal>
    </div>
  )
}
