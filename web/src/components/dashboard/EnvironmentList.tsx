"use client";

import { useEnvironments } from "@/hooks/useEnvironments";
import { useDashboardStore } from "@/store/dashboard";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { StatusChip } from "@/components/dashboard/StatusChip";
import Link from "next/link";

export function EnvironmentList() {
  const { data = [] } = useEnvironments();
  const { query, sortBy } = useDashboardStore();

  let filtered = (data ?? []).filter((e) => {
    if (query && !e.name.toLowerCase().includes(query.toLowerCase())) return false;
    return true;
  });

  filtered = filtered.slice().sort((a, b) => {
    if (sortBy === "name") {
      return a.name.localeCompare(b.name);
    }
    // Default: sort by created date (newest first)
    const ad = Date.parse(a.createdAt);
    const bd = Date.parse(b.createdAt);
    return bd - ad;
  });

  return (
    <div className="glass grain rounded-lg">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>URL</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filtered.map((e) => (
            <TableRow key={e.id}>
              <TableCell className="font-medium">{e.name}</TableCell>
              <TableCell className="uppercase text-xs text-muted-foreground">{e.type}</TableCell>
              <TableCell className="text-xs"><StatusChip status={mapStatus(e.status)} /></TableCell>
              <TableCell className="truncate max-w-[280px] text-xs text-muted-foreground">{e.url ?? "Not deployed"}</TableCell>
              <TableCell className="text-right space-x-2">
                <Button asChild variant="outline" size="sm" disabled={!e.url}><Link href={e.url ?? "#"}>Open</Link></Button>
                <Button asChild variant="outline" size="sm"><Link href={`/environments/${e.id}/logs`}>Logs</Link></Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function mapStatus(s: string): "Running" | "Stopped" | "Creating" | "Destroying" | "Error" | "Unknown" {
  switch (s) {
    case "active":
      return "Running";
    case "creating":
      return "Creating";
    case "destroying":
      return "Destroying";
    case "error":
      return "Error";
    default:
      return "Unknown";
  }
}


