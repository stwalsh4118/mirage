"use client";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { KpiStrip } from "@/components/dashboard/KpiStrip";
import { ControlsBar } from "@/components/dashboard/ControlsBar";
import { CommandMenu } from "@/components/dashboard/CommandMenu";
import { ProjectsAccordion } from "@/components/dashboard/ProjectsAccordion";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";

export default function DashboardPage() {
  const { data: projects = [] } = useRailwayProjectsDetails();
  console.log(projects)
  const totalProjects = projects.length;
  const totalServices = projects.reduce((sum, p) => sum + (p.services?.length ?? 0), 0);
  const totalPlugins = projects.reduce((sum, p) => sum + (p.plugins?.length ?? 0), 0);
  const totalEnvironments = projects.reduce((sum, p) => sum + (p.environments?.length ?? 0), 0);
  const kpis = [
    { title: "Projects", value: totalProjects },
    { title: "Services", value: totalServices },
    { title: "Plugins", value: totalPlugins },
    { title: "Environments", value: totalEnvironments },
  ];

  return (
    <div className="space-y-6">
      <DashboardHeader />
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        <KpiStrip items={kpis} />
        <ControlsBar />
        <ProjectsAccordion />
      </main>
      <CommandMenu />
    </div>
  );
}


