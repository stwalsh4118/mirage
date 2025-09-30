"use client";

import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDashboardStore } from "@/store/dashboard";

export function ControlsBar() {
  const { sortBy, setSortBy, view, setView, query, setQuery } = useDashboardStore();
  return (
    <div className="glass grain rounded-lg p-3 flex flex-col lg:flex-row gap-3 items-start lg:items-center justify-between">
      <div className="flex items-center gap-2" />
      <div className="flex-1 grid grid-cols-1 md:grid-cols-3 gap-3">
        <div className="flex flex-col gap-1">
          <Label htmlFor="filter">Filter projects…</Label>
          <Input id="filter" placeholder="Filter projects…" className="h-9" value={query} onChange={(e) => setQuery(e.target.value)} />
        </div>
        <div className="flex flex-col gap-1">
          <Label>Sort</Label>
          <Select value={sortBy} onValueChange={(v) => setSortBy(v as any)}>
            <SelectTrigger className="h-9"><SelectValue placeholder="Last updated" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="name">Name</SelectItem>
              <SelectItem value="services">Services</SelectItem>
              <SelectItem value="plugins">Plugins</SelectItem>
              <SelectItem value="environments">Environments</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <ToggleGroup type="single" value={view} onValueChange={(v) => v && setView(v as any)}>
          <ToggleGroupItem value="grid">▦</ToggleGroupItem>
          <ToggleGroupItem value="list">☰</ToggleGroupItem>
        </ToggleGroup>
      </div>
    </div>
  );
}


