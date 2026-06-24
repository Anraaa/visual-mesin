import { useState, useMemo } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Table, Input, Select, Button, Space, Typography, Tag, Drawer, Descriptions, Spin,
  Modal, Transfer, message, Tabs, Row, Col, DatePicker,
} from 'antd'
import dayjs from 'dayjs'
import {
  ArrowLeftOutlined, SearchOutlined, ReloadOutlined, EyeOutlined,
  DownloadOutlined, TableOutlined, DashboardOutlined,
} from '@ant-design/icons'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, AreaChart, Area, LineChart, Line, Legend,
} from 'recharts'
import api from '../services/api'

const { Text, Title } = Typography

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div style={{
        background: 'var(--bg-elevated)',
        padding: '12px 16px',
        borderRadius: 12,
        boxShadow: 'var(--shadow-lg)',
        border: '1px solid var(--border-color)',
      }}>
        <Text strong style={{ fontSize: 13, display: 'block', marginBottom: 4 }}>{label}</Text>
        {payload.map((p: any, i: number) => (
          <div key={i} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 13 }}>
            <span style={{ width: 8, height: 8, borderRadius: '50%', background: p.color }} />
            <span style={{ color: 'var(--text-secondary)' }}>{p.name}:</span>
            <span style={{ fontWeight: 700 }}>{p.value?.toFixed?.(2) ?? p.value}</span>
          </div>
        ))}
      </div>
    )
  }
  return null
}

const QualityTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    const total = payload.reduce((s: number, p: any) => s + (p.value || 0), 0)
    const ok = payload.find((p: any) => p.name === 'OK')?.value || 0
    return (
      <div style={{
        background: 'var(--bg-elevated)',
        padding: '12px 16px',
        borderRadius: 12,
        boxShadow: 'var(--shadow-lg)',
        border: '1px solid var(--border-color)',
      }}>
        <Text strong style={{ fontSize: 13, display: 'block', marginBottom: 6 }}>{label}</Text>
        {payload.map((p: any, i: number) => (
          <div key={i} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 13, marginBottom: 2 }}>
            <span style={{ width: 8, height: 8, borderRadius: '50%', background: p.color }} />
            <span style={{ color: 'var(--text-secondary)' }}>{p.name}:</span>
            <span style={{ fontWeight: 700 }}>{p.value}</span>
          </div>
        ))}
        <div style={{ borderTop: '1px solid var(--border-color)', marginTop: 6, paddingTop: 6, fontSize: 13 }}>
          <span style={{ color: 'var(--text-secondary)' }}>OK Rate: </span>
          <span style={{ fontWeight: 700, color: '#10b981' }}>{total > 0 ? `${(ok / total * 100).toFixed(1)}%` : '-'}</span>
        </div>
      </div>
    )
  }
  return null
}

const phaseColors = ['#1677ff', '#52c41a', '#722ed1', '#fa8c16', '#f5222d', '#13c2c2']
const specPieColors = ['#1677ff', '#52c41a', '#722ed1', '#fa8c16', '#f5222d', '#13c2c2', '#eb2f96', '#fa541c', '#2f54eb', '#faad14']

type ResourceGroup = 'building' | 'extruder-cyclic' | 'extruder-pcs' | 'extruder'

function getResourceGroup(resource?: string): ResourceGroup {
  if (!resource) return 'building'
  if (resource.startsWith('recorddatacyclic')) return 'extruder-cyclic'
  if (resource.startsWith('recorddatapcs')) return 'extruder-pcs'
  if (resource.startsWith('rteex') || resource.startsWith('datalog')) return 'extruder'
  return 'building'
}

function getSPCConfig(resource?: string) {
  const group = getResourceGroup(resource)
  if (group === 'extruder-cyclic') {
    return {
      time_col: 'timestamp_record',
      actual: 'aberat_control',
      target: 'sberat_control',
      tol_pp: 'sberatctrl_tol++',
      tol_mm: 'sberatctrl_tol--',
      tol_p: 'sberatctrl_tol+',
      tol_m: 'sberatctrl_tol-',
      title: 'Weight Control (SPC)',
      yLabel: 'Weight (kg)',
    }
  }
  if (group === 'extruder-pcs') {
    return {
      time_col: 'waktu',
      actual: 'aberat_finish',
      target: 'sberat_finish',
      tol_pp: 'sberat_finish++',
      tol_mm: 'sberat_finish--',
      tol_p: 'sberat_finish+',
      tol_m: 'sberat_finish-',
      title: 'Weight Finish (SPC)',
      yLabel: 'Weight (kg)',
    }
  }
  return null
}

interface ResourceWidgetConfig {
  sortCol: string
  trendCol: string | null
  trendTimeCol: string | null
  trendTitle?: string
  distCol: string | null
  distTitle: string
  hasCT: boolean
  statOKKey: string | null
  statNGKey: string | null
  statTotKey: string | null
  hasQualityTrend?: boolean
  qualityTitle?: string
  qualityStatusCol?: string
  qualityTimeCol?: string
  extraDistCol?: string
  extraDistTitle?: string
  durationCol?: string
  hasJudgmentSummary?: boolean
}

