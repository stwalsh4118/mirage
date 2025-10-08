"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	RailwayProject,
	RailwayProjectDetails,
	ImportRailwayEnvsRequest,
	ImportRailwayEnvsResponse,
	ProvisionProjectRequest,
	ProvisionProjectResponse,
	ProvisionEnvironmentRequest,
	ProvisionEnvironmentResponse,
	ProvisionServicesRequest,
	ProvisionServicesResponse,
	EnvironmentMetadata,
	ServiceBuildConfig,
	TemplateListItem,
	EnvironmentSnapshot,
} from "@/lib/api/railway";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";

const RAILWAY_POLL_INTERVAL_MS = 30_000;

export function useRailwayProjects(names: string[]) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<RailwayProject[]>({
    queryKey: ["railway-projects", names.slice().sort().join(",")],
    queryFn: async () => {
      if (names == null) names = [] as string[];
      if (!Array.isArray(names)) {
        throw new TypeError('names must be an array of strings');
      }
      const qs = new URLSearchParams();
      if (names.length) {
        qs.set('names', names.join(','));
      }
      const suffix = qs.toString();
      const path = suffix ? `/railway/projects?${suffix}` : `/railway/projects`;
      return fetch<RailwayProject[]>(`/api/v1${path}`);
    },
    enabled: names.length > 0,
    staleTime: RAILWAY_POLL_INTERVAL_MS,
    refetchInterval: RAILWAY_POLL_INTERVAL_MS,
    refetchOnWindowFocus: false,
  });
}

export function useRailwayProjectsDetails(names?: string[]) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<RailwayProjectDetails[]>({
    queryKey: ["railway-projects-details", (names ?? []).slice().sort().join(",")],
    queryFn: async () => {
      const nameList = names ?? [];
      const qs = new URLSearchParams();
      qs.set('details', '1');
      if (nameList.length) {
        qs.set('names', nameList.join(','));
      }
      return fetch<RailwayProjectDetails[]>(`/api/v1/railway/projects?${qs.toString()}`);
    },
    staleTime: RAILWAY_POLL_INTERVAL_MS,
    refetchInterval: RAILWAY_POLL_INTERVAL_MS,
    refetchOnWindowFocus: false,
  });
}

export function useImportRailwayEnvironments() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<ImportRailwayEnvsResponse, Error, ImportRailwayEnvsRequest>({
    mutationFn: async (body) => {
      return fetch<ImportRailwayEnvsResponse>(`/api/v1/railway/import/environments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["environments"] });
    },
  });
}

export function useProvisionProject() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<ProvisionProjectResponse, Error, ProvisionProjectRequest>({
    mutationFn: async (body) => {
      return fetch<ProvisionProjectResponse>(`/api/v1/provision/project`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["railway-projects"] });
    },
  });
}

export function useProvisionEnvironment() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<ProvisionEnvironmentResponse, Error, ProvisionEnvironmentRequest>({
    mutationFn: async (body) => {
      return fetch<ProvisionEnvironmentResponse>(`/api/v1/provision/environment`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useProvisionServices() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<ProvisionServicesResponse, Error, ProvisionServicesRequest>({
    mutationFn: async (body) => {
      return fetch<ProvisionServicesResponse>(`/api/v1/provision/services`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useDeleteRailwayEnvironment() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<void, Error, string>({
    mutationFn: async (railwayEnvironmentId: string) => {
      return fetch<void>(`/api/v1/railway/environment/${railwayEnvironmentId}`, {
        method: 'DELETE',
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
    },
  });
}

export function useDeleteRailwayProject() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<void, Error, string>({
    mutationFn: async (projectId: string) => {
      return fetch<void>(`/api/v1/railway/project/${projectId}`, {
        method: 'DELETE',
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["railway-projects"] });
    },
  });
}

export function useDeleteRailwayService() {
  const qc = useQueryClient();
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<void, Error, string>({
    mutationFn: async (railwayServiceId: string) => {
      return fetch<void>(`/api/v1/railway/service/${railwayServiceId}`, {
        method: 'DELETE',
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["railway-projects-details"] });
      await qc.invalidateQueries({ queryKey: ["environment-services"] });
    },
  });
}

// Environment Metadata hooks
export function useEnvironmentMetadata(environmentId: string | null | undefined, enabled = true) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<EnvironmentMetadata>({
    queryKey: ["environment-metadata", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return fetch<EnvironmentMetadata>(`/api/v1/environments/${environmentId}/metadata`);
    },
    enabled: enabled && !!environmentId,
    retry: false, // Don't retry for 404s (missing metadata)
  });
}

export function useEnvironmentServices(environmentId: string | null | undefined, enabled = true) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<ServiceBuildConfig[]>({
    queryKey: ["environment-services", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return fetch<ServiceBuildConfig[]>(`/api/v1/environments/${environmentId}/services`);
    },
    enabled: enabled && !!environmentId,
    retry: false,
  });
}

export function useServiceDetails(serviceId: string | null | undefined, enabled = true) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<ServiceBuildConfig>({
    queryKey: ["service-details", serviceId],
    queryFn: () => {
      if (!serviceId) throw new Error("Service ID is required");
      return fetch<ServiceBuildConfig>(`/api/v1/services/${serviceId}`);
    },
    enabled: enabled && !!serviceId,
    retry: false,
  });
}

export function useTemplates() {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<TemplateListItem[]>({
    queryKey: ["templates"],
    queryFn: () => fetch<TemplateListItem[]>(`/api/v1/templates`),
  });
}

export function useEnvironmentSnapshot(environmentId: string | null) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<EnvironmentSnapshot>({
    queryKey: ["environment-snapshot", environmentId],
    queryFn: () => {
      if (!environmentId) throw new Error("Environment ID is required");
      return fetch<EnvironmentSnapshot>(`/api/v1/environments/${environmentId}/snapshot`);
    },
    enabled: !!environmentId,
    retry: false, // Don't retry on errors (snapshot might not exist)
  });
}
