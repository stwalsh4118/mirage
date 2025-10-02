/**
 * Type definitions for Dockerfile discovery
 * Matches the API response from POST /api/v1/discovery/dockerfiles
 */

export type DiscoveredService = {
  name: string;
  dockerfilePath: string;
  buildContext: string;
  exposedPorts: number[];
  buildArgs: string[];
  baseImage: string;
};

export type DiscoveryResults = {
  services: DiscoveredService[];
  owner: string;
  repo: string;
  branch: string;
  cacheHit?: boolean;
};

export type DiscoveryRequest = {
  owner: string;
  repo: string;
  branch: string;
  userToken?: string;
  requestId?: string;
};

/**
 * Internal state for selected services with editable names
 */
export type SelectedService = DiscoveredService & {
  selected: boolean;
  editedName?: string;
};

