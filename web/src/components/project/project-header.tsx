import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Breadcrumbs } from "@/components/Breadcrumbs"

interface Project {
  id: string
  name: string
  description?: string
  status?: "active" | "inactive" | "error"
  createdAt?: string
  environments: any[]
}

interface ProjectHeaderProps {
  project: Project
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "bg-emerald-500/20 text-emerald-700 border-emerald-500/30"
      case "inactive":
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
      case "error":
        return "bg-red-500/20 text-red-700 border-red-500/30"
      default:
        return "bg-gray-500/20 text-gray-700 border-gray-500/30"
    }
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
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div className="space-y-2">
              <div className="flex items-center gap-3">
                <h1 className="text-2xl font-bold">{project.name}</h1>
                {project.status && (
                  <Badge variant="secondary" className={getStatusColor(project.status)}>{project.status}</Badge>
                )}
              </div>
              {project.description && (
                <p className="text-muted-foreground text-sm">{project.description}</p>
              )}
              <div className="flex items-center gap-4 text-sm text-muted-foreground">
                {project.createdAt && (
                  <>
                    <span>Created {new Date(project.createdAt).toLocaleDateString()}</span>
                    <span>•</span>
                  </>
                )}
                <span>{project.environments.length} environments</span>
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm">
                View Metrics
              </Button>
              <Button size="sm">New Environment</Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
