"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { listRailwayProjectsByNames, RailwayProject, listRailwayProjectsDetails, RailwayProjectDetails, importRailwayEnvironments, ImportRailwayEnvsRequest, ImportRailwayEnvsResponse, provisionProject, ProvisionProjectRequest, ProvisionProjectResponse, provisionEnvironment, ProvisionEnvironmentRequest, ProvisionEnvironmentResponse, provisionServices, ProvisionServicesRequest, ProvisionServicesResponse, deleteRailwayEnvironment, deleteRailwayProject, deleteRailwayService, getEnvironmentMetadata, EnvironmentMetadata, getEnvironmentServices, ServiceBuildConfig, getServiceDetails, listTemplates, TemplateListItem, getEnvironmentSnapshot, EnvironmentSnapshot } from "@/lib/api/railway";

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

export function useDeleteRailwayService() {
  const qc = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: (railwayServiceId: string) => deleteRailwayService(railwayServiceId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["environment-services"] });
    },
  });
}

// Environment Metadata hooks
export function useEnvironmentMetadata(environmentId: string | null | undefined, enabled = true) {
  return useQuery<EnvironmentMetadata>({
    queryKey: ["environment-metadata", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return getEnvironmentMetadata(environmentId);
    },
    enabled: enabled && !!environmentId,
    retry: false, // Don't retry for 404s (missing metadata)
  });
}

export function useEnvironmentServices(environmentId: string | null | undefined, enabled = true) {
  return useQuery<ServiceBuildConfig[]>({
    queryKey: ["environment-services", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return getEnvironmentServices(environmentId);
    },
    enabled: enabled && !!environmentId,
    retry: false,
  });
}

export function useServiceDetails(serviceId: string | null | undefined, enabled = true) {
  return useQuery<ServiceBuildConfig>({
    queryKey: ["service-details", serviceId],
    queryFn: () => {
      if (!serviceId) throw new Error("Service ID is required");
      return getServiceDetails(serviceId);
    },
    enabled: enabled && !!serviceId,
    retry: false,
  });
}

export function useTemplates() {
  return useQuery<TemplateListItem[]>({
    queryKey: ["templates"],
    queryFn: () => listTemplates(),
  });
}

export function useEnvironmentSnapshot(environmentId: string | null) {
  return useQuery<EnvironmentSnapshot>({
    queryKey: ["environment-snapshot", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return getEnvironmentSnapshot(environmentId);
    },
    enabled: !!environmentId,
    retry: false, // Don't retry on errors (snapshot might not exist)
  });
}
