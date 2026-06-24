import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Select, Table, Button, Typography, message, Tag, Space, Checkbox, Spin,
} from 'antd'
import { SaveOutlined } from '@ant-design/icons'
import api from '../services/api'
import { useAuthStore } from '../stores/authStore'

const { Title } = Typography

const actions = ['view', 'view-any', 'create', 'update', 'delete', 'delete-any'] as const

const modules = [
  {
    key: 'dashboard', label: 'Dashboard',
    actions: ['view', 'view-any'],
  },
  {
    key: 'user', label: 'User Management',
    actions: ['view', 'view-any', 'create', 'update', 'delete', 'delete-any'],
  },
  {
    key: 'role', label: 'Role Management',
    actions: ['view', 'view-any', 'create', 'update', 'delete', 'delete-any'],
  },
  {
    key: 'permission', label: 'Permission Management',
    actions: ['view', 'view-any', 'create', 'update', 'delete', 'delete-any'],
  },
  {
    key: 'resource-connection', label: 'Resource Connection',
    actions: ['view', 'view-any', 'create', 'update', 'delete', 'delete-any'],
  },
  {
    key: 'activity-log', label: 'Activity Log',
    actions: ['view', 'view-any'],
  },
  {
    key: 'data-produksi', label: 'Data Produksi',
    actions: ['view', 'view-any'],
  },
  {
    key: 'export', label: 'Export',
    actions: ['view', 'view-any', 'create', 'delete'],
  },
  {
    key: 'ai-chat', label: 'AI Chat',
    actions: ['view', 'view-any'],
  },
]

const allActionLabels: Record<string, string> = {
  view: 'View',
  'view-any': 'View Any',
  create: 'Create',
  update: 'Update',
  delete: 'Delete',
  'delete-any': 'Delete Any',
}

function permName(action: string, moduleKey: string) {
  return `${action}-${moduleKey}`
}

export default function RolePermissionPage() {
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null)
  const [checked, setChecked] = useState<Set<string>>(new Set())
  const queryClient = useQueryClient()

  const { data: rolesData } = useQuery({
    queryKey: ['roles'],
    queryFn: () => api.get('/api/v1/roles'),
  })

  const { data: roleDetail, isLoading: loadingRole } = useQuery({
    queryKey: ['role', selectedRoleId],
    queryFn: () => api.get(`/api/v1/roles/${selectedRoleId}`),
    enabled: !!selectedRoleId,
  })

  useEffect(() => {
    if (roleDetail?.data?.permissions) {
      setChecked(new Set(roleDetail.data.permissions.map((p: any) => p.name)))
    } else {
      setChecked(new Set())
    }
  }, [roleDetail])

  const refreshProfile = async () => {
    try {
      const res = await api.get<any>('/api/v1/auth/me')
      if (res?.data) useAuthStore.getState().updateUser(res.data)
    } catch { /* skip */ }
  }

  const syncMutation = useMutation({
    mutationFn: (perms: string[]) =>
      api.post(`/api/v1/roles/${selectedRoleId}/sync-permissions`, { permissions: perms }),
    onSuccess: () => {
      message.success('Permission berhasil disimpan')
      queryClient.invalidateQueries({ queryKey: ['role', selectedRoleId] })
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      refreshProfile()
    },
    onError: (err: any) => {
      message.error('Gagal: ' + (err.response?.data?.message || err.message))
    },
  })

  const togglePerm = (name: string) => {
    setChecked((prev) => {
      const next = new Set(prev)
      if (next.has(name)) next.delete(name)
      else next.add(name)
      return next
    })
  }

  const columns = [
    {
      title: 'Module', dataIndex: 'label', key: 'label', width: 200,
      render: (v: string) => <strong>{v}</strong>,
    },
    ...actions.map((action) => ({
      title: (
        <span style={{ fontSize: 12, whiteSpace: 'nowrap' }}>{allActionLabels[action]}</span>
      ),
      dataIndex: action,
      key: action,
      width: 80,
      render: (_: any, record: any) => {
        if (!record._actions.includes(action)) return null
        const name = permName(action, record.key)
        return (
          <Checkbox
            checked={checked.has(name)}
            onChange={() => togglePerm(name)}
          />
        )
      },
    })),
  ]

  const tableData = modules.map((m) => ({
    key: m.key,
    label: m.label,
    _actions: m.actions,
    ...Object.fromEntries(actions.map((a) => [a, null])),
  }))

  const handleSave = () => {
    if (!selectedRoleId) return
    syncMutation.mutate(Array.from(checked))
  }

  const selectedRole = rolesData?.data?.find((r: any) => r.id === selectedRoleId)

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>Role Permission</Title>
      </div>
      <Card style={{ marginBottom: 16 }}>
        <Space size="large" wrap>
          <div>
            <div style={{ marginBottom: 6, fontWeight: 500 }}>Pilih Role</div>
            <Select
              placeholder="Pilih role..."
              style={{ width: 280 }}
              value={selectedRoleId}
              onChange={setSelectedRoleId}
              options={(rolesData?.data || []).map((r: any) => ({
                value: r.id,
                label: <><Tag color="blue">{r.name}</Tag> ({r.permissions?.length || 0} permissions)</>,
              }))}
            />
          </div>
          {selectedRoleId && (
            <Space>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={handleSave}
                loading={syncMutation.isPending}
              >
                Simpan Permission
              </Button>
              <span style={{ color: '#888', fontSize: 13 }}>
                {checked.size} permission dipilih
              </span>
            </Space>
          )}
        </Space>
      </Card>

      {selectedRoleId ? (
        <Card>
          {loadingRole ? (
            <div style={{ textAlign: 'center', padding: 40 }}><Spin /></div>
          ) : (
            <Table
              dataSource={tableData}
              columns={columns}
              pagination={false}
              bordered
              size="small"
              scroll={{ x: 'max-content' }}
            />
          )}
        </Card>
      ) : (
        <Card>
          <div style={{ textAlign: 'center', padding: 40, color: '#888' }}>
            Pilih role terlebih dahulu untuk mengatur permission
          </div>
        </Card>
      )}
    </div>
  )
}
