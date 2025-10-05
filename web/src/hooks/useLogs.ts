"use client";

import { useQuery } from "@tanstack/react-query";
import { getServiceLogs, type GetServiceLogsParams, type ServiceLogsResponse } from "@/lib/api/logs";

export function useServiceLogs(
  params: GetServiceLogsParams,
  enabled = true
) {
  return useQuery<ServiceLogsResponse>({
    queryKey: ["service-logs", params.serviceId, params.limit, params.search, params.minSeverity],
    queryFn: () => getServiceLogs(params),
    enabled: enabled && !!params.serviceId,
    staleTime: 10_000, // Consider logs stale after 10 seconds
    retry: false,
  });
}

