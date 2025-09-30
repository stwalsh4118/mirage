"use client";

import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Environment } from "@/lib/api/environments";
import { MoreHorizontal } from "lucide-react";
import { StatusChip } from "./StatusChip";
import { Pill } from "./Pill";

function typeBadge(type: Environment["type"]) {
  switch (type) {
    case "prod":
      return <Badge variant="default">Prod</Badge>;
    case "staging":
      return <Badge variant="secondary">Staging</Badge>;
    case "ephemeral":
      return <Badge variant="secondary">Ephemeral</Badge>;
    case "dev":
    default:
      return <Badge variant="secondary">Dev</Badge>;
  }
}

export function EnvironmentCard({ env }: { env: Environment }) {
  return (
    <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <CardTitle className="text-base font-medium">{env.name}</CardTitle>
            <span className="text-[10px]">{typeBadge(env.type)}</span>
          </div>
          <StatusChip status={mapStatus(env.status)} />
        </div>
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" aria-label="More actions">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link href={env.url ?? "#"} aria-disabled={!env.url} className={!env.url ? "pointer-events-none opacity-60" : ""}>Open</Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href={`/environments/${env.id}/logs`}>View logs</Link>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <div>
          <div className="text-[11px] text-muted-foreground mb-1">URL</div>
          <div className="flex items-center gap-2">
            <div className="relative flex-1">
              <div className="h-6 rounded-md bg-muted/60 border border-border/50 px-2 text-xs flex items-center text-muted-foreground truncate">
                {env.url ?? "Not deployed"}
              </div>
            </div>
            <Button variant="ghost" size="icon" aria-label="Copy URL" disabled={!env.url}
              onClick={() => env.url && navigator.clipboard.writeText(env.url!)}>
              🔗
            </Button>
          </div>
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>3 services</span>
          <span>{new Date(env.createdAt).toLocaleString()}</span>
        </div>
        <div className="text-xs text-muted-foreground">Latest: <Pill color="neutral">ab23de</Pill></div>
      </CardContent>
      <CardFooter className="flex items-center justify-between">
        <div className="flex gap-2">
          <Button asChild variant="outline" size="sm" aria-label="Open environment" disabled={!env.url}>
            <Link href={env.url ?? "#"} className={!env.url ? "pointer-events-none" : ""}>Open</Link>
          </Button>
          <Button asChild variant="outline" size="sm" aria-label="View logs">
            <Link href={`/environments/${env.id}/logs`}>Logs</Link>
          </Button>
        </div>
        <Button variant="ghost" size="icon" aria-label="More actions" className="rounded-md">
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      </CardFooter>
    </Card>
  );
}


function mapStatus(s: Environment["status"]): "Running" | "Stopped" | "Creating" | "Destroying" | "Error" | "Unknown" {
  switch (s) {
    case "active":
      return "Running";
    case "creating":
      return "Creating";
    case "destroying":
      return "Destroying";
    case "error":
      return "Error";
    default:
      return "Unknown";
  }
}

