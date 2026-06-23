import api from './api'

export interface ChatMessage {
  id?: number
  question: string
  answer?: string
  detected_intent?: string
  generated_sql?: string
  sql_status?: string
  query_result?: {
    columns: string[]
    rows: Record<string, unknown>[]
    total: number
    latency: string
  }
  status: string
  created_at: string
  session_id: string
}

export interface ChatSession {
  session_id: string
  question: string
  created_at: string
}

export const chatApi = {
  send: async (question: string, sessionID?: string) => {
    const res = await api.post('/api/v1/ai/chat', {
      question,
      session_id: sessionID || '',
    })
    return res.data.data as {
      session_id: string
      answer: string
      detected_intent?: string
      generated_sql?: string
      sql_status: string
      query_result?: ChatMessage['query_result']
      latency?: string
    }
  },

  getHistory: async (sessionID: string, page = 1, limit = 50) => {
    const res = await api.get(`/api/v1/ai/chat/${sessionID}/history`, {
      params: { page, limit },
    })
    return res.data
  },

  getSessions: async () => {
    const res = await api.get('/api/v1/ai/chat/sessions')
    return res.data.data as ChatSession[]
  },

  deleteSession: async (sessionID: string) => {
    await api.delete(`/api/v1/ai/chat/${sessionID}`)
  },
}
