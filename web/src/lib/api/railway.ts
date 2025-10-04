import { fetchJSON } from '@/lib/api';

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

export function listRailwayProjectsByNames(names: unknown): Promise<RailwayProject[]> {
  if (names == null) names = [] as string[];
  if (!Array.isArray(names)) {
    throw new TypeError('listRailwayProjectsByNames: names must be an array of strings');
  }
  const qs = new URLSearchParams();
  if ((names as string[]).length) {
    qs.set('names', (names as string[]).join(','));
  }
  const suffix = qs.toString();
  const path = suffix ? `/railway/projects?${suffix}` : `/railway/projects`;
  return fetchJSON<RailwayProject[]>(`/api/v1${path}`);
}

export function listRailwayProjectsDetails(names?: unknown): Promise<RailwayProjectDetails[]> {
  if (names == null) names = [] as string[];
  if (!Array.isArray(names)) {
    throw new TypeError('listRailwayProjectsDetails: names must be an array of strings');
  }
  const qs = new URLSearchParams();
  qs.set('details', '1');
  if ((names as string[]).length) {
    qs.set('names', (names as string[]).join(','));
  }
  return fetchJSON<RailwayProjectDetails[]>(`/api/v1/railway/projects?${qs.toString()}`);
}

export type ImportRailwayEnvsRequest = {
  projectId: string;
  environmentIds: string[];
};

export type ImportRailwayEnvsResponse = {
  imported: number;
  skipped: number;
  items: { id: string; name: string; type: string; status: string; createdAt: string }[];
};

export function importRailwayEnvironments(body: ImportRailwayEnvsRequest): Promise<ImportRailwayEnvsResponse> {
  return fetchJSON<ImportRailwayEnvsResponse>(`/api/v1/railway/import/environments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

// Provision: Create Project
export type ProvisionProjectRequest = {
  defaultEnvironmentName?: string;
  name?: string;
  requestId: string;
};

export type ProvisionProjectResponse = {
  projectId: string;
  baseEnvironmentId: string;
  name: string;
};

export function provisionProject(body: ProvisionProjectRequest): Promise<ProvisionProjectResponse> {
  return fetchJSON<ProvisionProjectResponse>(`/api/v1/provision/project`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

// Provision: Create Environment
export type ProvisionEnvironmentRequest = {
  projectId: string;
  name: string;
  requestId: string;
  envType?: 'dev' | 'staging' | 'prod' | 'ephemeral';
  wizardInputs?: Record<string, unknown>;
};

export type ProvisionEnvironmentResponse = {
  environmentId: string;
};

export function provisionEnvironment(body: ProvisionEnvironmentRequest): Promise<ProvisionEnvironmentResponse> {
  return fetchJSON<ProvisionEnvironmentResponse>(`/api/v1/provision/environment`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

// Provision: Create Services
export type ProvisionServicesRequest = {
  projectId: string;
  environmentId: string;
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

export function provisionServices(body: ProvisionServicesRequest): Promise<ProvisionServicesResponse> {
  return fetchJSON<ProvisionServicesResponse>(`/api/v1/provision/services`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

// Delete Railway environment by Railway environment ID
export function deleteRailwayEnvironment(railwayEnvironmentId: string): Promise<void> {
  return fetchJSON<void>(`/api/v1/railway/environment/${railwayEnvironmentId}`, {
    method: 'DELETE',
  });
}

// Delete Railway project by Railway project ID
// WARNING: This is a destructive operation that deletes the project and all its resources
export function deleteRailwayProject(projectId: string): Promise<void> {
  return fetchJSON<void>(`/api/v1/railway/project/${projectId}`, {
    method: 'DELETE',
  });
}

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

// Get environment metadata
export function getEnvironmentMetadata(environmentId: string): Promise<EnvironmentMetadata> {
  return fetchJSON<EnvironmentMetadata>(`/api/v1/environments/${environmentId}/metadata`);
}

// Get services for an environment
export function getEnvironmentServices(environmentId: string): Promise<ServiceBuildConfig[]> {
  return fetchJSON<ServiceBuildConfig[]>(`/api/v1/environments/${environmentId}/services`);
}

// Get single service details
export function getServiceDetails(serviceId: string): Promise<ServiceBuildConfig> {
  return fetchJSON<ServiceBuildConfig>(`/api/v1/services/${serviceId}`);
}

// List environment templates
export type TemplateListItem = {
  id: string;
  templateName: string;
  environmentType?: string;
  serviceCount: number;
  createdAt: string;
};

export function listTemplates(): Promise<TemplateListItem[]> {
  return fetchJSON<TemplateListItem[]>(`/api/v1/templates`);
}
