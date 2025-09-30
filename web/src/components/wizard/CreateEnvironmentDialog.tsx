"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { WizardStepper } from "@/components/wizard/WizardStepper";
import { WizardFooter } from "@/components/wizard/WizardFooter";
import { useWizardStore } from "@/store/wizard";
import { StepProject, StepSource, StepConfig, StepReview } from "./steps";
import { useProvisionEnvironment, useProvisionProject } from "@/hooks/useRailway";

export function CreateEnvironmentDialog(props: { trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  const reset = useWizardStore((s) => s.reset);
  const { currentStepIndex } = useWizardStore();
  const { projectSelectionMode, defaultEnvironmentName, newProjectName, setField } = useWizardStore();
  const provisionProject = useProvisionProject();
  const provisionEnvironment = useProvisionEnvironment();

  const handleOpenChange = (next: boolean) => {
    setOpen(next);
    if (!next) {
      reset();
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {props.trigger ?? <Button variant="default">Create environment</Button>}
      </DialogTrigger>
      <DialogContent className={`glass grain rounded-lg border border-border/60 shadow-sm ${currentStepIndex === 2 ? "sm:max-w-4xl" : "sm:max-w-lg"}`}>
        <DialogHeader className="pb-0">
          <DialogTitle className="text-foreground/90">Create environment</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-2">
          <WizardStepper />
          <div className="min-h-[200px]">
            {currentStepIndex === 0 && <StepProject />}
            {currentStepIndex === 1 && <StepSource />}
            {currentStepIndex === 2 && <StepConfig />}
            {currentStepIndex === 3 && <StepReview />}
          </div>
          <WizardFooter
            isSubmitting={provisionProject.isPending || provisionEnvironment.isPending}
            onConfirm={() => {
              const state = useWizardStore.getState();
              // Intentionally logging full wizard state for inspection
              // eslint-disable-next-line no-console
              console.log("CreateEnvironment submission", {
                project: {
                  mode: state.projectSelectionMode,
                  existingProjectId: state.existingProjectId,
                  existingProjectName: state.existingProjectName,
                  newProjectName: state.newProjectName,
                  defaultEnvironmentName: state.defaultEnvironmentName,
                },
                source: {
                  repositoryUrl: state.repositoryUrl,
                  repositoryBranch: state.repositoryBranch,
                },
                config: {
                  environmentName: state.environmentName,
                  templateKind: state.templateKind,
                  ttlHours: state.ttlHours,
                  environmentVariables: state.environmentVariables,
                },
              });
              // Kick off project creation if needed
              if (projectSelectionMode === "new") {
                const requestId = crypto.randomUUID();
                provisionProject.mutate(
                  { requestId, defaultEnvironmentName, name: newProjectName || undefined },
                  {
                    onSuccess: (res) => {
                      setField("createdProjectId", res.projectId);
                    },
                  }
                );
              } else if (projectSelectionMode === "existing") {
                const requestId = crypto.randomUUID();
                const projId = state.existingProjectId;
                const envName = state.environmentName || state.defaultEnvironmentName;
                if (projId && envName) {
                  provisionEnvironment.mutate({ requestId, projectId: projId, name: envName });
                }
              }
            }}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}


