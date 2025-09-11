"use client";

import * as React from "react";
import { usePathname, useRouter } from "next/navigation";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import { useRailwayProjectsDetails } from "@/hooks/useRailway";
import { useDashboardStore } from "@/store/dashboard";

export function CommandMenu() {
  const [open, setOpen] = React.useState(false);
  const { data: projects = [] } = useRailwayProjectsDetails();
  const { setQuery, setSortBy, setView } = useDashboardStore();
  const router = useRouter();
  const pathname = usePathname();
  const isDashboard = pathname?.startsWith("/dashboard");
  const isProject = pathname?.startsWith("/project/");
  const currentUrl = typeof window !== "undefined" ? window.location.href : "";

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.key.toLowerCase() === "k" && (e.metaKey || e.ctrlKey))) {
        e.preventDefault();
        setOpen((o) => !o);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Type a command… (⌘K)" />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>
        <CommandGroup heading="Navigation">
          <CommandItem onSelect={() => { router.push("/dashboard"); setOpen(false); }}>Go to Dashboard</CommandItem>
        </CommandGroup>
        <CommandSeparator />
        <CommandGroup heading="Projects">
          {projects.slice(0, 10).map((p) => (
            <CommandItem key={p.id} value={p.name} onSelect={() => { router.push(`/project/${p.id}`); setOpen(false); }}>{p.name}</CommandItem>
          ))}
        </CommandGroup>
        <CommandSeparator />
        {isDashboard && (
          <>
            <CommandGroup heading="Dashboard">
              <CommandItem onSelect={() => { setSortBy("name" as any); setOpen(false); }}>Sort by name</CommandItem>
              <CommandItem onSelect={() => { setSortBy("services" as any); setOpen(false); }}>Sort by services</CommandItem>
              <CommandItem onSelect={() => { setSortBy("plugins" as any); setOpen(false); }}>Sort by plugins</CommandItem>
              <CommandItem onSelect={() => { setSortBy("environments" as any); setOpen(false); }}>Sort by environments</CommandItem>
              <CommandSeparator />
              <CommandItem onSelect={() => { setView("grid" as any); setOpen(false); }}>View: Grid</CommandItem>
              <CommandItem onSelect={() => { setView("list" as any); setOpen(false); }}>View: Table</CommandItem>
            </CommandGroup>
            <CommandSeparator />
          </>
        )}
        {isProject && (
          <>
            <CommandGroup heading="Project">
              <CommandItem onSelect={() => { if (currentUrl) { navigator.clipboard?.writeText(currentUrl).catch(() => {});} setOpen(false); }}>Copy page link</CommandItem>
              <CommandItem onSelect={() => { const url = new URL(currentUrl); if (!url.searchParams.has("demo")) url.searchParams.set("demo", "1"); router.push(url.pathname + "?" + url.searchParams.toString()); setOpen(false); }}>Enable demo data</CommandItem>
              <CommandItem onSelect={() => { const url = new URL(currentUrl); url.searchParams.delete("demo"); router.push(url.pathname + (url.searchParams.toString() ? ("?" + url.searchParams.toString()) : "")); setOpen(false); }}>Disable demo data</CommandItem>
            </CommandGroup>
            <CommandSeparator />
          </>
        )}
      </CommandList>
    </CommandDialog>
  );
}






