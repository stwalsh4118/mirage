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
  return useQuery<RailwayProjectDetails[]>({
    queryKey: ["railway-projects-details", (names ?? []).slice().sort().join(",")],
    queryFn: () => listRailwayProjectsDetails(names),
    staleTime: RAILWAY_POLL_INTERVAL_MS,
    refetchInterval: RAILWAY_POLL_INTERVAL_MS,
    refetchOnWindowFocus: false,
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
