"use client";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import Link from "next/link";
import type { RailwayEnvironmentWithServices } from "@/lib/api/railway";
import { Pill } from "./Pill";

export function RailwayEnvironmentCard({ env, href }: { env: RailwayEnvironmentWithServices; href?: string }) {
  const serviceCount = env.services?.length ?? 0;
  const content = (
    <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div className="space-y-1">
          <CardTitle className="text-base font-medium">{env.name}</CardTitle>
          <div className="text-[11px] text-muted-foreground">{env.id}</div>
        </div>
        <div className="flex items-center gap-2">
          <Pill color={serviceCount > 0 ? "green" : "amber"}>
            {serviceCount} service{serviceCount === 1 ? "" : "s"}
          </Pill>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="text-xs text-muted-foreground">Railway environment</div>
        {serviceCount > 0 && (
          <div className="flex flex-wrap gap-2 text-[11px]">
            {env.services!.slice(0, 4).map((s) => (
              <Pill key={s.id} color="blue">{s.name}</Pill>
            ))}
            {env.services!.length > 4 && (
              <Pill color="neutral">+{env.services!.length - 4} more</Pill>
            )}
          </div>
        )}
      </CardContent>
      <CardFooter className="flex items-center justify-between">
        <div className="text-[11px] text-muted-foreground">Read-only</div>
        <div className="text-[11px] text-muted-foreground">Services reflect latest fetch</div>
      </CardFooter>
    </Card>
  );
  if (href) {
    return (
      <Link href={href} className="block" aria-label={`Open project for ${env.name}`}>
        {content}
      </Link>
    );
  }
  return content;
}


