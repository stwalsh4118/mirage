"use client";

import { Button } from "@/components/ui/button";
import Link from "next/link";
import { SignedIn, SignedOut, SignInButton } from "@clerk/nextjs";

export function FinalCTA() {
  return (
    <section className="py-24 bg-background border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="glass grain sheen rounded-2xl p-12 lg:p-16 text-center max-w-4xl mx-auto">
          <h2 className="text-3xl lg:text-5xl font-semibold tracking-tight mb-6">Start managing your Railway projects</h2>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto leading-relaxed">Take control of your Railway infrastructure with a powerful, intuitive dashboard.</p>
          <SignedOut>
            <SignInButton mode="redirect" forceRedirectUrl={"/dashboard"} fallbackRedirectUrl={"/dashboard"}>
              <Button size="lg" className="font-medium">Get started for free</Button>
            </SignInButton>
          </SignedOut>
          <SignedIn>
            <Link href="/dashboard">
              <Button size="lg" className="font-medium">Go to Dashboard</Button>
            </Link>
          </SignedIn>
        </div>
      </div>
    </section>
  );
}



