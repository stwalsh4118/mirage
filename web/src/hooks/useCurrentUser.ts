"use client";

import { useUser } from "@clerk/nextjs";
import { useQuery } from "@tanstack/react-query";
import { useAuthenticatedFetch } from "./useAuthenticatedFetch";

/**
 * Mirage user data from the backend database
 * Synchronized from Clerk via webhooks
 */
export interface MirageUser {
	id: string;
	clerkUserId: string;
	email: string;
	firstName?: string;
	lastName?: string;
	profileImageUrl?: string;
	lastSeenAt?: string;
	isActive: boolean;
	createdAt: string;
	updatedAt: string;
}

/**
 * Hook to fetch the current user's Mirage user data from the backend
 * This is the user record in the Mirage database, synced from Clerk
 * 
 * @returns React Query result with Mirage user data
 * 
 * @example
 * const { data: user, isLoading } = useCurrentUser();
 * if (user) {
 *   console.log(user.email);
 * }
 */
export function useCurrentUser() {
	const { isSignedIn } = useUser();
	const { fetch } = useAuthenticatedFetch();

	return useQuery({
		queryKey: ["user", "me"],
		queryFn: () => fetch<MirageUser>("/api/v1/users/me"),
		enabled: isSignedIn,
		staleTime: 5 * 60 * 1000, // 5 minutes
		retry: (failureCount, error) => {
			// Don't retry on 401 (will redirect) or 404 (user not synced yet)
			if (
				error instanceof Error &&
				(error.message.includes("401") || error.message.includes("404"))
			) {
				return false;
			}
			return failureCount < 3;
		},
	});
}

