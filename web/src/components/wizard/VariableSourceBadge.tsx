import { Badge } from "@/components/ui/badge";
import { Globe, Package, Settings } from "lucide-react";

export type VariableSource = "global" | "service" | "system";

interface VariableSourceBadgeProps {
  source: VariableSource;
  className?: string;
}

export function VariableSourceBadge({ source, className }: VariableSourceBadgeProps) {
  const config = {
    global: {
      label: "Global",
      icon: Globe,
      variant: "secondary" as const,
      title: "Shared across all services in this environment",
    },
    service: {
      label: "Service",
      icon: Package,
      variant: "default" as const,
      title: "Specific to this service only",
    },
    system: {
      label: "System",
      icon: Settings,
      variant: "outline" as const,
      title: "Automatically managed by Railway",
    },
  };

  const { label, icon: Icon, variant, title } = config[source];

  return (
    <Badge variant={variant} className={className} title={title}>
      <Icon className="h-3 w-3" />
      {label}
    </Badge>
  );
}

