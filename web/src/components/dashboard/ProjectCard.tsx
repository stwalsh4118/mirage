"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { RailwayProjectDetails } from "@/lib/api/railway";

export function ProjectCard({ project }: { project: RailwayProjectDetails }) {
  const servicesCount = project.services?.length ?? 0;
  const pluginsCount = project.plugins?.length ?? 0;
  const environmentsCount = project.environments?.length ?? 0;

  return (
    <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div className="space-y-1">
          <CardTitle className="text-base font-medium">{project.name}</CardTitle>
          <div className="text-[11px] text-muted-foreground">{project.id}</div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>services</span>
          <span className="font-medium text-foreground/80">{servicesCount}</span>
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>plugins</span>
          <span className="font-medium text-foreground/80">{pluginsCount}</span>
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>environments</span>
          <span className="font-medium text-foreground/80">{environmentsCount}</span>
        </div>
      </CardContent>
      <CardFooter className="flex items-center justify-between">
        <Button asChild variant="outline" size="sm" aria-label="Open project details">
          <Link href={`/dashboard/projects/${project.id}`}>Open</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}


