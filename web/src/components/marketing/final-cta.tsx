"use client";

import { Button } from "@/components/ui/button";

export function FinalCTA() {
  const handleCreateEnvironment = () => alert("Ready to create your first environment! Sign up coming soon.");
  return (
    <section className="py-24 bg-background border-t border-border/40">
      <div className="container mx-auto px-4">
        <div className="glass grain sheen rounded-2xl p-12 lg:p-16 text-center max-w-4xl mx-auto">
          <h2 className="text-3xl lg:text-5xl font-semibold tracking-tight mb-6">Create your first environment</h2>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto leading-relaxed">Join thousands of developers who trust Mirage for their environment management.</p>
          <Button size="lg" onClick={handleCreateEnvironment} className="font-medium">Get started for free</Button>
        </div>
      </div>
    </section>
  );
}



