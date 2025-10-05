"use client"

import { memo, useState } from "react"
import { LogSeverityBadge } from "./LogSeverityBadge"
import { ServiceBadge } from "./ServiceBadge"
import { Button } from "@/components/ui/button"
import { Copy, ChevronDown, ChevronUp } from "lucide-react"
import { cn } from "@/lib/utils"
import { toast } from "sonner"
import type { Log } from "./types"

const MAX_MESSAGE_LENGTH = 200

interface LogLineProps {
  log: Log
  lineNumber: number
  index: number // Stable index from data array for alternating backgrounds
  searchQuery?: string
}

export const LogLine = memo(function LogLine({ log, lineNumber, index, searchQuery }: LogLineProps) {
  const [expanded, setExpanded] = useState(false)
  const messageToDisplay = log.message || log.rawLine || ""
  const isLongMessage = messageToDisplay.length > MAX_MESSAGE_LENGTH

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(log.rawLine)
      toast.success("Copied to clipboard")
    } catch (error) {
      console.error("Failed to copy to clipboard:", error)
      toast.error("Copy failed — try using keyboard shortcuts")
    }
  }

  const getSeverityBorderColor = (severity: string) => {
    switch (severity.toUpperCase()) {
      case "ERROR":
        return "border-l-red-500"
      case "WARN":
        return "border-l-yellow-500"
      case "INFO":
        return "border-l-blue-500"
      case "DEBUG":
        return "border-l-gray-500"
      case "TRACE":
        return "border-l-purple-500"
      default:
        return "border-l-transparent"
    }
  }

  const getSeverityBgColor = (severity: string) => {
    switch (severity.toUpperCase()) {
      case "ERROR":
        return "bg-red-500/5 hover:bg-red-500/10"
      case "WARN":
        return "bg-yellow-500/5 hover:bg-yellow-500/10"
      case "INFO":
        return "bg-blue-500/5 hover:bg-blue-500/10"
      case "DEBUG":
        return "bg-muted/30 hover:bg-muted/40"
      case "TRACE":
        return "bg-purple-500/5 hover:bg-purple-500/10"
      default:
        return "hover:bg-accent/10"
    }
  }

  const escapeRegExp = (str: string): string => {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  }

  const highlightSearch = (text: string, query?: string) => {
    if (!query) return text

    // Guard against overly long queries to prevent ReDoS
    const maxQueryLength = 100
    const safeQuery = query.length > maxQueryLength ? query.substring(0, maxQueryLength) : query

    try {
      const escapedQuery = escapeRegExp(safeQuery.trim())
      if (!escapedQuery) return text
      
      const parts = text.split(new RegExp(`(${escapedQuery})`, "gi"))
      return parts.map((part, i) =>
        part.toLowerCase() === safeQuery.toLowerCase() ? (
          <mark key={i} className="bg-yellow-300/60 text-foreground">
            {part}
          </mark>
        ) : (
          part
        )
      )
    } catch {
      // If regex is invalid, just return the text
      return text
    }
  }

  const displayMessage = expanded || !isLongMessage 
    ? messageToDisplay 
    : messageToDisplay.substring(0, MAX_MESSAGE_LENGTH)

  return (
    <div
      className={cn(
        "flex flex-col sm:flex-row sm:items-start gap-1 sm:gap-2 px-2 py-1.5 group border-l-2",
        getSeverityBorderColor(log.severity),
        getSeverityBgColor(log.severity),
        index % 2 === 0 && "bg-muted/20"
      )}
    >
      {/* Mobile: First row with metadata */}
      <div className="flex items-center gap-2 sm:contents">
        {/* Line number - hidden on mobile */}
        <span className="hidden sm:inline-block text-muted-foreground text-xs w-12 text-right flex-shrink-0 select-none font-mono">
          {lineNumber}
        </span>

        {/* Timestamp */}
        <span className="text-muted-foreground text-[10px] sm:text-xs sm:w-20 flex-shrink-0 font-mono">
          {new Date(log.timestamp).toLocaleTimeString()}
        </span>

        {/* Service badge - hidden on mobile */}
        <div className="hidden sm:block">
          <ServiceBadge serviceName={log.serviceName} />
        </div>

        {/* Severity badge - hidden on mobile, border color shows severity */}
        <div className="hidden sm:block">
          <LogSeverityBadge severity={log.severity} />
        </div>

        {/* Spacer on mobile to push action buttons to the right */}
        <div className="flex-1 sm:hidden" />

        {/* Action buttons - always visible on mobile for touch */}
        <div className="flex items-center gap-1 sm:contents">
          {/* Expand button for long messages */}
          {isLongMessage && (
            <Button
              variant="ghost"
              size="sm"
              className="flex-shrink-0 h-6 w-6 p-0 sm:px-2 sm:w-auto opacity-70 hover:opacity-100"
              onClick={() => setExpanded(!expanded)}
              aria-label={expanded ? "Collapse log line" : "Expand log line"}
              aria-expanded={expanded}
            >
              {expanded ? (
                <ChevronUp className="h-3 w-3" />
              ) : (
                <ChevronDown className="h-3 w-3" />
              )}
            </Button>
          )}

          {/* Copy button - always visible on mobile, hover on desktop */}
          <Button
            variant="ghost"
            size="sm"
            className="flex-shrink-0 h-6 w-6 p-0 sm:px-2 sm:w-auto sm:opacity-0 sm:group-hover:opacity-100"
            onClick={handleCopy}
            aria-label="Copy log line"
          >
            <Copy className="h-3 w-3" />
          </Button>
        </div>
      </div>

      {/* Message - full width on mobile, flex-1 on desktop */}
      <code className="w-full sm:flex-1 whitespace-pre-wrap break-all text-[11px] sm:text-xs font-mono text-foreground">
        {highlightSearch(displayMessage, searchQuery)}
        {isLongMessage && !expanded && "..."}
      </code>
    </div>
  )
})

