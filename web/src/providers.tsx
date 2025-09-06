"use client";

import { ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const defaultStaleTimeMs = 5 * 1000;

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: defaultStaleTimeMs,
			retry: 2,
			refetchOnWindowFocus: false,
		},
	},
});

export default function Providers({ children }: { children: ReactNode }) {
	return (
		<QueryClientProvider client={queryClient}>
			{children}
			<ReactQueryDevtools initialIsOpen={false} buttonPosition="bottom-left" />
		</QueryClientProvider>
	);
}
