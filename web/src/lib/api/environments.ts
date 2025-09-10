import { fetchJSON } from "@/lib/api";

export type EnvironmentType = "dev" | "prod";
export type EnvironmentStatus = "creating" | "active" | "error" | "destroying" | "unknown";

export interface Environment {
  id: string;
  name: string;
  type: EnvironmentType;
  status: EnvironmentStatus;
  url?: string;
  createdAt: string;
  description?: string;
}

export function listEnvironments(): Promise<Environment[]> {
  return fetchJSON<Environment[]>("/environments");
}

export function destroyEnvironment(id: string): Promise<{ ok: boolean }>{
  return fetchJSON<{ ok: boolean }>(`/environments/${id}`, { method: "DELETE" });
}





