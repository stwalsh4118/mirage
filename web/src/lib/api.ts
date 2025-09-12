export const DEFAULT_API_BASE_URL = "http://localhost:8090";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || DEFAULT_API_BASE_URL;

export type HealthResponse = {
    status: string;
};

export async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
    const url = new URL(path, API_BASE_URL).toString();
    const res = await fetch(url, {
        ...init,
        headers: init?.headers,
        cache: "no-store",
    });
    if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(`Request failed ${res.status}: ${body}`);
    }
    return (await res.json()) as T;
}

export function getHealth(): Promise<HealthResponse> {
    return fetchJSON<HealthResponse>("/api/v1/healthz");
}


