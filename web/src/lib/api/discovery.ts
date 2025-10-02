import { fetchJSON } from '@/lib/api';
import type { DiscoveryRequest, DiscoveryResults } from '@/components/discovery/types';

/**
 * Discover Dockerfiles in a GitHub repository
 * POST /api/v1/discovery/dockerfiles
 */
export function discoverDockerfiles(body: DiscoveryRequest): Promise<DiscoveryResults> {
  return fetchJSON<DiscoveryResults>(`/api/v1/discovery/dockerfiles`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

