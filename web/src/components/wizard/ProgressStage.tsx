"use client";

import { type StageStatus } from "@/store/wizard";
import { CheckCircle2, Circle, Loader2, AlertCircle } from "lucide-react";

interface ProgressStageProps {
  label: string;
  description: string;
  status: StageStatus;
  error?: string;
  duration?: string;
  isActive?: boolean;
}

export function ProgressStage({ label, description, status, error, duration, isActive }: ProgressStageProps) {
  const getIcon = () => {
    switch (status) {
      case "success":
        return <CheckCircle2 className="h-5 w-5 text-green-600" />;
      case "running":
        return <Loader2 className="h-5 w-5 animate-spin text-primary" />;
      case "error":
        return <AlertCircle className="h-5 w-5 text-destructive" />;
      default:
        return <Circle className="h-5 w-5 text-muted-foreground" />;
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case "success":
        return "text-green-600";
      case "running":
        return "text-primary";
      case "error":
        return "text-destructive";
      default:
        return "text-muted-foreground";
    }
  };

  return (
    <div
      className={`flex items-start gap-3 rounded-md p-3 transition-colors ${
        isActive ? "bg-secondary/20" : "bg-transparent"
      }`}
    >
      <div className="mt-0.5">{getIcon()}</div>
      <div className="flex-1 space-y-1">
        <div className="flex items-center justify-between">
          <h4 className={`text-sm font-medium ${getStatusColor()}`}>{label}</h4>
          {duration && <span className="text-xs text-muted-foreground">({duration}s)</span>}
        </div>
        <p className="text-xs text-muted-foreground">{description}</p>
        {error && (
          <p className="text-xs text-destructive mt-1 rounded bg-destructive/10 p-2 border border-destructive/20">
            {error}
          </p>
        )}
      </div>
    </div>
  );
}
