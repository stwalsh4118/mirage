"use client";

import { useQuery } from "@tanstack/react-query";
import { listRailwayProjectsByNames, RailwayProject, listRailwayProjectsDetails, RailwayProjectDetails } from "@/lib/api/railway";

export function useRailwayProjects(names: string[]) {
  return useQuery<RailwayProject[]>({
    queryKey: ["railway-projects", names.slice().sort().join(",")],
    queryFn: () => listRailwayProjectsByNames(names),
    enabled: names.length > 0,
    staleTime: 30_000,
    refetchOnWindowFocus: false,
  });
}

export function useRailwayProjectsDetails(names?: string[]) {
  return useQuery<RailwayProjectDetails[]>({
    queryKey: ["railway-projects-details", (names ?? []).slice().sort().join(",")],
    queryFn: () => listRailwayProjectsDetails(names),
    staleTime: 30_000,
    refetchOnWindowFocus: false,
  });
}
