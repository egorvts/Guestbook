import './style.css'
import { createMessage, listMessages, type Message } from './api'

const app = document.querySelector<HTMLDivElement>('#app')!

app.innerHTML = `
  <main class="page">
    <section class="input-panel" aria-label="Write a message">
      <form id="message-form">
        <h1>Thanks for your visit!<br />Leave a message</h1>
        <label for="name-input">Your name or alias</label>
        <input
          id="name-input"
          name="name"
          type="text"
          placeholder="Anonymous user"
          autocomplete="nickname"
          maxlength="255"
        />
        <label for="message-input">Your message</label>
        <textarea
          id="message-input"
          name="message"
          placeholder="Type your message..."
          maxlength="1024"
          rows="3"
          required
        ></textarea>
        <button type="submit" disabled>Publish</button>
        <p id="form-status" class="status" role="status"></p>
        <p class="publish-notice">Your name and message will be public.<br />Messages cannot be edited or deleted.</p>
      </form>
    </section>
    <section class="messages-panel" aria-label="Messages" tabindex="0">
      <h2>Message Feed</h2>
      <p id="feed-status" class="status" role="status">Loading messages…</p>
      <ul id="messages" aria-live="polite" aria-relevant="additions"></ul>
    </section>
  </main>
`

const form = document.querySelector<HTMLFormElement>('#message-form')!
const input = document.querySelector<HTMLTextAreaElement>('#message-input')!
const nameInput = document.querySelector<HTMLInputElement>('#name-input')!
const messages = document.querySelector<HTMLUListElement>('#messages')!
const button = form.querySelector<HTMLButtonElement>('button')!
const formStatus = document.querySelector<HTMLParagraphElement>('#form-status')!
const feedStatus = document.querySelector<HTMLParagraphElement>('#feed-status')!
const dateFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  hour: 'numeric',
  minute: '2-digit',
  hour12: true,
})

function setStatus(element: HTMLElement, text: string, isError = false): void {
  element.textContent = text
  element.classList.toggle('is-error', isError)
}

function renderMessage(entry: Message): HTMLLIElement {
  const message = document.createElement('li')
  const metadata = document.createElement('div')
  metadata.className = 'message-meta'

  const author = document.createElement('strong')
  author.className = 'message-author'
  author.textContent = entry.name

  const date = new Date(entry.date)
  const timestamp = document.createElement('time')
  timestamp.dateTime = entry.date
  timestamp.textContent = dateFormatter.format(date)
  timestamp.title = date.toLocaleString()

  metadata.append(author, timestamp)

  const body = document.createElement('p')
  body.className = 'message-body'
  body.textContent = entry.content
  message.append(body, metadata)
  return message
}

async function refreshMessages(): Promise<void> {
  const entries = await listMessages()
  const cards = entries.map(renderMessage)
  messages.replaceChildren(...cards)
  setStatus(feedStatus, entries.length ? '' : 'No messages yet. Leave the first one!')
}

form.addEventListener('submit', async (event) => {
  event.preventDefault()
  if (button.disabled) return

  const content = input.value.trim()
  const name = nameInput.value.trim() || 'Anonymous'
  if (!content) {
    setStatus(formStatus, 'Please enter a message.\nIt cannot be empty.', true)
    return
  }

  // The backend limits UTF-8 bytes, rather than characters.
  const encoder = new TextEncoder()
  if (encoder.encode(name).length >= 256 || encoder.encode(content).length >= 2048) {
    setStatus(formStatus, 'Your name or message is too long.\nPlease shorten it.', true)
    return
  }

  button.disabled = true
  button.textContent = 'Publishing…'
  setStatus(formStatus, '')

  try {
    await createMessage(name, content)
    if (input.value.trim() === content) input.value = ''
    setStatus(formStatus, 'Message published.')

    try {
      await refreshMessages()
    } catch {
      setStatus(feedStatus, 'Your message was published.\nReload to refresh the feed.', true)
    }
  } catch {
    setStatus(formStatus, 'Could not confirm your message was published.\nCheck the feed before trying again.', true)
  } finally {
    button.disabled = false
    button.textContent = 'Publish'
  }
})

async function initialize(): Promise<void> {
  try {
    await refreshMessages()
  } catch {
    setStatus(feedStatus, 'Could not load messages.\nPlease reload to try again.', true)
  } finally {
    button.disabled = false
  }
}

void initialize()
