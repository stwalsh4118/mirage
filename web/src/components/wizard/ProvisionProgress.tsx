"use client";

import { useRouter } from "next/navigation";
import { useWizardStore, type ProvisionStage, type StageInfo } from "@/store/wizard";
import { ProgressStage } from "./ProgressStage";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { CheckCircle2, AlertCircle } from "lucide-react";

interface StageConfig {
  id: ProvisionStage;
  label: string;
  description: string;
}

interface ProvisionProgressProps {
  onClose?: () => void;
}

export function ProvisionProgress({ onClose }: ProvisionProgressProps) {
  const router = useRouter();
  const {
    projectSelectionMode,
    currentStage,
    stages,
    serviceProgress,
    createdProjectId,
    existingProjectId,
  } = useWizardStore();

  // Define stages based on flow
  const stageConfigs: StageConfig[] = projectSelectionMode === "new"
    ? [
        { id: "creating-project", label: "Create Project", description: "Creating Railway project..." },
        { id: "creating-services", label: "Create Services", description: "Provisioning services..." },
        { id: "complete", label: "Complete", description: "Environment ready!" },
      ]
    : [
        { id: "creating-environment", label: "Create Environment", description: "Creating environment..." },
        { id: "creating-services", label: "Create Services", description: "Provisioning services..." },
        { id: "complete", label: "Complete", description: "Environment ready!" },
      ];

  const isComplete = currentStage === "complete";
  const isFailed = currentStage === "failed" || Object.values(stages).some((s) => s.status === "error");
  const hasError = Object.values(stages).some((s) => s.status === "error");

  return (
    <div className="space-y-6">
      <div className="space-y-3">
        {stageConfigs.map((config) => {
          const stageInfo: StageInfo = stages[config.id] || { status: "pending" };
          return (
            <ProgressStage
              key={config.id}
              label={config.label}
              description={config.description}
              status={stageInfo.status}
              error={stageInfo.error}
              duration={
                stageInfo.startedAt && stageInfo.completedAt
                  ? ((stageInfo.completedAt - stageInfo.startedAt) / 1000).toFixed(1)
                  : undefined
              }
              isActive={currentStage === config.id}
            />
          );
        })}

        {/* Show per-service progress when creating services */}
        {currentStage === "creating-services" && serviceProgress.length > 0 && (
          <div className="ml-8 space-y-2 border-l-2 border-border/40 pl-4">
            {serviceProgress.map((service, idx) => (
              <div key={idx} className="flex items-center gap-2 text-sm">
                {service.status === "success" && <CheckCircle2 className="h-4 w-4 text-green-600" />}
                {service.status === "running" && (
                  <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                )}
                {service.status === "error" && <AlertCircle className="h-4 w-4 text-destructive" />}
                {service.status === "pending" && <div className="h-4 w-4 rounded-full border-2 border-border" />}
                <span className={service.status === "error" ? "text-destructive" : "text-foreground/80"}>
                  {service.name}
                </span>
                {service.error && <span className="text-xs text-destructive">({service.error})</span>}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Success message */}
      {isComplete && !hasError && (
        <Alert className="glass grain border-border/60 bg-green-500/10">
          <CheckCircle2 className="h-4 w-4 text-green-600" />
          <AlertDescription className="text-foreground/90">
            Your environment is ready! You can now view it in the dashboard.
          </AlertDescription>
        </Alert>
      )}

      {/* Error message with retry placeholder */}
      {isFailed && (
        <Alert className="glass grain border-border/60 bg-destructive/10">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <AlertDescription className="text-foreground/90">
            Failed to create environment. Please try again or contact support.
          </AlertDescription>
        </Alert>
      )}

      {/* Actions */}
      <div className="flex items-center justify-end gap-2 border-t pt-4">
        {isComplete && (
          <>
            <Button
              variant="outline"
              onClick={() => {
                useWizardStore.getState().reset();
              }}
            >
              Create Another
            </Button>
            <Button
              onClick={() => {
                // Close dialog first
                onClose?.();
                
                // Reset wizard state
                useWizardStore.getState().reset();
                
                // Navigate to project
                const projectId = existingProjectId || createdProjectId;
                if (projectId) {
                  router.push(`/project/${projectId}`);
                } else {
                  // Fallback to dashboard if no project ID
                  router.push("/dashboard");
                }
              }}
            >
              View Project →
            </Button>
          </>
        )}
        {isFailed && (
          <>
            <Button
              variant="outline"
              onClick={() => {
                useWizardStore.getState().reset();
              }}
            >
              Start Over
            </Button>
            <Button
              variant="default"
              onClick={() => {
                // TODO: Implement retry logic
                console.log("Retry from failed stage");
              }}
            >
              Retry
            </Button>
            {/* Placeholder for cleanup - will be implemented later */}
            <Button
              variant="destructive"
              disabled
              title="Cleanup functionality coming soon"
            >
              Clean Up (Coming Soon)
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
