import { fetchJSON } from '@/lib/api';

export type RailwayProject = { id: string; name: string };
export type RailwayProjectItem = { id: string; name: string };
export type RailwayProjectDetails = {
  id: string;
  name: string;
  services: RailwayProjectItem[];
  plugins: RailwayProjectItem[];
  environments: RailwayProjectItem[];
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
  return fetchJSON<RailwayProject[]>(path);
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
  return fetchJSON<RailwayProjectDetails[]>(`/railway/projects?${qs.toString()}`);
}
