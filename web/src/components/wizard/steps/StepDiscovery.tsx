"use client";

import { useMemo } from "react";
import { useWizardStore } from "@/store/wizard";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Loader2, Search, AlertCircle, Package } from "lucide-react";
import { useDockerfileDiscovery } from "@/hooks/useDockerfileDiscovery";
import { DockerfileDiscoveryView } from "@/components/discovery";
import type { SelectedService } from "@/components/discovery/types";

export function StepDiscovery() {
  const {
    repositoryUrl,
    repositoryBranch,
    githubToken,
    deploymentSource,
    discoveryTriggered,
    discoveryLoading,
    discoveryError,
    discoveredServices,
    selectedServiceIndices,
    serviceNameOverrides,
    discoverySkipped,
    setField,
  } = useWizardStore();

  const discovery = useDockerfileDiscovery();

  // Parse GitHub owner and repo from URL
  const { owner, repo } = useMemo(() => {
    if (!repositoryUrl) return { owner: null, repo: null };
    
    const trimmed = repositoryUrl.trim();
    
    // Handle various GitHub URL formats
    // 1. Simple format: "owner/repo" (Railway format)
    // 2. Full URL: "github.com/owner/repo" or "https://github.com/owner/repo"
    // 3. SSH format: "git@github.com:owner/repo.git"
    
    // Try simple "owner/repo" format first (Railway's expected format)
    const simpleMatch = trimmed.match(/^([\w-]+)\/([\w.-]+?)(?:\.git)?$/i);
    if (simpleMatch) {
      return { owner: simpleMatch[1], repo: simpleMatch[2] };
    }
    
    // Try full URL format
    const urlMatch = trimmed.match(/github\.com[/:]([\w-]+)\/([\w.-]+)/i);
    if (urlMatch) {
      return { owner: urlMatch[1], repo: urlMatch[2].replace(/\.git$/, "") };
    }
    
    return { owner: null, repo: null };
  }, [repositoryUrl]);

  const canScan = owner && repo && repositoryBranch && !discoveryLoading;

  const handleScan = () => {
    if (!owner || !repo || !repositoryBranch) return;

    setField("discoveryTriggered", true);
    setField("discoveryLoading", true);
    setField("discoveryError", null);
    setField("discoverySkipped", false);

    discovery.mutate(
      {
        owner,
        repo,
        branch: repositoryBranch,
        userToken: githubToken || undefined,
      },
      {
        onSuccess: (data) => {
          setField("discoveryLoading", false);
          setField("discoveredServices", data.services);
          
          // Select all services by default
          setField(
            "selectedServiceIndices",
            data.services.map((_, i) => i)
          );
          
          // Initialize name overrides with original names
          const overrides: Record<number, string> = {};
          data.services.forEach((service, i) => {
            overrides[i] = service.name;
          });
          setField("serviceNameOverrides", overrides);
        },
        onError: (error) => {
          setField("discoveryLoading", false);
          setField("discoveryError", error.message);
        },
      }
    );
  };

  const handleSkip = () => {
    setField("discoverySkipped", true);
    setField("discoveredServices", []);
    setField("selectedServiceIndices", []);
    setField("serviceNameOverrides", {});
  };

  const handleSelectionChange = (selectedServices: SelectedService[]) => {
    // Find indices of selected services
    const indices = selectedServices
      .map((selected) => {
        return discoveredServices.findIndex(
          (ds) => ds.dockerfilePath === selected.dockerfilePath
        );
      })
      .filter((i) => i !== -1);

    setField("selectedServiceIndices", indices);

    // Update name overrides
    const overrides: Record<number, string> = {};
    selectedServices.forEach((selected) => {
      const index = discoveredServices.findIndex(
        (ds) => ds.dockerfilePath === selected.dockerfilePath
      );
      if (index !== -1 && selected.editedName) {
        overrides[index] = selected.editedName;
      }
    });
    setField("serviceNameOverrides", overrides);
  };

  // If deployment source is image, skip discovery automatically
  if (deploymentSource === "image") {
    return (
      <div className="space-y-4">
        <Alert className="bg-muted/30 border-border/60">
          <Package className="h-4 w-4" />
          <AlertTitle>Image-based Deployment</AlertTitle>
          <AlertDescription>
            Dockerfile discovery is not available for Docker image deployments. Continue to the next step to configure your deployment.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  // If no repository provided, show message
  if (!repositoryUrl || !owner || !repo) {
    return (
      <div className="space-y-4">
        <Alert className="bg-muted/30 border-border/60">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Repository Required</AlertTitle>
          <AlertDescription>
            Please provide a repository URL in the previous step to enable Dockerfile discovery.
          </AlertDescription>
        </Alert>
        <Button variant="outline" onClick={handleSkip}>
          Skip Discovery
        </Button>
      </div>
    );
  }

  // Initial state: Show scan button
  if (!discoveryTriggered && !discoverySkipped) {
    return (
      <div className="space-y-4">
        <Alert className="bg-muted/30 border-border/60">
          <Search className="h-4 w-4" />
          <AlertTitle>Scan for Dockerfiles</AlertTitle>
          <AlertDescription>
            Automatically discover services in your monorepo by scanning for Dockerfiles.
            This will help you quickly deploy multiple services from a single repository.
          </AlertDescription>
        </Alert>

        <div className="flex gap-3">
          <Button onClick={handleScan} disabled={!canScan}>
            <Search className="mr-2 h-4 w-4" />
            Scan Repository
          </Button>
          <Button variant="outline" onClick={handleSkip}>
            Skip Discovery
          </Button>
        </div>

        {!canScan && (
          <p className="text-sm text-muted-foreground">
            Ensure repository URL and branch are provided in the previous step.
          </p>
        )}
      </div>
    );
  }

  // Loading state
  if (discoveryLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-center py-12">
          <div className="text-center space-y-4">
            <Loader2 className="h-8 w-8 animate-spin mx-auto text-primary" />
            <div>
              <p className="text-sm font-medium">Scanning repository...</p>
              <p className="text-xs text-muted-foreground mt-1">
                Looking for Dockerfiles in {owner}/{repo} on {repositoryBranch}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (discoveryError) {
    return (
      <div className="space-y-4">
        <Alert className="bg-destructive/10 text-destructive-foreground border-destructive/40">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Discovery Failed</AlertTitle>
          <AlertDescription>{discoveryError}</AlertDescription>
        </Alert>

        <div className="flex gap-3">
          <Button onClick={handleScan}>
            <Search className="mr-2 h-4 w-4" />
            Retry Scan
          </Button>
          <Button variant="outline" onClick={handleSkip}>
            Skip Discovery
          </Button>
        </div>
      </div>
    );
  }

  // Skipped state
  if (discoverySkipped) {
    return (
      <div className="space-y-4">
        <Alert className="bg-muted/30 border-border/60">
          <AlertTitle>Discovery Skipped</AlertTitle>
          <AlertDescription>
            You can still manually configure services in the next step.
          </AlertDescription>
        </Alert>

        <Button variant="outline" onClick={() => setField("discoverySkipped", false)}>
          <Search className="mr-2 h-4 w-4" />
          Run Discovery
        </Button>
      </div>
    );
  }

  // Results state
  if (discoveredServices.length > 0) {
    // Convert to SelectedService format for the discovery view
    const servicesWithSelection = discoveredServices.map((service, index) => ({
      ...service,
      selected: selectedServiceIndices.includes(index),
      editedName: serviceNameOverrides[index] || service.name,
    }));

    return (
      <div className="space-y-4">
        <DockerfileDiscoveryView
          services={discoveredServices}
          owner={owner}
          repo={repo}
          branch={repositoryBranch}
          onSelectionChange={handleSelectionChange}
        />

        <div className="flex justify-between items-center pt-2">
          <Button variant="outline" size="sm" onClick={handleScan}>
            <Search className="mr-2 h-4 w-4" />
            Rescan
          </Button>
          <p className="text-sm text-muted-foreground">
            {selectedServiceIndices.length} service
            {selectedServiceIndices.length !== 1 ? "s" : ""} selected
          </p>
        </div>
      </div>
    );
  }

  // No services found
  return (
    <div className="space-y-4">
      <Alert className="bg-muted/30 border-border/60">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>No Dockerfiles Found</AlertTitle>
        <AlertDescription>
          No Dockerfiles were discovered in {owner}/{repo} on {repositoryBranch}.
          You can rescan or continue without discovered services.
        </AlertDescription>
      </Alert>

      <div className="flex gap-3">
        <Button variant="outline" onClick={handleScan}>
          <Search className="mr-2 h-4 w-4" />
          Rescan
        </Button>
        <Button variant="outline" onClick={handleSkip}>
          Continue Anyway
        </Button>
      </div>
    </div>
  );
}

