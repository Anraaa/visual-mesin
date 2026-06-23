import { useState, useRef, useEffect, useCallback } from 'react'
import {
  Layout, Input, Button, Typography, Space, List, Tag, Spin,
  Card, Collapse, Tooltip, Empty, Popconfirm, Avatar,
} from 'antd'
import {
  SendOutlined, RobotOutlined, UserOutlined, DeleteOutlined,
  PlusOutlined, ClockCircleOutlined, CodeOutlined,
  CheckCircleOutlined, CloseCircleOutlined,
} from '@ant-design/icons'
import type { ChatMessage, ChatSession } from '../services/chat'
import { chatApi } from '../services/chat'

const { Text, Title, Paragraph } = Typography
const { TextArea } = Input
const { Sider, Content } = Layout
const { Panel } = Collapse

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  detected_intent?: string
  generated_sql?: string
  sql_status?: string
  query_result?: ChatMessage['query_result']
  timestamp: Date
  status: 'sending' | 'sent' | 'error'
}

export default function AIChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [sending, setSending] = useState(false)
  const [sessionID, setSessionID] = useState<string>('')
  const [sessions, setSessions] = useState<ChatSession[]>([])
  const [loadingSessions, setLoadingSessions] = useState(true)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [])

  useEffect(() => {
    scrollToBottom()
  }, [messages, scrollToBottom])

  useEffect(() => {
    loadSessions()
  }, [])

  const loadSessions = async () => {
    try {
      setLoadingSessions(true)
      const data = await chatApi.getSessions()
      setSessions(data)
    } catch {
      // ignore
    } finally {
      setLoadingSessions(false)
    }
  }

  const loadSession = async (sid: string) => {
    try {
      setSessionID(sid)
      const res = await chatApi.getHistory(sid)
      const items = res.data as ChatMessage[]

      const msgs: Message[] = []
      for (const item of items) {
        msgs.push({
          id: `q-${item.id}`,
          role: 'user',
          content: item.question,
          timestamp: new Date(item.created_at),
          status: 'sent',
        })
        if (item.status === 'completed' && item.answer) {
          msgs.push({
            id: `a-${item.id}`,
            role: 'assistant',
            content: item.answer,
            detected_intent: item.detected_intent,
            generated_sql: item.generated_sql,
            sql_status: item.sql_status,
            query_result: item.query_result,
            timestamp: new Date(item.created_at),
            status: 'sent',
          })
        }
      }
      setMessages(msgs)
    } catch {
      // ignore
    }
  }

  const newChat = () => {
    setSessionID('')
    setMessages([])
    setInput('')
    inputRef.current?.focus()
  }

  const handleSend = async () => {
    const text = input.trim()
    if (!text || sending) return

    setInput('')
    const userMsg: Message = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: text,
      timestamp: new Date(),
      status: 'sending',
    }
    setMessages((prev) => [...prev, userMsg])
    setSending(true)

    try {
      const result = await chatApi.send(text, sessionID)

      if (!sessionID) {
        setSessionID(result.session_id)
      }

      setMessages((prev) => {
        const updated = [...prev]
        const last = updated[updated.length - 1]
        if (last.role === 'user') {
          last.status = 'sent'
        }
        updated.push({
          id: `ai-${Date.now()}`,
          role: 'assistant',
          content: result.answer,
          detected_intent: result.detected_intent,
          generated_sql: result.generated_sql,
          sql_status: result.sql_status,
          query_result: result.query_result,
          timestamp: new Date(),
          status: 'sent',
        })
        return updated
      })

      loadSessions()
    } catch {
      setMessages((prev) => {
        const updated = [...prev]
        const last = updated[updated.length - 1]
        if (last.role === 'user') {
          last.status = 'error'
        }
        return updated
      })
    } finally {
      setSending(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleDeleteSession = async (sid: string) => {
    try {
      await chatApi.deleteSession(sid)
      if (sessionID === sid) {
        newChat()
      }
      loadSessions()
    } catch {
      // ignore
    }
  }

  const sqlStatusIcon = (status?: string) => {
    switch (status) {
      case 'valid': return <CheckCircleOutlined style={{ color: '#52c41a' }} />
      case 'invalid':
      case 'rejected': return <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
      default: return null
    }
  }

  return (
    <Layout style={{ background: 'transparent', height: 'calc(100vh - 140px)' }}>
      <Sider
        width={280}
        style={{ background: 'transparent', borderRight: '1px solid #f0f0f0', overflow: 'auto' }}
        theme="light"
      >
        <div style={{ padding: '0 0 12px' }}>
          <Button type="primary" icon={<PlusOutlined />} block onClick={newChat}>
            Chat Baru
          </Button>
        </div>
        <Spin spinning={loadingSessions}>
          {sessions.length === 0 ? (
            <Empty description="Belum ada chat" image={Empty.PRESENTED_IMAGE_SIMPLE} />
          ) : (
            <List
              dataSource={sessions}
              renderItem={(item) => (
                <List.Item
                  onClick={() => loadSession(item.session_id)}
                  style={{
                    cursor: 'pointer',
                    background: sessionID === item.session_id ? '#e6f4ff' : undefined,
                    padding: '8px 12px',
                    borderRadius: 6,
                    marginBottom: 4,
                  }}
                  actions={[
                    <Popconfirm
                      title="Hapus sesi ini?"
                      onConfirm={() => handleDeleteSession(item.session_id)}
                    >
                      <DeleteOutlined style={{ color: '#999' }} />
                    </Popconfirm>,
                  ]}
                >
                  <List.Item.Meta
                    title={
                      <Text ellipsis style={{ maxWidth: 180 }}>
                        {item.question}
                      </Text>
                    }
                    description={
                      <Space size={4}>
                        <ClockCircleOutlined style={{ fontSize: 11 }} />
                        <Text type="secondary" style={{ fontSize: 11 }}>
                          {new Date(item.created_at).toLocaleDateString('id-ID')}
                        </Text>
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          )}
        </Spin>
      </Sider>

      <Content style={{ display: 'flex', flexDirection: 'column', paddingLeft: 24 }}>
        <div style={{ flex: 1, overflow: 'auto', paddingBottom: 16 }}>
          {messages.length === 0 ? (
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                height: '100%',
                color: '#999',
              }}
            >
              <RobotOutlined style={{ fontSize: 48, marginBottom: 16 }} />
              <Title level={4} type="secondary">AI Chat Assistant</Title>
              <Text type="secondary">
                Tanyakan tentang data produksi dalam bahasa Indonesia
              </Text>
              <div style={{ marginTop: 24, maxWidth: 400 }}>
                <Text type="secondary" style={{ fontSize: 13 }}>
                  Contoh: "Berapa total produksi hari ini?", "Tampilkan alarm terakhir",
                  "Data curing shift malam", "Rata-rata cycle time building"
                </Text>
              </div>
            </div>
          ) : (
            messages.map((msg) => (
              <div
                key={msg.id}
                style={{
                  display: 'flex',
                  justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                  marginBottom: 16,
                }}
              >
                <div style={{ maxWidth: '75%' }}>
                  <Space align="start" size={8}>
                    {msg.role === 'assistant' && (
                      <Avatar icon={<RobotOutlined />} style={{ backgroundColor: '#1677ff' }} />
                    )}
                    <div>
                      <Card
                        size="small"
                        style={{
                          background: msg.role === 'user' ? '#1677ff' : '#f5f5f5',
                          color: msg.role === 'user' ? '#fff' : undefined,
                          border: 'none',
                        }}
                      >
                        {msg.role === 'assistant' && msg.detected_intent && (
                          <div style={{ marginBottom: 8 }}>
                            {sqlStatusIcon(msg.sql_status)}
                            <Tag color="blue" style={{ marginLeft: 4 }}>
                              {msg.detected_intent}
                            </Tag>
                            {msg.sql_status && (
                              <Tag>{msg.sql_status}</Tag>
                            )}
                          </div>
                        )}

                        <Paragraph style={{ margin: 0, whiteSpace: 'pre-wrap', color: 'inherit' }}>
                          {msg.content}
                        </Paragraph>

                        {msg.role === 'assistant' && msg.generated_sql && (
                          <Collapse ghost size="small" style={{ marginTop: 8 }}>
                            <Panel
                              header={
                                <Space size={4}>
                                  <CodeOutlined />
                                  <Text style={{ fontSize: 12 }}>SQL Query</Text>
                                </Space>
                              }
                              key="sql"
                            >
                              <pre
                                style={{
                                  margin: 0,
                                  fontSize: 12,
                                  background: '#1a1a2e',
                                  color: '#e0e0e0',
                                  padding: 12,
                                  borderRadius: 4,
                                  overflow: 'auto',
                                }}
                              >
                                {msg.generated_sql}
                              </pre>
                            </Panel>
                          </Collapse>
                        )}

                        {msg.role === 'assistant' && msg.query_result && msg.query_result.rows.length > 0 && (
                          <div style={{ marginTop: 8, fontSize: 12, color: '#888' }}>
                            <Text type="secondary">
                              {msg.query_result.rows.length} dari {msg.query_result.total} baris
                              {msg.query_result.latency ? ` · ${msg.query_result.latency}` : ''}
                            </Text>
                          </div>
                        )}
                      </Card>
                      <div style={{ fontSize: 11, color: '#999', marginTop: 4, paddingLeft: 4 }}>
                        {msg.timestamp.toLocaleTimeString('id-ID')}
                        {msg.status === 'sending' && ' · mengirim...'}
                        {msg.status === 'error' && ' · gagal'}
                      </div>
                    </div>
                    {msg.role === 'user' && (
                      <Avatar icon={<UserOutlined />} style={{ backgroundColor: '#52c41a' }} />
                    )}
                  </Space>
                </div>
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>

        <div style={{ borderTop: '1px solid #f0f0f0', paddingTop: 16 }}>
          <Space.Compact style={{ width: '100%' }}>
            <TextArea
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Tanya tentang data produksi..."
              autoSize={{ minRows: 1, maxRows: 4 }}
              disabled={sending}
              style={{ flex: 1 }}
            />
            <Tooltip title="Kirim (Enter)">
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleSend}
                loading={sending}
                disabled={!input.trim()}
              />
            </Tooltip>
          </Space.Compact>
          <div style={{ marginTop: 4 }}>
            <Text type="secondary" style={{ fontSize: 11 }}>
              Enter untuk kirim · Shift+Enter untuk baris baru
            </Text>
          </div>
        </div>
      </Content>
    </Layout>
  )
}


