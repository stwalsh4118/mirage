import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { DeploymentTypeBadge, type DeploymentType } from "@/components/service/DeploymentTypeBadge"

interface Service {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  type: string
  // Deployment configuration (optional, defaults to source_repo)
  deploymentType?: DeploymentType
  sourceRepo?: string
  sourceBranch?: string
  dockerImage?: string
  imageRegistry?: string
  imageTag?: string
}

interface ServicesListProps {
  services: Service[]
}

export function ServicesList({ services }: ServicesListProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "running":
        return "bg-emerald-500/20 text-emerald-700 border-emerald-500/30"
      case "stopped":
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
      case "error":
        return "bg-red-500/20 text-red-700 border-red-500/30"
      default:
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case "web":
        return "🌐"
      case "api":
        return "⚡"
      case "database":
        return "🗄️"
      case "cache":
        return "💾"
      default:
        return "📦"
    }
  }

  // Helper to get deployment details text
  const getDeploymentDetails = (service: Service) => {
    const deploymentType = service.deploymentType || "source_repo"
    
    if (deploymentType === "docker_image") {
      if (service.dockerImage) {
        return service.dockerImage
      }
      // Build image reference from parts
      if (service.imageRegistry && service.imageRegistry !== "docker.io") {
        return `${service.imageRegistry}/${service.imageTag || "latest"}`
      }
      return service.imageTag || "latest"
    }
    
    // Source repo deployment
    if (service.sourceRepo && service.sourceBranch) {
      return `${service.sourceRepo}@${service.sourceBranch}`
    }
    if (service.sourceRepo) {
      return service.sourceRepo
    }
    return null
  }

  return (
    <div className="space-y-3">
      {services.map((service) => {
        const deploymentDetails = getDeploymentDetails(service)
        const deploymentType = service.deploymentType || "source_repo"
        
        return (
          <div
            key={service.id}
            className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-border/50"
          >
            <div className="flex items-center gap-3 flex-1 min-w-0">
              <span className="text-lg flex-shrink-0">{getTypeIcon(service.type)}</span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <div className="font-medium">{service.name}</div>
                  <DeploymentTypeBadge type={deploymentType} />
                </div>
                <div className="text-sm text-muted-foreground capitalize">
                  {service.type}
                </div>
                {deploymentDetails && (
                  <div className="text-xs text-muted-foreground truncate mt-0.5">
                    {deploymentDetails}
                  </div>
                )}
              </div>
            </div>

            <div className="flex items-center gap-3 flex-shrink-0">
              <Badge className={getStatusColor(service.status)}>{service.status}</Badge>
              <div className="flex gap-1">
                {service.status === "running" ? (
                  <Button variant="outline" size="sm">
                    Stop
                  </Button>
                ) : (
                  <Button size="sm">Start</Button>
                )}
                <Button variant="ghost" size="sm">
                  Logs
                </Button>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
