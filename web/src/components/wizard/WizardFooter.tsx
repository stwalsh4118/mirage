"use client";

import { Button } from "@/components/ui/button";
import { useWizardStore, WIZARD_STEPS_ORDER } from "@/store/wizard";

export function WizardFooter(props: { onConfirm?: () => void; isSubmitting?: boolean; canProceed?: boolean }) {
  const { goBack, goNext, currentStepIndex, projectSelectionMode, existingProjectId, newProjectName, deploymentSource, repositoryUrl, repositoryBranch, imageName, imageTag, imageDigest, useDigest, environmentName } = useWizardStore();
  const step0Valid = projectSelectionMode === "existing" ? Boolean(existingProjectId) : newProjectName.trim().length > 0;
  
  // Step 1 validation: Either valid repo OR valid image
  const repoValid = repositoryUrl.trim().length === 0 || (repositoryUrl.trim().length > 0 && repositoryBranch.trim().length > 0);
  
  // Image validation: name is required, and either tag or valid digest
  const imageNameValid = imageName.trim().length > 0;
  const tagOrDigestValid = useDigest 
    ? imageDigest.trim().length === 0 || /^sha256:[a-f0-9]{64}$/.test(imageDigest.trim())
    : imageTag.trim().length > 0;
  const imageValid = imageNameValid && tagOrDigestValid;
  
  const step1Valid = deploymentSource === "repository" ? repoValid : imageValid;
  
  const step2Valid = environmentName.trim().length > 0;
  const last = currentStepIndex >= WIZARD_STEPS_ORDER.length - 1;
  return (
    <div className="mt-6 flex items-center justify-between border-t pt-4">
      <div className="space-x-2">
        <Button variant="ghost" onClick={goBack} disabled={currentStepIndex === 0}>
          Back
        </Button>
      </div>
      <div className="space-x-2">
        {!last && (
          <Button onClick={goNext} disabled={props.canProceed === false || (currentStepIndex === 0 && !step0Valid) || (currentStepIndex === 1 && !step1Valid) || (currentStepIndex === 2 && !step2Valid)}>
            Next
          </Button>
        )}
        {last && (
          <Button onClick={props.onConfirm} disabled={props.isSubmitting || props.canProceed === false}>
            {props.isSubmitting ? "Creating…" : "Confirm"}
          </Button>
        )}
      </div>
    </div>
  );
}


