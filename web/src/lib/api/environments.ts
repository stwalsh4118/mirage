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


