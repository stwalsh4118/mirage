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
  const url = process.env.NEXT_PUBLIC_API_BASE_URL
    ? `${process.env.NEXT_PUBLIC_API_BASE_URL}${path}`
    : path;
  
  const headers = new Headers(init?.headers as HeadersInit);
  
  // Add Authorization header if token provided
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  
  // Add Content-Type for requests with body
  const hasBody = init?.body != null;
  if (hasBody && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  
  const res = await fetch(url, {
    ...init,
    headers,
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
    const text = await res.text().catch(() => "");
    throw new Error(`API ${res.status}: ${text}`);
  }
  
  return res.json() as Promise<T>;
}

