"use client";

import { useState, useEffect } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { WizardStepper } from "@/components/wizard/WizardStepper";
import { WizardFooter } from "@/components/wizard/WizardFooter";
import { ProvisionProgress } from "@/components/wizard/ProvisionProgress";
import { useWizardStore } from "@/store/wizard";
import { StepProject, StepSource, StepDiscovery, StepConfig, StepReview } from "./steps";
import { useProvisionEnvironment, useProvisionProject, useProvisionServices } from "@/hooks/useRailway";
import { generateUniqueServiceName } from "@/lib/serviceNaming";

const DEFAULT_SERVICE_NAME = "web-app";

export function CreateEnvironmentDialog(props: { trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  
  const {
    reset,
    currentStepIndex,
    projectSelectionMode,
    defaultEnvironmentName,
    environmentName,
    isProvisioning,
    currentStage,
    requestId,
    setField,
    setStageStatus,
    setServiceProgress,
    updateServiceProgress,
    startProvisioning,
  } = useWizardStore();

  const provisionProject = useProvisionProject();
  const provisionEnvironment = useProvisionEnvironment();
  const provisionServices = useProvisionServices();

  // Warn user if they try to close during provisioning
  useEffect(() => {
    if (isProvisioning) {
      const handleBeforeUnload = (e: BeforeUnloadEvent) => {
        e.preventDefault();
        e.returnValue = "";
      };
      window.addEventListener("beforeunload", handleBeforeUnload);
      return () => window.removeEventListener("beforeunload", handleBeforeUnload);
    }
  }, [isProvisioning]);

  const handleOpenChange = (next: boolean) => {
    if (isProvisioning && next === false) {
      // Show notification instead of closing
      // TODO: Add proper toast notification system
      console.log("Provisioning in progress - you can navigate away, but provisioning will continue");
      return;
    }
    setOpen(next);
    if (!next) {
      reset();
    }
  };

  const handleProvision = async () => {
    startProvisioning();
    const state = useWizardStore.getState();

    try {
      let finalProjectId = state.existingProjectId || "";
      let finalEnvironmentId = "";
      let finalRailwayEnvironmentId = ""; // Track Railway ID separately

      // Stage 1: Create Project (if new)
      if (projectSelectionMode === "new") {
        setStageStatus("creating-project", "running");
        
        const projectResult = await provisionProject.mutateAsync({
          requestId: requestId!,
          defaultEnvironmentName: state.defaultEnvironmentName,
          name: state.newProjectName || undefined,
        });

        setField("createdProjectId", projectResult.projectId);
        finalProjectId = projectResult.projectId;
        finalEnvironmentId = projectResult.baseEnvironmentId; // Mirage ID
        finalRailwayEnvironmentId = projectResult.railwayEnvironmentId; // Railway ID
        
        setStageStatus("creating-project", "success");
      }

      // Stage 2: Create Environment (if existing project)
      if (projectSelectionMode === "existing") {
        setStageStatus("creating-environment", "running");
        
        const envName = state.environmentName || state.defaultEnvironmentName;
        
        // Capture wizard inputs for metadata persistence
        const wizardInputs: Record<string, unknown> = {
          sourceType: state.deploymentSource,
          environmentType: state.templateKind, // dev or prod
        };
        
        // Add source-specific details
        if (state.deploymentSource === "repository") {
          wizardInputs.repositoryUrl = state.repositoryUrl;
          wizardInputs.branch = state.repositoryBranch;
          if (state.discoveredServices.length > 0 && state.selectedServiceIndices.length > 0) {
            wizardInputs.discoveredServices = state.selectedServiceIndices.map(idx => ({
              name: state.serviceNameOverrides[idx] || state.discoveredServices[idx].name,
              path: state.discoveredServices[idx].buildContext || "/",
              dockerfilePath: state.discoveredServices[idx].dockerfilePath,
            }));
          }
        } else if (state.deploymentSource === "image") {
          wizardInputs.dockerImage = state.imageName;
          wizardInputs.imageRegistry = state.imageRegistry;
          wizardInputs.imageTag = state.useDigest ? state.imageDigest : state.imageTag;
        }
        
        // Add TTL if set
        if (state.ttlHours) {
          wizardInputs.ttl = state.ttlHours;
        }
        
        const envResult = await provisionEnvironment.mutateAsync({
          requestId: requestId!,
          projectId: finalProjectId,
          name: envName,
          envType: state.templateKind as 'dev' | 'prod',
          wizardInputs,
        });

        setField("createdEnvironmentId", envResult.environmentId);
        finalEnvironmentId = envResult.environmentId; // Mirage ID
        finalRailwayEnvironmentId = envResult.railwayEnvironmentId; // Railway ID
        
        setStageStatus("creating-environment", "success");
      }

      // Stage 3: Create Services
      setStageStatus("creating-services", "running");
      
      // Build services list from discovered services or manual config
      const servicesToCreate: Array<{
        name: string;
        repo?: string;
        branch?: string;
        dockerfilePath?: string;
        imageName?: string;
        imageRegistry?: string;
        imageTag?: string;
        ports?: number[];
        environmentVariables?: Record<string, string>;
      }> = [];

      // Determine the environment name for unique service naming
      const envName = state.environmentName || state.defaultEnvironmentName;

      // Helper function to merge environment variables
      const mergeEnvironmentVariables = (serviceIndex?: number): Record<string, string> | undefined => {
        // Start with global variables
        const globalVars = Object.fromEntries(
          state.environmentVariables
            .filter((v) => v.key.trim().length > 0)
            .map((v) => [v.key, v.value])
        );

        // Merge with service-specific variables if provided
        let mergedVars = { ...globalVars };
        if (serviceIndex !== undefined && state.serviceEnvironmentVariables[serviceIndex]) {
          const serviceVars = Object.fromEntries(
            state.serviceEnvironmentVariables[serviceIndex]
              .filter((v) => v.key.trim().length > 0)
              .map((v) => [v.key, v.value])
          );
          mergedVars = { ...mergedVars, ...serviceVars };
        }

        // Return undefined if no variables, otherwise return the merged object
        return Object.keys(mergedVars).length > 0 ? mergedVars : undefined;
      };

      // Add discovered services
      if (state.selectedServiceIndices.length > 0) {
        state.selectedServiceIndices.forEach((index) => {
          const discoveredService = state.discoveredServices[index];
          if (discoveredService) {
            const baseName = state.serviceNameOverrides[index] || discoveredService.name;
            // Generate unique name by appending environment (e.g., "api-prod", "web-stg")
            const uniqueName = generateUniqueServiceName(baseName, envName);
            servicesToCreate.push({
              name: uniqueName,
              repo: state.repositoryUrl?.trim() || undefined,
              branch: state.repositoryBranch?.trim() || undefined,
              dockerfilePath: discoveredService.dockerfilePath,
              environmentVariables: mergeEnvironmentVariables(index),
            });
          }
        });
      } else {
        // Fallback: create single service with manual config
        // Generate unique name for the default service too
        const uniqueName = generateUniqueServiceName(DEFAULT_SERVICE_NAME, envName);
        
        const serviceConfig: {
          name: string;
          repo?: string;
          branch?: string;
          imageName?: string;
          imageRegistry?: string;
          imageTag?: string;
          ports?: number[];
          environmentVariables?: Record<string, string>;
        } = {
          name: uniqueName,
          environmentVariables: mergeEnvironmentVariables(), // Use global variables only
        };
        
        if (state.deploymentSource === "repository") {
          const repo = state.repositoryUrl?.trim() || undefined;
          const branch = state.repositoryBranch?.trim() || undefined;
          serviceConfig.repo = repo;
          serviceConfig.branch = branch;
        } else {
          // Image-based deployment
          const imageName = state.imageName?.trim() || undefined;
          const imageRegistry = state.imageRegistry?.trim() || undefined;
          const imageTag = state.imageTag?.trim() || undefined;
          const ports = state.imagePorts && state.imagePorts.length > 0 ? state.imagePorts : undefined;
          
          serviceConfig.imageName = imageName;
          if (imageRegistry) serviceConfig.imageRegistry = imageRegistry;
          if (imageTag) serviceConfig.imageTag = imageTag;
          if (ports) serviceConfig.ports = ports;
        }
        
        servicesToCreate.push(serviceConfig);
      }

      // Initialize service progress
      setServiceProgress(
        servicesToCreate.map((s) => ({ name: s.name, status: "pending" as const }))
      );

      // Create services sequentially with progress tracking
      const createdServiceIds: string[] = [];
      for (let i = 0; i < servicesToCreate.length; i++) {
        const service = servicesToCreate[i];
        updateServiceProgress(i, { status: "running" });

        try {
          const serviceResult = await provisionServices.mutateAsync({
            requestId: requestId!,
            projectId: finalProjectId,
            environmentId: finalEnvironmentId,                   // Mirage ID for database FK
            railwayEnvironmentId: finalRailwayEnvironmentId,    // Railway ID for Railway API
            services: [service],
          });

          createdServiceIds.push(...serviceResult.serviceIds);
          updateServiceProgress(i, { 
            status: "success",
            serviceId: serviceResult.serviceIds[0],
          });
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : "Unknown error";
          updateServiceProgress(i, { 
            status: "error",
            error: errorMessage,
          });
          throw error; // Propagate to outer catch
        }
      }

      setField("createdServiceIds", createdServiceIds);
      setStageStatus("creating-services", "success");

      // Stage 4: Apply Config (placeholder for future env vars implementation)
      // TODO: Implement env vars application when backend supports it
      // if (state.environmentVariables.length > 0) {
      //   setStageStatus("applying-config", "running");
      //   await applyEnvVars({ ... });
      //   setStageStatus("applying-config", "success");
      // }

      // Complete!
      setStageStatus("complete", "success");
      
      console.log(`Environment created: "${environmentName || defaultEnvironmentName}"`);

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Unknown error occurred";
      setStageStatus("failed", "error", errorMessage);
      
      console.error("Provisioning failed:", errorMessage);
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {props.trigger ?? <Button variant="default">Create environment</Button>}
      </DialogTrigger>
      <DialogContent 
        className={`glass grain rounded-lg border border-border/60 shadow-sm max-h-[90vh] flex flex-col ${
          currentStepIndex === 3 ? "sm:max-w-6xl" : currentStepIndex === 2 ? "sm:max-w-4xl" : "sm:max-w-lg"
        }`}
      >
        <DialogHeader className="pb-0">
          <DialogTitle className="text-foreground/90">
            {isProvisioning || currentStage === "complete" || currentStage === "failed" 
              ? "Creating Environment" 
              : "Create Environment"}
          </DialogTitle>
        </DialogHeader>
        
        <div className="space-y-4 pt-2">
          {/* Show progress screen during provisioning AND after completion */}
          {isProvisioning || currentStage === "complete" || currentStage === "failed" ? (
            <ProvisionProgress onClose={() => setOpen(false)} />
          ) : (
            <>
              <WizardStepper />
              <div className="min-h-[200px] overflow-y-auto flex-1">
                {currentStepIndex === 0 && <StepProject />}
                {currentStepIndex === 1 && <StepSource />}
                {currentStepIndex === 2 && <StepDiscovery />}
                {currentStepIndex === 3 && <StepConfig />}
                {currentStepIndex === 4 && <StepReview />}
              </div>
              <WizardFooter
                isSubmitting={false}
                onConfirm={handleProvision}
              />
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}