export interface Message {
  id: string
  name: string
  content: string
  date: string
}

const baseUrl = (import.meta.env?.VITE_API_BASE_URL || '/api').replace(/\/$/, '')

async function request(options?: RequestInit): Promise<Response> {
  const response = await fetch(`${baseUrl}/messages`, {
    ...options,
    signal: AbortSignal.timeout(10_000),
  })

  if (!response.ok) {
    throw new Error(`Message request failed (${response.status})`)
  }

  return response
}

export async function listMessages(): Promise<Message[]> {
  const response = await request({ cache: 'no-store' })
  return response.json()
}

export async function createMessage(name: string, content: string): Promise<void> {
  await request({
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, content }),
  })
}
