"use client";

import { useWizardStore } from "@/store/wizard";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { Label } from "@/components/ui/label";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { Card, CardContent } from "@/components/ui/card";

export function StepReview() {
  const {
    projectSelectionMode,
    existingProjectId,
    existingProjectName,
    newProjectName,
    repositoryUrl,
    repositoryBranch,
    environmentName,
    templateKind,
    ttlHours,
    environmentVariables,
    setField,
  } = useWizardStore();

  const summaryVars = environmentVariables.filter((v) => v.key.trim().length > 0);
  const { data: projects = [] } = useRailwayProjectsDetails();
  const resolvedProjectName = projectSelectionMode === "existing"
    ? (projects.find((p) => p.id === existingProjectId)?.name || existingProjectName || existingProjectId || "(none)")
    : (newProjectName || "(new project)");

  return (
    <div className="space-y-6">
      {/* Strategy toggle removed as not needed */}

      <Card className="glass grain border border-border/60">
        <CardContent className="pt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div className="space-y-2">
              <div>
                <div className="text-muted-foreground">Project</div>
                <div className="font-medium">{resolvedProjectName}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Source</div>
                <div className="font-medium">{repositoryUrl || "(none)"}{repositoryBranch ? `@${repositoryBranch}` : ""}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Environment</div>
                <div className="font-medium">{environmentName || "(none)"}</div>
              </div>
            </div>
            <div className="space-y-2">
              <div>
                <div className="text-muted-foreground">Template</div>
                <div className="font-medium">{templateKind}</div>
              </div>
              <div>
                <div className="text-muted-foreground">TTL</div>
                <div className="font-medium">{ttlHours ?? "(none)"}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Variables</div>
                <div className="font-medium">{summaryVars.length ? `${summaryVars.length} set` : "(none)"}</div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}


