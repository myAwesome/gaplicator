import { useState, useCallback } from 'react'
import { Copy, Check } from './icons.jsx'

const BRIDGE_URL = import.meta.env.VITE_GAPLICATOR_BRIDGE_URL || 'http://127.0.0.1:8787'

function outputName(projectName) {
  const cleaned = (projectName || 'app')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return cleaned || 'app'
}

export default function YamlPreview({ fullYaml, fullHighlighted, simpleYaml, simpleHighlighted, defaultTab, projectName }) {
  const [tab, setTab] = useState(defaultTab || 'full')
  const [copied, setCopied] = useState(false)
  const [sending, setSending] = useState(false)
  const [status, setStatus] = useState(null)
  const [statusText, setStatusText] = useState('')

  const yaml = tab === 'simple' ? simpleYaml : fullYaml
  const highlighted = tab === 'simple' ? simpleHighlighted : fullHighlighted

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(yaml)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      const el = document.createElement('textarea')
      el.value = yaml
      el.style.position = 'fixed'
      el.style.opacity = '0'
      document.body.appendChild(el)
      el.select()
      document.execCommand('copy')
      document.body.removeChild(el)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }, [yaml])

  const handleSend = useCallback(async () => {
    if (sending) return
    setSending(true)
    setStatus(null)
    setStatusText('')

    const out = `dist/${outputName(projectName)}`

    try {
      const res = await fetch(`${BRIDGE_URL}/build`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          yaml,
          output: out,
        }),
      })

      const payload = await res.json().catch(() => ({}))
      if (!res.ok) {
        const msg = Array.isArray(payload.errors) && payload.errors.length
          ? payload.errors.join('; ')
          : `Bridge error: ${res.status}`
        throw new Error(msg)
      }

      setStatus('ok')
      setStatusText(`Generated in ${payload.output || out}`)
    } catch (err) {
      setStatus('error')
      setStatusText(err?.message || 'Failed to send schema to gaplicator bridge')
    } finally {
      setSending(false)
    }
  }, [sending, projectName, yaml])

  return (
    <div className="preview-panel">
      <div className="preview-header">
        <div className="preview-tabs">
          <button
            className={`preview-tab ${tab === 'simple' ? 'active' : ''}`}
            onClick={() => setTab('simple')}
          >
            Simple
          </button>
          <button
            className={`preview-tab ${tab === 'full' ? 'active' : ''}`}
            onClick={() => setTab('full')}
          >
            Full
          </button>
        </div>
        <div className="preview-actions">
          <button
            className={`btn-copy ${copied ? 'copied' : ''}`}
            onClick={handleCopy}
            title="Copy to clipboard"
          >
            {copied ? <Check size={13} /> : <Copy size={13} />}
            {copied ? 'Copied!' : 'Copy'}
          </button>
          <button
            className={`btn-copy btn-send ${sending ? 'sending' : ''}`}
            onClick={handleSend}
            disabled={sending}
            title="Send YAML schema to local gaplicator bridge"
          >
            {sending ? 'Sending...' : 'Send to CLI'}
          </button>
        </div>
      </div>
      {status && (
        <div className={`preview-status ${status}`}>
          {statusText}
        </div>
      )}
      <div
        className="preview-code"
        dangerouslySetInnerHTML={{ __html: highlighted }}
      />
    </div>
  )
}
