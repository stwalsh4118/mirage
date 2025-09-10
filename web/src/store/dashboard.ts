import { create } from "zustand";
import { persist } from "zustand/middleware";

export type StatusFilter = "all" | "active" | "creating" | "error";
export type TypeFilter = "any" | "dev" | "prod";
export type SortBy = "updated" | "created" | "name" | "status";
export type ViewMode = "grid" | "list";

type DashboardState = {
  status: StatusFilter;
  type: TypeFilter;
  sortBy: SortBy;
  view: ViewMode;
  // removed density toggle per design
  query: string;
  setStatus: (s: StatusFilter) => void;
  setType: (t: TypeFilter) => void;
  setSortBy: (s: SortBy) => void;
  setView: (v: ViewMode) => void;
  setQuery: (q: string) => void;
};

export const useDashboardStore = create<DashboardState>()(
  persist(
    (set) => ({
      status: "all",
      type: "any",
      sortBy: "updated",
      view: "grid",
      
      query: "",
      setStatus: (status) => set({ status }),
      setType: (type) => set({ type }),
      setSortBy: (sortBy) => set({ sortBy }),
      setView: (view) => set({ view }),
      
      setQuery: (query) => set({ query }),
    }),
    { name: "mirage-dashboard-store" }
  )
);


