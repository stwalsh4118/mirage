import { fetchJSON } from "@/lib/api";

export type EnvironmentType = "dev" | "staging" | "prod" | "ephemeral";
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

export function destroyEnvironment(id: string): Promise<void> {
  return fetchJSON<void>(`/environments/${id}`, { method: "DELETE" });
}






