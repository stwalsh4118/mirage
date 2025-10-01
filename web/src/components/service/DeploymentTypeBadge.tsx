import { Badge } from "@/components/ui/badge"
import { GitBranch, Package } from "lucide-react"

export type DeploymentType = "source_repo" | "docker_image"

interface DeploymentTypeBadgeProps {
  type: DeploymentType
  className?: string
}

export function DeploymentTypeBadge({ type, className = "" }: DeploymentTypeBadgeProps) {
  if (type === "docker_image") {
    return (
      <Badge 
        variant="outline" 
        className={`bg-blue-500/10 text-blue-700 dark:text-blue-400 border-blue-500/30 ${className}`}
      >
        <Package className="h-3 w-3 mr-1" />
        Image
      </Badge>
    )
  }

  return (
    <Badge 
      variant="outline" 
      className={`bg-violet-500/10 text-violet-700 dark:text-violet-400 border-violet-500/30 ${className}`}
    >
      <GitBranch className="h-3 w-3 mr-1" />
      Source
    </Badge>
  )
}

