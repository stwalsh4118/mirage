"use client"

import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
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
      <DialogContent className="w-[70vw] max-w-[1800px] h-[85vh] flex flex-col gap-4 overflow-hidden sm:max-w-[95vw]">
        <DialogHeader className="flex-shrink-0">
          <DialogTitle>Service Logs: {serviceName}</DialogTitle>
          <DialogDescription>
            Viewing the latest {limit} log entries for this service
          </DialogDescription>
        </DialogHeader>

        <div className="flex gap-2 items-center flex-shrink-0">
          <div className="flex-1 relative">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search logs..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="pl-8"
            />
          </div>
          <Button onClick={() => refetch()} disabled={isLoading} variant="outline">
            {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
          </Button>
          <Button 
            onClick={() => setFollowLogs(!followLogs)} 
            variant={followLogs ? "default" : "outline"}
            size="sm"
          >
            <ArrowDown className="h-4 w-4 mr-1" />
            Follow
          </Button>
          <Button onClick={() => handleExport("json")} variant="outline" size="sm">
            <Download className="h-4 w-4 mr-1" />
            JSON
          </Button>
          <Button onClick={() => handleExport("csv")} variant="outline" size="sm">
            <Download className="h-4 w-4 mr-1" />
            CSV
          </Button>
          <Button onClick={() => handleExport("txt")} variant="outline" size="sm">
            <Download className="h-4 w-4 mr-1" />
            TXT
          </Button>
        </div>

        <LogViewer 
          logs={logs} 
          isLoading={isLoading}
          emptyMessage="No logs found for this service"
          followLogs={followLogs}
        />

        <div className="text-xs text-muted-foreground flex-shrink-0">
          Showing {logs.length} of {limit} logs
        </div>
      </DialogContent>
    </Dialog>
  )
}

