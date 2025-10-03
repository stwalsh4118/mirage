"use client";

import { useMemo, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { VariableSourceBadge } from "./VariableSourceBadge";
import { Trash2, Plus, Upload } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { EnvImportDialog } from "./EnvImportDialog";

interface Variable {
  key: string;
  value: string;
}

interface ServiceVariableEditorProps {
  serviceName: string;
  serviceIndex: number;
  globalVariables: Variable[];
  serviceVariables: Variable[];
  onServiceVariablesChange: (variables: Variable[]) => void;
  systemVariables?: Record<string, string>; // Optional system variables like RAILWAY_DOCKERFILE_PATH
}

export function ServiceVariableEditor({
  serviceName,
  serviceIndex,
  globalVariables,
  serviceVariables,
  onServiceVariablesChange,
  systemVariables = {},
}: ServiceVariableEditorProps) {
  const [importDialogOpen, setImportDialogOpen] = useState(false);

  const handleImport = (imported: Array<{ key: string; value: string }>) => {
    onServiceVariablesChange([...serviceVariables, ...imported]);
  };

  // Calculate effective variables (merged)
  const effectiveVariables = useMemo(() => {
    const result: Array<{ key: string; value: string; source: "global" | "service" | "system" }> = [];
    
    // Add global variables first
    globalVariables
      .filter((v) => v.key.trim().length > 0)
      .forEach((v) => {
        result.push({ ...v, source: "global" });
      });
    
    // Merge/override with service-specific variables
    serviceVariables
      .filter((v) => v.key.trim().length > 0)
      .forEach((sv) => {
        const existingIndex = result.findIndex((r) => r.key === sv.key);
        if (existingIndex >= 0) {
          // Override global variable
          result[existingIndex] = { ...sv, source: "service" };
        } else {
          // New service-specific variable
          result.push({ ...sv, source: "service" });
        }
      });
    
    // Add system variables last (these always override)
    Object.entries(systemVariables).forEach(([key, value]) => {
      const existingIndex = result.findIndex((r) => r.key === key);
      if (existingIndex >= 0) {
        result[existingIndex] = { key, value, source: "system" };
      } else {
        result.push({ key, value, source: "system" });
      }
    });
    
    return result;
  }, [globalVariables, serviceVariables, systemVariables]);

  const updateVar = (index: number, key: "key" | "value", value: string) => {
    const next = serviceVariables.slice();
    next[index] = { ...next[index], [key]: value };
    onServiceVariablesChange(next);
  };

  const addVar = () => {
    onServiceVariablesChange([...serviceVariables, { key: "", value: "" }]);
  };

  const removeVar = (index: number) => {
    onServiceVariablesChange(serviceVariables.filter((_, i) => i !== index));
  };

  const isOverridingGlobal = (key: string): boolean => {
    return globalVariables.some((gv) => gv.key === key && gv.key.trim().length > 0);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <Label className="text-base">Variables for {serviceName}</Label>
          <p className="text-xs text-muted-foreground mt-1">
            Add service-specific variables or override global variables
          </p>
        </div>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={() => setImportDialogOpen(true)}>
            <Upload className="h-4 w-4 mr-1" />
            Import
          </Button>
          <Button size="sm" variant="outline" onClick={addVar}>
            <Plus className="h-4 w-4 mr-1" />
            Add
          </Button>
        </div>
      </div>

      {/* Split layout: Variables on left, Preview on right */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Left side: Service-specific variables editor */}
        <div className="space-y-2">
          <Label className="text-sm">Service-Specific Variables</Label>
          <div className="space-y-2 max-h-[400px] overflow-y-auto rounded-md border border-border/60 p-3 bg-card/50">
            {serviceVariables.map((kv, i) => (
              <div key={i} className="grid grid-cols-12 gap-2 items-center">
                <div className="col-span-5 relative">
                  <Input
                    placeholder="KEY"
                    className="bg-card"
                    value={kv.key}
                    onChange={(e) => updateVar(i, "key", e.target.value)}
                  />
                  {isOverridingGlobal(kv.key) && kv.key.trim().length > 0 && (
                    <div className="absolute -top-1 -right-1">
                      <div className="bg-yellow-500 text-xs px-1 py-0.5 rounded text-black font-medium">
                        Override
                      </div>
                    </div>
                  )}
                </div>
                <Input
                  placeholder="value"
                  className="bg-card col-span-6"
                  value={kv.value}
                  onChange={(e) => updateVar(i, "value", e.target.value)}
                />
                <Button
                  variant="ghost"
                  size="icon"
                  className="col-span-1"
                  onClick={() => removeVar(i)}
                  aria-label="Remove variable"
                  title="Remove"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))}
            {serviceVariables.length === 0 && (
              <p className="text-sm text-muted-foreground">No service-specific variables added.</p>
            )}
          </div>
        </div>

        {/* Right side: Preview of effective variables */}
        <div className="space-y-2">
          <div>
            <Label className="text-sm">Effective Variables Preview</Label>
            <p className="text-xs text-muted-foreground mt-0.5">
              All variables this service will receive
            </p>
          </div>
          <ScrollArea className="h-[400px] rounded-md border border-border/60 p-3 bg-muted/20">
            {effectiveVariables.length > 0 ? (
              <div className="space-y-2">
                {effectiveVariables.map((v, i) => (
                  <div key={i} className="flex items-center justify-between text-xs">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <span className="font-mono font-semibold text-foreground/90">{v.key}</span>
                      <span className="text-muted-foreground">=</span>
                      <span className="font-mono text-foreground/70 truncate">{v.value}</span>
                    </div>
                    <VariableSourceBadge source={v.source} className="ml-2 shrink-0" />
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-xs text-muted-foreground">No variables configured</p>
            )}
          </ScrollArea>
        </div>
      </div>

      <EnvImportDialog
        open={importDialogOpen}
        onOpenChange={setImportDialogOpen}
        onImport={handleImport}
        title={`Import Variables for ${serviceName}`}
        description="Import service-specific environment variables"
      />
    </div>
  );
}

