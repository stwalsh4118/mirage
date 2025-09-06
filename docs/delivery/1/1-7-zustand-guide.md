# 1-7 Zustand Guide

Date: 2025-09-06

[Official Docs](https://zustand.docs.pmnd.rs/)

## Overview
Zustand provides minimal, unopinionated client-state with simple stores and selectors. We use it for UI-only state (e.g., theme mode, panel toggles), not for server-state.

## Installation
Already installed in `web`:

```bash
pnpm add zustand
```

## Basic Store Pattern

```ts
// src/store/ui.ts
import { create } from "zustand";

export type ThemeMode = "light" | "dark" | "system";

type UIState = { themeMode: ThemeMode; setThemeMode: (mode: ThemeMode) => void };

export const useUIStore = create<UIState>((set) => ({
  themeMode: "system",
  setThemeMode: (mode) => set({ themeMode: mode }),
}));
```

Usage with selector to avoid unnecessary re-renders:

```tsx
import { useUIStore } from "@/store/ui";

function ThemeToggle() {
  const themeMode = useUIStore((s) => s.themeMode);
  const setThemeMode = useUIStore((s) => s.setThemeMode);
  return (
    <select value={themeMode} onChange={(e) => setThemeMode(e.target.value as any)}>
      <option value="system">System</option>
      <option value="light">Light</option>
      <option value="dark">Dark</option>
    </select>
  );
}
```

## Recommendations
- Use selectors `(state) => state.slice` for component subscriptions.
- Keep server-state in TanStack Query; store only ephemeral UI in Zustand.
- Add middleware (persist, devtools) later as needed, not in MVP.


