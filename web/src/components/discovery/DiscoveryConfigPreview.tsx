"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Container, GitBranch, Folder, Network } from "lucide-react";
import type { SelectedService } from "./types";

interface DiscoveryConfigPreviewProps {
  selectedServices: SelectedService[];
  owner: string;
  repo: string;
  branch: string;
}

export function DiscoveryConfigPreview({
  selectedServices,
  owner,
  repo,
  branch,
}: DiscoveryConfigPreviewProps) {
  if (selectedServices.length === 0) {
    return (
      <Card className="glass grain">
        <CardContent className="flex flex-col items-center justify-center py-8 text-center">
          <p className="text-sm text-muted-foreground">
            No services selected. Select at least one service to see the preview.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="glass grain">
      <CardHeader>
        <CardTitle className="text-lg flex items-center gap-2">
          Deployment Preview
          <Badge variant="secondary" className="bg-accent/15 text-accent">
            {selectedServices.length} service{selectedServices.length !== 1 ? "s" : ""}
          </Badge>
        </CardTitle>
        <p className="text-sm text-muted-foreground">
          The following services will be created on Railway
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Repository Info */}
        <div className="space-y-2">
          <h4 className="text-sm font-semibold flex items-center gap-2">
            <GitBranch className="size-4" />
            Repository
          </h4>
          <div className="text-sm space-y-1 pl-6">
            <p className="font-mono text-foreground/80">
              {owner}/{repo}
            </p>
            <p className="text-muted-foreground">
              Branch: <span className="font-mono">{branch}</span>
            </p>
          </div>
        </div>

        <Separator />

        {/* Services List */}
        <div className="space-y-3">
          <h4 className="text-sm font-semibold flex items-center gap-2">
            <Container className="size-4" />
            Services
          </h4>
          <div className="space-y-3 pl-6">
            {selectedServices.map((service) => (
              <div
                key={service.dockerfilePath}
                className="space-y-2 pb-3 border-b border-border/50 last:border-0 last:pb-0"
              >
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-sm">
                      {service.editedName || service.name}
                    </p>
                    <p className="text-xs font-mono text-muted-foreground truncate">
                      {service.dockerfilePath}
                    </p>
                  </div>
                  {service.editedName && service.editedName !== service.name && (
                    <Badge variant="outline" className="text-xs shrink-0">
                      Renamed
                    </Badge>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-2 text-xs">
                  {/* Build Context */}
                  <div className="space-y-1">
                    <p className="text-muted-foreground flex items-center gap-1">
                      <Folder className="size-3" />
                      Build Context
                    </p>
                    <p className="font-mono text-foreground/80">
                      {service.buildContext}
                    </p>
                  </div>

                  {/* Ports */}
                  {service.exposedPorts.length > 0 && (
                    <div className="space-y-1">
                      <p className="text-muted-foreground flex items-center gap-1">
                        <Network className="size-3" />
                        Exposed Ports
                      </p>
                      <div className="flex flex-wrap gap-1">
                        {service.exposedPorts.map((port) => (
                          <Badge
                            key={port}
                            variant="outline"
                            className="text-xs font-mono"
                          >
                            {port}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {/* Base Image */}
                {service.baseImage && (
                  <div className="text-xs">
                    <p className="text-muted-foreground">Base Image</p>
                    <p className="font-mono text-foreground/80">
                      {service.baseImage}
                    </p>
                  </div>
                )}

                {/* Build Args */}
                {service.buildArgs.length > 0 && (
                  <div className="text-xs">
                    <p className="text-muted-foreground mb-1">Build Arguments</p>
                    <div className="flex flex-wrap gap-1">
                      {service.buildArgs.map((arg) => (
                        <Badge
                          key={arg}
                          variant="outline"
                          className="text-xs font-mono"
                        >
                          {arg}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Environment Variables Info */}
        <Separator />
        <div className="text-xs text-muted-foreground bg-muted/30 rounded-md p-3">
          <p className="font-medium mb-1">Note:</p>
          <ul className="list-disc list-inside space-y-1">
            <li>Each service will use its specified Dockerfile path via RAILWAY_DOCKERFILE_PATH</li>
            <li>Build context will be set relative to the repository root</li>
            <li>Build arguments can be configured in Railway after deployment</li>
          </ul>
        </div>
      </CardContent>
    </Card>
  );
}

