export const DEFAULT_API_BASE_URL = "http://localhost:8090";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || DEFAULT_API_BASE_URL;

export type HealthResponse = {
    status: string;
};

/**
 * Fetch JSON from API with optional authentication token
 * @param path - API path (e.g., /api/v1/users/me)
 * @param init - Fetch options
 * @param token - Optional JWT token for authentication
 * @returns Promise<T>
 */
export async function fetchJSON<T>(
    path: string,
    init?: RequestInit,
    token?: string | null
): Promise<T> {
    const url = new URL(path, API_BASE_URL).toString();
    
    const headers = new Headers(init?.headers);
    
    // Add Authorization header if token provided
    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }
    
    // Add Content-Type for requests with body
    if (init?.body && !headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
    }
    
    const res = await fetch(url, {
        ...init,
        headers,
        cache: "no-store",
    });
    
    // Handle 204 No Content
    if (res.status === 204) {
        return {} as T;
    }
    
    // Handle 401 Unauthorized - redirect to sign-in
    if (res.status === 401) {
        if (typeof window !== "undefined") {
            const currentUrl = window.location.pathname + window.location.search;
            const signInUrl = `/sign-in?redirect_url=${encodeURIComponent(currentUrl)}`;
            window.location.href = signInUrl;
        }
        throw new Error("Authentication required");
    }
    
    // Handle other errors
    if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(`Request failed ${res.status}: ${body}`);
    }
    
    return (await res.json()) as T;
}

/**
 * Fetch Blob from API with optional authentication token
 * @param path - API path
 * @param init - Fetch options
 * @param token - Optional JWT token for authentication
 * @returns Promise<Blob>
 */
export async function fetchBlob(
    path: string,
    init?: RequestInit,
    token?: string | null
): Promise<Blob> {
    const url = new URL(path, API_BASE_URL).toString();
    
    const headers = new Headers(init?.headers);
    
    // Add Authorization header if token provided
    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }
    
    const res = await fetch(url, {
        ...init,
        headers,
        cache: "no-store",
    });
    
    // Handle 401 Unauthorized - redirect to sign-in
    if (res.status === 401) {
        if (typeof window !== "undefined") {
            const currentUrl = window.location.pathname + window.location.search;
            const signInUrl = `/sign-in?redirect_url=${encodeURIComponent(currentUrl)}`;
            window.location.href = signInUrl;
        }
        throw new Error("Authentication required");
    }
    
    // Handle other errors
    if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(`Request failed ${res.status}: ${body}`);
    }
    
    return await res.blob();
}

export function getHealth(): Promise<HealthResponse> {
    return fetchJSON<HealthResponse>("/api/v1/healthz");
}


