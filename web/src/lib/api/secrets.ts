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

// GitHub Token Types
export interface GitHubTokenStatusResponse {
  configured: boolean;
  username?: string;
  scopes?: string[];
  last_validated?: string;
  message?: string;
}

export interface StoreGitHubTokenRequest {
  token: string;
}

export interface StoreGitHubTokenResponse {
  success: boolean;
  validated: boolean;
  stored_at: string;
  username: string;
  scopes: string[];
  message: string;
}

export interface ValidateGitHubTokenResponse {
  valid: boolean;
  username?: string;
  scopes?: string[];
  message: string;
}

export interface DeleteGitHubTokenResponse {
  success: boolean;
  message: string;
}

