"use client";

import { Environment } from "@/lib/api/environments";
import { useEnvironments } from "@/hooks/useEnvironments";
import { useDashboardStore } from "@/store/dashboard";
import { EnvironmentCard } from "./EnvironmentCard";
import { Skeleton } from "@/components/ui/skeleton";

function LoadingGrid() {
  return (
    <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {Array.from({ length: 8 }).map((_, i) => (
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

export function EnvironmentGrid() {
  const { data, isLoading, isError, refetch } = useEnvironments();
  const { query, sortBy } = useDashboardStore();

  if (isLoading) return <LoadingGrid />;
  if (isError)
    return (
      <div className="glass grain p-6 rounded-lg">
        <div className="text-foreground/80">Failed to load environments.</div>
        <button className="mt-3 underline" onClick={() => void refetch()}>Retry</button>
      </div>
    );

  let environments = (data ?? []) as Environment[];
  // Simple client-side filters for now
  environments = environments.filter((e) => {
    if (query && !e.name.toLowerCase().includes(query.toLowerCase())) return false;
    return true;
  });

  // Sorting
  environments = environments.slice().sort((a, b) => {
    if (sortBy === "name") {
      return a.name.localeCompare(b.name);
    }
    // Default: sort by created date (newest first)
    const ad = Date.parse(a.createdAt);
    const bd = Date.parse(b.createdAt);
    return bd - ad;
  });
  if (!environments.length)
    return (
      <div className="glass grain p-10 rounded-lg text-center">
        <div className="text-lg font-medium mb-2">No environments yet</div>
        <div className="text-muted-foreground">Create your first environment to get started.</div>
      </div>
    );

  return (
    <div className="grid gap-6 grid-cols-1 md:grid-cols-2 xl:grid-cols-3">
      {environments.map((env) => (
        <EnvironmentCard key={env.id} env={env} />
      ))}
    </div>
  );
}



