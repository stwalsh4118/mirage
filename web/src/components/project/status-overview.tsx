import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

interface Project {
  environments: Array<{
    status: string
    services: Array<{ status: string }>
  }>
}

interface StatusOverviewProps {
  project: Project
}

export function StatusOverview({ project }: StatusOverviewProps) {
  const totalServices = project.environments.reduce((acc, env) => acc + env.services.length, 0)
  const runningServices = project.environments.reduce(
    (acc, env) => acc + env.services.filter((s) => s.status === "running").length,
    0,
  )
  const operationalEnvs = project.environments.filter((env) => env.status === "operational").length

  const healthPercentage = totalServices > 0 ? (runningServices / totalServices) * 100 : 0

  return (
    <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
      <CardHeader>
        <CardTitle className="text-sm">Project Health</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span>Overall Health</span>
            <span>{Math.round(healthPercentage)}%</span>
          </div>
          <Progress value={healthPercentage} className="h-2" />
        </div>

        <div className="grid grid-cols-2 gap-4 text-center">
          <div className="space-y-1">
            <div className="text-2xl font-bold text-emerald-600">{operationalEnvs}</div>
            <div className="text-xs text-muted-foreground">Operational</div>
          </div>
          <div className="space-y-1">
            <div className="text-2xl font-bold text-accent">{runningServices}</div>
            <div className="text-xs text-muted-foreground">Running Services</div>
          </div>
        </div>

        <div className="pt-2 border-t border-border/50">
          <div className="text-xs text-muted-foreground space-y-1">
            <div className="flex justify-between">
              <span>Total Environments:</span>
              <span>{project.environments.length}</span>
            </div>
            <div className="flex justify-between">
              <span>Total Services:</span>
              <span>{totalServices}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
