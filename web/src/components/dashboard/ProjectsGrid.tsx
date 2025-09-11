"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { useDashboardStore } from "@/store/dashboard";
import { ProjectCard } from "./ProjectCard";

function LoadingGrid() {
  return (
    <div className="grid gap-6 grid-cols-1 md:grid-cols-2 xl:grid-cols-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="glass grain p-4 rounded-lg">
          <div className="space-y-3">
            <Skeleton className="h-5 w-1/2" />
            <Skeleton className="h-4 w-2/3" />
            <Skeleton className="h-4 w-1/3" />
            <Skeleton className="h-10 w-full" />
          </div>
        </div>
      ))}
    </div>
  );
}

export function ProjectsGrid() {
  const { data, isLoading, isError, refetch } = useRailwayProjectsDetails();
  const { query, sortBy, view } = useDashboardStore();

  if (isLoading) return <LoadingGrid />;
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
    case "plugins":
      projects.sort((a, b) => (b.plugins?.length ?? 0) - (a.plugins?.length ?? 0));
      break;
    case "environments":
      projects.sort((a, b) => (b.environments?.length ?? 0) - (a.environments?.length ?? 0));
      break;
    case "created":
    case "updated":
    default:
      // No created/updated metadata available from API yet; keep as-is.
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
    <div className={view === "grid" ? "grid gap-6 grid-cols-1 md:grid-cols-2 xl:grid-cols-3" : "space-y-3"}>
      {projects.map((p) => (
        <ProjectCard key={p.id} project={p} />
      ))}
    </div>
  );
}


