"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { WizardStepper } from "@/components/wizard/WizardStepper";
import { WizardFooter } from "@/components/wizard/WizardFooter";
import { useWizardStore } from "@/store/wizard";
import { StepProject, StepSource, StepConfig, StepReview } from "./steps";

export function CreateEnvironmentDialog(props: { trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  const reset = useWizardStore((s) => s.reset);
  const { currentStepIndex } = useWizardStore();

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
            }}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}


