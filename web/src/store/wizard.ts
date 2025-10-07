"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";

export type WizardStepId = "project" | "source" | "discovery" | "config" | "review";

export const WIZARD_STEPS_ORDER: WizardStepId[] = [
  "project",
  "source",
  "discovery",
  "config",
  "review",
];

export type ProjectSelectionMode = "existing" | "new";

export type ProvisionStage = 
  | 'idle'
  | 'creating-project'
  | 'creating-environment'
  | 'creating-services'
  | 'applying-config'
  | 'complete'
  | 'failed';

export type StageStatus = 'pending' | 'running' | 'success' | 'error';

export interface ServiceProgress {
  name: string;
  status: StageStatus;
  serviceId?: string;
  error?: string;
}

export interface StageInfo {
  status: StageStatus;
  error?: string;
  startedAt?: number;
  completedAt?: number;
}

export type WizardState = {
  currentStepIndex: number;
  requestId: string | null;

  // Clone mode
  sourceMode: 'new' | 'clone';
  cloneSourceEnvId: string | null;

  // Step 0: Project
  projectSelectionMode: ProjectSelectionMode;
  existingProjectId: string | null;
  existingProjectName: string | null;
  newProjectName: string;
  defaultEnvironmentName: string;

  // Step 1: Source
  deploymentSource: "repository" | "image";
  repositoryUrl: string;
  repositoryBranch: string;
  githubToken: string; // Optional GitHub token for private repos
  imageName: string;
  imageRegistry: string;
  imageTag: string;
  imageDigest: string;
  useDigest: boolean;
  imagePorts: number[];

  // Step 2: Service Discovery
  discoveryTriggered: boolean;
  discoveryLoading: boolean;
  discoveryError: string | null;
  discoveredServices: Array<{
    name: string;
    dockerfilePath: string;
    buildContext: string;
    exposedPorts: number[];
    buildArgs: string[];
    baseImage: string;
  }>;
  selectedServiceIndices: number[]; // Indices of selected services
  serviceNameOverrides: Record<number, string>; // Index -> edited name
  discoverySkipped: boolean;

  // Step 3: Config
  environmentName: string;
  templateKind: "dev" | "prod";
  ttlHours: number | null;
  environmentVariables: Array<{ key: string; value: string }>;
  serviceEnvironmentVariables: Record<number, Array<{ key: string; value: string }>>; // Per-service env vars keyed by service index

  // Step 4: Strategy
  deploymentStrategy: "sequential" | "parallel";

  // Provision state
  isProvisioning: boolean;
  currentStage: ProvisionStage;
  stages: Record<string, StageInfo>;
  serviceProgress: ServiceProgress[];

  // Derived IDs returned during submission (for resume)
  createdProjectId?: string;
  createdEnvironmentId?: string;
  createdServiceIds?: string[];

  // Actions
  goNext: () => void;
  goBack: () => void;
  goTo: (index: number) => void;
  setField: <K extends keyof WizardState>(key: K, value: WizardState[K]) => void;
  setStageStatus: (stage: ProvisionStage, status: StageStatus, error?: string) => void;
  setServiceProgress: (services: ServiceProgress[]) => void;
  updateServiceProgress: (index: number, update: Partial<ServiceProgress>) => void;
  startProvisioning: () => void;
  reset: () => void;
};

const initialState: Omit<WizardState,
  | "goNext"
  | "goBack"
  | "goTo"
  | "setField"
  | "setStageStatus"
  | "setServiceProgress"
  | "updateServiceProgress"
  | "startProvisioning"
  | "reset"
> = {
  currentStepIndex: 0,
  requestId: null,

  sourceMode: "new",
  cloneSourceEnvId: null,

  projectSelectionMode: "existing",
  existingProjectId: null,
  existingProjectName: null,
  newProjectName: "",
  defaultEnvironmentName: "production",

  deploymentSource: "repository",
  repositoryUrl: "",
  repositoryBranch: "main",
  githubToken: "",
  imageName: "",
  imageRegistry: "",
  imageTag: "latest",
  imageDigest: "",
  useDigest: false,
  imagePorts: [],

  discoveryTriggered: false,
  discoveryLoading: false,
  discoveryError: null,
  discoveredServices: [],
  selectedServiceIndices: [],
  serviceNameOverrides: {},
  discoverySkipped: false,

  environmentName: "",
  templateKind: "dev",
  ttlHours: null,
  environmentVariables: [],
  serviceEnvironmentVariables: {},

  deploymentStrategy: "sequential",

  isProvisioning: false,
  currentStage: "idle",
  stages: {},
  serviceProgress: [],

  createdProjectId: undefined,
  createdEnvironmentId: undefined,
  createdServiceIds: undefined,
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
      setStageStatus: (stage, status, error) => {
        const stages = { ...get().stages };
        const now = Date.now();
        stages[stage] = {
          ...stages[stage],
          status,
          error,
          startedAt: stages[stage]?.startedAt || (status === 'running' ? now : undefined),
          completedAt: status === 'success' || status === 'error' ? now : undefined,
        };
        set({ 
          stages,
          currentStage: stage,
          // Keep isProvisioning true until explicitly complete or failed
          isProvisioning: status === 'running' || status === 'pending',
        });
      },
      setServiceProgress: (services) => set({ serviceProgress: services }),
      updateServiceProgress: (index, update) => {
        const serviceProgress = [...get().serviceProgress];
        serviceProgress[index] = { ...serviceProgress[index], ...update };
        set({ serviceProgress });
      },
      startProvisioning: () => {
        set({ 
          isProvisioning: true,
          currentStage: 'idle',
          stages: {},
          serviceProgress: [],
          requestId: crypto.randomUUID(),
        });
      },
      reset: () => set({ ...initialState }),
    }),
    { name: "mirage-wizard-store" }
  )
);


