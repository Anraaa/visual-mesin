import { Card, Col, Row, Statistic, Typography } from 'antd'
import {
  CheckCircleOutlined, CloseCircleOutlined, DashboardOutlined,
  RiseOutlined,
} from '@ant-design/icons'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts'

const { Title } = Typography

const productionData = [
  { name: 'RTBA1', OK: 180, NG: 12 },
  { name: 'RTBA2', OK: 165, NG: 8 },
  { name: 'RTBA3', OK: 195, NG: 15 },
  { name: 'Extruder', OK: 220, NG: 5 },
  { name: 'Curing', OK: 190, NG: 10 },
  { name: 'Trimming', OK: 210, NG: 7 },
]

const pieData = [
  { name: 'OK', value: 85 },
  { name: 'NG', value: 15 },
]

const COLORS = ['#52c41a', '#ff4d4f']

export default function Dashboard() {
  return (
    <div>
      <Title level={4}>Dashboard Produksi</Title>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="Total Produksi Hari Ini" value={1250} prefix={<DashboardOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="OK" value={1180} prefix={<CheckCircleOutlined />} valueStyle={{ color: '#52c41a' }} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="NG" value={70} prefix={<CloseCircleOutlined />} valueStyle={{ color: '#ff4d4f' }} />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic title="Yield Rate" value={94.4} suffix="%" precision={1} prefix={<RiseOutlined />} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="Produksi per Mesin (OK/NG)">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={productionData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="OK" fill="#52c41a" radius={[4, 4, 0, 0]} />
                <Bar dataKey="NG" fill="#ff4d4f" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="Rasio OK / NG">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  dataKey="value"
                  label={({ name, percent }: { name: string; percent?: number }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
                >
                  {pieData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index]} />
                  ))}
                </Pie>
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
