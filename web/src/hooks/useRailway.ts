"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { listRailwayProjectsByNames, RailwayProject, listRailwayProjectsDetails, RailwayProjectDetails, importRailwayEnvironments, ImportRailwayEnvsRequest, ImportRailwayEnvsResponse, provisionProject, ProvisionProjectRequest, ProvisionProjectResponse, provisionEnvironment, ProvisionEnvironmentRequest, ProvisionEnvironmentResponse, provisionServices, ProvisionServicesRequest, ProvisionServicesResponse, deleteRailwayEnvironment, deleteRailwayProject } from "@/lib/api/railway";

const RAILWAY_POLL_INTERVAL_MS = 30_000;

export function useRailwayProjects(names: string[]) {
  return useQuery<RailwayProject[]>({
    queryKey: ["railway-projects", names.slice().sort().join(",")],
    queryFn: () => listRailwayProjectsByNames(names),
    enabled: names.length > 0,
    staleTime: RAILWAY_POLL_INTERVAL_MS,
    refetchInterval: RAILWAY_POLL_INTERVAL_MS,
    refetchOnWindowFocus: false,
  });
}

export function useRailwayProjectsDetails(names?: string[]) {
  const includeDemo = typeof window !== "undefined" && new URLSearchParams(window.location.search).has("demo");
  return useQuery<RailwayProjectDetails[]>({
    queryKey: ["railway-projects-details", (names ?? []).slice().sort().join(","), includeDemo ? "demo" : "nodemo"],
    queryFn: () => listRailwayProjectsDetails(names),
    staleTime: RAILWAY_POLL_INTERVAL_MS,
    refetchInterval: RAILWAY_POLL_INTERVAL_MS,
    refetchOnWindowFocus: false,
    select: (data) => {
      if (!includeDemo) return data;
      const demo: RailwayProjectDetails = {
        id: "starlight-orchestra",
        name: "Starlight Orchestra",
        services: [
          { id: "api-gateway", name: "api-gateway" },
          { id: "web-app", name: "web-app" },
          { id: "worker-bus", name: "worker-bus" },
          { id: "postgres", name: "postgres" },
          { id: "redis", name: "redis" },
          { id: "vector-db", name: "vector-db" },
          { id: "image-cdn", name: "image-cdn" },
        ],
        plugins: [
          { id: "grafana", name: "grafana" },
          { id: "sentry", name: "sentry" },
          { id: "datadog", name: "datadog" },
          { id: "stripe", name: "stripe" },
        ],
        environments: [
          {
            id: "production",
            name: "production",
            services: [
              { id: "web-app", name: "web-app" },
              { id: "api-gateway", name: "api-gateway" },
              { id: "postgres", name: "postgres" },
              { id: "redis", name: "redis" },
              { id: "image-cdn", name: "image-cdn" },
            ],
          },
          {
            id: "staging",
            name: "staging",
            services: [
              { id: "web-app", name: "web-app" },
              { id: "api-gateway", name: "api-gateway" },
              { id: "postgres", name: "postgres" },
              { id: "vector-db", name: "vector-db" },
            ],
          },
          {
            id: "preview",
            name: "preview",
            services: [
              { id: "web-app", name: "web-app" },
              { id: "worker-bus", name: "worker-bus" },
            ],
          },
        ],
      };
      // Put demo first for visibility
      const existing = data ?? [];
      const withoutDup = existing.filter((p) => p.id !== demo.id);
      return [demo, ...withoutDup];
    },
  });
}

export function useImportRailwayEnvironments() {
  const qc = useQueryClient();
  return useMutation<ImportRailwayEnvsResponse, Error, ImportRailwayEnvsRequest>({
    mutationFn: (body) => importRailwayEnvironments(body),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["environments"] });
    },
  });
}

export function useProvisionProject() {
  const qc = useQueryClient();
  return useMutation<ProvisionProjectResponse, Error, ProvisionProjectRequest>({
    mutationFn: (body) => provisionProject(body),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["railway-projects"] });
    },
  });
}

export function useProvisionEnvironment() {
  const qc = useQueryClient();
  return useMutation<ProvisionEnvironmentResponse, Error, ProvisionEnvironmentRequest>({
    mutationFn: (body) => provisionEnvironment(body),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useProvisionServices() {
  const qc = useQueryClient();
  return useMutation<ProvisionServicesResponse, Error, ProvisionServicesRequest>({
    mutationFn: (body) => provisionServices(body),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useDeleteRailwayEnvironment() {
  const qc = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: (railwayEnvironmentId: string) => deleteRailwayEnvironment(railwayEnvironmentId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useDeleteRailwayProject() {
  const qc = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: (projectId: string) => deleteRailwayProject(projectId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["railway-projects"] });
    },
  });
}
