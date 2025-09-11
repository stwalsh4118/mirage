"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { listRailwayProjectsByNames, RailwayProject, listRailwayProjectsDetails, RailwayProjectDetails, importRailwayEnvironments, ImportRailwayEnvsRequest, ImportRailwayEnvsResponse } from "@/lib/api/railway";

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
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["environments"] });
    },
  });
}
