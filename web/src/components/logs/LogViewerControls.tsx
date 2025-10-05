"use client"

import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { ArrowUp, ArrowDown } from "lucide-react"

export interface LogViewerControlsProps {
  autoScroll?: boolean
  onToggleAutoScroll?: () => void
  onJumpToTop?: () => void
  onJumpToBottom?: () => void
  logCount?: number
}

export function LogViewerControls({
  autoScroll = false,
  onToggleAutoScroll,
  onJumpToTop,
  onJumpToBottom,
  logCount = 0,
}: LogViewerControlsProps) {
  return (
    <div className="flex items-center gap-2 sm:gap-4 flex-wrap">
      {/* Log count */}
      <span className="text-[10px] sm:text-xs text-muted-foreground font-mono">
        {logCount.toLocaleString()} {logCount === 1 ? "log" : "logs"}
      </span>

      {/* Auto-scroll toggle */}
      {onToggleAutoScroll && (
        <div className="flex items-center gap-2">
          <Switch
            id="auto-scroll"
            checked={autoScroll}
            onCheckedChange={onToggleAutoScroll}
            className="scale-75 sm:scale-100"
          />
          <Label htmlFor="auto-scroll" className="text-[10px] sm:text-xs cursor-pointer whitespace-nowrap">
            Auto-scroll
          </Label>
        </div>
      )}

      {/* Jump buttons */}
      <div className="flex items-center gap-1">
        {onJumpToTop && (
          <Button
            variant="outline"
            size="sm"
            className="h-6 sm:h-7 px-1.5 sm:px-2"
            onClick={onJumpToTop}
          >
            <ArrowUp className="h-3 w-3" />
            <span className="hidden sm:inline ml-1">Top</span>
          </Button>
        )}
        {onJumpToBottom && (
          <Button
            variant="outline"
            size="sm"
            className="h-6 sm:h-7 px-1.5 sm:px-2"
            onClick={onJumpToBottom}
          >
            <ArrowDown className="h-3 w-3" />
            <span className="hidden sm:inline ml-1">Bottom</span>
          </Button>
        )}
      </div>
    </div>
  )
}

