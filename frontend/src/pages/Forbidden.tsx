import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'

export default function Forbidden() {
  const navigate = useNavigate()

  return (
    <Result
      status="403"
      title="403"
      subTitle="Anda tidak memiliki akses ke halaman ini"
      extra={
        <Button type="primary" onClick={() => navigate('/dashboard')}>
          Kembali ke Dashboard
        </Button>
      }
    />
  )
}