function getWidgetConfig(resource?: string): ResourceWidgetConfig {
  const group = getResourceGroup(resource)

  if (group === 'extruder-cyclic' || group === 'extruder-pcs') {
    return { sortCol: 'id', trendCol: null, trendTimeCol: null, distCol: null, distTitle: 'Recipe', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null }
  }

  if (group === 'extruder') {
    if (resource === 'rteex1') {
      return { sortCol: '', trendCol: 'ActWS', trendTimeCol: 'trxtime', trendTitle: 'Actual WS Trend', distCol: 'machine', distTitle: 'Machine', hasCT: false, statOKKey: 'OK', statNGKey: 'NG', statTotKey: 'Total' }
    }
    if (resource === 'rteex2') {
      return { sortCol: 'recid', trendCol: 'actual', trendTimeCol: 'trxtime', trendTitle: 'Actual Weight Trend', distCol: 'mesin', distTitle: 'Machine', hasCT: false, statOKKey: 'ok', statNGKey: 'ng', statTotKey: 'total' }
    }
    if (resource === 'rteex3head') {
      return { sortCol: 'recid', trendCol: 'TempHeadExtUp', trendTimeCol: 'trxtime', trendTitle: 'Ext. Head Temp Trend', distCol: 'machine', distTitle: 'Machine', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null }
    }
    if (resource === 'datalog') {
      return { sortCol: 'id', trendCol: 'act_runscale', trendTimeCol: 'datetime', trendTitle: 'Run Scale Trend', distCol: 'recipe1', distTitle: 'Recipe', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null }
    }
    return { sortCol: 'id', trendCol: null, trendTimeCol: null, distCol: null, distTitle: 'Category', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null }
  }

  if (resource?.startsWith('rtba') || resource?.startsWith('rtbc') || resource?.startsWith('rtbe')) {
    return { sortCol: 'recid', trendCol: 'BDD_CT', trendTimeCol: 'Timestamp', trendTitle: 'BDD CT Trend', distCol: 'specification', distTitle: 'Specification', hasCT: true, statOKKey: null, statNGKey: null, statTotKey: null }
  }

  if (resource === 'curtire') {
    return { sortCol: 'recid', trendCol: null, trendTimeCol: 'eventdate', distCol: 'finaljdg', distTitle: 'Quality', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null, hasQualityTrend: true, qualityTitle: 'Quality Trend (OK/NG)', qualityStatusCol: 'finaljdg', qualityTimeCol: 'eventdate', extraDistCol: 'finaldfc', extraDistTitle: 'Defect Code' }
  }

  if (resource === 'item_measurement') {
    return { sortCol: 'recid', trendCol: null, trendTimeCol: 'eventdate', distCol: 'jdg', distTitle: 'Final Judgment', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null, hasQualityTrend: true, qualityTitle: 'Quality Trend (OK/NG)', qualityStatusCol: 'jdg', qualityTimeCol: 'eventdate', extraDistCol: 'codespec', extraDistTitle: 'Code Spec', hasJudgmentSummary: true }
  }

  if (resource === 'trimming') {
    return { sortCol: 'id', trendCol: 'Duration_Process', trendTimeCol: 'Start_Trimming', trendTitle: 'Duration Trend', distCol: 'Machine_number', distTitle: 'Machine', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null, extraDistCol: 'Tirecode', extraDistTitle: 'Tire Code', durationCol: 'Duration_Process' }
  }

  if (resource === 'rtc-tr1') {
    return { sortCol: 'recid', trendCol: 'Duration_TrimingProcess', trendTimeCol: 'Start_Triming', trendTitle: 'Duration Trend', distCol: 'Trimming_MachineNumber', distTitle: 'Machine', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null, extraDistCol: 'Tire_Code', extraDistTitle: 'Tire Code', durationCol: 'Duration_TrimingProcess' }
  }

  const idCol = resource === 'alarm_history' ? 'id' : 'recid'
  return { sortCol: idCol, trendCol: null, trendTimeCol: null, distCol: null, distTitle: 'Category', hasCT: false, statOKKey: null, statNGKey: null, statTotKey: null }
}

function getDateColumn(widgetCfg: ResourceWidgetConfig): string {
  if (widgetCfg.trendTimeCol) return widgetCfg.trendTimeCol
  if (widgetCfg.qualityTimeCol) return widgetCfg.qualityTimeCol
  return 'Timestamp'
}

function getStatValue(summary: any, keys: string[]): number | undefined {
  for (const key of keys) {
    const v = summary?.[key]
    if (v != null && v !== 0) return v
  }
  return undefined
}

const SPCTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div style={{
        background: 'var(--bg-elevated)',
        padding: '12px 16px',
        borderRadius: 12,
        boxShadow: 'var(--shadow-lg)',
        border: '1px solid var(--border-color)',
        maxHeight: 250,
        overflow: 'auto',
      }}>
        <Text strong style={{ fontSize: 12, display: 'block', marginBottom: 6 }}>{label}</Text>
        {payload.map((p: any, i: number) => (
          p.value != null && (
            <div key={i} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 12, padding: '2px 0' }}>
              <span style={{
                width: 8, height: 8, borderRadius: '50%',
                background: p.color,
                border: p.payload?.borderDash ? '1px dashed #666' : 'none',
                flexShrink: 0,
              }} />
              <span style={{ color: 'var(--text-secondary)', flex: 1 }}>{p.name}</span>
              <span style={{ fontWeight: 700 }}>{typeof p.value === 'number' ? p.value.toFixed(3) : p.value}</span>
            </div>
          )
        ))}
      </div>
    )
  }
  return null
}

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
  const [activeTab, setActiveTab] = useState('dashboard')
  const queryClient = useQueryClient()

  const group = getResourceGroup(resource)
  const spcCfg = getSPCConfig(resource)
  const widgetCfg = getWidgetConfig(resource)
  const isBuilding = group === 'building'
  const dateColumn = getDateColumn(widgetCfg)
  const [startDate, setStartDate] = useState<dayjs.Dayjs | null>(null)
  const [endDate, setEndDate] = useState<dayjs.Dayjs | null>(null)

  const { data, isLoading, isFetching, refetch } = useQuery({
    queryKey: ['resource-data', resource, page, search, searchBy, sortBy, sortDir, startDate, endDate, dateColumn],
    queryFn: () => {
      const params: Record<string, any> = { page, limit: 25, search, search_by: searchBy, sort_by: sortBy, sort_dir: sortDir }
      if (startDate) params.start_date = startDate.format('YYYY-MM-DD HH:mm:ss')
      if (endDate) {
        if (endDate.format('HH:mm:ss') === '00:00:00') {
          params.end_date = endDate.format('YYYY-MM-DD') + ' 23:59:59'
        } else {
          params.end_date = endDate.format('YYYY-MM-DD HH:mm:ss')
        }
      }
      if (dateColumn) params.date_column = dateColumn
      return api.get(`/api/v1/resources/${resource}`, { params })
    },
    enabled: activeTab === 'table',
  })

  const { data: statsData, isLoading: statsLoading } = useQuery({
    queryKey: ['resource-stats', resource, widgetCfg.durationCol],
    queryFn: () => {
      const params: Record<string, any> = {}
      if (widgetCfg.durationCol) params.duration_col = widgetCfg.durationCol
      return api.get(`/api/v1/resources/${resource}/stats`, { params })
    },
    enabled: activeTab === 'dashboard',
    refetchInterval: 30_000,
  })

  const { data: trendData, isLoading: trendLoading } = useQuery({
    queryKey: ['resource-trend', resource, widgetCfg.trendCol],
    queryFn: () => {
      if (!widgetCfg.trendCol) return Promise.resolve({ data: [] })
      const params: Record<string, string> = { column: widgetCfg.trendCol }
      if (widgetCfg.trendTimeCol) params.time_column = widgetCfg.trendTimeCol
      return api.get(`/api/v1/resources/${resource}/trend`, { params })
    },
    enabled: activeTab === 'dashboard' && !!widgetCfg.trendCol,
  })

  const { data: spcData, isLoading: spcLoading } = useQuery({
    queryKey: ['resource-spc', resource],
    queryFn: () =>
      api.get(`/api/v1/resources/${resource}/spc`, { params: spcCfg || {} }),
    enabled: activeTab === 'dashboard' && !!spcCfg,
    refetchInterval: 30_000,
  })

  const { data: distSpec, isLoading: distSpecLoading } = useQuery({
    queryKey: ['resource-dist-spec', resource, widgetCfg.distCol],
    queryFn: () => {
      if (!widgetCfg.distCol) return Promise.resolve({ data: [] })
      return api.get(`/api/v1/resources/${resource}/distribution`, { params: { column: widgetCfg.distCol } })
    },
    enabled: activeTab === 'dashboard' && !!widgetCfg.distCol,
  })

  const { data: extraDist, isLoading: extraDistLoading } = useQuery({
    queryKey: ['resource-extra-dist', resource, widgetCfg.extraDistCol],
    queryFn: () => {
      if (!widgetCfg.extraDistCol) return Promise.resolve({ data: [] })
      return api.get(`/api/v1/resources/${resource}/distribution`, { params: { column: widgetCfg.extraDistCol } })
    },
    enabled: activeTab === 'dashboard' && !!widgetCfg.extraDistCol,
  })

  const { data: qualityTrend, isLoading: qualityTrendLoading } = useQuery({
    queryKey: ['resource-quality-trend', resource],
    queryFn: () => {
      if (!widgetCfg.hasQualityTrend) return Promise.resolve({ data: [] })
      return api.get(`/api/v1/resources/${resource}/quality-trend`, {
        params: { time_col: widgetCfg.qualityTimeCol, status_col: widgetCfg.qualityStatusCol },
      })
    },
    enabled: activeTab === 'dashboard' && !!widgetCfg.hasQualityTrend,
    refetchInterval: 30_000,
  })

  const { data: judgmentData, isLoading: judgmentLoading } = useQuery({
    queryKey: ['resource-judgment', resource],
    queryFn: () => api.get(`/api/v1/resources-judgment/${resource}`),
    enabled: activeTab === 'dashboard' && !!widgetCfg.hasJudgmentSummary,
    refetchInterval: 30_000,
  })

  const { data: recentData } = useQuery({
    queryKey: ['resource-recent', resource, widgetCfg.sortCol],
    queryFn: () => {
      const params: Record<string, any> = { page: 1, limit: 8 }
      if (widgetCfg.sortCol) {
        params.sort_by = widgetCfg.sortCol
        params.sort_dir = 'desc'
      }
      return api.get(`/api/v1/resources/${resource}`, { params })
    },
    enabled: activeTab === 'dashboard',
  })

  const rows = data?.data || []
  const total = data?.meta?.total || 0
  const allColumns = rows.length > 0 ? Object.keys(rows[0]) : []
  const recentRows = recentData?.data || []
  const stats = statsData?.data
  const trendPoints = (trendData?.data || []) as any[]
  const specDist = (distSpec?.data || []) as any[]
  const extraDistData = (extraDist?.data || []) as any[]
  const qualityTrendData = (qualityTrend?.data || []) as any[]
  const judgmentSummaryData = (judgmentData?.data || []) as any[]
  const spcResult = spcData?.data

  const exportMutation = useMutation({
    mutationFn: (body: any) => api.post('/api/v1/exports', body),
    onSuccess: () => {
      message.success('Export job submitted! Cek halaman Export.')
      setExportModalOpen(false)
      queryClient.invalidateQueries({ queryKey: ['exports'] })
    },
    onError: (err: any) => {
      message.error('Gagal submit export: ' + (err.response?.data?.message || err.message))
    },
  })

  const dataColumns = allColumns.slice(0, 10).map((key) => ({
    title: key,
    dataIndex: key,
    key,
    ellipsis: true,
    width: key === 'id' || key === 'recid' ? 80 : undefined,
    sorter: true,
  }))

  const columns: any[] = [
    ...dataColumns,
    {
      title: 'Aksi',
      key: 'action',
      width: 80,
      render: (_: any, record: any) => (
        <Button
          size="small"
          icon={<EyeOutlined />}
          onClick={() => { setSelectedRow(record); setDrawerOpen(true) }}
          style={{ borderRadius: 8 }}
        />
      ),
    },
  ]

  const handleTableChange = (pagination: any, _filters: any, sorter: any) => {
    setPage(pagination.current)
    if (sorter.field) {
      setSortBy(sorter.field)
      setSortDir(sorter.order === 'ascend' ? 'asc' : 'desc')
    }
  }

  const handleExport = () => {
    const filters: Record<string, string> = {}
    if (startDate) filters.start_date = startDate.format('YYYY-MM-DD HH:mm:ss')
    if (endDate) {
      if (endDate.format('HH:mm:ss') === '00:00:00') {
        filters.end_date = endDate.format('YYYY-MM-DD') + ' 23:59:59'
      } else {
        filters.end_date = endDate.format('YYYY-MM-DD HH:mm:ss')
      }
    }
    if (dateColumn) filters.date_column = dateColumn
    exportMutation.mutate({
      resource_name: resource,
      columns: selectedColumns.length > 0 ? selectedColumns : allColumns,
      format: 'csv',
      filters: Object.keys(filters).length > 0 ? filters : undefined,
    })
  }

  const ctSummary = useMemo(() => {
    const summary = stats?.ct_summary || {}
    return [
      { name: 'PUD_CT', value: summary.PUD_CT || 0 },
      { name: 'BTD_CT', value: summary.BTD_CT || 0 },
      { name: 'BDD_CT', value: summary.BDD_CT || 0 },
    ]
  }, [stats])

  const recentCols = useMemo(() => {
    if (recentRows.length === 0) return []
    const keys = Object.keys(recentRows[0]).slice(0, 8)
    return keys.map(key => ({
      title: key,
      dataIndex: key,
      key,
      ellipsis: true,
      width: key === 'recid' || key === 'id' ? 70 : undefined,
      render: (v: any) => typeof v === 'number' ? (v % 1 === 0 ? v : v.toFixed?.(2) ?? v) : v,
    }))
  }, [recentRows])

  const statCards = useMemo(() => {
    const cards = [
      { label: 'Total Records', value: stats?.total_records ?? '-', gradient: 'gradient-blue' as const },
    ]
    if (widgetCfg.hasCT) {
      const summary = stats?.ct_summary || {}
      cards.push(
        { label: 'Avg Cycle Time', value: summary.PUD_CT != null ? `${Number(summary.PUD_CT).toFixed(1)}s` : '-', gradient: 'gradient-purple' as const },
      )
    }
    if (widgetCfg.durationCol) {
      const summary = stats?.ct_summary || {}
      if (summary.duration_avg != null && summary.duration_avg !== 0) {
        cards.push({ label: 'Avg Duration', value: `${Number(summary.duration_avg).toFixed(1)}s`, gradient: 'gradient-purple' as const })
        cards.push({ label: 'Min Duration', value: `${Number(summary.duration_min).toFixed(1)}s`, gradient: 'gradient-green' as const })
        cards.push({ label: 'Max Duration', value: `${Number(summary.duration_max).toFixed(1)}s`, gradient: 'gradient-red' as const })
      }
    }
    if (!isBuilding && !widgetCfg.durationCol) {
      const summary = stats?.ct_summary || {}
      const okVal = getStatValue(summary, ['prod_OK', 'OK', 'ok'])
      const ngVal = getStatValue(summary, ['prod_NG', 'NG', 'ng'])
      const totVal = getStatValue(summary, ['prod_Tot', 'Total', 'total'])
      if (okVal != null) cards.push({ label: 'Total OK', value: String(okVal), gradient: 'gradient-green' as const })
      if (ngVal != null) cards.push({ label: 'Total NG', value: String(ngVal), gradient: 'gradient-red' as const })
      if (totVal != null) cards.push({ label: 'Total Prod', value: String(totVal), gradient: 'gradient-orange' as const })
    }
    return cards
  }, [stats, isBuilding, widgetCfg.hasCT, widgetCfg.durationCol])

  const renderDashboard = () => {
    if (spcCfg) {
      return (
        <div style={{ marginTop: 4 }}>
          {/* Stat cards */}
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            {[
              { label: 'Total Records', value: stats?.total_records ?? '-', gradient: 'gradient-blue' },
              { label: 'Avg Actual', value: spcResult?.datasets?.[0]?.data ? 
                (spcResult.datasets[0].data.filter((v: any) => v != null).reduce((a: number, b: number) => a + b, 0) / 
                 Math.max(spcResult.datasets[0].data.filter((v: any) => v != null).length, 1)).toFixed(3) : '-', 
                gradient: 'gradient-green' },
              { label: 'Latest Actual', value: spcResult?.datasets?.[0]?.data ? 
                (spcResult.datasets[0].data[spcResult.datasets[0].data.length - 1] ?? '-').toFixed?.(3) ?? '-' : '-', 
                gradient: 'gradient-purple' },
              { label: 'Deviation', value: spcResult?.datasets?.[0]?.data && spcResult?.datasets?.[1]?.data ? 
                ((spcResult.datasets[0].data[spcResult.datasets[0].data.length - 1] ?? 0) - 
                 (spcResult.datasets[1].data[spcResult.datasets[1].data.length - 1] ?? 0)).toFixed(3) : '-', 
                gradient: spcResult?.datasets?.[0]?.data && spcResult?.datasets?.[1]?.data &&
                  Math.abs((spcResult.datasets[0].data[spcResult.datasets[0].data.length - 1] ?? 0) - 
                    (spcResult.datasets[1].data[spcResult.datasets[1].data.length - 1] ?? 0)) > 0.1
                  ? 'gradient-red' : 'gradient-orange' },
            ].map((card, i) => (
              <Col xs={24} sm={12} lg={6} key={i}>
                <Card className={`stat-card ${card.gradient}`} bordered={false}>
                  <div style={{ position: 'relative', zIndex: 1 }}>
                    <div className="stat-value">{card.value}</div>
                    <div className="stat-label">{card.label}</div>
                  </div>
                </Card>
              </Col>
            ))}
          </Row>

          {/* SPC Chart - Full width */}
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            <Col xs={24}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#10b981' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>{spcCfg.title}</span>
                  </Space>
                }
                styles={{ body: { padding: '20px 16px 8px' } }}
              >
                <Spin spinning={spcLoading && !spcResult}>
                  {spcResult?.labels?.length > 0 ? (
                    <ResponsiveContainer width="100%" height={360}>
                      <LineChart data={(() => {
                        const ds = spcResult.datasets
                        return spcResult.labels.map((l: string, i: number) => ({
                          label: l,
                          'Actual Weight': ds[0]?.data?.[i] ?? null,
                          'Set Weight': ds[1]?.data?.[i] ?? null,
                          'USL': ds[2]?.data?.[i] ?? null,
                          'LSL': ds[3]?.data?.[i] ?? null,
                          'UCL (Warning)': ds[4]?.data?.[i] ?? null,
                          'LCL (Warning)': ds[5]?.data?.[i] ?? null,
                        }))
                      })()} margin={{ top: 8, right: 24, bottom: 0, left: 0 }}
                      >
                        <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                        <XAxis
                          dataKey="label"
                          tick={{ fill: 'var(--text-secondary)', fontSize: 11 }}
                          axisLine={false}
                          tickLine={false}
                          tickFormatter={(v) => {
                            if (!v) return v
                            const parts = v.split(' ')
                            return parts[parts.length - 1]?.substring(0, 8) || v
                          }}
                        />
                        <YAxis tick={{ fill: 'var(--text-secondary)', fontSize: 11 }} axisLine={false} tickLine={false} />
                        <Tooltip content={<SPCTooltip />} />
                        <Legend
                          wrapperStyle={{ fontSize: 12, paddingTop: 12 }}
                        />
                        <Line type="monotone" dataKey="Actual Weight" stroke="#10b981" strokeWidth={2} dot={false} activeDot={{ r: 4 }} connectNulls />
                        <Line type="monotone" dataKey="Set Weight" stroke="#9ca3af" strokeWidth={1.5} strokeDasharray="5 5" dot={false} connectNulls />
                        <Line type="monotone" dataKey="USL" stroke="#ef4444" strokeWidth={1} strokeDasharray="6 3" dot={false} connectNulls />
                        <Line type="monotone" dataKey="LSL" stroke="#ef4444" strokeWidth={1} strokeDasharray="6 3" dot={false} connectNulls />
                        <Line type="monotone" dataKey="UCL (Warning)" stroke="#f59e0b" strokeWidth={1} strokeDasharray="2 2" dot={false} connectNulls />
                        <Line type="monotone" dataKey="LCL (Warning)" stroke="#f59e0b" strokeWidth={1} strokeDasharray="2 2" dot={false} connectNulls />
                      </LineChart>
                    </ResponsiveContainer>
                  ) : (
                    <div style={{ textAlign: 'center', padding: 60, color: 'var(--text-tertiary)' }}>
                      {spcLoading ? 'Loading...' : 'No SPC data available'}
                    </div>
                  )}
                </Spin>
              </Card>
            </Col>
          </Row>

          {/* Second Row */}
          <Row gutter={[16, 16]}>
            {widgetCfg.distCol && (
              <Col xs={24} lg={8}>
                <Card
                  className="modern-card"
                  title={
                    <Space>
                      <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#52c41a' }} />
                      <span style={{ fontWeight: 600, fontSize: 14 }}>{widgetCfg.distTitle}</span>
                    </Space>
                  }
                  styles={{ body: { padding: '16px' } }}
                >
                  <Spin spinning={distSpecLoading}>
                    {specDist.length > 0 ? (
                      <ResponsiveContainer width="100%" height={240}>
                        <PieChart>
                          <Pie
                            data={specDist}
                            cx="50%" cy="50%"
                            innerRadius={55} outerRadius={90}
                            dataKey="value" nameKey="label"
                            paddingAngle={2}
                          >
                            {specDist.map((_: any, idx: number) => (
                              <Cell key={idx} fill={specPieColors[idx % specPieColors.length]} stroke="none" />
                            ))}
                          </Pie>
                          <Tooltip content={<CustomTooltip />} />
                        </PieChart>
                      </ResponsiveContainer>
                    ) : (
                      <div style={{ textAlign: 'center', padding: 40, color: 'var(--text-tertiary)' }}>No data</div>
                    )}
                    {specDist.length > 0 && (
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center', marginTop: 8 }}>
                        {specDist.slice(0, 6).map((item: any, idx: number) => (
                          <div key={item.label} style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
                            <span style={{ width: 8, height: 8, borderRadius: '50%', background: specPieColors[idx % specPieColors.length] }} />
                            <span style={{ color: 'var(--text-secondary)' }}>{item.label}</span>
                            <span style={{ fontWeight: 600 }}>{item.value}</span>
                          </div>
                        ))}
                      </div>
                    )}
                  </Spin>
                </Card>
              </Col>
            )}
            <Col xs={24} {...(widgetCfg.distCol ? { lg: 16 } : { lg: 24 })}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#13c2c2' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>Recent Records</span>
                  </Space>
                }
                styles={{ body: { padding: 0 } }}
              >
                <Table
                  dataSource={recentRows}
                  columns={recentCols}
                  rowKey={(r: any) => r.recid || r.id || Math.random()}
                  pagination={false}
                  size="small"
                  scroll={{ x: 'max-content' }}
                  locale={{ emptyText: 'No recent records' }}
                />
              </Card>
            </Col>
          </Row>
        </div>
      )
    }

    return (
      <div style={{ marginTop: 4 }}>
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          {statCards.map((card, i) => (
            <Col xs={24} sm={12} lg={6} key={i}>
              <Card className={`stat-card ${card.gradient}`} bordered={false}>
                <div style={{ position: 'relative', zIndex: 1 }}>
                  <div className="stat-value">{card.value}</div>
                  <div className="stat-label">{card.label}</div>
                </div>
              </Card>
            </Col>
          ))}
        </Row>

        {(widgetCfg.hasCT || widgetCfg.trendCol) && (
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            {widgetCfg.hasCT && (
              <Col xs={24} {...(widgetCfg.trendCol ? { lg: 12 } : { lg: 24 })}>
                <Card
                  className="modern-card"
                  title={
                    <Space>
                      <div style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--primary-color)' }} />
                      <span style={{ fontWeight: 600, fontSize: 14 }}>Cycle Time per Phase</span>
                    </Space>
                  }
                  styles={{ body: { padding: '16px 16px 8px' } }}
                >
                  <ResponsiveContainer width="100%" height={260}>
                    <BarChart data={ctSummary} margin={{ top: 8, right: 8, bottom: 0, left: -16 }}>
                      <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                      <XAxis dataKey="name" tick={{ fill: 'var(--text-secondary)', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis tick={{ fill: 'var(--text-secondary)', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <Tooltip content={<CustomTooltip />} cursor={{ fill: 'var(--bg-hover)' }} />
                      <Bar dataKey="value" radius={[6, 6, 0, 0]} maxBarSize={48}>
                        {ctSummary.map((_: any, idx: number) => (
                          <Cell key={idx} fill={phaseColors[idx % phaseColors.length]} />
                        ))}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                </Card>
              </Col>
            )}
            {widgetCfg.trendCol && (
              <Col xs={24} {...(widgetCfg.hasCT ? { lg: 12 } : { lg: 24 })}>
                <Card
                  className="modern-card"
                  title={
                    <Space>
                      <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#722ed1' }} />
                      <span style={{ fontWeight: 600, fontSize: 14 }}>{widgetCfg.trendTitle || widgetCfg.trendCol}</span>
                    </Space>
                  }
                  styles={{ body: { padding: '16px 16px 8px' } }}
                >
                  <Spin spinning={trendLoading && trendPoints.length === 0}>
                    <ResponsiveContainer width="100%" height={260}>
                      <AreaChart data={trendPoints.slice(-24)} margin={{ top: 8, right: 8, bottom: 0, left: -16 }}>
                        <defs>
                          <linearGradient id="trendGrad" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#722ed1" stopOpacity={0.3} />
                            <stop offset="95%" stopColor="#722ed1" stopOpacity={0} />
                          </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                        <XAxis
                          dataKey="label"
                          tick={{ fill: 'var(--text-secondary)', fontSize: 11 }}
                          axisLine={false} tickLine={false}
                          tickFormatter={(v) => v?.split(' ')?.[1]?.substring(0, 5) ?? v}
                        />
                        <YAxis tick={{ fill: 'var(--text-secondary)', fontSize: 12 }} axisLine={false} tickLine={false} />
                        <Tooltip content={<CustomTooltip />} cursor={{ fill: 'var(--bg-hover)' }} />
                        <Area type="monotone" dataKey="value" stroke="#722ed1" fill="url(#trendGrad)" strokeWidth={2} dot={false} activeDot={{ r: 4 }} />
                      </AreaChart>
                    </ResponsiveContainer>
                  </Spin>
                </Card>
              </Col>
            )}
          </Row>
        )}

        {widgetCfg.hasQualityTrend && (
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            <Col xs={24}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#10b981' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>{widgetCfg.qualityTitle || 'Quality Trend'}</span>
                  </Space>
                }
                styles={{ body: { padding: '20px 16px 8px' } }}
              >
                <Spin spinning={qualityTrendLoading && qualityTrendData.length === 0}>
                  {qualityTrendData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <AreaChart data={qualityTrendData} margin={{ top: 8, right: 16, bottom: 0, left: -8 }}>
                        <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" vertical={false} />
                        <XAxis
                          dataKey="label"
                          tick={{ fill: 'var(--text-secondary)', fontSize: 11 }}
                          axisLine={false} tickLine={false}
                          tickFormatter={(v) => v?.split(' ')?.[1]?.substring(0, 5) ?? v}
                        />
                        <YAxis tick={{ fill: 'var(--text-secondary)', fontSize: 11 }} axisLine={false} tickLine={false} />
                        <Tooltip content={<QualityTooltip />} />
                        <Legend wrapperStyle={{ fontSize: 12, paddingTop: 8 }} />
                        <Area type="monotone" dataKey="ok" name="OK" stroke="#10b981" fill="#10b981" fillOpacity={0.25} strokeWidth={2} dot={false} stackId="1" />
                        <Area type="monotone" dataKey="ng" name="NG" stroke="#ef4444" fill="#ef4444" fillOpacity={0.25} strokeWidth={2} dot={false} stackId="1" />
                      </AreaChart>
                    </ResponsiveContainer>
                  ) : (
                    <div style={{ textAlign: 'center', padding: 60, color: 'var(--text-tertiary)' }}>
                      {qualityTrendLoading ? 'Loading...' : 'No quality data available'}
                    </div>
                  )}
                </Spin>
              </Card>
            </Col>
          </Row>
        )}

        {widgetCfg.hasJudgmentSummary && (
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            <Col xs={24} lg={12}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#f59e0b' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>Parameter Failure Rate (NG)</span>
                  </Space>
                }
                styles={{ body: { padding: '16px 16px 8px' } }}
              >
                <Spin spinning={judgmentLoading && judgmentSummaryData.length === 0}>
                  {judgmentSummaryData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={judgmentSummaryData} layout="vertical" margin={{ top: 8, right: 24, bottom: 0, left: 0 }} barCategoryGap={8}>
                        <CartesianGrid strokeDasharray="4 4" stroke="var(--border-color)" horizontal={false} />
                        <XAxis type="number" tick={{ fill: 'var(--text-secondary)', fontSize: 11 }} axisLine={false} tickLine={false} />
                        <YAxis type="category" dataKey="parameter" tick={{ fill: 'var(--text-secondary)', fontSize: 12 }} axisLine={false} tickLine={false} width={80} />
                        <Tooltip content={<CustomTooltip />} cursor={{ fill: 'var(--bg-hover)' }} />
                        <Legend wrapperStyle={{ fontSize: 12, paddingTop: 8 }} />
                        <Bar dataKey="left_ng" name="Left NG" stackId="a" radius={[0, 4, 4, 0]} fill="#1677ff" />
                        <Bar dataKey="right_ng" name="Right NG" stackId="a" radius={[0, 4, 4, 0]} fill="#f5222d" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: 'var(--text-tertiary)' }}>
                      {judgmentLoading ? 'Loading...' : 'No data'}
                    </div>
                  )}
                </Spin>
              </Card>
            </Col>
            <Col xs={24} lg={12}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#10b981' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>NG Summary</span>
                  </Space>
                }
                styles={{ body: { padding: '16px' } }}
              >
                <Spin spinning={judgmentLoading && judgmentSummaryData.length === 0}>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                    {judgmentSummaryData.map((item: any) => (
                      <div key={item.parameter} style={{
                        display: 'flex', alignItems: 'center', gap: 12,
                        padding: '10px 14px',
                        borderRadius: 10,
                        background: 'var(--bg-elevated)',
                        border: '1px solid var(--border-color)',
                      }}>
                        <div style={{ flex: 1, fontWeight: 600, fontSize: 13, textTransform: 'capitalize' }}>{item.parameter}</div>
                        <div style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
                          Left: <span style={{ fontWeight: 700, color: '#1677ff' }}>{item.left_ng}</span>
                        </div>
                        <div style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
                          Right: <span style={{ fontWeight: 700, color: '#f5222d' }}>{item.right_ng}</span>
                        </div>
                        <div style={{
                          fontSize: 13, fontWeight: 700, padding: '2px 12px',
                          borderRadius: 20,
                          background: item.total_ng > 0 ? '#fef2f2' : '#f0fdf4',
                          color: item.total_ng > 0 ? '#ef4444' : '#10b981',
                        }}>
                          {item.total_ng}
                        </div>
                      </div>
                    ))}
                  </div>
                </Spin>
              </Card>
            </Col>
          </Row>
        )}

        <Row gutter={[16, 16]}>
          {widgetCfg.distCol && (
            <Col xs={24} lg={widgetCfg.extraDistCol ? 8 : 8}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#52c41a' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>{widgetCfg.distTitle}</span>
                  </Space>
                }
                styles={{ body: { padding: '16px' } }}
              >
                <Spin spinning={distSpecLoading}>
                  {specDist.length > 0 ? (
                    <ResponsiveContainer width="100%" height={240}>
                      <PieChart>
                        <Pie
                          data={specDist}
                          cx="50%" cy="50%" innerRadius={55} outerRadius={90}
                          dataKey="value" nameKey="label" paddingAngle={2}
                        >
                          {specDist.map((_: any, idx: number) => (
                            <Cell key={idx} fill={specPieColors[idx % specPieColors.length]} stroke="none" />
                          ))}
                        </Pie>
                        <Tooltip content={<CustomTooltip />} />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: 'var(--text-tertiary)' }}>No data</div>
                  )}
                  {specDist.length > 0 && (
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center', marginTop: 8 }}>
                      {specDist.slice(0, 6).map((item: any, idx: number) => (
                        <div key={item.label} style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
                          <span style={{ width: 8, height: 8, borderRadius: '50%', background: specPieColors[idx % specPieColors.length] }} />
                          <span style={{ color: 'var(--text-secondary)' }}>{item.label}</span>
                          <span style={{ fontWeight: 600 }}>{item.value}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </Spin>
              </Card>
            </Col>
          )}
          {widgetCfg.extraDistCol && (
            <Col xs={24} lg={8}>
              <Card
                className="modern-card"
                title={
                  <Space>
                    <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#f59e0b' }} />
                    <span style={{ fontWeight: 600, fontSize: 14 }}>{widgetCfg.extraDistTitle}</span>
                  </Space>
                }
                styles={{ body: { padding: '16px' } }}
              >
                <Spin spinning={extraDistLoading}>
                  {extraDistData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={240}>
                      <PieChart>
                        <Pie
                          data={extraDistData}
                          cx="50%" cy="50%" innerRadius={55} outerRadius={90}
                          dataKey="value" nameKey="label" paddingAngle={2}
                        >
                          {extraDistData.map((_: any, idx: number) => (
                            <Cell key={idx} fill={specPieColors[idx % specPieColors.length]} stroke="none" />
                          ))}
                        </Pie>
                        <Tooltip content={<CustomTooltip />} />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: 'var(--text-tertiary)' }}>No data</div>
                  )}
                  {extraDistData.length > 0 && (
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center', marginTop: 8 }}>
                      {extraDistData.slice(0, 6).map((item: any, idx: number) => (
                        <div key={item.label} style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
                          <span style={{ width: 8, height: 8, borderRadius: '50%', background: specPieColors[idx % specPieColors.length] }} />
                          <span style={{ color: 'var(--text-secondary)' }}>{item.label}</span>
                          <span style={{ fontWeight: 600 }}>{item.value}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </Spin>
              </Card>
            </Col>
          )}
          <Col xs={24} {...(widgetCfg.distCol && widgetCfg.extraDistCol ? { lg: 8 } : widgetCfg.distCol || widgetCfg.extraDistCol ? { lg: 16 } : { lg: 24 })}>
            <Card
              className="modern-card"
              title={
                <Space>
                  <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#13c2c2' }} />
                  <span style={{ fontWeight: 600, fontSize: 14 }}>Recent Records</span>
                </Space>
              }
              styles={{ body: { padding: 0 } }}
            >
              <Table
                dataSource={recentRows}
                columns={recentCols}
                rowKey={(r: any) => r.recid || r.id || Math.random()}
                pagination={false}
                size="small"
                scroll={{ x: 'max-content' }}
                locale={{ emptyText: 'No recent records' }}
              />
            </Card>
          </Col>
        </Row>
      </div>
    )
  }

  return (
    <div className="page-enter">
      {/* Header */}
      <div style={{
        display: 'flex', alignItems: 'center', gap: 16,
        marginBottom: 20, paddingBottom: 16,
        borderBottom: '1px solid var(--border-color)', flexWrap: 'wrap',
      }}>
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/data')}
          style={{ borderRadius: 10, display: 'flex', alignItems: 'center', gap: 4 }}
        >
          Kembali
        </Button>
        <div style={{
          width: 36, height: 36, borderRadius: 10,
          background: 'var(--primary-bg)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          color: 'var(--primary-color)', fontSize: 18,
        }}>
          <TableOutlined />
        </div>
        <div style={{ flex: 1, minWidth: 120 }}>
          <Title level={4} style={{ margin: 0, fontSize: 18 }}>{resource?.toUpperCase()}</Title>
          {stats && (
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              {stats.total_records} total records
            </Typography.Text>
          )}
        </div>
        <Tag color="blue" style={{ fontSize: 13, padding: '4px 14px', fontWeight: 600 }}>
          {stats?.total_records ?? '-'} records
        </Tag>
        <Button
          icon={<DownloadOutlined />}
          onClick={() => { setSelectedColumns(allColumns); setExportModalOpen(true) }}
          style={{ borderRadius: 10 }}
        >
          Export CSV
        </Button>
      </div>

      {/* Tabs */}
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        style={{ marginBottom: 0 }}
        size="large"
        items={[
          {
            key: 'dashboard',
            label: <span><DashboardOutlined style={{ marginRight: 6 }} />Dashboard</span>,
            children: (
              <Spin spinning={statsLoading && !statsData}>
                {renderDashboard()}
              </Spin>
            ),
          },
          {
            key: 'table',
            label: <span><TableOutlined style={{ marginRight: 6 }} />Data Table</span>,
            children: (
              <Card className="modern-card" bodyStyle={{ padding: 0 }}>
                <div style={{
                  display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                  padding: '12px 16px', borderBottom: '1px solid var(--border-color)', flexWrap: 'wrap', gap: 8,
                }}>
                  <Space>
                    <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>All Records</span>
                  </Space>
                  <Space wrap>
                    <DatePicker
                      showTime={{ format: 'HH:mm:ss', defaultValue: dayjs('00:00:00', 'HH:mm:ss') }}
                      format="YYYY-MM-DD HH:mm:ss"
                      value={startDate}
                      onChange={(d) => setStartDate(d)}
                      style={{ borderRadius: 10, width: 220 }}
                      placeholder="Dari Tanggal & Jam"
                    />
                    <span style={{ color: 'var(--text-tertiary)' }}>—</span>
                    <DatePicker
                      showTime={{ format: 'HH:mm:ss', defaultValue: dayjs('23:59:59', 'HH:mm:ss') }}
                      format="YYYY-MM-DD HH:mm:ss"
                      value={endDate}
                      onChange={(d) => setEndDate(d)}
                      style={{ borderRadius: 10, width: 220 }}
                      placeholder="Sampai Tanggal & Jam"
                    />
                    <Select
                      allowClear placeholder="Cari di kolom"
                      style={{ width: 150 }}
                      value={searchBy || undefined}
                      onChange={(v) => setSearchBy(v || '')}
                      options={allColumns.map((c) => ({ value: c, label: c }))}
                    />
                    <Input
                      placeholder="Cari..."
                      prefix={<SearchOutlined style={{ color: 'var(--text-tertiary)' }} />}
                      style={{ width: 200, borderRadius: 10 }}
                      value={search}
                      onChange={(e) => setSearch(e.target.value)}
                      onPressEnter={() => refetch()}
                    />
                    <Button icon={<ReloadOutlined />} onClick={() => refetch()} style={{ borderRadius: 10 }} />
                  </Space>
                </div>
                <Table
                  dataSource={rows}
                  columns={columns}
                  rowKey={(record: any) => record.id || record.recid || Math.random()}
                  loading={isLoading || isFetching}
                  onChange={handleTableChange}
                  pagination={{ current: page, pageSize: 25, total, showSizeChanger: false, showTotal: (t) => `Total ${t} records` }}
                  scroll={{ x: 'max-content' }}
                  size="small"
                />
              </Card>
            ),
          },
        ]}
      />

      <Drawer
        title={<Space><EyeOutlined style={{ color: 'var(--primary-color)' }} /><span style={{ fontWeight: 600 }}>Detail Record</span></Space>}
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
        title={<Space><DownloadOutlined style={{ color: 'var(--primary-color)' }} /><span style={{ fontWeight: 600 }}>Export CSV — Pilih Kolom</span></Space>}
        open={exportModalOpen}
        onCancel={() => setExportModalOpen(false)}
        onOk={handleExport}
        confirmLoading={exportMutation.isPending}
        width={600}
        okText="Export"
        cancelText="Batal"
      >
        <div style={{ marginBottom: 12 }}>
          <Typography.Text type="secondary">Pilih kolom yang ingin diexport ke CSV</Typography.Text>
        </div>
        <Transfer
          dataSource={allColumns.map((c) => ({ key: c, title: c }))}
          titles={['Tersedia', 'Dipilih']}
          targetKeys={selectedColumns}
          onChange={(keys) => setSelectedColumns(keys as string[])}
          render={(item) => item.title}
          listStyle={{ width: 250, height: 400 }}
        />
      </Modal>
    </div>
  )
}
