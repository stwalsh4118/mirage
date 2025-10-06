"use client"

import { useRef, useEffect, useState } from "react"
import { useVirtualizer } from "@tanstack/react-virtual"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import { LogLine } from "./LogLine"
import { LogViewerControls } from "./LogViewerControls"
import { EmptyLogState } from "./EmptyLogState"
import type { LogViewerProps } from "./types"
import { cn } from "@/lib/utils"

const DEFAULT_MAX_HEIGHT = "600px"
const ESTIMATED_LOG_LINE_HEIGHT = 32
const SCROLL_THRESHOLD = 10 // Pixels threshold for scroll position detection

export function LogViewer({
  logs,
  loading = false,
  onLoadMore,
  searchQuery,
  autoScroll = false,
  onToggleAutoScroll,
  maxHeight = DEFAULT_MAX_HEIGHT,
  hideHeader = false,
}: LogViewerProps) {
  const parentRef = useRef<HTMLDivElement>(null)
  const [userScrolled, setUserScrolled] = useState(false)

  const virtualizer = useVirtualizer({
    count: logs.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ESTIMATED_LOG_LINE_HEIGHT,
    overscan: 10,
  })

  // Auto-scroll to bottom when new logs arrive (if auto-scroll is enabled and user hasn't manually scrolled)
  useEffect(() => {
    if (autoScroll && !userScrolled && parentRef.current && logs.length > 0) {
      // Use double RAF to ensure virtualizer has measured and rendered
      // First RAF: Let React commit the changes
      requestAnimationFrame(() => {
        // Second RAF: Let the browser paint before scrolling
        requestAnimationFrame(() => {
          if (parentRef.current) {
            parentRef.current.scrollTop = parentRef.current.scrollHeight
          }
        })
      })
    }
  }, [logs.length, autoScroll, userScrolled, logs])

  // Reset userScrolled when auto-scroll is toggled on
  useEffect(() => {
    if (autoScroll) {
      setUserScrolled(false)
    }
  }, [autoScroll])

  // Detect when user manually scrolls
  const handleScroll = () => {
    if (!parentRef.current) return
    
    const { scrollTop, scrollHeight, clientHeight } = parentRef.current
    const distanceFromBottom = scrollHeight - clientHeight - scrollTop
    const isAtBottom = distanceFromBottom <= SCROLL_THRESHOLD
    
    if (!isAtBottom && autoScroll) {
      setUserScrolled(true)
    }

    // Trigger load more when scrolling to top
    if (scrollTop <= SCROLL_THRESHOLD && onLoadMore) {
      onLoadMore()
    }
  }

  // Jump to top/bottom handlers
  const handleJumpToTop = () => {
    if (parentRef.current) {
      parentRef.current.scrollTo({ top: 0, behavior: "smooth" })
      setUserScrolled(true)
    }
  }

  const handleJumpToBottom = () => {
    if (parentRef.current) {
      parentRef.current.scrollTo({ top: parentRef.current.scrollHeight, behavior: "smooth" })
      setUserScrolled(false)
    }
  }

  return (
    <Card className="glass grain flex flex-col h-full">
      {!hideHeader && (
        <CardHeader className="pb-2 sm:pb-3 flex-shrink-0 px-3 sm:px-6 pt-3 sm:pt-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-0">
            <CardTitle className="text-sm sm:text-base text-foreground/90">Logs</CardTitle>
            <LogViewerControls
              autoScroll={autoScroll}
              onToggleAutoScroll={onToggleAutoScroll}
              onJumpToTop={handleJumpToTop}
              onJumpToBottom={handleJumpToBottom}
              logCount={logs.length}
            />
          </div>
        </CardHeader>
      )}
      <CardContent className="flex-1 flex flex-col min-h-0 p-0">
        {loading && logs.length === 0 ? (
          <div 
            className="flex items-center justify-center bg-muted/30 rounded-b-lg border-t border-border/50"
            style={{ height: maxHeight }}
          >
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : logs.length === 0 ? (
          <div 
            className="flex items-center justify-center bg-muted/30 rounded-b-lg border-t border-border/50"
            style={{ height: maxHeight }}
          >
            <EmptyLogState />
          </div>
        ) : (
          <div
            ref={parentRef}
            className={cn(
              "overflow-auto bg-muted/30 rounded-b-lg border-t border-border/50 font-mono text-sm",
              "scrollbar-thin scrollbar-thumb-muted-foreground/20 scrollbar-track-transparent"
            )}
            style={{ height: maxHeight }}
            onScroll={handleScroll}
          >
            <div
              style={{
                height: `${virtualizer.getTotalSize()}px`,
                width: "100%",
                position: "relative",
              }}
            >
              {virtualizer.getVirtualItems().map((virtualRow) => {
                const log = logs[virtualRow.index]
                return (
                  <div
                    key={virtualRow.index}
                    data-index={virtualRow.index}
                    ref={virtualizer.measureElement}
                    style={{
                      position: "absolute",
                      top: 0,
                      left: 0,
                      width: "100%",
                      transform: `translateY(${virtualRow.start}px)`,
                    }}
                  >
                    <LogLine
                      log={log}
                      lineNumber={virtualRow.index + 1}
                      index={virtualRow.index}
                      searchQuery={searchQuery}
                    />
                  </div>
                )
              })}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

