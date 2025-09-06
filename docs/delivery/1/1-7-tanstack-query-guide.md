# 1-7 TanStack Query Guide

Date: 2025-09-06

[Official Docs](https://tanstack.com/query/latest/docs/framework/react/overview)

## Overview
TanStack Query manages server-state with caching, deduping, retries, background refetching, and devtools. In this project it powers API data fetching for the Next.js App Router UI.

## Installation
Already installed in `web`:

```bash
pnpm add @tanstack/react-query @tanstack/react-query-devtools
```

## Setup
Provider is centralized in `src/providers.tsx` and attached in `src/app/layout.tsx`.

```tsx
// src/providers.tsx
"use client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 5_000, retry: 2, refetchOnWindowFocus: false } },
});

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} buttonPosition="bottom-left" />
    </QueryClientProvider>
  );
}
```

Attach in the root layout:

```tsx
// src/app/layout.tsx
import Providers from "@/providers";
...
<Providers>{children}</Providers>
```

## Usage Pattern
Define API helpers in `src/lib/api.ts` and consume with `useQuery` in client components.

```tsx
import { useQuery } from "@tanstack/react-query";
import { getHealth } from "@/lib/api";

export function Health() {
  const { data, isLoading, isError, error } = useQuery({ queryKey: ["health"], queryFn: getHealth });
  if (isLoading) return <p>Loading…</p>;
  if (isError) return <p>Error: {(error as Error).message}</p>;
  return <pre>{JSON.stringify(data, null, 2)}</pre>;
}
```

## Recommendations
- Use descriptive `queryKey` arrays: `["environments", envId]`.
- Keep fetchers side-effect free and throw on non-200 responses.
- Prefer `enabled: false` and `refetch()` for imperative refresh buttons.
- Co-locate queries with components; move to hooks for reuse.


