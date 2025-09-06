"use client";

import { useQuery } from "@tanstack/react-query";
import { getHealth } from "@/lib/api";

export default function HealthPage() {
    const { data, isLoading, isError, error, refetch, isFetching } = useQuery({
        queryKey: ["health"],
        queryFn: getHealth,
    });

    return (
        <div className="space-y-4">
            <h1 className="text-xl font-semibold">Backend Health</h1>
            <div className="flex items-center gap-3">
                <button
                    onClick={() => refetch()}
                    className="px-3 py-1.5 rounded border border-black/10 dark:border-white/15 text-sm hover:bg-black/[.03] dark:hover:bg-white/[.06]"
                >
                    Refresh
                </button>
                {isFetching && <span className="text-xs text-black/60 dark:text-white/60">Fetching…</span>}
            </div>
            {isLoading && <p className="text-black/70 dark:text-white/70">Loading…</p>}
            {isError && (
                <p className="text-red-600">{(error as Error)?.message ?? "Failed to load health"}</p>
            )}
            {data && (
                <pre className="rounded bg-black/[.03] dark:bg-white/[.06] p-3 text-sm">
{JSON.stringify(data, null, 2)}
                </pre>
            )}
        </div>
    );
}


