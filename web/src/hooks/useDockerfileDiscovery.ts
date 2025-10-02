"use client";

import { useMutation } from "@tanstack/react-query";
import { discoverDockerfiles } from "@/lib/api/discovery";
import type { DiscoveryRequest, DiscoveryResults } from "@/components/discovery/types";

/**
 * Hook for triggering Dockerfile discovery in a repository
 * Uses mutation since it's a POST request that actively scans
 */
export function useDockerfileDiscovery() {
  return useMutation<DiscoveryResults, Error, DiscoveryRequest>({
    mutationFn: (request) => discoverDockerfiles(request),
  });
}

