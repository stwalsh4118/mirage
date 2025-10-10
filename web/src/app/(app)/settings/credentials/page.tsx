"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Key, Github, Package, FolderLock } from "lucide-react";
import { RailwayTokenTab } from "@/components/credentials/RailwayTokenTab";
import { GitHubTokenTab } from "@/components/credentials/GitHubTokenTab";

// Force dynamic rendering - don't prerender this page at build time
export const dynamic = 'force-dynamic';

export default function CredentialsPage() {
  return (
    <div className="space-y-6">
      <main className="max-w-screen-2xl mx-auto px-8 space-y-6">
        {/* Page Header */}
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Credentials</h1>
          <p className="text-muted-foreground">
            Manage your API tokens, registry credentials, and secrets securely
          </p>
        </div>

        {/* Tabbed Interface */}
        <Tabs defaultValue="railway" className="w-full">
          <TabsList className="grid w-full grid-cols-4 lg:w-auto lg:inline-flex">
            <TabsTrigger value="railway" className="gap-2">
              <Key className="h-4 w-4" />
              <span className="hidden sm:inline">Railway Token</span>
              <span className="sm:hidden">Railway</span>
            </TabsTrigger>
            <TabsTrigger value="github" className="gap-2">
              <Github className="h-4 w-4" />
              <span className="hidden sm:inline">GitHub Token</span>
              <span className="sm:hidden">GitHub</span>
            </TabsTrigger>
            <TabsTrigger value="docker" className="gap-2">
              <Package className="h-4 w-4" />
              <span className="hidden sm:inline">Docker Registries</span>
              <span className="sm:hidden">Docker</span>
            </TabsTrigger>
            <TabsTrigger value="custom" className="gap-2">
              <FolderLock className="h-4 w-4" />
              <span className="hidden sm:inline">Custom Secrets</span>
              <span className="sm:hidden">Custom</span>
            </TabsTrigger>
          </TabsList>

          {/* Railway Token Tab */}
          <TabsContent value="railway" className="mt-6 space-y-4">
            <RailwayTokenTab />
          </TabsContent>

          {/* GitHub Token Tab */}
          <TabsContent value="github" className="mt-6 space-y-4">
            <GitHubTokenTab />
          </TabsContent>

          {/* Docker Registries Tab */}
          <TabsContent value="docker" className="mt-6 space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Docker Registry Credentials</CardTitle>
                <CardDescription>
                  Manage credentials for Docker Hub, GHCR, and other container registries
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-center py-12 text-muted-foreground">
                  Docker registry management will be implemented here
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Custom Secrets Tab */}
          <TabsContent value="custom" className="mt-6 space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Custom Secrets</CardTitle>
                <CardDescription>
                  Store and manage custom secrets with version history and rollback capabilities
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-center py-12 text-muted-foreground">
                  Custom secrets management will be implemented here
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}

