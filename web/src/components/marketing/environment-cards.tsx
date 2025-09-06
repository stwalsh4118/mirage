"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

interface Environment {
  name: string;
  type: "Dev" | "Prod" | "Staging";
  status: "Running" | "Building" | "Stopped";
  url: string;
  createdAt: string;
}

export function EnvironmentCards() {
  const [loading, setLoading] = useState(true);
  const [environments] = useState<Environment[]>([
    { name: "api-service", type: "Dev", status: "Running", url: "api-dev.mirage.app", createdAt: "2 hours ago" },
    { name: "frontend-app", type: "Prod", status: "Running", url: "app.mirage.com", createdAt: "1 day ago" },
    { name: "worker-queue", type: "Staging", status: "Building", url: "worker-staging.mirage.app", createdAt: "5 minutes ago" },
    { name: "analytics-db", type: "Dev", status: "Stopped", url: "analytics-dev.mirage.app", createdAt: "3 days ago" },
  ]);

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 1500);
    return () => clearTimeout(timer);
  }, []);

  const statusColor = (status: string) => (status === "Running" ? "bg-green-500/10 text-green-700 dark:text-green-400" : status === "Building" ? "bg-yellow-500/10 text-yellow-700 dark:text-yellow-400" : "bg-gray-500/10 text-gray-700 dark:text-gray-400");
  const typeColor = (type: string) => (type === "Prod" ? "bg-red-500/10 text-red-700 dark:text-red-400" : type === "Staging" ? "bg-yellow-500/10 text-yellow-700 dark:text-yellow-400" : "bg-blue-500/10 text-blue-700 dark:text-blue-400");

  return (
    <section className="py-24 section-soft-alt border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl lg:text-4xl font-semibold tracking-tight mb-4">Live environment dashboard</h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">See how your environments look in the Mirage dashboard with real-time status updates.</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 max-w-6xl mx-auto">
          {environments.map((env, index) => (
            <Card key={index} className="glass grain hover:scale-[1.01] hover:-translate-y-1 transition-all duration-200">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-base font-medium truncate">{loading ? <Skeleton className="h-4 w-24" /> : env.name}</CardTitle>
                  {loading ? <Skeleton className="h-5 w-12" /> : <Badge variant="secondary" className={typeColor(env.type)}>{env.type}</Badge>}
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Status</span>
                  {loading ? <Skeleton className="h-5 w-16" /> : <Badge variant="secondary" className={statusColor(env.status)}>{env.status}</Badge>}
                </div>
                <div>
                  <span className="text-sm text-muted-foreground">URL</span>
                  <div className="text-sm font-mono truncate mt-1">{loading ? <Skeleton className="h-4 w-full" /> : env.url}</div>
                </div>
                <div className="flex items-center justify-between text-xs text-muted-foreground"><span>Created {loading ? <Skeleton className="inline-block h-3 w-16" /> : env.createdAt}</span></div>
                <div className="flex gap-2 pt-2">
                  <Button size="sm" variant="outline" className="flex-1 text-xs bg-transparent">View</Button>
                  <Button size="sm" variant="outline" className="flex-1 text-xs bg-transparent">Logs</Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}



