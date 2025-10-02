/**
 * Utility functions for generating unique service names in Railway
 * 
 * Railway requires unique service names within a project, even across environments.
 * This module provides functions to append environment suffixes to service names.
 */

/**
 * Common environment name abbreviations for cleaner service names
 */
const ENV_ABBREVIATIONS: Record<string, string> = {
  'production': 'prod',
  'staging': 'stg',
  'development': 'dev',
  'testing': 'test',
  'preview': 'prev',
  'qa': 'qa',
};

/**
 * Generates a unique service name by appending the environment name
 * 
 * Examples:
 * - generateUniqueServiceName("api", "production") -> "api-prod"
 * - generateUniqueServiceName("web", "staging") -> "web-stg"
 * - generateUniqueServiceName("worker", "feature-branch") -> "worker-feature-branch"
 * 
 * @param baseName - The base service name (e.g., "api", "web", "worker")
 * @param environmentName - The environment name (e.g., "production", "staging")
 * @returns The unique service name with environment suffix
 */
export function generateUniqueServiceName(baseName: string, environmentName: string): string {
  // Clean the base name (remove any existing environment suffixes)
  const cleanBaseName = baseName.trim();
  
  // Clean and normalize the environment name
  const cleanEnvName = environmentName.trim().toLowerCase();
  
  // Use abbreviation if available, otherwise use the full name
  const envSuffix = ENV_ABBREVIATIONS[cleanEnvName] || cleanEnvName;
  
  // Combine with hyphen separator
  return `${cleanBaseName}-${envSuffix}`;
}

/**
 * Generates unique service names for multiple services
 * 
 * @param services - Array of service base names
 * @param environmentName - The environment name
 * @returns Array of unique service names
 */
export function generateUniqueServiceNames(
  services: string[],
  environmentName: string
): string[] {
  return services.map(baseName => generateUniqueServiceName(baseName, environmentName));
}

