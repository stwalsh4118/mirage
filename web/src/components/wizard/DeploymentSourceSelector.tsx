"use client";

import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { GitBranch, Package } from "lucide-react";

export function DeploymentSourceSelector() {
  const { deploymentSource, setField } = useWizardStore();

  return (
    <div className="space-y-4">
      <div>
        <Label className="text-base">Deployment Source</Label>
        <p className="text-sm text-muted-foreground mt-1">
          Choose how you want to deploy your services
        </p>
      </div>

      <RadioGroup
        value={deploymentSource}
        onValueChange={(value) => setField("deploymentSource", value as "repository" | "image")}
        className="gap-4"
      >
        <div className="flex items-start space-x-3">
          <RadioGroupItem value="repository" id="source-repository" className="mt-1" />
          <div className="flex-1">
            <label
              htmlFor="source-repository"
              className="flex items-center gap-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
            >
              <GitBranch className="h-4 w-4" />
              From Source Repository
            </label>
            <p className="text-sm text-muted-foreground mt-1.5">
              Deploy services from a Git repository. Railway will build and deploy your code automatically.
            </p>
          </div>
        </div>

        <div className="flex items-start space-x-3">
          <RadioGroupItem value="image" id="source-image" className="mt-1" />
          <div className="flex-1">
            <label
              htmlFor="source-image"
              className="flex items-center gap-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
            >
              <Package className="h-4 w-4" />
              From Docker Image
            </label>
            <p className="text-sm text-muted-foreground mt-1.5">
              Deploy pre-built Docker images from registries like Docker Hub, GHCR, or private registries.
            </p>
          </div>
        </div>
      </RadioGroup>
    </div>
  );
}

