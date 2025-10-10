"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";
import type {
  GitHubTokenStatusResponse,
  StoreGitHubTokenRequest,
  StoreGitHubTokenResponse,
  ValidateGitHubTokenResponse,
  DeleteGitHubTokenResponse,
} from "@/lib/api/secrets";

/**
 * Hook to get GitHub token status (configured, username, scopes, etc.)
 */
export function useGitHubTokenStatus() {
  const { fetch } = useAuthenticatedFetch();

  return useQuery<GitHubTokenStatusResponse>({
    queryKey: ["github-token", "status"],
    queryFn: () => fetch<GitHubTokenStatusResponse>("/api/v1/secrets/github/status"),
    staleTime: 30_000, // 30 seconds
    retry: 1,
  });
}

/**
 * Hook to store/update GitHub token
 */
export function useStoreGitHubToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<StoreGitHubTokenResponse, Error, StoreGitHubTokenRequest>({
    mutationFn: (request) =>
      fetch<StoreGitHubTokenResponse>("/api/v1/secrets/github", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(request),
      }),
    onSuccess: () => {
      // Invalidate status query to refresh the UI
      queryClient.invalidateQueries({ queryKey: ["github-token", "status"] });
    },
  });
}

/**
 * Hook to validate GitHub token
 */
export function useValidateGitHubToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<ValidateGitHubTokenResponse, Error, void>({
    mutationFn: () =>
      fetch<ValidateGitHubTokenResponse>("/api/v1/secrets/github/validate", {
        method: "POST",
      }),
    onSuccess: () => {
      // Invalidate status query to update last_validated timestamp
      queryClient.invalidateQueries({ queryKey: ["github-token", "status"] });
    },
  });
}

/**
 * Hook to delete GitHub token
 */
export function useDeleteGitHubToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<DeleteGitHubTokenResponse, Error, void>({
    mutationFn: () =>
      fetch<DeleteGitHubTokenResponse>("/api/v1/secrets/github", {
        method: "DELETE",
      }),
    onSuccess: () => {
      // Invalidate status query to show unconfigured state
      queryClient.invalidateQueries({ queryKey: ["github-token", "status"] });
    },
  });
}

