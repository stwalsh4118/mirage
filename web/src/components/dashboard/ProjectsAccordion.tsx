"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { useDashboardStore } from "@/store/dashboard";
import type { RailwayProjectDetails } from "@/lib/api/railway";
import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { RailwayEnvironmentCard } from "./RailwayEnvironmentCard";

function LoadingList() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="glass grain p-4 rounded-lg">
          <div className="flex items-center justify-between">
            <div className="space-y-2 w-2/3">
              <Skeleton className="h-5 w-1/2" />
              <Skeleton className="h-4 w-2/3" />
            </div>
            <div className="flex items-center gap-3 w-1/3 justify-end">
              <Skeleton className="h-5 w-16" />
              <Skeleton className="h-5 w-16" />
              <Skeleton className="h-5 w-16" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export function ProjectsAccordion() {
  const { data, isLoading, isError, refetch } = useRailwayProjectsDetails();
  const { query, sortBy } = useDashboardStore();

  if (isLoading) return <LoadingList />;
  if (isError)
    return (
      <div className="glass grain p-6 rounded-lg">
        <div className="text-foreground/80">Failed to load projects.</div>
        <button className="mt-3 underline" onClick={() => void refetch()}>Retry</button>
      </div>
    );

  let projects = (data ?? []).slice();
  if (query) {
    const q = query.toLowerCase();
    projects = projects.filter((p) => p.name.toLowerCase().includes(q) || p.id.toLowerCase().includes(q));
  }
  switch (sortBy) {
    case "name":
      projects.sort((a, b) => a.name.localeCompare(b.name));
      break;
    case "services":
      projects.sort((a, b) => (b.services?.length ?? 0) - (a.services?.length ?? 0));
      break;
    case "environments":
      projects.sort((a, b) => (b.environments?.length ?? 0) - (a.environments?.length ?? 0));
      break;
    case "created":
    case "updated":
    default:
      break;
  }

  if (!projects.length)
    return (
      <div className="glass grain p-10 rounded-lg text-center">
        <div className="text-lg font-medium mb-2">No projects found</div>
        <div className="text-muted-foreground">Ensure Railway integration is configured.</div>
      </div>
    );

  return (
    <div className="space-y-3">
      {projects.map((p) => (
        <ProjectAccordionItem key={p.id} project={p} />
      ))}
    </div>
  );
}

function ProjectAccordionItem({ project }: { project: RailwayProjectDetails }) {
  const [isOpen, setIsOpen] = useState(false);
  const [maxHeight, setMaxHeight] = useState(0);
  const contentRef = useRef<HTMLDivElement>(null);

  const services = project.services?.length ?? 0;
  const envs = project.environments?.length ?? 0;

  useEffect(() => {
    if (contentRef.current) {
      setMaxHeight(contentRef.current.scrollHeight);
    }
  }, [project.environments]);

  useEffect(() => {
    if (isOpen && contentRef.current) {
      // Re-measure on open in case content height changed
      setMaxHeight(contentRef.current.scrollHeight);
    }
  }, [isOpen]);

  return (
    <details className="group glass grain rounded-lg border border-border/60 open:ring-1 open:ring-primary/20" onToggle={(e) => setIsOpen((e.currentTarget as HTMLDetailsElement).open)}>
      <summary className="list-none cursor-pointer select-none p-4 flex items-center justify-between rounded-lg bg-gradient-to-r from-muted/60 to-transparent hover:from-muted/70 transition-colors">
        <div className="flex items-center gap-3">
          <div className="h-6 w-1 rounded bg-primary/50" style={{ backgroundColor: isOpen ? "rgb(59 130 246)" : undefined }} />
          <div>
            <div className="text-sm font-medium">{project.name}</div>
            <div className="text-[11px] text-muted-foreground">{project.id}</div>
          </div>
        </div>
        <div className="flex items-center gap-3 text-xs">
          <div className="rounded-md border px-2 py-1 bg-primary/10 text-primary border-primary/20">services: <span className="font-medium">{services}</span></div>
          <div className="rounded-md border px-2 py-1 bg-emerald-500/10 text-emerald-600 dark:text-emerald-300 border-emerald-500/20">envs: <span className="font-medium">{envs}</span></div>
          <Button asChild variant="outline" size="sm" className="ml-2 bg-transparent" onClick={(e) => e.stopPropagation()}>
            <Link href={`/project/${project.id}`}>Open</Link>
          </Button>
          <span className="ml-2 text-muted-foreground transition-transform" style={{ transform: isOpen ? "rotate(180deg)" : undefined }}>▾</span>
        </div>
      </summary>
      <div className="transition-[max-height] duration-300 overflow-hidden" style={{ maxHeight: isOpen ? maxHeight : 0 }}>
        <div ref={contentRef} className={`px-4 pt-4 pb-4 transition-opacity duration-300 ${isOpen ? "opacity-100" : "opacity-0"}`}>
          {envs === 0 ? (
            <div className="text-xs text-muted-foreground">No environments</div>
          ) : (
            <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {project.environments?.map((e) => (
                <RailwayEnvironmentCard key={e.id} env={e} href={`/project/${project.id}`} />
              ))}
            </div>
          )}
        </div>
      </div>
    </details>
  );
}


