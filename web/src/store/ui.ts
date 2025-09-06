import { create } from "zustand";

export type ThemeMode = "light" | "dark" | "system";

type UIState = {
    themeMode: ThemeMode;
    setThemeMode: (mode: ThemeMode) => void;
};

export const useUIStore = create<UIState>((set) => ({
    themeMode: "system",
    setThemeMode: (mode) => set({ themeMode: mode }),
}));


