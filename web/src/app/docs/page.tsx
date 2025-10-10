import Link from "next/link";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { BookOpen, Rocket, Wrench, AlertCircle } from "lucide-react";

export default function DocsHomePage() {
  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="space-y-4">
        <h1 className="text-4xl font-bold tracking-tight">
          Mirage Documentation
        </h1>
        <p className="text-xl text-muted-foreground max-w-3xl">
          Learn how to use Mirage to manage your Railway environments with ease. 
          From getting started to advanced features, find everything you need here.
        </p>
      </div>

      {/* Main Documentation Sections */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="glass grain hover:border-primary transition-colors">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-primary/10">
                <Rocket className="h-6 w-6 text-primary" />
              </div>
              <CardTitle>Getting Started</CardTitle>
            </div>
            <CardDescription>
              New to Mirage? Start here to learn the basics and create your first environment.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Link href="/docs/getting-started/introduction" className="block text-sm text-primary hover:underline">
                Introduction to Mirage
              </Link>
              <Link href="/docs/getting-started/prerequisites" className="block text-sm text-primary hover:underline">
                Prerequisites
              </Link>
              <Link href="/docs/getting-started/setup" className="block text-sm text-primary hover:underline">
                Setup Guide
              </Link>
              <Link href="/docs/getting-started/first-environment" className="block text-sm text-primary hover:underline">
                Create Your First Environment
              </Link>
            </div>
            <Button variant="ghost" size="sm" className="mt-4 w-full" asChild>
              <Link href="/docs/getting-started">View All →</Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="glass grain hover:border-primary transition-colors">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-accent/10">
                <BookOpen className="h-6 w-6 text-accent" />
              </div>
              <CardTitle>Features</CardTitle>
            </div>
            <CardDescription>
              Explore Mirage&apos;s powerful features for managing Railway environments.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Link href="/docs/features/railway" className="block text-sm text-primary hover:underline">
                Railway Integration
              </Link>
              <Link href="/docs/features/environments" className="block text-sm text-primary hover:underline">
                Environment Management
              </Link>
              <Link href="/docs/features/dashboard" className="block text-sm text-primary hover:underline">
                Dashboard Overview
              </Link>
            </div>
            <Button variant="ghost" size="sm" className="mt-4 w-full" asChild>
              <Link href="/docs/features">View All →</Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="glass grain hover:border-primary transition-colors">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-chart-1/20">
                <Wrench className="h-6 w-6 text-chart-1" />
              </div>
              <CardTitle>How-To Guides</CardTitle>
            </div>
            <CardDescription>
              Step-by-step guides for common tasks and workflows.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Link href="/docs/how-to/create-dev-environment" className="block text-sm text-primary hover:underline">
                Create a Dev Environment
              </Link>
              <Link href="/docs/how-to/manage-environment-variables" className="block text-sm text-primary hover:underline">
                Manage Environment Variables
              </Link>
              <Link href="/docs/how-to/use-templates-effectively" className="block text-sm text-primary hover:underline">
                Use Templates Effectively
              </Link>
            </div>
            <Button variant="ghost" size="sm" className="mt-4 w-full" asChild>
              <Link href="/docs/how-to">View All →</Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="glass grain hover:border-primary transition-colors">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-destructive/10">
                <AlertCircle className="h-6 w-6 text-destructive" />
              </div>
              <CardTitle>Troubleshooting</CardTitle>
            </div>
            <CardDescription>
              Find solutions to common issues and errors.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Link href="/docs/troubleshooting/connection-issues" className="block text-sm text-primary hover:underline">
                Connection Issues
              </Link>
              <Link href="/docs/troubleshooting/authentication-errors" className="block text-sm text-primary hover:underline">
                Authentication Errors
              </Link>
              <Link href="/docs/troubleshooting/environment-creation-failures" className="block text-sm text-primary hover:underline">
                Environment Creation Failures
              </Link>
            </div>
            <Button variant="ghost" size="sm" className="mt-4 w-full" asChild>
              <Link href="/docs/troubleshooting">View All →</Link>
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Quick Links Section */}
      <div className="mt-12 p-6 rounded-lg border border-border glass grain">
        <h2 className="text-2xl font-semibold mb-4">Quick Links</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <h3 className="font-medium mb-2">Popular Topics</h3>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li>
                <Link href="/docs/getting-started/first-environment" className="hover:text-primary">
                  Creating your first environment
                </Link>
              </li>
              <li>
                <Link href="/docs/features/railway/connecting-account" className="hover:text-primary">
                  Connecting Railway account
                </Link>
              </li>
              <li>
                <Link href="/docs/features/environments/templates" className="hover:text-primary">
                  Using templates
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="font-medium mb-2">Resources</h3>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li>
                <Link href="https://railway.app/docs" className="hover:text-primary" target="_blank" rel="noopener noreferrer">
                  Railway Documentation ↗
                </Link>
              </li>
              <li>
                <Link href="/dashboard" className="hover:text-primary">
                  Go to Dashboard
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="font-medium mb-2">Need Help?</h3>
            <ul className="space-y-1 text-sm text-muted-foreground">
              <li>
                <Link href="/docs/troubleshooting" className="hover:text-primary">
                  Troubleshooting Guide
                </Link>
              </li>
              <li>
                <Link href="/docs/troubleshooting/getting-help" className="hover:text-primary">
                  Getting Support
                </Link>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

