// Railway API Types - used by hooks in hooks/useRailway.ts

export type RailwayProject = { id: string; name: string };
export type RailwayProjectItem = { id: string; name: string };

export type ServiceSource = {
  image?: string | null;
  repo?: string | null;
};

export type LatestDeployment = {
  canRedeploy?: boolean | null;
  canRollback?: boolean | null;
  createdAt?: string | null;
  deploymentStopped?: boolean | null;
  environmentId?: string | null;
  id?: string | null;
  meta?: Record<string, unknown> | null;
  projectId?: string | null;
  serviceId?: string | null;
  snapshotId?: string | null;
  staticUrl?: string | null;
  status?: string | null;
  statusUpdatedAt?: string | null;
  suggestAddServiceDomain?: boolean | null;
  updatedAt?: string | null;
  url?: string | null;
};

export type RailwayServiceInstance = {
  id: string;
  serviceId: string;
  serviceName: string;
  environmentId: string;
  buildCommand?: string | null;
  builder?: string | null;
  createdAt?: string | null;
  cronSchedule?: string | null;
  deletedAt?: string | null;
  drainingSeconds?: number | null;
  healthcheckPath?: string | null;
  healthcheckTimeout?: number | null;
  isUpdatable?: boolean | null;
  nextCronRunAt?: string | null;
  nixpacksPlan?: string | null;
  numReplicas?: number | null;
  overlapSeconds?: number | null;
  preDeployCommand?: string | null;
  railpackInfo?: string | null;
  railwayConfigFile?: string | null;
  region?: string | null;
  restartPolicyMaxRetries?: number | null;
  restartPolicyType?: string | null;
  rootDirectory?: string | null;
  sleepApplication?: boolean | null;
  startCommand?: string | null;
  updatedAt?: string | null;
  upstreamUrl?: string | null;
  watchPatterns?: string[];
  source?: ServiceSource;
  latestDeployment?: LatestDeployment;
};

export type RailwayEnvironmentWithServices = {
  id: string;
  name: string;
  services: RailwayServiceInstance[];
};

export type RailwayProjectDetails = {
  id: string;
  name: string;
  services: RailwayProjectItem[];
  environments: RailwayEnvironmentWithServices[];
};


export type ImportRailwayEnvsRequest = {
  projectId: string;
  environmentIds: string[];
};

export type ImportRailwayEnvsResponse = {
  imported: number;
  skipped: number;
  items: { id: string; name: string; type: string; status: string; createdAt: string }[];
};


// Provision: Create Project
export type ProvisionProjectRequest = {
  defaultEnvironmentName?: string;
  name?: string;
  requestId: string;
};

export type ProvisionProjectResponse = {
  projectId: string;
  baseEnvironmentId: string;      // Mirage internal ID for foreign keys
  railwayEnvironmentId: string;   // Railway ID for Railway API calls
  name: string;
};


// Provision: Create Environment
export type ProvisionEnvironmentRequest = {
  projectId: string;
  name: string;
  requestId: string;
  envType?: 'dev' | 'staging' | 'prod' | 'ephemeral';
  wizardInputs?: Record<string, unknown>;
};

export type ProvisionEnvironmentResponse = {
  environmentId: string;        // Mirage internal ID for foreign keys
  railwayEnvironmentId: string; // Railway ID for Railway API calls
};


// Provision: Create Services
export type ProvisionServicesRequest = {
  projectId: string;
  environmentId: string;        // Mirage internal ID for database FK
  railwayEnvironmentId?: string; // Railway ID for Railway API calls (optional for backward compat)
  services: {
    name: string;
    // Repository-based deployment
    repo?: string;
    branch?: string;
    // Image-based deployment
    imageName?: string;
    imageRegistry?: string;
    imageTag?: string;
    // Dockerfile path for monorepo services
    dockerfilePath?: string;
    environmentVariables?: Record<string, string>;
    registryUsername?: string;
    registryPassword?: string;
  }[];
  requestId: string;
};

export type ProvisionServicesResponse = {
  serviceIds: string[];
};


// Environment Metadata Types
export type WizardInputs = {
  sourceType?: 'repository' | 'docker_image';
  repositoryUrl?: string;
  branch?: string;
  dockerImage?: string;
  imageRegistry?: string;
  imageTag?: string;
  discoveredServices?: Array<{
    name: string;
    path: string;
    dockerfilePath?: string;
  }>;
  environmentType?: string;
  ttl?: number;
  [key: string]: unknown;
};

export type ProvisionOutputs = {
  railwayProjectId?: string;
  railwayEnvironmentId?: string;
  serviceIds?: string[];
  [key: string]: unknown;
};

export type EnvironmentMetadata = {
  id: string;
  environmentId: string;
  wizardInputs?: WizardInputs;
  provisionOutputs?: ProvisionOutputs;
  isTemplate: boolean;
  templateName?: string;
  createdAt: string;
  updatedAt: string;
};

export type ServiceBuildConfig = {
  id: string;
  environmentId: string;
  name: string;
  status: string;
  deploymentType: 'repository' | 'docker_image';
  // Repository-based deployment
  sourceRepo?: string;
  sourceBranch?: string;
  dockerfilePath?: string;
  buildContext?: string;
  // Docker image-based deployment
  dockerImage?: string;
  imageRegistry?: string;
  imageTag?: string;
  // Build configuration
  buildArgs?: Record<string, string>;
  exposedPorts?: number[];
  // Railway IDs
  railwayServiceId?: string;
  // Timestamps
  createdAt: string;
  updatedAt: string;
};


// List environment templates
export type TemplateListItem = {
  id: string;
  templateName: string;
  environmentType?: string;
  serviceCount: number;
  createdAt: string;
};


// Environment Snapshot Types
export type ServiceVariablesSnapshot = {
  serviceId: string;
  serviceName: string;
  variables: Record<string, string>;
};

export type EnvironmentSnapshot = {
  environment: {
    id: string;
    name: string;
    type: string;
    status: string;
    sourceRepo: string;
    sourceBranch: string;
    sourceCommit: string;
    railwayProjectId: string;
    railwayEnvironmentId: string;
    ttlSeconds?: number;
    parentEnvironmentId?: string;
    createdAt: string;
    updatedAt: string;
  };
  services: Array<{
    id: string;
    environmentId: string;
    name: string;
    path: string;
    status: string;
    railwayServiceId: string;
    deploymentType: 'source_repo' | 'docker_image';
    sourceRepo: string;
    sourceBranch: string;
    dockerfilePath?: string;
    buildContext?: string;
    rootDirectory?: string;
    targetStage?: string;
    dockerImage: string;
    imageRegistry: string;
    imageName: string;
    imageTag: string;
    imageDigest: string;
    imageAuthStored: boolean;
    exposedPortsJson: string;
    healthCheckPath?: string;
    startCommand?: string;
    createdAt: string;
    updatedAt: string;
  }>;
  environmentVariables: Record<string, string>;
  serviceVariables: ServiceVariablesSnapshot[];
};

// Get environment snapshot for cloning
