"use client";

import { useMutation } from "@tanstack/react-query";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";
import type { DiscoveryRequest, DiscoveryResults } from "@/components/discovery/types";

/**
 * Hook for triggering Dockerfile discovery in a repository
 * Uses mutation since it's a POST request that actively scans
 * Now includes authentication via JWT token
 */
export function useDockerfileDiscovery() {
  const { fetch } = useAuthenticatedFetch();
  
  return useMutation<DiscoveryResults, Error, DiscoveryRequest>({
    mutationFn: (request) => 
      fetch<DiscoveryResults>('/api/v1/discovery/dockerfiles', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      }),
  });
}

