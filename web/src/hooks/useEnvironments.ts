"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { destroyEnvironment, Environment, listEnvironments } from "@/lib/api/environments";

export const ENV_POLL_INTERVAL_MS = 5000 as const;

export function useEnvironments() {
  return useQuery<Environment[]>({
    queryKey: ["environments"],
    queryFn: () => listEnvironments(),
    staleTime: 0,
    retry: 2,
    refetchOnWindowFocus: false,
    refetchInterval: ENV_POLL_INTERVAL_MS,
  });
}

export function useDestroyEnvironment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => destroyEnvironment(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["environments"] });
    },
  });
}

// Temporary mocked data to visualize the dashboard while backend wiring is in progress
export const MOCK_ENVIRONMENTS: Environment[] = [
  {
    id: "env-1",
    name: "analytics-db",
    type: "dev",
    status: "error",
    url: "https://analytics-dev.mirage.app",
    createdAt: new Date(Date.now() - 3 * 24 * 3600 * 1000).toISOString(),
    description: undefined,
  },
  {
    id: "env-2",
    name: "api-service",
    type: "dev",
    status: "active",
    url: "https://api-dev.mirage.app",
    createdAt: new Date(Date.now() - 2 * 3600 * 1000).toISOString(),
    description: undefined,
  },
  {
    id: "env-3",
    name: "auth-service",
    type: "dev",
    status: "unknown",
    url: undefined,
    createdAt: new Date(Date.now() - 7 * 24 * 3600 * 1000).toISOString(),
    description: "Not deployed",
  },
  {
    id: "env-4",
    name: "frontend-app",
    type: "prod",
    status: "active",
    url: "https://app.mirage.com",
    createdAt: new Date(Date.now() - 24 * 3600 * 1000).toISOString(),
    description: undefined,
  },
  {
    id: "env-5",
    name: "worker-queue",
    type: "dev",
    status: "creating",
    url: undefined,
    createdAt: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    description: undefined,
  },
];



