"use client";
import { EnvironmentGrid } from "@/components/dashboard/EnvironmentGrid";
import { EnvironmentList } from "@/components/dashboard/EnvironmentList";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { KpiStrip } from "@/components/dashboard/KpiStrip";
import { ControlsBar } from "@/components/dashboard/ControlsBar";
import { useDashboardStore } from "@/store/dashboard";
import { CommandMenu } from "@/components/dashboard/CommandMenu";

export default function DashboardPage() {
  const { view } = useDashboardStore();
  const kpis = [
    { title: "Total", value: 5, delta: 3 },
    { title: "Active", value: 2, delta: 2 },
    { title: "Creating", value: 1 },
    { title: "Errors", value: 1, delta: -1 },
  ];
  return (
    <div className="space-y-6">
      <DashboardHeader />
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        <KpiStrip items={kpis} />
        <ControlsBar />
        {view === "grid" ? <EnvironmentGrid /> : <EnvironmentList />}
      </main>
      <CommandMenu />
    </div>
  );
}


