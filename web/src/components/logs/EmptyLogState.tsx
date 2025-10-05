"use client"

import { FileText } from "lucide-react"

interface EmptyLogStateProps {
  message?: string
}

export function EmptyLogState({ message = "No logs available" }: EmptyLogStateProps) {
  return (
    <div className="flex flex-col items-center justify-center h-full min-h-[400px] text-muted-foreground">
      <FileText className="h-12 w-12 mb-4 opacity-50" />
      <p className="text-sm">{message}</p>
      <p className="text-xs mt-2 opacity-70">
        Logs will appear here when your services generate output
      </p>
    </div>
  )
}

