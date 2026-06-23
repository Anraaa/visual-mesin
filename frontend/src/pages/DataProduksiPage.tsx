import { useNavigate } from 'react-router-dom'
import { Card, Col, Row, Typography, Tag } from 'antd'
import {
  ToolOutlined, ExperimentOutlined, FireOutlined, ScissorOutlined,
  DashboardOutlined, AlertOutlined, FileTextOutlined, BuildOutlined,
} from '@ant-design/icons'

const { Title } = Typography

const resourceGroups = [
  {
    label: 'Building (RTBA)',
    resources: ['rtba1', 'rtba2', 'rtba3'],
    icon: <BuildOutlined style={{ fontSize: 24 }} />,
    color: '#1677ff',
  },
  {
    label: 'Building Quality (RTBC)',
    resources: ['rtbc1', 'rtbc2', 'rtbc3', 'rtbc4'],
    icon: <ExperimentOutlined style={{ fontSize: 24 }} />,
    color: '#52c41a',
  },
  {
    label: 'Building Evaluation (RTBE)',
    resources: ['rtbe1', 'rtbe2'],
    icon: <DashboardOutlined style={{ fontSize: 24 }} />,
    color: '#722ed1',
  },
  {
    label: 'Extruder',
    resources: ['rteex1', 'rteex2', 'rteex3head', 'recorddatacyclic', 'recorddatapcs', 'datalog'],
    icon: <ToolOutlined style={{ fontSize: 24 }} />,
    color: '#fa8c16',
  },
  {
    label: 'Curing',
    resources: ['curtire', 'item_measurement'],
    icon: <FireOutlined style={{ fontSize: 24 }} />,
    color: '#f5222d',
  },
  {
    label: 'Trimming',
    resources: ['trimming', 'rtc-tr1'],
    icon: <ScissorOutlined style={{ fontSize: 24 }} />,
    color: '#eb2f96',
  },
  {
    label: 'Monitoring & Yield',
    resources: ['monitoringtl1', 'rtl-tl1', 'rtltl1'],
    icon: <DashboardOutlined style={{ fontSize: 24 }} />,
    color: '#13c2c2',
  },
  {
    label: 'Alarm & Material',
    resources: ['alarm_history', 'material'],
    icon: <AlertOutlined style={{ fontSize: 24 }} />,
    color: '#fa541c',
  },
  {
    label: 'Recipe & Order',
    resources: ['recipe1', 'recipe1queue', 'recipe_history', 'order_report', 'batch_report'],
    icon: <FileTextOutlined style={{ fontSize: 24 }} />,
    color: '#2f54eb',
  },
  {
    label: 'Supporting',
    resources: ['mastermcn', 'bpbl', 'rsc_pc1', 'gtentire'],
    icon: <BuildOutlined style={{ fontSize: 24 }} />,
    color: '#faad14',
  },
]

export default function DataProduksiPage() {
  const navigate = useNavigate()

  return (
    <div>
      <Title level={4}>Data Produksi</Title>
      <p style={{ marginBottom: 24, color: '#666' }}>
        Pilih tabel resource untuk melihat data produksi
      </p>
      {resourceGroups.map((group) => (
        <div key={group.label} style={{ marginBottom: 24 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
            <span style={{ color: group.color }}>{group.icon}</span>
            <Title level={5} style={{ margin: 0 }}>{group.label}</Title>
          </div>
          <Row gutter={[12, 12]}>
            {group.resources.map((resource) => (
              <Col key={resource}>
                <Card
                  hoverable
                  size="small"
                  style={{ width: 180, cursor: 'pointer' }}
                  onClick={() => navigate(`/data/${resource}`)}
                >
                  <Tag color={group.color} style={{ margin: 0 }}>{resource}</Tag>
                </Card>
              </Col>
            ))}
          </Row>
        </div>
      ))}
    </div>
  )
}
