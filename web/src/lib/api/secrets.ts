/**
 * API client for secrets management (Railway tokens, GitHub PATs, Docker credentials, etc.)
 * All endpoints require authentication via JWT token
 */

// Railway Token Types
export interface RailwayTokenStatusResponse {
  configured: boolean;
  last_validated?: string;
  needs_rotation: boolean;
  message?: string;
}

export interface StoreRailwayTokenRequest {
  token: string;
}

export interface StoreRailwayTokenResponse {
  success: boolean;
  validated: boolean;
  stored_at: string;
  message: string;
}

export interface ValidateRailwayTokenResponse {
  valid: boolean;
  message: string;
}

export interface DeleteRailwayTokenResponse {
  success: boolean;
  message: string;
}

export interface RotateRailwayTokenRequest {
  new_token: string;
}

export interface RotateRailwayTokenResponse {
  success: boolean;
  message: string;
}

