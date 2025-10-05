import { fetchJSON, fetchBlob } from '@/lib/api';

export type LogEntry = {
  timestamp: string;
  serviceName: string;
  severity: string;
  message: string;
  rawLine: string;
};

// Alias for component compatibility
export type { LogEntry as Log };

export type ServiceLogsResponse = {
  logs: LogEntry[];
  count: number;
};

export type GetServiceLogsParams = {
  serviceId: string;
  limit?: number;
  search?: string;
  minSeverity?: string;
};

export async function getServiceLogs({
  serviceId,
  limit = 100,
  search,
  minSeverity,
}: GetServiceLogsParams): Promise<ServiceLogsResponse> {
  const params = new URLSearchParams({ limit: limit.toString() });
  if (search) params.append('search', search);
  if (minSeverity) params.append('minSeverity', minSeverity);

  return fetchJSON(`/api/v1/services/${serviceId}/logs?${params.toString()}`);
}

export type ExportFormat = 'json' | 'csv' | 'txt';

export async function exportServiceLogs(
  serviceId: string,
  format: ExportFormat,
  limit = 1000
): Promise<Blob> {
  const params = new URLSearchParams({
    serviceId,
    format,
    limit: limit.toString(),
  });

  return fetchBlob(`/api/v1/logs/export?${params.toString()}`);
}

