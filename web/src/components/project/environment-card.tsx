"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

export interface Service {
  id: string
  name: string
  status: "running" | "stopped" | "error"
  type: string
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
}

export function EnvironmentCard({ environment, isSelected, onSelect }: EnvironmentCardProps) {
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

  const runningServices = environment.services.filter((s) => s.status === "running").length

  return (
    <Card
      className={`glass grain cursor-pointer transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01] ${
        isSelected ? "ring-2 ring-accent" : ""
      }`}
      onClick={onSelect}
    >
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{environment.name}</CardTitle>
          <Badge variant="secondary" className={getStatusColor(environment.status)}>{environment.status}</Badge>
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
              // Handle direct action
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
  )
}
