"use client";
import { KpiStrip } from "@/components/dashboard/KpiStrip";
import { ControlsBar } from "@/components/dashboard/ControlsBar";
import { ProjectsAccordion } from "@/components/dashboard/ProjectsAccordion";
import { ProjectsTable } from "@/components/dashboard/ProjectsTable";
import { useDashboardStore } from "@/store/dashboard";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";

export default function DashboardPage() {
  const { data: projects = [] } = useRailwayProjectsDetails();
  const totalProjects = projects.length;
  const totalServices = projects.reduce((sum, p) => sum + (p.services?.length ?? 0), 0);
  const totalEnvironments = projects.reduce((sum, p) => sum + (p.environments?.length ?? 0), 0);
  const kpis = [
    { title: "Projects", value: totalProjects },
    { title: "Services", value: totalServices },
    { title: "Environments", value: totalEnvironments },
  ];

  const { view } = useDashboardStore();
  return (
    <div className="space-y-6">
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        <KpiStrip items={kpis} />
        <ControlsBar />
        {view === "grid" ? <ProjectsAccordion /> : <ProjectsTable />}
      </main>
      {/* CommandMenu moved to shared layout */}
    </div>
  );
}


