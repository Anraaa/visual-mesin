import { Component } from 'react'
import type { ErrorInfo, ReactNode } from 'react'
import { Button, Result } from 'antd'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('ErrorBoundary caught:', error, info)
  }

  render() {
    if (this.state.hasError) {
      return (
        <Result
          status="error"
          title="Terjadi Kesalahan"
          subTitle={this.state.error?.message || 'Terjadi error yang tidak terduga'}
          extra={
            <Button type="primary" onClick={() => window.location.reload()}>
              Muat Ulang
            </Button>
          }
        />
      )
    }

    return this.props.children
  }
}
