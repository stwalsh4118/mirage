"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";
import type {
  RailwayTokenStatusResponse,
  StoreRailwayTokenRequest,
  StoreRailwayTokenResponse,
  ValidateRailwayTokenResponse,
  DeleteRailwayTokenResponse,
  RotateRailwayTokenRequest,
  RotateRailwayTokenResponse,
} from "@/lib/api/secrets";

/**
 * Hook to get Railway token status (configured, last validated, etc.)
 */
export function useRailwayTokenStatus() {
  const { fetch } = useAuthenticatedFetch();

  return useQuery<RailwayTokenStatusResponse>({
    queryKey: ["railway-token", "status"],
    queryFn: () => fetch<RailwayTokenStatusResponse>("/api/v1/secrets/railway/status"),
    staleTime: 30_000, // 30 seconds
    retry: 1,
  });
}

/**
 * Hook to store/update Railway token
 */
export function useStoreRailwayToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<StoreRailwayTokenResponse, Error, StoreRailwayTokenRequest>({
    mutationFn: (request) =>
      fetch<StoreRailwayTokenResponse>("/api/v1/secrets/railway", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(request),
      }),
    onSuccess: () => {
      // Invalidate status query to refresh the UI
      queryClient.invalidateQueries({ queryKey: ["railway-token", "status"] });
    },
  });
}

/**
 * Hook to validate Railway token
 */
export function useValidateRailwayToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<ValidateRailwayTokenResponse, Error, void>({
    mutationFn: () =>
      fetch<ValidateRailwayTokenResponse>("/api/v1/secrets/railway/validate", {
        method: "POST",
      }),
    onSuccess: () => {
      // Invalidate status query to update last_validated timestamp
      queryClient.invalidateQueries({ queryKey: ["railway-token", "status"] });
    },
  });
}

/**
 * Hook to delete Railway token
 */
export function useDeleteRailwayToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<DeleteRailwayTokenResponse, Error, void>({
    mutationFn: () =>
      fetch<DeleteRailwayTokenResponse>("/api/v1/secrets/railway", {
        method: "DELETE",
      }),
    onSuccess: () => {
      // Invalidate status query to show unconfigured state
      queryClient.invalidateQueries({ queryKey: ["railway-token", "status"] });
    },
  });
}

/**
 * Hook to rotate Railway token
 */
export function useRotateRailwayToken() {
  const { fetch } = useAuthenticatedFetch();
  const queryClient = useQueryClient();

  return useMutation<RotateRailwayTokenResponse, Error, RotateRailwayTokenRequest>({
    mutationFn: (request) =>
      fetch<RotateRailwayTokenResponse>("/api/v1/secrets/railway/rotate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(request),
      }),
    onSuccess: () => {
      // Invalidate status query to refresh the UI
      queryClient.invalidateQueries({ queryKey: ["railway-token", "status"] });
    },
  });
}

