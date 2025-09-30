export async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const url = process.env.NEXT_PUBLIC_API_BASE_URL ? `${process.env.NEXT_PUBLIC_API_BASE_URL}${path}` : path;
  const merged = init ? { ...init } : {};
  const hdrs = new Headers(init?.headers as HeadersInit);
  const hasBody = merged.body != null;
  if (hasBody && !hdrs.has('Content-Type')) {
    hdrs.set('Content-Type', 'application/json');
  }
  merged.headers = hdrs;
  const res = await fetch(url, merged);
  if (res.status === 204) {
    return { } as T;
  }
  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`API ${res.status}: ${text}`);
  }
  return res.json() as Promise<T>;
}

