import { EnvironmentGrid } from "@/components/dashboard/EnvironmentGrid";

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Dashboard</h1>
      </div>
      <EnvironmentGrid />
    </div>
  );
}


