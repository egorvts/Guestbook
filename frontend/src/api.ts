export interface Message {
  id: string
  name: string
  content: string
  date: string
}

export interface AppInfo {
  version: string
  'infrastructure used': string
}

const baseUrl = (import.meta.env?.VITE_API_BASE_URL || '/api').replace(/\/$/, '')

async function request(path: string, options?: RequestInit): Promise<Response> {
  const response = await fetch(`${baseUrl}${path}`, {
    ...options,
    signal: AbortSignal.timeout(10_000),
  })

  if (!response.ok) {
    throw new Error(`API request failed (${response.status})`)
  }

  return response
}

export async function listMessages(): Promise<Message[]> {
  const response = await request('/messages', { cache: 'no-store' })
  return response.json()
}

export async function getAppInfo(): Promise<AppInfo> {
  const response = await request('/info', { cache: 'no-store' })
  return response.json()
}

export async function createMessage(name: string, content: string): Promise<void> {
  await request('/messages', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, content }),
  })
}
