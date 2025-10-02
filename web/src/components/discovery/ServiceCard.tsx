"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Copy, Check, Container } from "lucide-react";
import { cn } from "@/lib/utils";
import type { DiscoveredService } from "./types";

const MAX_VISIBLE_ITEMS = 3;

interface ServiceCardProps {
  service: DiscoveredService;
  selected: boolean;
  editedName?: string;
  onToggleSelection: () => void;
  onNameChange: (name: string) => void;
}

export function ServiceCard({
  service,
  selected,
  editedName,
  onToggleSelection,
  onNameChange,
}: ServiceCardProps) {
  const [isEditingName, setIsEditingName] = useState(false);
  const [copied, setCopied] = useState(false);
  const displayName = editedName || service.name;

  const handleCopyPath = async () => {
    await navigator.clipboard.writeText(service.dockerfilePath);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleNameBlur = () => {
    setIsEditingName(false);
    if (!editedName?.trim()) {
      onNameChange(service.name);
    }
  };

  const handleNameKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleNameBlur();
    } else if (e.key === "Escape") {
      onNameChange(service.name);
      setIsEditingName(false);
    }
  };

  return (
    <Card
      className={cn(
        "glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]",
        selected && "ring-2 ring-primary/50"
      )}
    >
      <CardHeader className="pb-3">
        <div className="flex items-start gap-3">
          <Checkbox
            checked={selected}
            onCheckedChange={onToggleSelection}
            className="mt-1"
          />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <Container className="size-4 text-muted-foreground shrink-0" />
              {isEditingName ? (
                <Input
                  value={displayName}
                  onChange={(e) => onNameChange(e.target.value)}
                  onBlur={handleNameBlur}
                  onKeyDown={handleNameKeyDown}
                  className="h-7 text-sm font-semibold"
                  autoFocus
                />
              ) : (
                <button
                  onClick={() => setIsEditingName(true)}
                  className="text-sm font-semibold text-foreground hover:text-foreground/80 transition-colors text-left"
                  title="Click to edit service name"
                >
                  {displayName}
                </button>
              )}
            </div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <span className="truncate">{service.dockerfilePath}</span>
              <Button
                variant="ghost"
                size="sm"
                className="h-5 w-5 p-0 shrink-0"
                onClick={handleCopyPath}
                title="Copy Dockerfile path"
              >
                {copied ? (
                  <Check className="size-3 text-green-500" />
                ) : (
                  <Copy className="size-3" />
                )}
              </Button>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0 space-y-3">
        {/* Base Image */}
        {service.baseImage && (
          <div>
            <p className="text-xs font-medium text-muted-foreground mb-1">
              Base Image
            </p>
            <Badge variant="secondary" className="font-mono text-xs">
              {service.baseImage}
            </Badge>
          </div>
        )}

        {/* Build Context */}
        <div>
          <p className="text-xs font-medium text-muted-foreground mb-1">
            Build Context
          </p>
          <p className="text-xs font-mono text-foreground/80">
            {service.buildContext}
          </p>
        </div>

        {/* Exposed Ports */}
        {service.exposedPorts.length > 0 && (
          <div>
            <p className="text-xs font-medium text-muted-foreground mb-1">
              Exposed Ports
            </p>
            <div className="flex flex-wrap gap-1.5">
              {service.exposedPorts.slice(0, MAX_VISIBLE_ITEMS).map((port) => (
                <Badge key={port} variant="outline" className="font-mono">
                  {port}
                </Badge>
              ))}
              {service.exposedPorts.length > MAX_VISIBLE_ITEMS && (
                <Badge variant="outline" className="text-muted-foreground">
                  +{service.exposedPorts.length - MAX_VISIBLE_ITEMS} more
                </Badge>
              )}
            </div>
          </div>
        )}

        {/* Build Arguments */}
        {service.buildArgs.length > 0 && (
          <div>
            <p className="text-xs font-medium text-muted-foreground mb-1">
              Build Arguments
            </p>
            <div className="flex flex-wrap gap-1.5">
              {service.buildArgs.slice(0, MAX_VISIBLE_ITEMS).map((arg) => (
                <Badge key={arg} variant="outline" className="font-mono">
                  {arg}
                </Badge>
              ))}
              {service.buildArgs.length > MAX_VISIBLE_ITEMS && (
                <Badge variant="outline" className="text-muted-foreground">
                  +{service.buildArgs.length - MAX_VISIBLE_ITEMS} more
                </Badge>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

