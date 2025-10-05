"use client"

import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

interface ServiceBadgeProps {
  serviceName: string
}

// Generate consistent color for service name
function getServiceColor(serviceName: string): string {
  const colors = [
    "bg-cyan-500/20 text-cyan-700 border-cyan-500/40",
    "bg-emerald-500/20 text-emerald-700 border-emerald-500/40",
    "bg-violet-500/20 text-violet-700 border-violet-500/40",
    "bg-pink-500/20 text-pink-700 border-pink-500/40",
    "bg-orange-500/20 text-orange-700 border-orange-500/40",
    "bg-teal-500/20 text-teal-700 border-teal-500/40",
    "bg-indigo-500/20 text-indigo-700 border-indigo-500/40",
  ]
  
  // Simple hash function for consistent color assignment
  let hash = 0
  for (let i = 0; i < serviceName.length; i++) {
    hash = serviceName.charCodeAt(i) + ((hash << 5) - hash)
  }
  
  return colors[Math.abs(hash) % colors.length]
}

export function ServiceBadge({ serviceName }: ServiceBadgeProps) {
  if (!serviceName || serviceName === "unknown") {
    return (
      <Badge 
        variant="outline" 
        className="text-xs font-mono flex-shrink-0 px-1.5 py-0 bg-muted/50 text-muted-foreground"
      >
        unknown
      </Badge>
    )
  }

  return (
    <Badge 
      variant="outline" 
      className={cn(
        "text-xs font-mono flex-shrink-0 px-1.5 py-0",
        getServiceColor(serviceName)
      )}
    >
      {serviceName}
    </Badge>
  )
}

