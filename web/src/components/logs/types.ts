export interface Log {
  timestamp: string
  severity: string
  message: string
  serviceName: string
  rawLine: string
}

export interface LogViewerProps {
  logs: Log[]
  loading?: boolean
  onLoadMore?: () => void
  searchQuery?: string
  autoScroll?: boolean
  onToggleAutoScroll?: () => void
  maxHeight?: string
  hideHeader?: boolean
}

