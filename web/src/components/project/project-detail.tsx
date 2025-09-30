"use client"

import { useMemo, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { ProjectHeader } from "./project-header"
import { EnvironmentCard, type Environment as UIEnvironment, type Service as UIService } from "./environment-card"
import { ServicesList } from "./services-list"
import { StatusOverview } from "./status-overview"

import { useRailwayProjectsDetails } from "@/hooks/useRailway"
import type { RailwayProjectDetails } from "@/lib/api/railway"

interface ProjectDetailProps {
  projectId: string
}


export function ProjectDetail({ projectId }: ProjectDetailProps) {
  const [selectedEnvironment, setSelectedEnvironment] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState("")

  const { data: projects = [], isLoading, isError, refetch } = useRailwayProjectsDetails()

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
        services: e.services.map((s) => ({ id: s.id, name: s.name, status: "running", type: "web" } as UIService)),
      }))
    }, [project])

  const selectedEnv = environments.find((env) => env.id === selectedEnvironment)

  const filteredEnvironments = environments.filter((env) =>
    env.name.toLowerCase().includes(searchQuery.toLowerCase()),
  )

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
              />
            ))}
          </div>

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
                <ServicesList services={selectedEnv.services} />
              </CardContent>
            </Card>
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
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
