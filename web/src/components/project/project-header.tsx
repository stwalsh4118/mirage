"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Breadcrumbs } from "@/components/Breadcrumbs"
import { ExternalLink, Layers, Server } from "lucide-react"

interface Environment {
  id: string;
  name: string;
  services?: { id: string }[];
}

interface Project {
  id: string
  name: string
  description?: string
  status?: "active" | "inactive" | "error"
  createdAt?: string
  environments: Environment[]
}

interface ProjectHeaderProps {
  project: Project
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  const totalServices = project.environments.reduce((sum, env) => sum + (env.services?.length || 0), 0)
  const runningServices = totalServices // In the future, could calculate based on actual status
  
  const handleOpenInRailway = () => {
    window.open(`https://railway.app/project/${project.id}`, '_blank', 'noopener,noreferrer')
  }

  return (
    <div className="space-y-4">
      {/* Breadcrumb */}
      <Breadcrumbs items={[
        { label: "Dashboard", href: "/dashboard" },
        { label: "Projects", href: "/dashboard" },
        { label: project.name },
      ]} />

      {/* Project Overview */}
      <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
        <CardContent className="p-6">
          <div className="flex flex-col gap-6">
            {/* Header Row */}
            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
              <div className="space-y-1">
                <h1 className="text-3xl font-bold">{project.name}</h1>
                {project.description && (
                  <p className="text-muted-foreground">{project.description}</p>
                )}
              </div>
              <Button 
                variant="outline" 
                size="sm"
                onClick={handleOpenInRailway}
                className="bg-transparent shrink-0"
              >
                <ExternalLink className="mr-2 h-4 w-4" />
                Open in Railway
              </Button>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
              <div className="space-y-1">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Layers className="h-4 w-4" />
                  <span>Environments</span>
                </div>
                <div className="text-2xl font-bold">{project.environments.length}</div>
              </div>
              
              <div className="space-y-1">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Server className="h-4 w-4" />
                  <span>Total Services</span>
                </div>
                <div className="text-2xl font-bold">{totalServices}</div>
              </div>
              
              <div className="space-y-1">
                <div className="text-sm text-muted-foreground">Status</div>
                <div className="flex items-center gap-2">
                  <div className="text-2xl font-bold text-green-600 dark:text-green-400">{runningServices}</div>
                  <Badge variant="secondary" className="bg-green-500/20 text-green-700 dark:text-green-400 border-green-500/30">
                    Healthy
                  </Badge>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
