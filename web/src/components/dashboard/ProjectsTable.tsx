"use client";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { useRailwayTokenStatus } from "@/hooks/useRailwayToken";
import { useDashboardStore } from "@/store/dashboard";
import { Pill } from "./Pill";

export function ProjectsTable() {
  const { data, isLoading, isError, refetch } = useRailwayProjectsDetails();
  const { data: tokenStatus } = useRailwayTokenStatus();
  const { query, sortBy } = useDashboardStore();

  if (isLoading)
    return (
      <div className="glass grain p-6 rounded-lg text-sm text-muted-foreground">Loading projects…</div>
    );
  
  // Don't show error if no token is configured - onboarding banner will guide them
  if (isError && tokenStatus?.configured) {
    return (
      <div className="glass grain p-6 rounded-lg">
        <div className="text-foreground/80">Failed to load projects.</div>
        <button className="mt-3 underline" onClick={() => void refetch()}>Retry</button>
      </div>
    );
  }
  
  // If no token configured, don't show anything (onboarding banner handles it)
  if (!tokenStatus?.configured) {
    return null;
  }

  let projects = (data ?? []).slice();
  if (query) {
    const q = query.toLowerCase();
    projects = projects.filter((p) => p.name.toLowerCase().includes(q) || p.id.toLowerCase().includes(q));
  }
  switch (sortBy) {
    case "services":
      projects.sort((a, b) => (b.services?.length ?? 0) - (a.services?.length ?? 0));
      break;
    case "environments":
      projects.sort((a, b) => (b.environments?.length ?? 0) - (a.environments?.length ?? 0));
      break;
    case "name":
    default:
      projects.sort((a, b) => a.name.localeCompare(b.name));
      break;
  }

  if (!projects.length)
    return (
      <div className="glass grain p-10 rounded-lg text-center space-y-3">
        <div className="text-lg font-medium mb-2">No Railway projects found</div>
        <div className="text-muted-foreground max-w-md mx-auto">
          {query ? "No projects match your search. Try a different query." : "Your Railway account doesn't have any projects yet. Create one in Railway to get started."}
        </div>
      </div>
    );

  return (
    <div className="glass grain rounded-lg border border-border/60 shadow-sm">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="bg-muted/40">Name</TableHead>
            <TableHead className="bg-muted/40">ID</TableHead>
            <TableHead className="text-right bg-muted/40">Services</TableHead>
            <TableHead className="text-right bg-muted/40">Environments</TableHead>
            <TableHead className="text-right bg-muted/40">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {projects.map((p) => {
            const services = p.services?.length ?? 0;
            const envs = p.environments?.length ?? 0;
            return (
              <TableRow key={p.id} className="cursor-pointer" onClick={(e) => { if ((e.target as HTMLElement).closest("a,button")) return; }}>
                <TableCell className="font-medium">
                  <Link href={`/project/${p.id}`} className="hover:underline">{p.name}</Link>
                </TableCell>
                <TableCell className="text-xs text-muted-foreground">{p.id}</TableCell>
                <TableCell className="text-right">
                  <Pill color={services > 0 ? "green" : "amber"}>{services}</Pill>
                </TableCell>
                <TableCell className="text-right">
                  <Pill color={envs > 0 ? "accent" : "neutral"}>{envs}</Pill>
                </TableCell>
                <TableCell className="text-right">
                  <Button asChild variant="outline" size="sm" className="bg-transparent">
                    <Link href={`/project/${p.id}`}>Open</Link>
                  </Button>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}


