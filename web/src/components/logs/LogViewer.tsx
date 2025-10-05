"use client"

import { useEffect, useRef } from "react"
import { Badge } from "@/components/ui/badge"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Loader2 } from "lucide-react"

export interface LogEntry {
  timestamp: string
  message: string
  severity: string
  serviceName?: string
}

interface LogViewerProps {
  logs: LogEntry[]
  isLoading?: boolean
  emptyMessage?: string
  showServiceName?: boolean
  followLogs?: boolean
}

export function LogViewer({ 
  logs, 
  isLoading = false, 
  emptyMessage = "No logs found",
  showServiceName = false,
  followLogs = false
}: LogViewerProps) {
  const scrollAreaRef = useRef<HTMLDivElement>(null)
  const logsEndRef = useRef<HTMLDivElement>(null)

  // Auto-scroll to bottom when followLogs is enabled and logs change
  useEffect(() => {
    if (followLogs && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [logs, followLogs])

  const getSeverityColor = (severity: string) => {
    switch (severity.toUpperCase()) {
      case "ERROR":
        return "bg-red-500/20 text-red-700 border-red-500/30"
      case "WARN":
        return "bg-yellow-500/20 text-yellow-700 border-yellow-500/30"
      case "INFO":
        return "bg-blue-500/20 text-blue-700 border-blue-500/30"
      case "DEBUG":
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
      default:
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
    }
  }

  return (
    <ScrollArea className="flex-1 border rounded-md min-h-0" ref={scrollAreaRef}>
      {isLoading ? (
        <div className="flex items-center justify-center h-full min-h-[400px]">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : logs.length === 0 ? (
        <div className="flex items-center justify-center h-full min-h-[400px] text-muted-foreground">
          {emptyMessage}
        </div>
      ) : (
        <div className="p-4 space-y-2 font-mono text-sm">
          {logs.map((log, index) => (
            <div key={index} className="flex gap-3 items-start hover:bg-muted/50 p-2 rounded">
              <span className="text-xs text-muted-foreground whitespace-nowrap">
                {new Date(log.timestamp).toLocaleTimeString()}
              </span>
              {showServiceName && log.serviceName && (
                <span className="text-xs text-muted-foreground whitespace-nowrap px-2 py-1 bg-muted rounded">
                  {log.serviceName}
                </span>
              )}
              <Badge className={getSeverityColor(log.severity)} variant="outline">
                {log.severity}
              </Badge>
              <span className="flex-1 break-all">{log.message}</span>
            </div>
          ))}
          <div ref={logsEndRef} />
        </div>
      )}
    </ScrollArea>
  )
}

