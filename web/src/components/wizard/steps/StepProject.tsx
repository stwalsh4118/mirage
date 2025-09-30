"use client";

import { useMemo, useState } from "react";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";

export function StepProject() {
  const { projectSelectionMode, existingProjectId, newProjectName, defaultEnvironmentName, setField } = useWizardStore();

  const { data: projects = [], isLoading } = useRailwayProjectsDetails();
  const selectedProject = useMemo(() => projects.find((p) => p.id === existingProjectId), [projects, existingProjectId]);
  const [open, setOpen] = useState(false);

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <Label>Project</Label>
        <ToggleGroup
          type="single"
          value={projectSelectionMode}
          onValueChange={(v) => v && setField("projectSelectionMode", v as "existing" | "new")}
          variant="outline"
          className=""
        >
          <ToggleGroupItem value="existing" aria-label="Use existing project" className="hover:bg-primary/70 hover:text-foreground data-[state=on]:bg-primary/70 data-[state=on]:text-foreground">Existing</ToggleGroupItem>
          <ToggleGroupItem value="new" aria-label="Create new project" className="hover:bg-primary/70 hover:text-foreground data-[state=on]:bg-primary/70 data-[state=on]:text-foreground">New</ToggleGroupItem>
        </ToggleGroup>
      </div>

      {projectSelectionMode === "existing" ? (
        <div className="space-y-2">
          <Label htmlFor="project">Select a project</Label>
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={open}
                className="w-full justify-between bg-card hover:bg-muted/40 hover:text-foreground data-[state=open]:bg-muted/50"
              >
                {selectedProject ? (
                  <span className="flex items-center gap-2">
                    <span>{selectedProject.name}</span>
                    <Badge variant="secondary" className="font-mono text-xs">{selectedProject.id}</Badge>
                  </span>
                ) : (
                  <span className="text-muted-foreground">{isLoading ? "Loading projects…" : "Search projects…"}</span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="p-0 w-[420px] bg-card border border-border/60 shadow-md">
              <Command>
                <CommandInput placeholder="Search projects by name or id" />
                <CommandList>
                  <CommandEmpty>No projects found.</CommandEmpty>
                  <CommandGroup>
                    <ScrollArea className="h-64">
                      {projects.map((p) => (
                        <CommandItem
                          key={p.id}
                          value={`${p.name} ${p.id}`}
                          onSelect={() => {
                            setField("existingProjectId", p.id);
                            setOpen(false);
                            setField("existingProjectName", p.name);
                          }}
                        >
                          <div className="flex flex-col">
                            <span className="font-medium">{p.name}</span>
                            <span className="text-xs text-muted-foreground font-mono">{p.id}</span>
                          </div>
                        </CommandItem>
                      ))}
                    </ScrollArea>
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 items-start">
          <div className="space-y-2">
            <Label htmlFor="projectName">Project name</Label>
            <Input
              id="projectName"
              placeholder="e.g. Starlight Orchestra"
              value={newProjectName}
              onChange={(e) => setField("newProjectName", e.target.value)}
              className="bg-card"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="defaultEnv">Default environment name</Label>
            <Input
              id="defaultEnv"
              placeholder="production"
              value={defaultEnvironmentName}
              onChange={(e) => setField("defaultEnvironmentName", e.target.value)}
              className="bg-card"
            />
          </div>
        </div>
      )}
    </div>
  );
}


