import { useEffect, useRef, useState, useCallback } from 'react'
import type { Log } from '@/components/logs/types'

const DEFAULT_MAX_BUFFER_SIZE = 5000
const MAX_RECONNECT_ATTEMPTS = 5
const BASE_RECONNECT_DELAY = 1000
const MAX_RECONNECT_DELAY = 30000

interface LogStreamOptions {
  serviceId: string
  filters?: {
    search?: string
    levels?: string[]
  }
  enabled?: boolean
  maxBufferSize?: number
}

interface LogStreamResult {
  logs: Log[]
  status: 'connecting' | 'connected' | 'disconnected' | 'error'
  error: string | null
  connect: () => void
  disconnect: () => void
  clear: () => void
}

/**
 * Custom hook for managing WebSocket connection to stream deployment logs in real-time
 * Connects to service-specific log stream endpoint
 * Handles connection lifecycle, reconnection with exponential backoff, and buffer management
 */
export function useLogStream(options: LogStreamOptions): LogStreamResult {
  const [logs, setLogs] = useState<Log[]>([])
  const [status, setStatus] = useState<'connecting' | 'connected' | 'disconnected' | 'error'>('disconnected')
  const [error, setError] = useState<string | null>(null)
  
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const lastTimestampRef = useRef<string | null>(null)
  const maxBufferSize = options.maxBufferSize || DEFAULT_MAX_BUFFER_SIZE
  
  // Store latest options in ref to avoid stale closures
  const optionsRef = useRef(options)
  useEffect(() => {
    optionsRef.current = options
  }, [options])

  const buildWebSocketUrl = useCallback(() => {
    const opts = optionsRef.current
    
    // Determine WebSocket protocol based on current page protocol
    const protocol = typeof window !== 'undefined' && window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    
    // Get API base URL from environment or use current host
    const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL
    let host: string
    
    if (apiBaseUrl) {
      // Extract host from API base URL (e.g., "http://localhost:8090" -> "localhost:8090")
      host = apiBaseUrl.replace(/^https?:\/\//, '')
    } else if (typeof window !== 'undefined') {
      host = window.location.host
    } else {
      host = 'localhost:8090' // Fallback for SSR
    }
    
    let url = `${protocol}//${host}/api/v1/services/${opts.serviceId}/logs/stream`
    
    // Build query parameters
    const params = new URLSearchParams()
    if (opts.filters?.search) {
      params.set('search', opts.filters.search)
    }
    if (opts.filters?.levels?.length) {
      params.set('levels', opts.filters.levels.join(','))
    }
    if (lastTimestampRef.current) {
      params.set('lastTimestamp', lastTimestampRef.current)
    }
    
    const queryString = params.toString()
    if (queryString) {
      url += `?${queryString}`
    }
    
    return url
  }, [])

  const connect = useCallback(() => {
    // Don't connect if already connected or connecting
    if (wsRef.current?.readyState === WebSocket.OPEN || 
        wsRef.current?.readyState === WebSocket.CONNECTING) {
      return
    }

    // Clear any pending reconnect timeout
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    setStatus('connecting')
    setError(null)

    try {
      const ws = new WebSocket(buildWebSocketUrl())
      wsRef.current = ws

      ws.onopen = () => {
        console.log('[useLogStream] WebSocket connected')
        setStatus('connected')
        reconnectAttemptsRef.current = 0
        // Clear logs on fresh connection to avoid duplicates
        setLogs([])
      }

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)

          switch (message.type) {
            case 'log':
              setLogs((prevLogs) => {
                const newLog = message.data as Log
                lastTimestampRef.current = newLog.timestamp

                // Circular buffer: remove old logs if exceeding limit
                const updatedLogs = [...prevLogs, newLog]
                if (updatedLogs.length > maxBufferSize) {
                  return updatedLogs.slice(-maxBufferSize)
                }
                return updatedLogs
              })
              break

            case 'error':
              console.error('[useLogStream] Stream error:', message.data)
              setError(typeof message.data === 'string' ? message.data : 'Stream error occurred')
              break

            case 'status':
              console.log('[useLogStream] Stream status:', message.data)
              break

            case 'ping':
              // Heartbeat - no action needed
              break

            default:
              console.warn('[useLogStream] Unknown message type:', message.type)
          }
        } catch (err) {
          console.error('[useLogStream] Failed to parse message:', err)
        }
      }

      ws.onerror = (event) => {
        console.error('[useLogStream] WebSocket error:', event)
        setStatus('error')
        setError('Connection error occurred')
      }

      ws.onclose = (event) => {
        console.log('[useLogStream] WebSocket closed:', event.code, event.reason)
        setStatus('disconnected')
        wsRef.current = null

        // Don't reconnect if this is a manual disconnect (code 1000) or if disabled
        const isManualDisconnect = event.code === 1000
        const shouldReconnect = !isManualDisconnect && 
                               reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS && 
                               optionsRef.current.enabled !== false

        if (shouldReconnect) {
          const delay = Math.min(
            BASE_RECONNECT_DELAY * Math.pow(2, reconnectAttemptsRef.current),
            MAX_RECONNECT_DELAY
          )
          reconnectAttemptsRef.current++
          console.log(`[useLogStream] Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current}/${MAX_RECONNECT_ATTEMPTS})`)
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, delay)
        } else if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
          setStatus('error')
          setError('Maximum reconnection attempts reached')
        } else if (isManualDisconnect) {
          console.log('[useLogStream] Manual disconnect, not reconnecting')
        }
      }
    } catch (err) {
      console.error('[useLogStream] Failed to create WebSocket:', err)
      setStatus('error')
      setError('Failed to connect to log stream')
    }
  }, [buildWebSocketUrl, maxBufferSize])

  const disconnect = useCallback(() => {
    console.log('[useLogStream] disconnect() called')
    
    // Clear any pending reconnect timeout
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
      console.log('[useLogStream] Closing WebSocket connection')
      // Use normal closure code to prevent reconnection
      wsRef.current.close(1000, 'Client disconnect')
      wsRef.current = null
    }
    setStatus('disconnected')
    reconnectAttemptsRef.current = 0
  }, [])

  const clear = useCallback(() => {
    setLogs([])
    lastTimestampRef.current = null
  }, [])

  // Connect/disconnect based on enabled flag
  useEffect(() => {
    console.log('[useLogStream] Effect running, enabled:', options.enabled)
    if (options.enabled) {
      connect()
    } else {
      disconnect()
    }

    // Cleanup on unmount
    return () => {
      console.log('[useLogStream] Effect cleanup, disconnecting')
      disconnect()
    }
  }, [options.enabled, connect, disconnect])

  // Handle page visibility changes to save resources
  // Only disconnect after being hidden for a while (not on quick tab switches)
  useEffect(() => {
    let visibilityTimeout: NodeJS.Timeout | null = null

    const handleVisibilityChange = () => {
      if (document.hidden) {
        // Set a timer to disconnect after 30 seconds of being hidden
        visibilityTimeout = setTimeout(() => {
          console.log('[useLogStream] Page hidden for 30s, disconnecting to save resources')
          disconnect()
        }, 30000) // 30 seconds
      } else {
        // Page became visible again - cancel disconnect timer
        if (visibilityTimeout) {
          console.log('[useLogStream] Page visible again, cancelling disconnect timer')
          clearTimeout(visibilityTimeout)
          visibilityTimeout = null
        }
        
        // Reconnect if we were supposed to be streaming but got disconnected
        if (optionsRef.current.enabled && (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN)) {
          console.log('[useLogStream] Page visible and not connected, reconnecting')
          connect()
        }
      }
    }

    if (typeof document !== 'undefined') {
      document.addEventListener('visibilitychange', handleVisibilityChange)
      return () => {
        document.removeEventListener('visibilitychange', handleVisibilityChange)
        if (visibilityTimeout) {
          clearTimeout(visibilityTimeout)
        }
      }
    }
  }, [connect, disconnect])

  return {
    logs,
    status,
    error,
    connect,
    disconnect,
    clear,
  }
}

