"use client";

import { useMemo } from "react";
import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Trash2 } from "lucide-react";

export function StepConfig() {
  const { environmentName, templateKind, ttlHours, environmentVariables, setField } = useWizardStore();

  const preview = useMemo(() => {
    const lines = environmentVariables
      .filter((kv) => kv.key.trim().length > 0)
      .map((kv) => `${kv.key}=${kv.value}`);
    return lines.join("\n");
  }, [environmentVariables]);

  const updateVar = (index: number, key: "key" | "value", value: string) => {
    const next = environmentVariables.slice();
    next[index] = { ...next[index], [key]: value };
    setField("environmentVariables", next);
  };

  const addVar = () => setField("environmentVariables", [...environmentVariables, { key: "", value: "" }]);
  const removeVar = (index: number) => setField("environmentVariables", environmentVariables.filter((_, i) => i !== index));

  const importFromPaste = (text: string) => {
    const rows = text
      .split(/\r?\n/)
      .map((l) => l.trim())
      .filter(Boolean)
      .map((l) => {
        const eq = l.indexOf("=");
        if (eq === -1) return { key: l, value: "" };
        return { key: l.slice(0, eq), value: l.slice(eq + 1) };
      });
    setField("environmentVariables", [...environmentVariables, ...rows]);
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
        <div className="flex items-center justify-between">
          <Label>Environment variables</Label>
          <div className="space-x-2">
            <Button size="sm" variant="outline" onClick={addVar}>Add variable</Button>
            <Button size="sm" variant="outline" onClick={() => importFromPaste(prompt("Paste KEY=VALUE lines:") || "")}>Import</Button>
          </div>
        </div>

        <ResizablePanelGroup direction="horizontal" className="rounded-md border border-border/60">
          <ResizablePanel defaultSize={60} minSize={40} className="p-3">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {environmentVariables.map((kv, i) => (
                <div key={i} className="col-span-1 md:col-span-2 grid grid-cols-12 gap-3 items-center">
                  <Input
                    placeholder="KEY"
                    className="bg-card col-span-5"
                    value={kv.key}
                    onChange={(e) => updateVar(i, "key", e.target.value)}
                  />
                  <Input
                    placeholder="value"
                    className="bg-card col-span-6"
                    value={kv.value}
                    onChange={(e) => updateVar(i, "value", e.target.value)}
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    className="col-span-1 justify-self-end"
                    onClick={() => removeVar(i)}
                    aria-label="Remove variable"
                    title="Remove"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
              {environmentVariables.length === 0 && (
                <p className="text-sm text-muted-foreground">No variables added.</p>
              )}
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle />
          <ResizablePanel defaultSize={40} minSize={30} className="p-3 bg-muted/20">
            <Label className="text-sm">Inheritance preview</Label>
            <ScrollArea className="h-48 mt-2">
              <pre className="text-xs whitespace-pre-wrap leading-5 text-foreground/90">{preview || "(empty)"}</pre>
            </ScrollArea>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}


