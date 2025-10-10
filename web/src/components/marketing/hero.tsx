"use client";

import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";
import { SignedIn, SignedOut, SignInButton } from "@clerk/nextjs";

export function Hero() {
  const [currentEnv, setCurrentEnv] = useState(0);
  const environments = ["production", "staging", "development", "preview"];

  useEffect(() => {
    const interval = setInterval(() => setCurrentEnv((prev) => (prev + 1) % environments.length), 2000);
    return () => clearInterval(interval);
  }, [environments.length]);

  return (
    <section className="relative min-h-[79vh] flex items-center overflow-hidden">
      <div className="absolute inset-0">
        <svg className="w-full h-full" viewBox="0 0 1200 800" fill="none">
          <defs>
            <radialGradient id="storm1" cx="20%" cy="30%" r="50%">
              <stop offset="0%" stopColor="rgba(218, 165, 32, 0.3)" />
              <stop offset="50%" stopColor="rgba(205, 133, 63, 0.15)" />
              <stop offset="100%" stopColor="transparent" />
            </radialGradient>
            <radialGradient id="storm2" cx="80%" cy="70%" r="45%">
              <stop offset="0%" stopColor="rgba(244, 164, 96, 0.25)" />
              <stop offset="100%" stopColor="transparent" />
            </radialGradient>
            <linearGradient id="dust" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="rgba(218, 165, 32, 0.1)" />
              <stop offset="50%" stopColor="rgba(205, 133, 63, 0.05)" />
              <stop offset="100%" stopColor="rgba(244, 164, 96, 0.1)" />
            </linearGradient>
          </defs>
          <rect width="100%" height="100%" fill="url(#dust)" />
          <ellipse cx="240" cy="240" rx="400" ry="300" fill="url(#storm1)" className="animate-pulse" />
          <ellipse cx="960" cy="560" rx="350" ry="250" fill="url(#storm2)" className="animate-pulse" style={{ animationDelay: "1s" }} />
        </svg>
      </div>

      <div className="relative container mx-auto px-4">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div className="relative">
            <div className="bg-card/90 backdrop-blur-2xl border border-border/70 shadow-2xl rounded-3xl p-8 lg:p-12 transform -rotate-1 hover:rotate-0 transition-all duration-500">
              <div className="bg-accent/20 backdrop-blur-sm rounded-lg px-3 py-1 text-xs font-medium text-accent mb-6 inline-block">Railway Environment Management</div>
              <div className="mb-6 pb-1 flex items-center gap-3">
                <Image src="/mirage_logo.png" alt="Mirage" width={600} height={140} className="h-14 md:h-16 w-auto drop-shadow-sm" priority />
                <span className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight leading-tight text-foreground/90">Mirage</span>
              </div>
              <p className="text-xl text-muted-foreground text-pretty mb-8 leading-relaxed">Your command center for Railway infrastructure. Browse projects, manage environments, and deploy services with an intuitive dashboard.</p>
              <div className="flex flex-col sm:flex-row gap-4">
                <SignedOut>
                  <SignInButton mode="redirect" forceRedirectUrl={"/dashboard"} fallbackRedirectUrl={"/dashboard"}>
                    <Button size="lg" className="font-semibold text-lg px-8 py-6 bg-gradient-to-r from-accent to-accent/80 hover:from-accent/90 hover:to-accent/70 shadow-lg">Start Building</Button>
                  </SignInButton>
                </SignedOut>
                <SignedIn>
                  <Link href="/dashboard">
                    <Button size="lg" className="font-semibold text-lg px-8 py-6 bg-gradient-to-r from-accent to-accent/80 hover:from-accent/90 hover:to-accent/70 shadow-lg">Go to Dashboard</Button>
                  </Link>
                </SignedIn>
                <Button size="lg" variant="outline" className="font-medium text-lg px-8 py-6 bg-card/60 backdrop-blur-sm border-border/50 hover:bg-card/80">Explore Docs</Button>
              </div>
            </div>
            <div className="absolute -bottom-6 -right-6 bg-card/95 backdrop-blur-xl border border-border/60 rounded-2xl p-6 shadow-xl transform rotate-2 hover:rotate-0 transition-all duration-300">
              <div className="text-sm text-muted-foreground mb-1">Live Updates</div>
              <div className="text-2xl font-bold text-accent">Real-time</div>
              <div className="text-xs text-muted-foreground">status monitoring</div>
            </div>
          </div>
          <div className="relative lg:pl-8">
            <div className="bg-card/95 backdrop-blur-2xl border border-border/70 rounded-2xl shadow-2xl overflow-hidden transform rotate-1 hover:rotate-0 transition-all duration-500">
              <div className="bg-muted/50 backdrop-blur-sm border-b border-border/30 px-6 py-4 flex items-center gap-3">
                <div className="flex gap-2"><div className="w-3 h-3 rounded-full bg-red-400"></div><div className="w-3 h-3 rounded-full bg-yellow-400"></div><div className="w-3 h-3 rounded-full bg-green-400"></div></div>
                <div className="text-sm font-mono text-muted-foreground lg:pl-2 pl-10">mirage dashboard</div>
              </div>
              <div className="p-6 text-sm space-y-4 min-h-[300px]">
                <div className="space-y-2">
                  <div className="text-xs text-muted-foreground uppercase font-semibold">Projects</div>
                  <div className="space-y-2">
                    <div className="flex items-center gap-3 p-2 rounded bg-muted/30">
                      <div className="w-2 h-2 rounded-full bg-green-400"></div>
                      <span className="text-foreground font-medium">production-api</span>
                      <span className="text-xs text-muted-foreground ml-auto">3 services</span>
                    </div>
                    <div className="flex items-center gap-3 p-2 rounded bg-muted/30">
                      <div className="w-2 h-2 rounded-full bg-blue-400"></div>
                      <span className="text-foreground font-medium">{environments[currentEnv]}-web</span>
                      <span className="text-xs text-muted-foreground ml-auto">2 services</span>
                    </div>
                  </div>
                </div>
                <div className="pt-2 border-t border-border/30">
                  <div className="text-xs text-muted-foreground uppercase font-semibold mb-2">Recent Activity</div>
                  <div className="text-xs text-accent">✓ Environment created in {environments[currentEnv]}-web</div>
                  <div className="text-xs text-muted-foreground mt-1">2 minutes ago</div>
                </div>
              </div>
            </div>
            <div className="absolute -top-4 -left-4 bg-card/90 backdrop-blur-xl border border-border/50 rounded-xl p-4 shadow-lg transform -rotate-3 hover:-rotate-1 transition-all duration-300"><div className="text-xs text-muted-foreground mb-1">Railway Projects</div><div className="text-lg font-bold text-accent">8</div></div>
            <div className="absolute -bottom-8 right-8 bg-card/90 backdrop-blur-xl border border-border/50 rounded-xl p-4 shadow-lg transform rotate-3 hover:rotate-1 transition-all duration-300"><div className="text-xs text-muted-foreground mb-1">Total Services</div><div className="text-lg font-bold text-green-400">24</div></div>
          </div>
        </div>
      </div>
    </section>
  );
}


