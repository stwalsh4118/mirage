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

export function listRailwayProjectsByNames(names: string[]): Promise<RailwayProject[]> {
  const qs = new URLSearchParams();
  if (names.length) {
    qs.set('names', names.join(','));
  }
  const suffix = qs.toString();
  const path = suffix ? `/railway/projects?${suffix}` : `/railway/projects`;
  return fetchJSON<RailwayProject[]>(path);
}

export function listRailwayProjectsDetails(names?: string[]): Promise<RailwayProjectDetails[]> {
  const qs = new URLSearchParams();
  qs.set('details', '1');
  if (names && names.length) {
    qs.set('names', names.join(','));
  }
  return fetchJSON<RailwayProjectDetails[]>(`/railway/projects?${qs.toString()}`);
}
