"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { useDeleteRailwayEnvironment } from "@/hooks/useRailway"
import { MoreVertical, Trash2, Copy } from "lucide-react"
import { toast } from "sonner"
import { CreateEnvironmentDialog } from "@/components/wizard/CreateEnvironmentDialog"

export interface Service {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  type: string
  // Railway service ID for delete operations
  railwayServiceId?: string
  // Deployment configuration (optional, defaults to source_repo)
  deploymentType?: "source_repo" | "docker_image"
  sourceRepo?: string
  sourceBranch?: string
  dockerImage?: string
  imageRegistry?: string
  imageTag?: string
  // Railway service instance details (optional)
  buildCommand?: string | null
  builder?: string | null
  startCommand?: string | null
  rootDirectory?: string | null
  region?: string | null
}

export interface Environment {
  id: string
  name: string
  status: "operational" | "maintenance" | "error"
  url: string
  services: Service[]
}

interface EnvironmentCardProps {
  environment: Environment
  isSelected: boolean
  onSelect: () => void
  projectId?: string
}

export function EnvironmentCard({ environment, isSelected, onSelect, projectId }: EnvironmentCardProps) {
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [showCloneDialog, setShowCloneDialog] = useState(false)
  const deleteEnvironment = useDeleteRailwayEnvironment()

  const getStatusColor = (status: string) => {
    switch (status) {
      case "operational":
        return "bg-emerald-500/20 text-emerald-700 border-emerald-500/30"
      case "maintenance":
        return "bg-amber-500/20 text-amber-700 border-amber-500/30"
      case "error":
        return "bg-red-500/20 text-red-700 border-red-500/30"
      default:
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
    }
  }

  const handleDelete = async () => {
    try {
      await deleteEnvironment.mutateAsync(environment.id)
      toast.success(`Environment "${environment.name}" deleted successfully`)
      setShowDeleteDialog(false)
    } catch (error) {
      toast.error(`Failed to delete environment: ${error instanceof Error ? error.message : "Unknown error"}`)
    }
  }

  const runningServices = environment.services.filter((s) => s.status === "running").length

  return (
    <>
      <Card
        className={`glass grain cursor-pointer transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01] ${
          isSelected ? "ring-2 ring-accent" : ""
        }`}
        onClick={onSelect}
      >
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg">{environment.name}</CardTitle>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className={getStatusColor(environment.status)}>{environment.status}</Badge>
              <DropdownMenu>
                <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
                  <Button variant="ghost" size="icon" className="h-8 w-8" aria-label="More actions">
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" onClick={(e) => e.stopPropagation()}>
                  <DropdownMenuItem
                    onSelect={(e) => {
                      e.stopPropagation()
                      onSelect()
                    }}
                  >
                    View details
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onSelect={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      setShowCloneDialog(true)
                    }}
                  >
                    <Copy className="h-4 w-4 mr-2" />
                    Clone environment
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="text-destructive focus:text-destructive"
                    onSelect={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      setShowDeleteDialog(true)
                    }}
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    Delete environment
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="text-sm text-muted-foreground">
            <div className="truncate">{environment.url}</div>
          </div>

          <div className="flex items-center justify-between">
            <div className="text-sm">
              <span className="font-medium">{runningServices}</span>
              <span className="text-muted-foreground">/{environment.services.length} services</span>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="bg-transparent"
              onClick={(e) => {
                e.stopPropagation()
                onSelect()
              }}
            >
              Manage
            </Button>
          </div>

          {/* Service Status Indicators */}
          <div className="flex gap-1">
            {environment.services.map((service) => (
              <div
                key={service.id}
                className={`w-2 h-2 rounded-full ${
                  service.status === "running"
                    ? "bg-emerald-500"
                    : service.status === "error"
                      ? "bg-red-500"
                      : "bg-gray-400"
                }`}
                title={`${service.name}: ${service.status}`}
              />
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Clone Dialog - hidden trigger */}
      <div className="hidden">
        <CreateEnvironmentDialog
          trigger={
            <button
              ref={(el) => {
                if (el && showCloneDialog) {
                  el.click();
                  setShowCloneDialog(false);
                }
              }}
            />
          }
          defaultSourceMode="clone"
          defaultCloneSourceId={environment.id}
          defaultProjectId={projectId}
        />
      </div>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent onClick={(e) => e.stopPropagation()}>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Environment</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete <strong>{environment.name}</strong>? This will permanently delete the environment{environment.services.length > 0 ? ` and all ${environment.services.length} service${environment.services.length === 1 ? "" : "s"}` : ""} from Railway and the database. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteEnvironment.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.preventDefault()
                e.stopPropagation()
                void handleDelete()
              }}
              disabled={deleteEnvironment.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteEnvironment.isPending ? "Deleting..." : "Delete Environment"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
