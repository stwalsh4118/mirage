"use client";
import { EnvironmentGrid } from "@/components/dashboard/EnvironmentGrid";
import { EnvironmentList } from "@/components/dashboard/EnvironmentList";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { KpiStrip } from "@/components/dashboard/KpiStrip";
import { ControlsBar } from "@/components/dashboard/ControlsBar";
import { useDashboardStore } from "@/store/dashboard";
import { CommandMenu } from "@/components/dashboard/CommandMenu";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";

export default function DashboardPage() {
  const { view } = useDashboardStore();
  const kpis = [
    { title: "Total", value: 5, delta: 3 },
    { title: "Active", value: 2, delta: 2 },
    { title: "Creating", value: 1 },
    { title: "Errors", value: 1, delta: -1 },
  ];

  const { data: projects, isLoading, isError } = useRailwayProjectsDetails();
  console.log(projects)

  return (
    <div className="space-y-6">
      <DashboardHeader />
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        <div className="glass grain p-3 rounded-lg">
          <div className="text-xs text-muted-foreground">Railway Projects Probe</div>
          {isLoading && <div className="text-sm">Loading…</div>}
          {isError && <div className="text-sm text-red-600">Failed to fetch projects</div>}
          {projects && (
            <div className="mt-2 space-y-2">
              <div className="text-xs">Found {projects.length} project(s)</div>
              <ul className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
                {projects.map((p) => (
                  <li key={p.id} className="text-sm">
                    <div className="font-medium">{p.name} <span className="text-muted-foreground">({p.id})</span></div>
                    <div className="text-xs text-muted-foreground mt-1">
                      services: {(p.services?.length ?? 0)} · plugins: {(p.plugins?.length ?? 0)} · envs: {(p.environments?.length ?? 0)}
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
        <KpiStrip items={kpis} />
        <ControlsBar />
        {view === "grid" ? <EnvironmentGrid /> : <EnvironmentList />}
      </main>
      <CommandMenu />
    </div>
  );
}


