"use client"

import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

export interface LogSeverityBadgeProps {
  severity: string
}

export function LogSeverityBadge({ severity }: LogSeverityBadgeProps) {
  const severityUpper = severity.toUpperCase()
  
  const severityConfig = {
    ERROR: {
      className: "bg-red-500/20 text-red-700 border-red-500/40 hover:bg-red-500/30",
      label: "ERROR"
    },
    WARN: {
      className: "bg-yellow-500/20 text-yellow-700 border-yellow-500/40 hover:bg-yellow-500/30",
      label: "WARN"
    },
    INFO: {
      className: "bg-blue-500/20 text-blue-700 border-blue-500/40 hover:bg-blue-500/30",
      label: "INFO"
    },
    DEBUG: {
      className: "bg-gray-500/20 text-gray-700 border-gray-500/40 hover:bg-gray-500/30",
      label: "DEBUG"
    },
    TRACE: {
      className: "bg-purple-500/20 text-purple-700 border-purple-500/40 hover:bg-purple-500/30",
      label: "TRACE"
    },
    UNKNOWN: {
      className: "bg-muted/50 text-muted-foreground border-muted hover:bg-muted/70",
      label: "UNKNOWN"
    }
  }

  const config = severityConfig[severityUpper as keyof typeof severityConfig] || severityConfig.UNKNOWN

  return (
    <Badge 
      variant="outline" 
      className={cn("text-xs font-mono flex-shrink-0 px-1.5 py-0", config.className)}
    >
      {config.label}
    </Badge>
  )
}

