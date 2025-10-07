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
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";

export function StepProject() {
  const { 
    sourceMode, 
    cloneSourceEnvId, 
    projectSelectionMode, 
    existingProjectId, 
    newProjectName, 
    defaultEnvironmentName, 
    setField 
  } = useWizardStore();

  const { data: projects = [], isLoading } = useRailwayProjectsDetails();
  const selectedProject = useMemo(() => projects.find((p) => p.id === existingProjectId), [projects, existingProjectId]);
  const [open, setOpen] = useState(false);
  const [envComboboxOpen, setEnvComboboxOpen] = useState(false);

  // Flatten all environments from all projects for the clone dropdown
  const allEnvironments = useMemo(() => {
    return projects.flatMap((project) =>
      project.environments.map((env) => ({
        id: env.id,
        name: env.name,
        projectName: project.name,
        projectId: project.id,
        serviceCount: env.services.length,
      }))
    );
  }, [projects]);

  return (
    <div className="space-y-6">
      {/* Source Mode Selection */}
      <div className="space-y-3">
        <Label>How would you like to create this environment?</Label>
        <RadioGroup
          value={sourceMode}
          onValueChange={(value: 'new' | 'clone') => setField('sourceMode', value)}
        >
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="new" id="new" />
            <Label htmlFor="new" className="font-normal cursor-pointer">
              Create from scratch
            </Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="clone" id="clone" />
            <Label htmlFor="clone" className="font-normal cursor-pointer">
              Clone from existing environment
            </Label>
          </div>
        </RadioGroup>
      </div>

      {/* Environment Picker (shown only in clone mode) */}
      {sourceMode === 'clone' && (
        <div className="space-y-2">
          <Label>Source Environment</Label>
          <Popover open={envComboboxOpen} onOpenChange={setEnvComboboxOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={envComboboxOpen}
                className="w-full justify-between bg-card"
              >
                {cloneSourceEnvId ? (
                  (() => {
                    const selectedEnv = allEnvironments.find((env) => env.id === cloneSourceEnvId);
                    return selectedEnv ? (
                      <div className="flex items-center gap-2">
                        <span>{selectedEnv.name}</span>
                        <Badge variant="secondary" className="text-xs">
                          {selectedEnv.projectName}
                        </Badge>
                        <span className="text-xs text-muted-foreground">
                          ({selectedEnv.serviceCount} {selectedEnv.serviceCount === 1 ? 'service' : 'services'})
                        </span>
                      </div>
                    ) : "Select environment to clone";
                  })()
                ) : "Select environment to clone"}
                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[500px] p-0" align="start">
              <Command>
                <CommandInput placeholder="Search environments..." />
                <CommandList>
                  <CommandEmpty>
                    {isLoading ? "Loading environments..." : "No environments found."}
                  </CommandEmpty>
                  <CommandGroup>
                    {allEnvironments.map((env) => (
                      <CommandItem
                        key={env.id}
                        value={`${env.name} ${env.projectName}`}
                        onSelect={() => {
                          setField('cloneSourceEnvId', env.id);
                          setEnvComboboxOpen(false);
                        }}
                      >
                        <Check
                          className={cn(
                            "mr-2 h-4 w-4",
                            cloneSourceEnvId === env.id ? "opacity-100" : "opacity-0"
                          )}
                        />
                        <div className="flex items-center gap-2">
                          <span>{env.name}</span>
                          <Badge variant="secondary" className="text-xs">
                            {env.projectName}
                          </Badge>
                          <span className="text-xs text-muted-foreground">
                            ({env.serviceCount} {env.serviceCount === 1 ? 'service' : 'services'})
                          </span>
                        </div>
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          {cloneSourceEnvId && (
            <p className="text-sm text-muted-foreground">
              The wizard will be pre-filled with this environment&apos;s configuration
            </p>
          )}
        </div>
      )}

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


