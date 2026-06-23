import { Card, Col, Row, Statistic, Typography } from 'antd'
import { CheckCircleOutlined, CloseCircleOutlined, DashboardOutlined } from '@ant-design/icons'

const { Title } = Typography

export default function Dashboard() {
  return (
    <div>
      <Title level={4}>Dashboard</Title>
      <Row gutter={[16, 16]}>
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
            <Statistic title="Yield Rate" value={94.4} suffix="%" precision={1} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
