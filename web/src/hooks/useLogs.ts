"use client";

import { useQuery } from "@tanstack/react-query";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";
import type { GetServiceLogsParams, ServiceLogsResponse } from "@/lib/api/logs";

export function useServiceLogs(
  params: GetServiceLogsParams,
  enabled = true
) {
  const { fetch } = useAuthenticatedFetch();
  
  return useQuery<ServiceLogsResponse>({
    queryKey: ["service-logs", params.serviceId, params.limit, params.search, params.minSeverity],
    queryFn: async () => {
      const queryParams = new URLSearchParams({ limit: (params.limit || 100).toString() });
      if (params.search) queryParams.append('search', params.search);
      if (params.minSeverity) queryParams.append('minSeverity', params.minSeverity);
      
      return fetch<ServiceLogsResponse>(`/api/v1/services/${params.serviceId}/logs?${queryParams.toString()}`);
    },
    enabled: enabled && !!params.serviceId,
    staleTime: 10_000, // Consider logs stale after 10 seconds
    retry: false,
  });
}

