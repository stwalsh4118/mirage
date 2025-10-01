"use client";

import { useState, useEffect } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { WizardStepper } from "@/components/wizard/WizardStepper";
import { WizardFooter } from "@/components/wizard/WizardFooter";
import { ProvisionProgress } from "@/components/wizard/ProvisionProgress";
import { useWizardStore } from "@/store/wizard";
import { StepProject, StepSource, StepConfig, StepReview } from "./steps";
import { useProvisionEnvironment, useProvisionProject, useProvisionServices } from "@/hooks/useRailway";

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
        finalEnvironmentId = projectResult.baseEnvironmentId;
        
        setStageStatus("creating-project", "success");
      }

      // Stage 2: Create Environment (if existing project)
      if (projectSelectionMode === "existing") {
        setStageStatus("creating-environment", "running");
        
        const envName = state.environmentName || state.defaultEnvironmentName;
        const envResult = await provisionEnvironment.mutateAsync({
          requestId: requestId!,
          projectId: finalProjectId,
          name: envName,
        });

        setField("createdEnvironmentId", envResult.environmentId);
        finalEnvironmentId = envResult.environmentId;
        
        setStageStatus("creating-environment", "success");
      }

      // Stage 3: Create Services
      setStageStatus("creating-services", "running");
      
      // Initialize service progress - for now just one service, but built for multiple
      const servicesToCreate = [{ name: DEFAULT_SERVICE_NAME }];
      setServiceProgress(
        servicesToCreate.map((s) => ({ name: s.name, status: "pending" as const }))
      );

      // Create services sequentially with progress tracking
      const createdServiceIds: string[] = [];
      for (let i = 0; i < servicesToCreate.length; i++) {
        const service = servicesToCreate[i];
        updateServiceProgress(i, { status: "running" });

        try {
          // Build service config based on deployment source
          const serviceConfig: {
            name: string;
            repo?: string;
            branch?: string;
            imageName?: string;
            imageRegistry?: string;
            imageTag?: string;
          } = { name: service.name };
          
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
            
            serviceConfig.imageName = imageName;
            if (imageRegistry) serviceConfig.imageRegistry = imageRegistry;
            if (imageTag) serviceConfig.imageTag = imageTag;
          }

          const serviceResult = await provisionServices.mutateAsync({
            requestId: requestId!,
            projectId: finalProjectId,
            environmentId: finalEnvironmentId,
            services: [serviceConfig],
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
        className={`glass grain rounded-lg border border-border/60 shadow-sm ${
          currentStepIndex === 2 ? "sm:max-w-4xl" : "sm:max-w-lg"
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
            <ProvisionProgress />
          ) : (
            <>
              <WizardStepper />
              <div className="min-h-[200px]">
                {currentStepIndex === 0 && <StepProject />}
                {currentStepIndex === 1 && <StepSource />}
                {currentStepIndex === 2 && <StepConfig />}
                {currentStepIndex === 3 && <StepReview />}
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