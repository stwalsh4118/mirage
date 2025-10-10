"use client";
import { KpiStrip } from "@/components/dashboard/KpiStrip";
import { ControlsBar } from "@/components/dashboard/ControlsBar";
import { ProjectsAccordion } from "@/components/dashboard/ProjectsAccordion";
import { ProjectsTable } from "@/components/dashboard/ProjectsTable";
import { RailwayTokenOnboarding } from "@/components/dashboard/RailwayTokenOnboarding";
import { useDashboardStore } from "@/store/dashboard";
import { useRailwayTokenStatus } from "@/hooks/useRailwayToken";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";

export default function DashboardPage() {
  const { data: tokenStatus } = useRailwayTokenStatus();
  const { view } = useDashboardStore();
  
  // Only fetch projects if Railway token is configured
  const { data: projects = [] } = useRailwayProjectsDetails(undefined, {
    enabled: tokenStatus?.configured === true,
  });
  
  const totalProjects = projects.length;
  const totalServices = projects.reduce((sum, p) => sum + (p.services?.length ?? 0), 0);
  const totalEnvironments = projects.reduce((sum, p) => sum + (p.environments?.length ?? 0), 0);
  const kpis = [
    { title: "Projects", value: totalProjects },
    { title: "Services", value: totalServices },
    { title: "Environments", value: totalEnvironments },
  ];

  return (
    <div className="space-y-6">
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        <RailwayTokenOnboarding />
        <KpiStrip items={kpis} />
        <ControlsBar />
        {view === "grid" ? <ProjectsAccordion /> : <ProjectsTable />}
      </main>
      {/* CommandMenu moved to shared layout */}
    </div>
  );
}


