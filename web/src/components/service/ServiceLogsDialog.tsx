"use client"

import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Download, Search, RefreshCw, Loader2, ArrowDown } from "lucide-react"
import { toast } from "sonner"
import { useServiceLogs } from "@/hooks/useLogs"
import { exportServiceLogs, type ExportFormat } from "@/lib/api/logs"
import { LogViewer } from "@/components/logs/LogViewer"

interface ServiceLogsDialogProps {
  serviceId: string
  serviceName: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ServiceLogsDialog({ serviceId, serviceName, open, onOpenChange }: ServiceLogsDialogProps) {
  const [searchQuery, setSearchQuery] = useState("")
  const [searchInput, setSearchInput] = useState("")
  const [limit, setLimit] = useState(100)
  const [followLogs, setFollowLogs] = useState(false)

  const { data, isLoading, refetch } = useServiceLogs(
    {
      serviceId,
      limit,
      search: searchQuery,
    },
    open
  )

  const logs = data?.logs || []

  const handleSearch = () => {
    setSearchQuery(searchInput)
  }

  const handleExport = async (format: ExportFormat) => {
    try {
      const blob = await exportServiceLogs(serviceId, format)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement("a")
      a.href = url
      a.download = `${serviceName}-logs-${new Date().toISOString().split("T")[0]}.${format}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      
      toast.success(`Logs exported as ${format.toUpperCase()}`)
    } catch (error) {
      console.error("Failed to export logs:", error)
      toast.error(error instanceof Error ? error.message : "Failed to export logs")
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-[85vw] max-[640px]:w-[95vw] max-w-[1800px] sm:max-w-[85vw] md:max-w-[1800px] h-[85vh] max-[640px]:h-[90vh] flex flex-col gap-2 overflow-hidden p-4 max-[640px]:p-3">
        <div className="flex items-center justify-between gap-4 flex-shrink-0 pb-2 border-b border-border/50">
          <div className="flex items-baseline gap-3">
            <DialogTitle className="text-base font-semibold">Service Logs: {serviceName}</DialogTitle>
            <span className="text-xs text-muted-foreground font-mono">
              {logs.length} logs
            </span>
          </div>
          <div className="flex gap-1.5 items-center mr-8">
            <Button onClick={() => refetch()} disabled={isLoading} variant="ghost" size="sm" className="h-7 px-2">
              {isLoading ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <RefreshCw className="h-3.5 w-3.5" />}
            </Button>
            <Button 
              onClick={() => setFollowLogs(!followLogs)} 
              variant={followLogs ? "default" : "ghost"}
              size="sm"
              className="h-7 px-2"
            >
              <ArrowDown className="h-3.5 w-3.5" />
            </Button>
            <Button onClick={() => handleExport("json")} variant="ghost" size="sm" className="hidden sm:inline-flex h-7 px-2">
              <Download className="h-3.5 w-3.5 mr-1" />
              JSON
            </Button>
            <Button onClick={() => handleExport("csv")} variant="ghost" size="sm" className="hidden sm:inline-flex h-7 px-2">
              <Download className="h-3.5 w-3.5 mr-1" />
              CSV
            </Button>
            <Button onClick={() => handleExport("txt")} variant="ghost" size="sm" className="hidden sm:inline-flex h-7 px-2">
              <Download className="h-3.5 w-3.5 mr-1" />
              TXT
            </Button>
          </div>
        </div>

        <div className="flex gap-2 items-center flex-shrink-0">
          <div className="flex-1 relative">
            <Search className="absolute left-2 top-2 h-3.5 w-3.5 text-muted-foreground" />
            <Input
              placeholder="Search logs..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="pl-7 h-8 text-sm"
            />
          </div>
        </div>

        <div className="flex-1 min-h-0">
          <LogViewer 
            logs={logs} 
            loading={isLoading}
            searchQuery={searchQuery}
            autoScroll={followLogs}
            onToggleAutoScroll={() => setFollowLogs(!followLogs)}
            maxHeight="100%"
            hideHeader={true}
          />
        </div>
      </DialogContent>
    </Dialog>
  )
}

