"use client";

import { useMemo } from "react";
import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Trash2, Plus, Info } from "lucide-react";
import { Switch } from "@/components/ui/switch";

const REGISTRY_PRESETS = {
  "docker.io": "Docker Hub",
  "ghcr.io": "GitHub Container Registry",
  custom: "Custom Registry",
} as const;

type RegistryPreset = keyof typeof REGISTRY_PRESETS;

const PORT_SUGGESTIONS: Record<string, number[]> = {
  nginx: [80, 443],
  postgres: [5432],
  redis: [6379],
  mysql: [3306],
  mongodb: [27017],
  rabbitmq: [5672, 15672],
};

export function DockerImageForm() {
  const { 
    imageName, 
    imageRegistry, 
    imageTag, 
    imageDigest,
    useDigest,
    imagePorts, 
    setField 
  } = useWizardStore();

  // Determine current registry preset
  const registryPreset = useMemo((): RegistryPreset => {
    if (!imageRegistry || imageRegistry === "docker.io") return "docker.io";
    if (imageRegistry === "ghcr.io") return "ghcr.io";
    return "custom";
  }, [imageRegistry]);

  // Generate full image reference preview
  const fullImageReference = useMemo(() => {
    const name = imageName.trim();
    if (!name) return "";

    const registry = imageRegistry.trim() || "docker.io";
    const registryPrefix = registry === "docker.io" ? "" : `${registry}/`;
    
    if (useDigest) {
      const digest = imageDigest.trim();
      return digest ? `${registryPrefix}${name}@${digest}` : `${registryPrefix}${name}`;
    } else {
      const tag = imageTag.trim() || "latest";
      return `${registryPrefix}${name}:${tag}`;
    }
  }, [imageName, imageRegistry, imageTag, imageDigest, useDigest]);

  // Suggest ports based on image name
  const suggestedPorts = useMemo(() => {
    const name = imageName.toLowerCase().split("/").pop() || "";
    for (const [key, ports] of Object.entries(PORT_SUGGESTIONS)) {
      if (name.includes(key)) return ports;
    }
    return [];
  }, [imageName]);

  const handleRegistryPresetChange = (preset: string) => {
    if (preset === "docker.io") {
      setField("imageRegistry", "");
    } else if (preset === "ghcr.io") {
      setField("imageRegistry", "ghcr.io");
    } else {
      // Custom - keep current value or clear if it was a preset
      if (imageRegistry === "" || imageRegistry === "ghcr.io") {
        setField("imageRegistry", "");
      }
    }
  };

  const addPort = (port?: number) => {
    const newPort = port || 80;
    if (!imagePorts.includes(newPort)) {
      setField("imagePorts", [...imagePorts, newPort]);
    }
  };

  const removePort = (index: number) => {
    setField("imagePorts", imagePorts.filter((_, i) => i !== index));
  };

  const updatePort = (index: number, value: string) => {
    const num = parseInt(value, 10);
    if (Number.isFinite(num) && num > 0 && num <= 65535) {
      const updated = [...imagePorts];
      updated[index] = num;
      setField("imagePorts", updated);
    }
  };

  const applySuggestedPorts = () => {
    setField("imagePorts", suggestedPorts);
  };

  // Validation
  const imageNameValid = imageName.trim().length > 0;
  const tagOrDigestValid = useDigest 
    ? imageDigest.trim().length === 0 || /^sha256:[a-f0-9]{64}$/.test(imageDigest.trim())
    : imageTag.trim().length > 0;

  return (
    <div className="space-y-6">
      {/* Registry Selection */}
      <div className="space-y-2">
        <Label htmlFor="registry">Registry</Label>
        <Select 
          value={registryPreset} 
          onValueChange={handleRegistryPresetChange}
        >
          <SelectTrigger id="registry" className="bg-card">
            <SelectValue placeholder="Select registry" />
          </SelectTrigger>
          <SelectContent className="bg-card">
            {Object.entries(REGISTRY_PRESETS).map(([value, label]) => (
              <SelectItem key={value} value={value}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {registryPreset === "custom" && (
          <Input
            placeholder="registry.example.com"
            className="bg-card mt-2"
            value={imageRegistry}
            onChange={(e) => setField("imageRegistry", e.target.value)}
            spellCheck={false}
          />
        )}
        <p className="text-xs text-muted-foreground">
          Choose the container registry hosting your image
        </p>
      </div>

      {/* Image Name */}
      <div className="space-y-2">
        <Label htmlFor="imageName">Image Name</Label>
        <Input
          id="imageName"
          className="bg-card"
          placeholder="nginx or owner/app"
          value={imageName}
          onChange={(e) => setField("imageName", e.target.value)}
          spellCheck={false}
        />
        <p className="text-xs text-muted-foreground">
          For Docker Hub: <code className="text-xs">nginx</code> or <code className="text-xs">owner/image</code>. 
          For other registries: <code className="text-xs">org/image</code>
        </p>
      </div>

      {/* Tag/Digest Toggle */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label>Version Specification</Label>
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">
              {useDigest ? "Using Digest" : "Using Tag"}
            </span>
            <Switch
              checked={useDigest}
              onCheckedChange={(checked) => setField("useDigest", checked)}
            />
          </div>
        </div>

        {useDigest ? (
          <div className="space-y-2">
            <Label htmlFor="imageDigest">Image Digest</Label>
            <Input
              id="imageDigest"
              className="bg-card font-mono text-xs"
              placeholder="sha256:abc123..."
              value={imageDigest}
              onChange={(e) => setField("imageDigest", e.target.value)}
              spellCheck={false}
            />
            <p className="text-xs text-muted-foreground">
              Use digest for immutable deployments (e.g., <code className="text-xs">sha256:abc123...</code>)
            </p>
            {imageDigest.trim() && !tagOrDigestValid && (
              <p className="text-xs text-destructive">
                Invalid digest format. Must be sha256:[64 hex chars]
              </p>
            )}
          </div>
        ) : (
          <div className="space-y-2">
            <Label htmlFor="imageTag">Image Tag</Label>
            <Input
              id="imageTag"
              className="bg-card"
              placeholder="latest"
              value={imageTag}
              onChange={(e) => setField("imageTag", e.target.value)}
              spellCheck={false}
            />
            <p className="text-xs text-muted-foreground">
              Common tags: <code className="text-xs">latest</code>, <code className="text-xs">stable</code>, <code className="text-xs">v1.0.0</code>
            </p>
          </div>
        )}
      </div>

      {/* Image Reference Preview */}
      {fullImageReference && (
        <Alert className="bg-muted/30 border-border/60">
          <Info className="h-4 w-4" />
          <AlertTitle>Full Image Reference</AlertTitle>
          <AlertDescription>
            <code className="text-sm font-mono break-all">{fullImageReference}</code>
          </AlertDescription>
        </Alert>
      )}

      {/* Port Configuration */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label>Exposed Ports</Label>
          <div className="flex gap-2">
            {suggestedPorts.length > 0 && imagePorts.length === 0 && (
              <Button 
                type="button" 
                size="sm" 
                variant="outline" 
                onClick={applySuggestedPorts}
              >
                Use Suggested ({suggestedPorts.join(", ")})
              </Button>
            )}
            <Button 
              type="button" 
              size="sm" 
              variant="outline" 
              onClick={() => addPort()}
            >
              <Plus className="h-4 w-4 mr-1" />
              Add Port
            </Button>
          </div>
        </div>

        {imagePorts.length > 0 ? (
          <div className="space-y-2">
            {imagePorts.map((port, index) => (
              <div key={index} className="flex gap-2 items-center">
                <Input
                  type="number"
                  min="1"
                  max="65535"
                  className="bg-card"
                  value={port}
                  onChange={(e) => updatePort(index, e.target.value)}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => removePort(index)}
                  aria-label="Remove port"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">
            No ports configured. Add ports that your container exposes.
          </p>
        )}
        <p className="text-xs text-muted-foreground">
          Specify which ports your container listens on (1-65535)
        </p>
      </div>

      {/* Validation Alerts */}
      {!imageNameValid && (
        <Alert className="bg-destructive/10 text-destructive-foreground border-destructive/40">
          <AlertTitle>Image Name Required</AlertTitle>
          <AlertDescription>
            Please provide a valid Docker image name to continue.
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}


