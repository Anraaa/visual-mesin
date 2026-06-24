import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import { HomeOutlined } from '@ant-design/icons'

export default function NotFound() {
  const navigate = useNavigate()

  return (
    <div className="page-enter" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '60vh' }}>
      <Result
        status="404"
        title={<span style={{ fontSize: 72, fontWeight: 800, color: 'var(--primary-color)', lineHeight: 1 }}>404</span>}
        subTitle="Halaman yang Anda cari tidak ditemukan"
        extra={
          <Button
            type="primary"
            icon={<HomeOutlined />}
            onClick={() => navigate('/dashboard')}
            size="large"
            className="btn-glow"
            style={{ borderRadius: 10, height: 44, paddingInline: 28 }}
          >
            Kembali ke Dashboard
          </Button>
        }
      />
    </div>
  )
}
