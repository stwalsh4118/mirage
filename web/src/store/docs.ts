import { create } from "zustand";
import { persist } from "zustand/middleware";

type DocsState = {
  openSections: string[];
  initialized: boolean;
  toggleSection: (sectionTitle: string) => void;
  isOpen: (sectionTitle: string) => boolean;
  setOpenSections: (sections: string[]) => void;
  setInitialized: (initialized: boolean) => void;
};

export const useDocsStore = create<DocsState>()(
  persist(
    (set, get) => ({
      openSections: [],
      initialized: false,
      
      toggleSection: (sectionTitle) =>
        set((state) => ({
          openSections: state.openSections.includes(sectionTitle)
            ? state.openSections.filter((s) => s !== sectionTitle)
            : [...state.openSections, sectionTitle],
        })),
      
      isOpen: (sectionTitle) => get().openSections.includes(sectionTitle),
      
      setOpenSections: (sections) => set({ openSections: sections, initialized: true }),
      
      setInitialized: (initialized) => set({ initialized }),
    }),
    { name: "mirage-docs-store" }
  )
);

