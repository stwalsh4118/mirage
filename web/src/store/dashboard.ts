import { create } from "zustand";
import { persist } from "zustand/middleware";

export type SortBy = "updated" | "created" | "name" | "services" | "environments";
export type ViewMode = "grid" | "list";

type DashboardState = {
  sortBy: SortBy;
  view: ViewMode;
  // removed density toggle per design
  query: string;
  setSortBy: (s: SortBy) => void;
  setView: (v: ViewMode) => void;
  setQuery: (q: string) => void;
};

export const useDashboardStore = create<DashboardState>()(
  persist(
    (set) => ({
      sortBy: "name",
      view: "grid",
      
      query: "",
      setSortBy: (sortBy) => set({ sortBy }),
      setView: (view) => set({ view }),
      
      setQuery: (query) => set({ query }),
    }),
    { name: "mirage-dashboard-store" }
  )
);


