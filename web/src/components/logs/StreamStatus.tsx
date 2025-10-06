import { Button } from "@/components/ui/button"
import { Loader2, Circle, CircleDot, CircleAlert } from "lucide-react"
import { cn } from "@/lib/utils"

interface StreamStatusProps {
  status: 'connecting' | 'connected' | 'disconnected' | 'error'
  error?: string | null
  onRetry?: () => void
}

/**
 * Displays the current WebSocket connection status with visual feedback
 * Shows appropriate icons and colors for each connection state
 */
export function StreamStatus({ status, error, onRetry }: StreamStatusProps) {
  const statusConfig = {
    connecting: {
      icon: Loader2,
      text: 'Connecting...',
      className: 'text-yellow-600 dark:text-yellow-500',
      iconClassName: 'animate-spin',
    },
    connected: {
      icon: CircleDot,
      text: 'Live',
      className: 'text-green-600 dark:text-green-500',
      iconClassName: 'animate-pulse',
    },
    disconnected: {
      icon: Circle,
      text: 'Disconnected',
      className: 'text-gray-500 dark:text-gray-400',
      iconClassName: '',
    },
    error: {
      icon: CircleAlert,
      text: error || 'Error',
      className: 'text-red-600 dark:text-red-500',
      iconClassName: '',
    },
  }

  const config = statusConfig[status]
  const Icon = config.icon

  return (
    <div className="flex items-center gap-2">
      <div className={cn("flex items-center gap-1.5 text-sm font-medium", config.className)}>
        <Icon className={cn("h-4 w-4", config.iconClassName)} />
        <span>{config.text}</span>
      </div>
      {status === 'error' && onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry} className="h-7 px-2">
          Retry
        </Button>
      )}
    </div>
  )
}

