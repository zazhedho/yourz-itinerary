import { useEffect, useRef, useState } from 'react'

import { getGoogleClientId } from '../../utils/runtimeConfig'

const googleScriptId = 'google-identity-services-script'

const ensureGoogleScript = () => {
  if (window.google?.accounts?.id) return Promise.resolve()

  const existing = document.getElementById(googleScriptId)
  if (existing?.dataset.loaded === 'true') return Promise.resolve()

  return new Promise((resolve, reject) => {
    const script = existing || document.createElement('script')
    script.id = googleScriptId
    script.src = 'https://accounts.google.com/gsi/client'
    script.async = true
    script.defer = true
    script.onload = () => {
      script.dataset.loaded = 'true'
      resolve()
    }
    script.onerror = () => reject(new Error('Google Identity gagal dimuat'))

    if (!existing) document.head.appendChild(script)
  })
}

const GoogleIdentityButton = ({ disabled = false, label, onCredential, onError, text = 'continue_with' }) => {
  const containerRef = useRef(null)
  const onCredentialRef = useRef(onCredential)
  const onErrorRef = useRef(onError)
  const [available, setAvailable] = useState(false)
  const clientId = getGoogleClientId()

  useEffect(() => {
    onCredentialRef.current = onCredential
  }, [onCredential])

  useEffect(() => {
    onErrorRef.current = onError
  }, [onError])

  useEffect(() => {
    let cancelled = false

    const init = async () => {
      if (!clientId) {
        setAvailable(false)
        return
      }

      try {
        await ensureGoogleScript()
      } catch (err) {
        if (!cancelled) {
          setAvailable(false)
          onErrorRef.current?.(err instanceof Error ? err.message : 'Google Identity gagal dimuat')
        }
        return
      }

      if (cancelled) return
      const api = window.google?.accounts?.id
      if (!api || !containerRef.current) {
        setAvailable(false)
        onErrorRef.current?.('Google Identity tidak tersedia.')
        return
      }

      api.initialize({
        client_id: clientId,
        callback: (response) => {
          const token = String(response.credential || '').trim()
          if (!token) {
            onErrorRef.current?.('Credential Google tidak tersedia.')
            return
          }
          void onCredentialRef.current(token)
        },
      })

      containerRef.current.innerHTML = ''
      api.renderButton(containerRef.current, {
        type: 'standard',
        theme: 'outline',
        text,
        shape: 'pill',
        size: 'large',
        width: 320,
        logo_alignment: 'left',
      })
      setAvailable(true)
    }

    init()

    return () => {
      cancelled = true
    }
  }, [clientId, text])

  return (
    <div className="google-auth-block">
      <div className="auth-divider">
        <span>{label}</span>
      </div>
      <div
        ref={containerRef}
        aria-disabled={!available || disabled}
        className={`google-identity-button ${!available || disabled ? 'disabled' : ''}`}
      />
    </div>
  )
}

export default GoogleIdentityButton
