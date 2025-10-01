"use client";

import { useWizardStore } from "@/store/wizard";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { Card, CardContent } from "@/components/ui/card";

export function StepReview() {
  const {
    projectSelectionMode,
    existingProjectId,
    existingProjectName,
    newProjectName,
    deploymentSource,
    repositoryUrl,
    repositoryBranch,
    imageName,
    imageRegistry,
    imageTag,
    imageDigest,
    useDigest,
    imagePorts,
    environmentName,
    templateKind,
    ttlHours,
    environmentVariables,
  } = useWizardStore();

  const summaryVars = environmentVariables.filter((v) => v.key.trim().length > 0);
  const { data: projects = [] } = useRailwayProjectsDetails();
  const resolvedProjectName = projectSelectionMode === "existing"
    ? (projects.find((p) => p.id === existingProjectId)?.name || existingProjectName || existingProjectId || "(none)")
    : (newProjectName || "(new project)");

  // Build full image reference for display
  const fullImageReference = (() => {
    if (deploymentSource !== "image" || !imageName) return null;
    const registry = imageRegistry || "docker.io";
    const registryPrefix = registry === "docker.io" ? "" : `${registry}/`;
    if (useDigest && imageDigest) {
      return `${registryPrefix}${imageName}@${imageDigest}`;
    }
    const tag = imageTag || "latest";
    return `${registryPrefix}${imageName}:${tag}`;
  })();

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
                <div className="text-muted-foreground">
                  {deploymentSource === "repository" ? "Source Repository" : "Docker Image"}
                </div>
                <div className="font-medium break-all">
                  {deploymentSource === "repository" 
                    ? `${repositoryUrl || "(none)"}${repositoryBranch ? `@${repositoryBranch}` : ""}` 
                    : fullImageReference || "(none)"}
                </div>
              </div>
              {deploymentSource === "image" && imagePorts.length > 0 && (
                <div>
                  <div className="text-muted-foreground">Exposed Ports</div>
                  <div className="font-medium">{imagePorts.join(", ")}</div>
                </div>
              )}
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


