import { useEffect, useRef, useCallback } from 'react'
import { useAuthStore } from '../stores/authStore'

type MessageHandler = (data: any) => void

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'

export function useWebSocket(handlers: Record<string, MessageHandler>) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectRef = useRef<number>(0)
  const token = useAuthStore((s) => s.token)

  const connect = useCallback(() => {
    if (!token) return

    const url = `${WS_URL}?token=${token}`
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => {
      reconnectRef.current = 0
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        const handler = handlers[msg.type]
        if (handler) {
          handler(msg.payload)
        }
      } catch {
        // ignore parse errors
      }
    }

    ws.onclose = () => {
      wsRef.current = null
      const delay = Math.min(1000 * 2 ** reconnectRef.current, 30000)
      reconnectRef.current++
      setTimeout(connect, delay)
    }

    ws.onerror = () => {
      ws.close()
    }
  }, [token, handlers])

  useEffect(() => {
    connect()
    return () => {
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [connect])

  const send = useCallback((type: string, payload: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, payload }))
    }
  }, [])

  return { send }
}
