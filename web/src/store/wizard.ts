"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";

export type WizardStepId = "project" | "source" | "config" | "review";

export const WIZARD_STEPS_ORDER: WizardStepId[] = [
  "project",
  "source",
  "config",
  "review",
];

export type ProjectSelectionMode = "existing" | "new";

export type WizardState = {
  currentStepIndex: number;
  requestId: string | null;

  // Step 0: Project
  projectSelectionMode: ProjectSelectionMode;
  existingProjectId: string | null;
  existingProjectName: string | null;
  newProjectName: string;
  defaultEnvironmentName: string;

  // Step 1: Source
  repositoryUrl: string;
  repositoryBranch: string;

  // Step 2: Config
  environmentName: string;
  templateKind: "dev" | "prod";
  ttlHours: number | null;
  environmentVariables: Array<{ key: string; value: string }>;

  // Step 3: Strategy
  deploymentStrategy: "sequential" | "parallel";

  // Derived IDs returned during submission (for resume)
  createdProjectId?: string;
  createdEnvironmentId?: string;

  // Actions
  goNext: () => void;
  goBack: () => void;
  goTo: (index: number) => void;
  setField: <K extends keyof WizardState>(key: K, value: WizardState[K]) => void;
  reset: () => void;
};

const initialState: Omit<WizardState,
  | "goNext"
  | "goBack"
  | "goTo"
  | "setField"
  | "reset"
> = {
  currentStepIndex: 0,
  requestId: null,

  projectSelectionMode: "existing",
  existingProjectId: null,
  existingProjectName: null,
  newProjectName: "",
  defaultEnvironmentName: "production",

  repositoryUrl: "",
  repositoryBranch: "main",

  environmentName: "",
  templateKind: "dev",
  ttlHours: null,
  environmentVariables: [],

  deploymentStrategy: "sequential",

  createdProjectId: undefined,
  createdEnvironmentId: undefined,
};

export const useWizardStore = create<WizardState>()(
  persist(
    (set, get) => ({
      ...initialState,
      goNext: () => {
        const next = Math.min(get().currentStepIndex + 1, WIZARD_STEPS_ORDER.length - 1);
        set({ currentStepIndex: next });
      },
      goBack: () => {
        const prev = Math.max(get().currentStepIndex - 1, 0);
        set({ currentStepIndex: prev });
      },
      goTo: (index) => {
        const bounded = Math.max(0, Math.min(index, WIZARD_STEPS_ORDER.length - 1));
        set({ currentStepIndex: bounded });
      },
      setField: (key, value) => set({ [key]: value } as Partial<WizardState>),
      reset: () => set({ ...initialState }),
    }),
    { name: "mirage-wizard-store" }
  )
);


