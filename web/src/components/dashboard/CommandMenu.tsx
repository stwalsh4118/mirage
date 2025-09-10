"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import { useEnvironments } from "@/hooks/useEnvironments";
import { useDashboardStore } from "@/store/dashboard";

export function CommandMenu() {
  const [open, setOpen] = React.useState(false);
  const { data = [] } = useEnvironments();
  const { setQuery, setStatus, setType } = useDashboardStore();
  const router = useRouter();

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
      <CommandInput placeholder="Search environments or type a command…" />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>
        <CommandGroup heading="Environments">
          {data.slice(0, 8).map((e) => (
            <CommandItem key={e.id} value={e.name} onSelect={() => { setQuery(e.name); setOpen(false); }}>
              {e.name}
            </CommandItem>
          ))}
        </CommandGroup>
        <CommandSeparator />
        <CommandGroup heading="Filters">
          <CommandItem onSelect={() => { setStatus("all"); setOpen(false); }}>All</CommandItem>
          <CommandItem onSelect={() => { setStatus("active"); setOpen(false); }}>Active</CommandItem>
          <CommandItem onSelect={() => { setStatus("creating"); setOpen(false); }}>Creating</CommandItem>
          <CommandItem onSelect={() => { setStatus("error"); setOpen(false); }}>Error</CommandItem>
          <CommandItem onSelect={() => { setType("dev"); setOpen(false); }}>Type: Dev</CommandItem>
          <CommandItem onSelect={() => { setType("prod"); setOpen(false); }}>Type: Prod</CommandItem>
          <CommandItem onSelect={() => { setType("any"); setOpen(false); }}>Type: Any</CommandItem>
        </CommandGroup>
        <CommandSeparator />
        <CommandGroup heading="Navigation">
          <CommandItem onSelect={() => { router.push("/dashboard"); setOpen(false); }}>Go to Dashboard</CommandItem>
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}




