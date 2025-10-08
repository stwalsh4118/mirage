"use client"

import { useMemo, useState } from "react"
import { useRouter } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
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
import { ProjectHeader } from "./project-header"
import { EnvironmentCard, type Environment as UIEnvironment, type Service as UIService } from "./environment-card"
import { ServicesList } from "./services-list"
import { StatusOverview } from "./status-overview"
import { EnvironmentMetadata } from "@/components/environment/EnvironmentMetadata"
import { ServiceBuildInfo } from "@/components/environment/ServiceBuildInfo"

import { useRailwayProjectsDetails, useDeleteRailwayProject, useEnvironmentServices } from "@/hooks/useRailway"
import { toast } from "sonner"
import type { RailwayProjectDetails } from "@/lib/api/railway"

interface ProjectDetailProps {
  projectId: string
}


export function ProjectDetail({ projectId }: ProjectDetailProps) {
  const [selectedEnvironment, setSelectedEnvironment] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState("")
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)

  const { data: projects = [], isLoading, isError, refetch } = useRailwayProjectsDetails()
  const deleteProject = useDeleteRailwayProject()
  const router = useRouter()
  
  // Fetch service build configurations for the selected environment
  const { data: serviceBuildConfigs = [] } = useEnvironmentServices(selectedEnvironment)

  const project: RailwayProjectDetails | undefined = useMemo(
    () => projects.find((p) => p.id === projectId || p.name === projectId),
    [projects, projectId],
  )

  const environments: UIEnvironment[] =
    useMemo(() => {
      if (!project) return []
      // Augment with plausible defaults for fields our UI expects (status/url)
      return project.environments.map((e) => ({
        ...e,
        status: "operational" as const,
        url: `${project.name}-${e.name}.mirage.dev`,
        services: (e.services || []).map((s) => {
          // Determine deployment type from source field
          const deploymentType = s.source?.image ? "docker_image" : "source_repo"
          
          return {
            id: s.id, 
            name: s.serviceName, 
            status: "running", 
            type: "web",
            // Railway service ID for deletion operations
            railwayServiceId: s.serviceId,
            // Deployment type and configuration
            deploymentType,
            sourceRepo: s.source?.repo,
            dockerImage: s.source?.image,
            // Pass through detailed service instance fields
            buildCommand: s.buildCommand,
            builder: s.builder,
            startCommand: s.startCommand,
            rootDirectory: s.rootDirectory,
            region: s.region,
          } as UIService
        }),
      }))
    }, [project])

  const selectedEnv = environments.find((env) => env.id === selectedEnvironment)

  const filteredEnvironments = environments.filter((env) =>
    env.name.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const handleDelete = async () => {
    if (!project) return
    try {
      await deleteProject.mutateAsync(project.id)
      toast.success(`"${project.name}" has been permanently deleted from Railway.`)
      setShowDeleteDialog(false)
      // Navigate after invalidation completes
      router.push("/dashboard")
    } catch (error) {
      toast.error((error as Error).message || "Could not delete project. Please try again.")
    }
  }

  return (
    <div className="container mx-auto px-4 py-6 space-y-6">
      {isLoading && (
        <div className="glass grain p-6 rounded-lg text-sm text-muted-foreground">Loading project…</div>
      )}
      {isError && (
        <div className="glass grain p-6 rounded-lg text-sm text-red-600">Failed to load project. <button className="underline" onClick={() => void refetch()}>Retry</button></div>
      )}
      {project && (
        <ProjectHeader project={{ id: project.id, name: project.name, environments }} />
      )}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-3 space-y-6">
          {/* Search and Filters */}
          <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
            <CardContent className="p-4">
              <div className="flex flex-col sm:flex-row gap-4 items-center">
                <Input
                  placeholder="Search environments..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="flex-1"
                />
                <div className="flex gap-2">
                  <Button variant="outline" size="sm" className="bg-transparent">
                    All
                  </Button>
                  <Button variant="outline" size="sm" className="bg-transparent">
                    Operational
                  </Button>
                  <Button variant="outline" size="sm" className="bg-transparent">
                    Maintenance
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Environments Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {filteredEnvironments.map((environment) => (
              <EnvironmentCard
                key={environment.id}
                environment={environment}
                isSelected={selectedEnvironment === environment.id}
                onSelect={() => setSelectedEnvironment(selectedEnvironment === environment.id ? null : environment.id)}
                projectId={project?.id}
              />
            ))}
          </div>

          {/* Environment Metadata */}
          {selectedEnv && (
            <EnvironmentMetadata 
              environmentId={selectedEnv.id}
              onSaveAsTemplate={() => {
                toast.info("Save as Template feature coming in PBI 15")
              }}
            />
          )}

          {/* Services Detail */}
          {selectedEnv && (
            <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  Services in {selectedEnv.name}
                  <Badge variant="secondary" className="bg-accent/15 text-accent">
                    {selectedEnv.services.length}
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ServicesList 
                  services={selectedEnv.services} 
                  mirageServices={serviceBuildConfigs}
                />
              </CardContent>
            </Card>
          )}

          {/* Service Build Configurations */}
          {selectedEnv && serviceBuildConfigs.length > 0 && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Service Build Configurations</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {serviceBuildConfigs.map((service) => (
                  <ServiceBuildInfo key={service.id} service={service} />
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {project && (
            <StatusOverview project={{ environments }} />
          )}

          <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
            <CardHeader>
              <CardTitle className="text-sm">Quick Actions</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <Button className="w-full" size="sm">
                Deploy All
              </Button>
              <Button variant="outline" className="w-full bg-transparent" size="sm">
                View Logs
              </Button>
              <Button variant="outline" className="w-full bg-transparent" size="sm">
                Settings
              </Button>
              <Button
                variant="outline"
                className="w-full bg-transparent text-destructive hover:text-destructive border-destructive/30 hover:bg-destructive/10"
                size="sm"
                onClick={() => setShowDeleteDialog(true)}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete Project
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      {project && (
        <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Delete Project</AlertDialogTitle>
              <AlertDialogDescription>
                Are you sure you want to delete <strong>{project.name}</strong>? This will permanently delete the project, all {environments.length} environment{environments.length === 1 ? "" : "s"}, and all associated services from Railway and the database. This action cannot be undone.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={deleteProject.isPending}>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={handleDelete}
                disabled={deleteProject.isPending}
                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              >
                {deleteProject.isPending ? "Deleting..." : "Delete Project"}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </div>
  )
}
