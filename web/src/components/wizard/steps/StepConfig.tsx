"use client";

import { useState, useEffect } from "react";
import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Trash2, Upload } from "lucide-react";
import { EnvImportDialog } from "../EnvImportDialog";
import { VariableSourceBadge } from "../VariableSourceBadge";

export function StepConfig() {
  const {
    environmentName,
    templateKind,
    ttlHours,
    environmentVariables,
    serviceEnvironmentVariables,
    discoveredServices,
    selectedServiceIndices,
    serviceNameOverrides,
    sourceMode,
    setField,
  } = useWizardStore();

  const [selectedServiceIndex, setSelectedServiceIndex] = useState<number | null>(null);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [serviceImportDialogOpen, setServiceImportDialogOpen] = useState(false);

  // Auto-select first service when in clone mode with pre-populated service variables
  useEffect(() => {
    if (sourceMode === 'clone' && 
        selectedServiceIndices.length > 0 && 
        Object.keys(serviceEnvironmentVariables).length > 0 &&
        selectedServiceIndex === null) {
      // Select the first service that has variables
      const firstServiceWithVars = selectedServiceIndices.find(idx => 
        serviceEnvironmentVariables[idx] && serviceEnvironmentVariables[idx].length > 0
      );
      if (firstServiceWithVars !== undefined) {
        setSelectedServiceIndex(firstServiceWithVars);
      } else {
        // If no service has vars yet, just select the first service
        setSelectedServiceIndex(selectedServiceIndices[0]);
      }
    }
  }, [sourceMode, selectedServiceIndices, serviceEnvironmentVariables, selectedServiceIndex]);

  // Ensure there's always an empty row at the end for easy adding (global vars)
  useEffect(() => {
    const lastVar = environmentVariables[environmentVariables.length - 1];
    const hasEmptyRow = lastVar && lastVar.key === "" && lastVar.value === "";
    if (!hasEmptyRow) {
      setField("environmentVariables", [...environmentVariables, { key: "", value: "" }]);
    }
  }, [environmentVariables, setField]);
  // Ensure there's always an empty row for selected service vars
  useEffect(() => {
    if (selectedServiceIndex === null) return;
    
    const serviceVars = serviceEnvironmentVariables[selectedServiceIndex] || [];
    const lastVar = serviceVars[serviceVars.length - 1];
    const hasEmptyRow = lastVar && lastVar.key === "" && lastVar.value === "";
    
    if (!hasEmptyRow) {
      const updated = { ...serviceEnvironmentVariables };
      updated[selectedServiceIndex] = [...serviceVars, { key: "", value: "" }];
      setField("serviceEnvironmentVariables", updated);
    }
  }, [selectedServiceIndex, serviceEnvironmentVariables, setField]);

  const updateVar = (index: number, key: "key" | "value", value: string) => {
    const next = environmentVariables.slice();
    next[index] = { ...next[index], [key]: value };
    
    // Auto-add new empty row when user types in the last row
    if (index === environmentVariables.length - 1 && value.trim()) {
      next.push({ key: "", value: "" });
    }
    
    setField("environmentVariables", next);
  };

  const removeVar = (index: number) => {
    // Don't allow removing the last empty row
    if (index === environmentVariables.length - 1 && 
        environmentVariables[index].key === "" && 
        environmentVariables[index].value === "") {
      return;
    }
    setField("environmentVariables", environmentVariables.filter((_, i) => i !== index));
  };

  const handleImport = (imported: Array<{ key: string; value: string }>) => {
    // Remove empty rows and add imported variables
    const existingNonEmpty = environmentVariables.filter((v) => v.key.trim() || v.value.trim());
    setField("environmentVariables", [...existingNonEmpty, ...imported, { key: "", value: "" }]);
  };

  const handleServiceImport = (imported: Array<{ key: string; value: string }>) => {
    if (selectedServiceIndex === null) return;
    
    const serviceVars = serviceEnvironmentVariables[selectedServiceIndex] || [];
    const existingNonEmpty = serviceVars.filter((v) => v.key.trim() || v.value.trim());
    
    const updated = { ...serviceEnvironmentVariables };
    updated[selectedServiceIndex] = [...existingNonEmpty, ...imported, { key: "", value: "" }];
    setField("serviceEnvironmentVariables", updated);
  };

  // Parse pasted content to detect .env format
  const parseEnvPaste = (text: string): Array<{ key: string; value: string }> | null => {
    const lines = text.split(/\r?\n/).filter((l) => l.trim() && !l.trim().startsWith("#"));
    
    // Only parse if we have multiple lines or at least one line with =
    if (lines.length === 0) return null;
    
    // Check if it looks like env format (at least one line with =)
    const hasEnvFormat = lines.some((line) => line.includes("="));
    if (!hasEnvFormat) return null;
    
    const parsed = lines
      .map((line) => {
        const eqIndex = line.indexOf("=");
        if (eqIndex === -1) return null;
        
        const key = line.slice(0, eqIndex).trim();
        let value = line.slice(eqIndex + 1).trim();
        
        // Remove surrounding quotes
        if ((value.startsWith('"') && value.endsWith('"')) || 
            (value.startsWith("'") && value.endsWith("'"))) {
          value = value.slice(1, -1);
        }
        
        return key ? { key, value } : null;
      })
      .filter((v): v is { key: string; value: string } => v !== null);
    
    return parsed.length > 0 ? parsed : null;
  };

  // Handle paste for global variables
  const handleGlobalPaste = (e: React.ClipboardEvent) => {
    const text = e.clipboardData.getData("text");
    const parsed = parseEnvPaste(text);
    
    if (parsed && parsed.length > 0) {
      e.preventDefault();
      handleImport(parsed);
    }
  };

  // Handle paste for service variables
  const handleServicePaste = (e: React.ClipboardEvent) => {
    const text = e.clipboardData.getData("text");
    const parsed = parseEnvPaste(text);
    
    if (parsed && parsed.length > 0) {
      e.preventDefault();
      handleServiceImport(parsed);
    }
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="space-y-2 md:col-span-2">
          <Label htmlFor="envName">Environment name</Label>
          <Input
            id="envName"
            className="bg-card"
            placeholder="staging"
            value={environmentName}
            onChange={(e) => setField("environmentName", e.target.value)}
          />
        </div>
        <div className="space-y-2 md:col-span-2">
          <Label htmlFor="template">Template</Label>
          <Select value={templateKind} onValueChange={(v) => setField("templateKind", v as "dev" | "prod")}> 
            <SelectTrigger id="template" className="bg-card w-full pr-6 gap-1">
              <SelectValue placeholder="Choose template" />
            </SelectTrigger>
            <SelectContent className="bg-card">
              <SelectItem value="dev">dev — minimal resources, fast iteration</SelectItem>
              <SelectItem value="prod">prod — higher resources, durability</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="space-y-2">
          <Label htmlFor="ttl">TTL (hours, optional)</Label>
          <Input
            id="ttl"
            className="bg-card"
            inputMode="numeric"
            placeholder="24"
            value={ttlHours ?? ""}
            onChange={(e) => {
              const v = e.target.value;
              if (v === "") setField("ttlHours", null);
              else {
                const num = Number(v.replace(/[^0-9]/g, ""));
                setField("ttlHours", Number.isFinite(num) ? num : null);
              }
            }}
          />
          <p className="text-xs text-muted-foreground">Leave blank for no auto-destroy; bounds validated on submit.</p>
        </div>
      </div>

      <Separator />

      <div className="space-y-3">
        <Label className="text-base">Environment Variables</Label>
        {selectedServiceIndices.length > 0 ? (
          <p className="text-xs text-muted-foreground">
            Configure global variables (shared across all services) and service-specific overrides
          </p>
        ) : (
          <p className="text-xs text-muted-foreground">
            Configure variables that will be shared across all services
          </p>
        )}

        {selectedServiceIndices.length > 0 ? (
          <div className="space-y-6">
            {/* Two-column layout for variable editors */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Left: Global Variables */}
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <Label>Global Variables</Label>
                  <Button size="sm" variant="outline" onClick={() => setImportDialogOpen(true)} className="h-8">
                    <Upload className="h-3.5 w-3.5 mr-1.5" />
                    Import
                  </Button>
                </div>
                <div className="rounded-md border border-border/60 p-3 bg-card/50 space-y-2.5 max-h-[400px] overflow-y-auto">
                  {environmentVariables.map((kv, i) => {
                    const isLastEmpty = i === environmentVariables.length - 1 && kv.key === "" && kv.value === "";
                    return (
                      <div key={i} className="grid grid-cols-12 gap-2 items-center">
                      <Input
                        placeholder={isLastEmpty ? "Type KEY..." : "KEY"}
                        className="bg-card col-span-5 h-9 text-sm"
                        value={kv.key}
                        onChange={(e) => updateVar(i, "key", e.target.value)}
                        onPaste={handleGlobalPaste}
                      />
                      <Input
                        placeholder="value"
                        className="bg-card col-span-6 h-9 text-sm"
                        value={kv.value}
                        onChange={(e) => updateVar(i, "value", e.target.value)}
                        onPaste={handleGlobalPaste}
                      />
                        <Button
                          variant="ghost"
                          size="icon"
                          className="col-span-1 h-9 w-9"
                          onClick={() => removeVar(i)}
                          disabled={isLastEmpty}
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                    );
                  })}
                </div>
              </div>

              {/* Right: Service-Specific Variables */}
              <div className="space-y-3">
                <div className="flex items-center justify-between gap-2">
                  <Label>Service Variables</Label>
                  <div className="flex items-center gap-2">
                    <Button 
                      size="sm" 
                      variant="outline" 
                      onClick={() => setServiceImportDialogOpen(true)} 
                      className="h-8"
                      disabled={selectedServiceIndex === null}
                    >
                      <Upload className="h-3.5 w-3.5 mr-1.5" />
                      Import
                    </Button>
                    <Select
                      value={selectedServiceIndex !== null ? String(selectedServiceIndex) : undefined}
                      onValueChange={(v) => setSelectedServiceIndex(Number(v))}
                    >
                      <SelectTrigger className="w-[180px] h-8 text-sm">
                        <SelectValue placeholder="Select service" />
                      </SelectTrigger>
                      <SelectContent>
                        {selectedServiceIndices.map((index) => {
                          const service = discoveredServices[index];
                          if (!service) return null;
                          const serviceName = serviceNameOverrides[index] || service.name;
                          const varCount = (serviceEnvironmentVariables[index] || []).filter((v) => v.key.trim()).length;
                          return (
                            <SelectItem key={index} value={String(index)}>
                              {serviceName} {varCount > 0 && `(${varCount})`}
                            </SelectItem>
                          );
                        })}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                
                {selectedServiceIndex !== null ? (
                  <div className="rounded-md border border-border/60 p-3 bg-card/50 space-y-2.5 max-h-[400px] overflow-y-auto">
                    {(serviceEnvironmentVariables[selectedServiceIndex] || []).map((kv, i) => {
                      const isOverriding = environmentVariables.some((gv) => gv.key === kv.key && gv.key.trim());
                      const serviceVars = serviceEnvironmentVariables[selectedServiceIndex] || [];
                      const isLastEmpty = i === serviceVars.length - 1 && kv.key === "" && kv.value === "";
                      
                      return (
                        <div key={i} className="grid grid-cols-12 gap-2 items-center">
                          <div className="col-span-5 relative">
                            <Input
                              placeholder={isLastEmpty ? "Type KEY..." : "KEY"}
                              className="bg-card h-9 text-sm"
                              value={kv.key}
                              onChange={(e) => {
                                const next = [...serviceVars];
                                next[i] = { ...next[i], key: e.target.value };
                                
                                // Auto-add new empty row when user types in the last row
                                if (i === serviceVars.length - 1 && e.target.value.trim()) {
                                  next.push({ key: "", value: "" });
                                }
                                
                                const updated = { ...serviceEnvironmentVariables };
                                updated[selectedServiceIndex] = next;
                                setField("serviceEnvironmentVariables", updated);
                              }}
                              onPaste={handleServicePaste}
                            />
                            {isOverriding && kv.key.trim() && (
                              <div className="absolute -top-1 -right-1 bg-yellow-500 text-[10px] px-1.5 py-0.5 rounded text-black font-medium">
                                ⚠
                              </div>
                            )}
                          </div>
                          <Input
                            placeholder="value"
                            className="bg-card col-span-6 h-9 text-sm"
                            value={kv.value}
                            onChange={(e) => {
                              const next = [...serviceVars];
                              next[i] = { ...next[i], value: e.target.value };
                              
                              // Auto-add new empty row when user types in the last row
                              if (i === serviceVars.length - 1 && e.target.value.trim()) {
                                next.push({ key: "", value: "" });
                              }
                              
                              const updated = { ...serviceEnvironmentVariables };
                              updated[selectedServiceIndex] = next;
                              setField("serviceEnvironmentVariables", updated);
                            }}
                            onPaste={handleServicePaste}
                          />
                          <Button
                            variant="ghost"
                            size="icon"
                            className="col-span-1 h-9 w-9"
                            disabled={isLastEmpty}
                            onClick={() => {
                              if (isLastEmpty) return;
                              const next = serviceVars.filter((_, idx) => idx !== i);
                              const updated = { ...serviceEnvironmentVariables };
                              updated[selectedServiceIndex] = next;
                              setField("serviceEnvironmentVariables", updated);
                            }}
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div className="rounded-md border border-border/60 p-4 bg-card/50 h-[400px] flex items-center justify-center">
                    <p className="text-sm text-muted-foreground">Select a service to configure</p>
                  </div>
                )}
              </div>
            </div>

            {/* Bottom: Effective Variables Preview - Full Width */}
            <div className="space-y-3">
              <Label>Effective Preview for Selected Service</Label>
              <ScrollArea className="h-[200px] rounded-md border border-border/60 p-4 bg-muted/20">
                {selectedServiceIndex !== null ? (
                  (() => {
                    const service = discoveredServices[selectedServiceIndex];
                    const serviceVars = serviceEnvironmentVariables[selectedServiceIndex] || [];
                    const systemVars: Record<string, string> = {};
                    if (service?.dockerfilePath) {
                      systemVars.RAILWAY_DOCKERFILE_PATH = service.dockerfilePath;
                    }

                    // Merge variables
                    const effective: Array<{ key: string; value: string; source: "global" | "service" | "system" }> = [];
                    
                    // Global first
                    environmentVariables
                      .filter((v) => v.key.trim())
                      .forEach((v) => effective.push({ ...v, source: "global" }));
                    
                    // Service overrides
                    serviceVars
                      .filter((v) => v.key.trim())
                      .forEach((sv) => {
                        const idx = effective.findIndex((e) => e.key === sv.key);
                        if (idx >= 0) {
                          effective[idx] = { ...sv, source: "service" };
                        } else {
                          effective.push({ ...sv, source: "service" });
                        }
                      });
                    
                    // System last
                    Object.entries(systemVars).forEach(([key, value]) => {
                      const idx = effective.findIndex((e) => e.key === key);
                      if (idx >= 0) {
                        effective[idx] = { key, value, source: "system" };
                      } else {
                        effective.push({ key, value, source: "system" });
                      }
                    });

                    return effective.length > 0 ? (
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                        {effective.map((v, i) => (
                          <div key={i} className="flex items-center gap-2 text-sm p-2 rounded-md bg-background/50">
                            <div className="flex items-center gap-1.5 flex-1 min-w-0">
                              <span className="font-mono font-semibold text-foreground/90 text-xs">{v.key}</span>
                              <span className="text-muted-foreground text-xs">=</span>
                              <span className="font-mono text-foreground/70 truncate text-xs">{v.value}</span>
                            </div>
                            <VariableSourceBadge source={v.source} className="shrink-0" />
                          </div>
                        ))}
                      </div>
                    ) : (
                      <p className="text-sm text-muted-foreground">No variables configured</p>
                    );
                  })()
                ) : (
                  <p className="text-sm text-muted-foreground">Select a service to preview</p>
                )}
              </ScrollArea>
            </div>
          </div>
        ) : (
          // Simple layout when no services discovered
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label className="text-sm">Global Variables</Label>
              <Button size="sm" variant="outline" onClick={() => setImportDialogOpen(true)}>
                <Upload className="h-4 w-4 mr-2" />
                Import .env
              </Button>
            </div>
            <div className="rounded-md border border-border/60 p-3 bg-card/50 space-y-2 max-h-[300px] overflow-y-auto">
              {environmentVariables.map((kv, i) => {
                const isLastEmpty = i === environmentVariables.length - 1 && kv.key === "" && kv.value === "";
                return (
                  <div key={i} className="grid grid-cols-12 gap-2 items-center">
                    <Input
                      placeholder={isLastEmpty ? "Type KEY to add..." : "KEY"}
                      className="bg-card col-span-5"
                      value={kv.key}
                      onChange={(e) => updateVar(i, "key", e.target.value)}
                      onPaste={handleGlobalPaste}
                    />
                    <Input
                      placeholder="value"
                      className="bg-card col-span-6"
                      value={kv.value}
                      onChange={(e) => updateVar(i, "value", e.target.value)}
                      onPaste={handleGlobalPaste}
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      className="col-span-1"
                      onClick={() => removeVar(i)}
                      disabled={isLastEmpty}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>

      <EnvImportDialog
        open={importDialogOpen}
        onOpenChange={setImportDialogOpen}
        onImport={handleImport}
        title="Import Global Variables"
        description="Import environment variables that will be shared across all services"
      />

      <EnvImportDialog
        open={serviceImportDialogOpen}
        onOpenChange={setServiceImportDialogOpen}
        onImport={handleServiceImport}
        title={`Import Variables for ${selectedServiceIndex !== null ? (serviceNameOverrides[selectedServiceIndex] || discoveredServices[selectedServiceIndex]?.name || 'Service') : 'Service'}`}
        description="Import service-specific environment variables"
      />
    </div>
  );
}


