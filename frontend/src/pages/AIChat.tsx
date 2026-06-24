import { useState, useRef, useEffect, useCallback } from 'react'
import {
  Layout, Input, Button, Typography, Space, List, Tag, Spin,
  Card, Collapse, Tooltip, Empty, Popconfirm, Avatar, Divider, Drawer,
} from 'antd'
import {
  SendOutlined, RobotOutlined, UserOutlined, DeleteOutlined,
  PlusOutlined, ClockCircleOutlined, CodeOutlined,
  CheckCircleOutlined, CloseCircleOutlined, MessageOutlined,
  ClearOutlined, ThunderboltOutlined, MenuOutlined,
} from '@ant-design/icons'
import type { ChatMessage, ChatSession } from '../services/chat'
import { chatApi } from '../services/chat'
import { useThemeStore } from '../stores/themeStore'

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
  const inputRef = useRef<any>(null)
  const { darkMode } = useThemeStore()
  const [sessionsOpen, setSessionsOpen] = useState(false)
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768)

  useEffect(() => {
    const handleResize = () => setIsMobile(window.innerWidth < 768)
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

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
      //
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
      //
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
      //
    }
  }

  const sqlStatusIcon = (status?: string) => {
    switch (status) {
      case 'valid': return <CheckCircleOutlined style={{ color: '#52c41a', fontSize: 12 }} />
      case 'invalid':
      case 'rejected': return <CloseCircleOutlined style={{ color: '#ff4d4f', fontSize: 12 }} />
      default: return null
    }
  }

  return (
    <div className="page-enter" style={{ height: isMobile ? 'calc(100vh - 108px)' : 'calc(100vh - 140px)' }}>
      {/* Chat container */}
      <div style={{
        display: 'flex',
        height: '100%',
        background: darkMode ? 'var(--bg-container)' : '#fff',
        borderRadius: 16,
        border: '1px solid var(--border-color)',
        overflow: 'hidden',
        boxShadow: 'var(--shadow-sm)',
      }}>
        {/* Sessions sidebar - desktop */}
        <div style={{
          width: 280,
          borderRight: '1px solid var(--border-color)',
          display: isMobile ? 'none' : 'flex',
          flexDirection: 'column',
          background: darkMode ? 'var(--bg-container)' : '#fafafa',
          flexShrink: 0,
        }}>
          <div style={{ padding: 16, borderBottom: '1px solid var(--border-color)' }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              block
              onClick={() => { newChat(); if (isMobile) setSessionsOpen(false) }}
              className="btn-glow"
              style={{ height: 40, fontWeight: 600 }}
            >
              Chat Baru
            </Button>
          </div>
          <div style={{ flex: 1, overflow: 'auto', padding: '8px' }}>
            <Spin spinning={loadingSessions}>
              {sessions.length === 0 ? (
                <div style={{ padding: 24, textAlign: 'center' }}>
                  <MessageOutlined style={{ fontSize: 32, color: 'var(--text-tertiary)', marginBottom: 8 }} />
                  <Text type="secondary" style={{ fontSize: 13 }}>Belum ada chat</Text>
                </div>
              ) : (
                <List
                  dataSource={sessions}
                  renderItem={(item) => (
                    <List.Item
                      onClick={() => loadSession(item.session_id)}
                      style={{
                        cursor: 'pointer',
                        background: sessionID === item.session_id
                          ? 'var(--primary-bg)'
                          : 'transparent',
                        padding: '10px 12px',
                        borderRadius: 10,
                        marginBottom: 2,
                        border: sessionID === item.session_id
                          ? '1px solid var(--border-color)'
                          : '1px solid transparent',
                        transition: 'all 0.2s',
                      }}
                      onMouseEnter={(e) => {
                        if (sessionID !== item.session_id) {
                          e.currentTarget.style.background = 'var(--bg-hover)'
                        }
                      }}
                      onMouseLeave={(e) => {
                        if (sessionID !== item.session_id) {
                          e.currentTarget.style.background = 'transparent'
                        }
                      }}
                      actions={[
                        <Popconfirm
                          title="Hapus sesi ini?"
                          onConfirm={() => handleDeleteSession(item.session_id)}
                          key="del"
                        >
                          <Button
                            type="text"
                            size="small"
                            icon={<DeleteOutlined />}
                            style={{ color: 'var(--text-tertiary)', opacity: 0, transition: 'opacity 0.2s' }}
                            className="delete-btn"
                          />
                        </Popconfirm>,
                      ]}
                      onMouseOver={(e) => {
                        const btn = e.currentTarget.querySelector('.delete-btn') as HTMLElement
                        if (btn) btn.style.opacity = '1'
                      }}
                      onMouseOut={(e) => {
                        const btn = e.currentTarget.querySelector('.delete-btn') as HTMLElement
                        if (btn) btn.style.opacity = '0'
                      }}
                    >
                      <List.Item.Meta
                        avatar={
                          <Avatar
                            icon={<MessageOutlined />}
                            size={28}
                            style={{
                              background: sessionID === item.session_id
                                ? 'var(--primary-color)'
                                : 'var(--bg-hover)',
                              color: sessionID === item.session_id ? '#fff' : 'var(--text-secondary)',
                            }}
                          />
                        }
                        title={
                          <Text
                            ellipsis
                            style={{
                              maxWidth: 160,
                              fontSize: 13,
                              fontWeight: sessionID === item.session_id ? 600 : 400,
                            }}
                          >
                            {item.question}
                          </Text>
                        }
                        description={
                          <Space size={4}>
                            <ClockCircleOutlined style={{ fontSize: 10, color: 'var(--text-tertiary)' }} />
                            <Text type="secondary" style={{ fontSize: 10 }}>
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
          </div>
        </div>

        {/* Sessions drawer - mobile */}
        <Drawer
          placement="left"
          width={280}
          open={isMobile && sessionsOpen}
          onClose={() => setSessionsOpen(false)}
          styles={{ body: { padding: 0, background: darkMode ? 'var(--bg-container)' : '#fafafa' } }}
          closable={false}
        >
          <div style={{ padding: 16, borderBottom: '1px solid var(--border-color)' }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              block
              onClick={() => { newChat(); setSessionsOpen(false) }}
              style={{ height: 40, fontWeight: 600 }}
            >
              Chat Baru
            </Button>
          </div>
          <div style={{ flex: 1, overflow: 'auto', padding: '8px', height: 'calc(100vh - 80px)' }}>
            <Spin spinning={loadingSessions}>
              {sessions.length === 0 ? (
                <div style={{ padding: 24, textAlign: 'center' }}>
                  <MessageOutlined style={{ fontSize: 32, color: 'var(--text-tertiary)', marginBottom: 8 }} />
                  <Text type="secondary" style={{ fontSize: 13 }}>Belum ada chat</Text>
                </div>
              ) : (
                <List
                  dataSource={sessions}
                  renderItem={(item) => (
                    <List.Item
                      onClick={() => { loadSession(item.session_id); setSessionsOpen(false) }}
                      style={{
                        cursor: 'pointer',
                        background: sessionID === item.session_id
                          ? 'var(--primary-bg)'
                          : 'transparent',
                        padding: '10px 12px',
                        borderRadius: 10,
                        marginBottom: 2,
                        border: sessionID === item.session_id
                          ? '1px solid var(--border-color)'
                          : '1px solid transparent',
                      }}
                    >
                      <List.Item.Meta
                        avatar={
                          <Avatar
                            icon={<MessageOutlined />}
                            size={28}
                            style={{
                              background: sessionID === item.session_id
                                ? 'var(--primary-color)'
                                : 'var(--bg-hover)',
                              color: sessionID === item.session_id ? '#fff' : 'var(--text-secondary)',
                            }}
                          />
                        }
                        title={
                          <Text
                            ellipsis
                            style={{ maxWidth: 180, fontSize: 13, fontWeight: sessionID === item.session_id ? 600 : 400 }}
                          >
                            {item.question}
                          </Text>
                        }
                        description={
                          <Space size={4}>
                            <ClockCircleOutlined style={{ fontSize: 10, color: 'var(--text-tertiary)' }} />
                            <Text type="secondary" style={{ fontSize: 10 }}>
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
          </div>
        </Drawer>

        {/* Main chat area */}
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
          {/* Chat header */}
          <div style={{
            padding: isMobile ? '10px 12px' : '12px 20px',
            borderBottom: '1px solid var(--border-color)',
            display: 'flex',
            alignItems: 'center',
            gap: 10,
          }}>
            {isMobile && (
              <Button
                type="text"
                icon={<MenuOutlined />}
                onClick={() => setSessionsOpen(true)}
                style={{ fontSize: 18, color: 'var(--text-secondary)', width: 32, height: 32 }}
              />
            )}
            <div style={{
              width: 32, height: 32, borderRadius: 8,
              background: 'linear-gradient(135deg, var(--primary-color), var(--primary-hover))',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 16, color: '#fff',
              flexShrink: 0,
            }}>
              <ThunderboltOutlined />
            </div>
            <div>
              <Text strong style={{ fontSize: 14 }}>AI Chat Assistant</Text>
              <br />
              <Text type="secondary" style={{ fontSize: 11 }}>Tanyakan tentang data produksi dalam bahasa Indonesia</Text>
            </div>
            {messages.length > 0 && (
              <Button
                type="text"
                size="small"
                icon={<ClearOutlined />}
                onClick={newChat}
                style={{ marginLeft: 'auto', color: 'var(--text-tertiary)' }}
              >
                Clear
              </Button>
            )}
          </div>

          {/* Messages */}
          <div style={{ flex: 1, overflow: 'auto', padding: isMobile ? '12px' : '20px 24px' }}>
            {messages.length === 0 ? (
              <div style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                height: '100%',
                gap: 12,
              }}>
                <div style={{
                  width: 72, height: 72, borderRadius: 20,
                  background: 'linear-gradient(135deg, var(--primary-color), #8b5cf6)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontSize: 32, color: '#fff',
                  boxShadow: '0 8px 32px rgba(22,119,255,0.25)',
                  marginBottom: 8,
                }}>
                  <RobotOutlined />
                </div>
                <Title level={4} style={{ color: 'var(--text-primary)', fontWeight: 600, margin: 0 }}>
                  AI Chat Assistant
                </Title>
                <Text type="secondary" style={{ fontSize: 14 }}>
                  Tanyakan tentang data produksi dalam bahasa Indonesia
                </Text>
                <div style={{ marginTop: 8, maxWidth: 480, textAlign: 'center' }}>
                  <div style={{
                    display: 'flex', flexWrap: 'wrap', gap: 8, justifyContent: 'center',
                  }}>
                    {[
                      'Berapa total produksi hari ini?',
                      'Tampilkan alarm terakhir',
                      'Data curing shift malam',
                      'Rata-rata cycle time building',
                    ].map((suggestion) => (
                      <div
                        key={suggestion}
                        onClick={() => {
                          setInput(suggestion)
                          inputRef.current?.focus()
                        }}
                        style={{
                          padding: '8px 14px',
                          borderRadius: 10,
                          background: 'var(--bg-hover)',
                          border: '1px solid var(--border-color)',
                          cursor: 'pointer',
                          fontSize: 12,
                          color: 'var(--text-secondary)',
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.background = 'var(--primary-bg)'
                          e.currentTarget.style.borderColor = 'var(--primary-color)'
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.background = 'var(--bg-hover)'
                          e.currentTarget.style.borderColor = 'var(--border-color)'
                        }}
                      >
                        {suggestion}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ) : (
              messages.map((msg) => (
                <div
                  key={msg.id}
                  style={{
                    display: 'flex',
                    justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                    marginBottom: 20,
                    animation: 'slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
                  }}
                >
                  <div style={{ maxWidth: isMobile ? '92%' : '80%', minWidth: 0 }}>
                    <Space align="start" size={10}>
                      {msg.role === 'assistant' && (
                        <Avatar
                          icon={<RobotOutlined />}
                          style={{
                            background: 'linear-gradient(135deg, var(--primary-color), #8b5cf6)',
                            boxShadow: '0 2px 8px rgba(22,119,255,0.3)',
                            marginTop: 4,
                          }}
                          size={32}
                        />
                      )}
                      <div>
                        <Card
                          size="small"
                          className={msg.role === 'user' ? 'chat-bubble-user' : 'chat-bubble-assistant'}
                          style={{
                            background: msg.role === 'user'
                              ? 'linear-gradient(135deg, var(--primary-color), var(--primary-active))'
                              : darkMode ? 'rgba(255,255,255,0.05)' : '#f0f0f5',
                            color: msg.role === 'user' ? '#fff' : undefined,
                            boxShadow: msg.role === 'user'
                              ? '0 4px 16px rgba(22,119,255,0.25)'
                              : '0 1px 4px rgba(0,0,0,0.04)',
                            padding: '12px 16px',
                            border: msg.role === 'user' ? 'none' : '1px solid var(--border-color)',
                            maxWidth: '100%',
                          }}
                        >
                          {msg.role === 'assistant' && msg.detected_intent && (
                            <div style={{ marginBottom: 8, display: 'flex', alignItems: 'center', gap: 6, flexWrap: 'wrap' }}>
                              {sqlStatusIcon(msg.sql_status)}
                              <Tag
                                color="blue"
                                style={{
                                  margin: 0,
                                  fontSize: 11,
                                  fontWeight: 600,
                                }}
                              >
                                {msg.detected_intent}
                              </Tag>
                              {msg.sql_status && (
                                <Tag style={{ margin: 0, fontSize: 11 }}>{msg.sql_status}</Tag>
                              )}
                            </div>
                          )}

                          <Paragraph style={{
                            margin: 0,
                            whiteSpace: 'pre-wrap',
                            color: msg.role === 'user' ? '#fff' : 'var(--text-primary)',
                            fontSize: 14,
                            lineHeight: 1.6,
                          }}>
                            {msg.content}
                          </Paragraph>

                          {msg.role === 'assistant' && msg.generated_sql && (
                            <Collapse ghost size="small" style={{ marginTop: 10 }}>
                              <Panel
                                header={
                                  <Space size={4}>
                                    <CodeOutlined style={{ fontSize: 12 }} />
                                    <Text style={{ fontSize: 12, color: 'var(--text-secondary)' }}>SQL Query</Text>
                                  </Space>
                                }
                                key="sql"
                                style={{ border: 'none' }}
                              >
                                <pre
                                  style={{
                                    margin: 0,
                                    fontSize: 12,
                                    background: darkMode ? '#0a0a0f' : '#1a1a2e',
                                    color: '#e0e0e0',
                                    padding: 14,
                                    borderRadius: 8,
                                    overflow: 'auto',
                                    lineHeight: 1.5,
                                    fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                                  }}
                                >
                                  {msg.generated_sql}
                                </pre>
                              </Panel>
                            </Collapse>
                          )}

                          {msg.role === 'assistant' && msg.query_result && msg.query_result.rows.length > 0 && (
                            <div style={{
                              marginTop: 8,
                              padding: '6px 10px',
                              background: darkMode ? 'rgba(255,255,255,0.03)' : 'rgba(0,0,0,0.02)',
                              borderRadius: 6,
                              display: 'flex',
                              alignItems: 'center',
                              gap: 8,
                              fontSize: 12,
                              color: 'var(--text-tertiary)',
                            }}>
                              <CheckCircleOutlined style={{ color: '#52c41a', fontSize: 12 }} />
                              <span>{msg.query_result.rows.length} baris </span>
                              {msg.query_result.total && (
                                <span>dari {msg.query_result.total}</span>
                              )}
                              {msg.query_result.latency && (
                                <span>· {msg.query_result.latency}</span>
                              )}
                            </div>
                          )}
                        </Card>
                        <div style={{
                          fontSize: 11,
                          color: 'var(--text-tertiary)',
                          marginTop: 4,
                          paddingLeft: msg.role === 'user' ? 0 : 4,
                          textAlign: msg.role === 'user' ? 'right' : 'left',
                        }}>
                          {msg.timestamp.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })}
                          {msg.status === 'sending' && ' · mengirim...'}
                          {msg.status === 'error' && ' · gagal'}
                        </div>
                      </div>
                      {msg.role === 'user' && (
                        <Avatar
                          icon={<UserOutlined />}
                          style={{
                            background: 'linear-gradient(135deg, #52c41a, #389e0d)',
                            boxShadow: '0 2px 8px rgba(82,196,26,0.3)',
                            marginTop: 4,
                          }}
                          size={32}
                        />
                      )}
                    </Space>
                  </div>
                </div>
              ))
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Input area */}
          <div style={{
            padding: isMobile ? '12px' : '16px 20px',
            borderTop: '1px solid var(--border-color)',
            background: darkMode ? 'var(--bg-container)' : '#fff',
          }}>
            <div style={{
              display: 'flex',
              gap: 8,
              background: darkMode ? 'rgba(255,255,255,0.04)' : 'var(--bg-hover)',
              borderRadius: 12,
              padding: 4,
              border: '1px solid var(--border-color)',
              transition: 'border-color 0.2s',
            }}
              onFocusCapture={(e) => {
                e.currentTarget.style.borderColor = 'var(--primary-color)'
              }}
              onBlurCapture={(e) => {
                e.currentTarget.style.borderColor = 'var(--border-color)'
              }}
            >
              <TextArea
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Tanya tentang data produksi..."
                autoSize={{ minRows: 1, maxRows: 4 }}
                disabled={sending}
                variant="borderless"
                style={{
                  flex: 1,
                  fontSize: 14,
                  padding: '8px 12px',
                  background: 'transparent',
                  border: 'none',
                  boxShadow: 'none',
                  resize: 'none',
                }}
              />
              <Button
                type="primary"
                icon={<SendOutlined />}
                onClick={handleSend}
                loading={sending}
                disabled={!input.trim()}
                style={{
                  width: 40,
                  height: 40,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  flexShrink: 0,
                  borderRadius: 10,
                }}
              />
            </div>
            <div style={{ marginTop: 6, display: 'flex', justifyContent: 'space-between' }}>
              <Text type="secondary" style={{ fontSize: 11 }}>
                Enter untuk kirim · Shift+Enter untuk baris baru
              </Text>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
