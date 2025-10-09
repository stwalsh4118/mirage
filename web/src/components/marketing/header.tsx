"use client";

import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { SignedIn, SignedOut, SignInButton, UserButton } from "@clerk/nextjs";

export function Header() {
  return (
    <header className="sticky top-0 z-50 glass border-b">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Link href="/">
            <Image
              src="/mirage_logo.png"
              alt="Mirage logo"
              width={112}
              height={32}
              className="h-7 w-auto select-none self-center relative top-[1px]"
              priority
            />
          </Link>
          <span className="text-2xl font-semibold tracking-tight leading-none text-foreground">Mirage</span>
        </div>
        <nav className="hidden md:flex items-center space-x-6">
          <a href="#" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Docs</a>
          <a href="#" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">GitHub</a>
          <SignedOut>
            <Link href="/sign-in" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Sign in</Link>
          </SignedOut>
          <SignedIn>
            <Link href="/dashboard" className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">Dashboard</Link>
          </SignedIn>
        </nav>
        <div className="flex items-center space-x-2">
          <ModeToggle />
          <SignedOut>
            <SignInButton mode="redirect" forceRedirectUrl={"/dashboard"} fallbackRedirectUrl={"/dashboard"}>
              <Button className="font-medium">Get started</Button>
            </SignInButton>
          </SignedOut>
          <SignedIn>
            <UserButton
              appearance={{
                elements: {
                  avatarBox: "w-8 h-8",
                  userButtonPopoverCard: "shadow-xl border border-border rounded-lg",
                },
              }}
              afterSignOutUrl="/"
            />
          </SignedIn>
        </div>
      </div>
    </header>
  );
}



