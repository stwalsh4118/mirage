"use client";

import { useState, useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Search, CheckSquare, Square } from "lucide-react";
import { ServiceCard } from "./ServiceCard";
import type { DiscoveredService, SelectedService } from "./types";

interface DockerfileDiscoveryViewProps {
  services: DiscoveredService[];
  owner: string;
  repo: string;
  branch: string;
  onSelectionChange?: (selectedServices: SelectedService[]) => void;
}

export function DockerfileDiscoveryView({
  services,
  owner,
  repo,
  branch,
  onSelectionChange,
}: DockerfileDiscoveryViewProps) {
  // Initialize all services as selected with editable state
  const [serviceStates, setServiceStates] = useState<SelectedService[]>(() =>
    services.map((service) => ({
      ...service,
      selected: true,
      editedName: service.name,
    }))
  );

  const [searchQuery, setSearchQuery] = useState("");

  // Filter services based on search query
  const filteredServices = useMemo(() => {
    if (!searchQuery.trim()) return serviceStates;

    const query = searchQuery.toLowerCase();
    return serviceStates.filter(
      (service) =>
        service.name.toLowerCase().includes(query) ||
        service.dockerfilePath.toLowerCase().includes(query) ||
        service.buildContext.toLowerCase().includes(query) ||
        service.baseImage.toLowerCase().includes(query)
    );
  }, [serviceStates, searchQuery]);

  const selectedCount = serviceStates.filter((s) => s.selected).length;
  const allSelected = serviceStates.length > 0 && selectedCount === serviceStates.length;
  const noneSelected = selectedCount === 0;

  const handleToggleAll = () => {
    const newSelected = !allSelected;
    const updated = serviceStates.map((service) => ({
      ...service,
      selected: newSelected,
    }));
    setServiceStates(updated);
    onSelectionChange?.(updated.filter((s) => s.selected));
  };

  const handleToggleService = (index: number) => {
    const updated = [...serviceStates];
    updated[index] = {
      ...updated[index],
      selected: !updated[index].selected,
    };
    setServiceStates(updated);
    onSelectionChange?.(updated.filter((s) => s.selected));
  };

  const handleNameChange = (index: number, newName: string) => {
    const updated = [...serviceStates];
    updated[index] = {
      ...updated[index],
      editedName: newName,
    };
    setServiceStates(updated);
    onSelectionChange?.(updated.filter((s) => s.selected));
  };

  if (services.length === 0) {
    return (
      <Card className="glass grain">
        <CardContent className="flex flex-col items-center justify-center py-12 text-center">
          <div className="rounded-full bg-muted/50 p-4 mb-4">
            <Square className="size-8 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">No Dockerfiles Found</h3>
          <p className="text-sm text-muted-foreground max-w-md">
            No Dockerfiles were discovered in <span className="font-mono">{owner}/{repo}</span> on branch{" "}
            <span className="font-mono">{branch}</span>.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header with controls */}
      <Card className="glass grain">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg flex items-center gap-2">
                Discovered Services
                <Badge variant="secondary" className="bg-accent/15 text-accent">
                  {selectedCount} / {services.length} selected
                </Badge>
              </CardTitle>
              <p className="text-sm text-muted-foreground mt-1">
                Found in{" "}
                <span className="font-mono">
                  {owner}/{repo}
                </span>{" "}
                on <span className="font-mono">{branch}</span>
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleToggleAll}
                className="bg-transparent"
              >
                {allSelected ? (
                  <>
                    <CheckSquare className="size-4 mr-2" />
                    Deselect All
                  </>
                ) : (
                  <>
                    <Square className="size-4 mr-2" />
                    Select All
                  </>
                )}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              placeholder="Search services..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
        </CardContent>
      </Card>

      {/* Services Grid */}
      {filteredServices.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {filteredServices.map((service, index) => {
            // Find original index in serviceStates array
            const originalIndex = serviceStates.findIndex(
              (s) => s.dockerfilePath === service.dockerfilePath
            );
            return (
              <ServiceCard
                key={service.dockerfilePath}
                service={service}
                selected={service.selected}
                editedName={service.editedName}
                onToggleSelection={() => handleToggleService(originalIndex)}
                onNameChange={(name) => handleNameChange(originalIndex, name)}
              />
            );
          })}
        </div>
      ) : (
        <Card className="glass grain">
          <CardContent className="flex flex-col items-center justify-center py-8 text-center">
            <p className="text-sm text-muted-foreground">
              No services match your search query
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

