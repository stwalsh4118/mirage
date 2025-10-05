import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { DeploymentTypeBadge, type DeploymentType } from "@/components/service/DeploymentTypeBadge"
import { ServiceLogsDialog } from "@/components/service/ServiceLogsDialog"
import { 
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Trash2 } from "lucide-react"
import { useDeleteRailwayService } from "@/hooks/useRailway"
import { useState } from "react"
import { toast } from "sonner"

interface Service {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  type: string
  // Railway service ID for delete operations
  railwayServiceId?: string
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
  const [deletingServiceId, setDeletingServiceId] = useState<string | null>(null)
  const [serviceToDelete, setServiceToDelete] = useState<Service | null>(null)
  const [logsServiceId, setLogsServiceId] = useState<string | null>(null)
  const [logsServiceName, setLogsServiceName] = useState<string>("")
  const deleteService = useDeleteRailwayService()

  const handleDeleteService = async () => {
    if (!serviceToDelete?.railwayServiceId) return

    setDeletingServiceId(serviceToDelete.id)
    try {
      await deleteService.mutateAsync(serviceToDelete.railwayServiceId)
      toast.success(`Service "${serviceToDelete.name}" deleted successfully`)
    } catch (error) {
      console.error("Failed to delete service:", error)
      toast.error(`Failed to delete service: ${error instanceof Error ? error.message : "Unknown error"}`)
    } finally {
      setDeletingServiceId(null)
      setServiceToDelete(null)
    }
  }

  const handleDeleteClick = (service: Service) => {
    if (!service.railwayServiceId) {
      toast.error("Cannot delete service: Railway ID not found")
      return
    }
    setServiceToDelete(service)
  }

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
    <>
      <AlertDialog open={!!serviceToDelete} onOpenChange={(open) => !open && setServiceToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Service</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete <strong>{serviceToDelete?.name}</strong>? This will remove the service from Railway and the database. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteService}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete Service
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {logsServiceId && (
        <ServiceLogsDialog
          serviceId={logsServiceId}
          serviceName={logsServiceName}
          open={!!logsServiceId}
          onOpenChange={(open) => !open && setLogsServiceId(null)}
        />
      )}

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
                <Button 
                  variant="ghost" 
                  size="sm"
                  onClick={() => {
                    if (!service.railwayServiceId) {
                      toast.error("Cannot view logs: Railway service ID not found")
                      return
                    }
                    setLogsServiceId(service.railwayServiceId)
                    setLogsServiceName(service.name)
                  }}
                >
                  Logs
                </Button>
                <Button 
                  variant="ghost" 
                  size="sm"
                  onClick={() => handleDeleteClick(service)}
                  disabled={deletingServiceId === service.id}
                  className="text-destructive hover:text-destructive hover:bg-destructive/10"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        )
        })}
      </div>
    </>
  )
}
