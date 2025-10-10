"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertCircle, ExternalLink, Key } from "lucide-react";
import Link from "next/link";
import { useRailwayTokenStatus } from "@/hooks/useRailwayToken";

export function RailwayTokenOnboarding() {
  const { data: status, isLoading } = useRailwayTokenStatus();

  // Don't show if loading or token is already configured
  if (isLoading || status?.configured) {
    return null;
  }

  return (
    <Card className="glass grain border-2 border-accent/30 bg-accent/5">
      <CardContent className="p-6">
        <div className="flex items-start gap-4">
          <div className="rounded-full bg-accent/20 p-3 mt-1">
            <Key className="h-6 w-6 text-accent" />
          </div>
          <div className="flex-1 space-y-4">
            <div>
              <h3 className="text-lg font-semibold mb-2">Welcome to Mirage! 🎉</h3>
              <p className="text-muted-foreground leading-relaxed">
                To get started managing your Railway infrastructure, you&apos;ll need to connect your Railway API token.
                This allows Mirage to securely access and manage your Railway projects.
              </p>
            </div>

            <div className="bg-muted/30 rounded-lg p-4 space-y-2">
              <div className="flex items-start gap-2">
                <AlertCircle className="h-4 w-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="text-sm text-muted-foreground">
                  <p className="font-medium mb-1">Quick setup (2 steps):</p>
                  <ol className="list-decimal list-inside space-y-1">
                    <li>Get your Railway API token from your account settings</li>
                    <li>Add it to Mirage in the Credentials page</li>
                  </ol>
                </div>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-3">
              <Link href="/settings/credentials">
                <Button size="lg" className="w-full sm:w-auto">
                  <Key className="mr-2 h-4 w-4" />
                  Add Railway Token
                </Button>
              </Link>
              <a
                href="https://railway.app/account/tokens"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="lg" variant="outline" className="w-full sm:w-auto bg-transparent">
                  <ExternalLink className="mr-2 h-4 w-4" />
                  Get Token from Railway
                </Button>
              </a>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

