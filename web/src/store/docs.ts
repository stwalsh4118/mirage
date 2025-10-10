import { create } from "zustand";
import { persist } from "zustand/middleware";

type DocsState = {
  openSections: string[];
  toggleSection: (sectionTitle: string) => void;
  isOpen: (sectionTitle: string) => boolean;
  setOpenSections: (sections: string[]) => void;
};

export const useDocsStore = create<DocsState>()(
  persist(
    (set, get) => ({
      openSections: [],
      
      toggleSection: (sectionTitle) =>
        set((state) => ({
          openSections: state.openSections.includes(sectionTitle)
            ? state.openSections.filter((s) => s !== sectionTitle)
            : [...state.openSections, sectionTitle],
        })),
      
      isOpen: (sectionTitle) => get().openSections.includes(sectionTitle),
      
      setOpenSections: (sections) => set({ openSections: sections }),
    }),
    { name: "mirage-docs-store" }
  )
);

