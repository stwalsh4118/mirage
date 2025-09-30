import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

interface Service {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  type: string
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

  return (
    <div className="space-y-3">
      {services.map((service) => (
        <div
          key={service.id}
          className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-border/50"
        >
          <div className="flex items-center gap-3">
            <span className="text-lg">{getTypeIcon(service.type)}</span>
            <div>
              <div className="font-medium">{service.name}</div>
              <div className="text-sm text-muted-foreground capitalize">{service.type}</div>
            </div>
          </div>

          <div className="flex items-center gap-3">
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
      ))}
    </div>
  )
}
