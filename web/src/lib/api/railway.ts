import { fetchJSON } from '@/lib/api';

export type RailwayProject = { id: string; name: string };
export type RailwayProjectItem = { id: string; name: string };
export type RailwayEnvironmentWithServices = {
  id: string;
  name: string;
  services: RailwayProjectItem[];
};

export type RailwayProjectDetails = {
  id: string;
  name: string;
  services: RailwayProjectItem[];
  plugins: RailwayProjectItem[];
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
