"use client";

import Image from "next/image";
import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";

export function Header() {
  const handleGetStarted = () => {
    alert("Welcome to Mirage! Sign up coming soon.");
  };

  return (
    <header className="sticky top-0 z-50 glass border-b">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Image
            src="/mirage_logo.png"
            alt="Mirage logo"
            width={112}
            height={32}
            className="h-7 w-auto select-none self-center relative top-[1px]"
            priority
          />
          <span className="text-2xl font-semibold tracking-tight leading-none text-foreground">Mirage</span>
        </div>
        <nav className="hidden md:flex items-center space-x-6">
          <a href="#" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Docs</a>
          <a href="#" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">GitHub</a>
          <a href="#" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Sign in</a>
        </nav>
        <div className="flex items-center space-x-2">
          <ModeToggle />
          <Button onClick={handleGetStarted} className="font-medium">Get started</Button>
        </div>
      </div>
    </header>
  );
}



