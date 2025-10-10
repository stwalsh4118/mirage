"use client";

import { useWizardStore } from "@/store/wizard";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Folder, GitBranch, Package, Server, Settings } from "lucide-react";

export function StepReview() {
  const {
    projectSelectionMode,
    existingProjectId,
    existingProjectName,
    newProjectName,
    defaultEnvironmentName,
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
    environmentVariables,
    serviceEnvironmentVariables,
    discoveredServices,
    selectedServiceIndices,
    serviceNameOverrides,
  } = useWizardStore();

  const summaryVars = environmentVariables.filter((v) => v.key.trim().length > 0);
  
  // Calculate service-specific variable counts
  const serviceVarCounts = selectedServiceIndices.map((index) => {
    const serviceVars = serviceEnvironmentVariables[index] || [];
    return serviceVars.filter((v) => v.key.trim().length > 0).length;
  });
  const totalServiceVars = serviceVarCounts.reduce((sum, count) => sum + count, 0);
  
  const { data: projects = [] } = useRailwayProjectsDetails();
  const resolvedProjectName = projectSelectionMode === "existing"
    ? (projects.find((p) => p.id === existingProjectId)?.name || existingProjectName || existingProjectId || "Unknown")
    : newProjectName;

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

  // Get selected services with their display names
  const selectedServices = selectedServiceIndices.map((index) => {
    const service = discoveredServices[index];
    const displayName = serviceNameOverrides[index] || service?.name || `Service ${index + 1}`;
    return {
      name: displayName,
      dockerfilePath: service?.dockerfilePath,
      exposedPorts: service?.exposedPorts || [],
    };
  });

  // Determine environment name to show
  const effectiveEnvironmentName = projectSelectionMode === "existing" 
    ? environmentName 
    : defaultEnvironmentName;

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h3 className="text-lg font-semibold">Review Configuration</h3>
        <p className="text-sm text-muted-foreground">
          Verify your settings before creating the environment
        </p>
      </div>

      {/* Project Section */}
      <Card className="glass grain border border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Folder className="h-4 w-4" />
            Project
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Type</span>
            <Badge variant="secondary" className="font-medium">
              {projectSelectionMode === "existing" ? "Existing Project" : "New Project"}
            </Badge>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Name</span>
            <span className="font-medium">{resolvedProjectName}</span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Environment</span>
            <span className="font-medium">{effectiveEnvironmentName || "production"}</span>
          </div>
        </CardContent>
      </Card>

      {/* Source Section */}
      <Card className="glass grain border border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            {deploymentSource === "repository" ? (
              <><GitBranch className="h-4 w-4" /> Repository</>
            ) : (
              <><Package className="h-4 w-4" /> Docker Image</>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          {deploymentSource === "repository" ? (
            <>
              <div className="flex items-start justify-between gap-4">
                <span className="text-muted-foreground">URL</span>
                <span className="font-medium font-mono text-xs break-all text-right">{repositoryUrl || "Not specified"}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">Branch</span>
                <span className="font-medium font-mono text-xs">{repositoryBranch || "main"}</span>
              </div>
            </>
          ) : (
            <>
              <div className="flex items-start justify-between gap-4">
                <span className="text-muted-foreground">Image</span>
                <span className="font-medium font-mono text-xs break-all text-right">{fullImageReference || "Not specified"}</span>
              </div>
              {imagePorts.length > 0 && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Ports</span>
                  <span className="font-medium">{imagePorts.join(", ")}</span>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* Services Section */}
      {selectedServices.length > 0 && (
        <Card className="glass grain border border-border/60">
          <CardHeader className="pb-3">
            <CardTitle className="text-base flex items-center gap-2">
              <Server className="h-4 w-4" />
              Services
              <Badge variant="secondary" className="ml-auto">{selectedServices.length}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {selectedServices.map((service, idx) => (
                <div key={idx} className="flex items-start justify-between gap-4 text-sm">
                  <span className="font-medium">{service.name}</span>
                  <div className="text-right space-y-0.5">
                    {service.dockerfilePath && (
                      <div className="text-xs text-muted-foreground font-mono">{service.dockerfilePath}</div>
                    )}
                    {service.exposedPorts.length > 0 && (
                      <div className="text-xs text-muted-foreground">Ports: {service.exposedPorts.join(", ")}</div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Environment Variables Section */}
      <Card className="glass grain border border-border/60">
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Settings className="h-4 w-4" />
            Environment Variables
          </CardTitle>
        </CardHeader>
        <CardContent>
          {summaryVars.length > 0 || totalServiceVars > 0 ? (
            <div className="space-y-2 text-sm">
              {summaryVars.length > 0 && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Global variables</span>
                  <Badge variant="secondary">{summaryVars.length}</Badge>
                </div>
              )}
              {totalServiceVars > 0 && (
                <div className="flex items-center justify-between">
                  <span className="text-muted-foreground">Service-specific variables</span>
                  <Badge variant="secondary">{totalServiceVars}</Badge>
                </div>
              )}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">No environment variables configured</p>
          )}
        </CardContent>
      </Card>

      <Separator />

      <div className="rounded-lg bg-muted/50 p-4 text-sm text-muted-foreground">
        <p>
          <strong>Note:</strong> This will create {projectSelectionMode === "new" ? "a new Railway project with a default environment" : "a new environment in the selected project"} and provision {selectedServices.length > 0 ? `${selectedServices.length} service${selectedServices.length > 1 ? 's' : ''}` : 'the configured services'} from your specified source.
        </p>
      </div>
    </div>
  );
}


