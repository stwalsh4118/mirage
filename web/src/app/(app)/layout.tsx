import type { ReactNode } from "react";
import { DashboardHeader } from "@/components/dashboard/DashboardHeader";
import { CommandMenu } from "@/components/dashboard/CommandMenu";

export default function AppLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background sandstorm-bg">
      <DashboardHeader />
      <div className="pt-6">
        {children}
      </div>
      <CommandMenu />
    </div>
  );
}


