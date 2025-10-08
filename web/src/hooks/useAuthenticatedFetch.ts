"use client";

import { useAuth } from "@clerk/nextjs";
import { useCallback } from "react";
import { fetchJSON, fetchBlob } from "@/lib/api";

/**
 * Hook to make authenticated API requests with automatic JWT token injection
 * Uses Clerk's getToken() to retrieve the current session token
 * 
 * @returns Object with fetch and fetchBlob methods
 * 
 * @example
 * const { fetch } = useAuthenticatedFetch();
 * const data = await fetch<MyType>('/api/v1/users/me');
 */
export function useAuthenticatedFetch() {
	const { getToken } = useAuth();

	const authenticatedFetch = useCallback(
		async <T>(path: string, init?: RequestInit): Promise<T> => {
			// Get JWT token from Clerk (automatically refreshed if needed)
			const token = await getToken();

			// Call fetchJSON with token
			return fetchJSON<T>(path, init, token);
		},
		[getToken]
	);

	const authenticatedFetchBlob = useCallback(
		async (path: string, init?: RequestInit): Promise<Blob> => {
			// Get JWT token from Clerk
			const token = await getToken();

			// Call fetchBlob with token
			return fetchBlob(path, init, token);
		},
		[getToken]
	);

	return {
		fetch: authenticatedFetch,
		fetchBlob: authenticatedFetchBlob,
	};
}

