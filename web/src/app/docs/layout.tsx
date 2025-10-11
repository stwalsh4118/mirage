import type { ReactNode } from "react";
import Link from "next/link";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { DocsSidebar } from "@/components/docs/DocsSidebar";
import { Breadcrumbs } from "@/components/docs/Breadcrumbs";
import { WipBanner } from "@/components/docs/WipBanner";

export default function DocsLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background sandstorm-bg">
      {/* Documentation Header */}
      <div className="glass grain sticky top-0 z-40">
        <div className="max-w-screen-2xl mx-auto px-8 py-3 flex items-center gap-3">
          <div className="flex items-center gap-2 pr-2">
            <Link href="/" aria-label="Go to Home" className="inline-flex items-center">
              <Image src="/mirage_logo.png" alt="Mirage" width={36} height={36} className="h-9 w-auto" />
            </Link>
            <span className="text-sm font-medium text-muted-foreground">Docs</span>
          </div>
          <Separator orientation="vertical" className="mx-1 h-6" />
          <div className="flex-1" />
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" asChild>
              <Link href="/dashboard">Dashboard</Link>
            </Button>
            <Button variant="ghost" size="sm" asChild>
              <Link href="/sign-in">Sign In</Link>
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="max-w-screen-2xl mx-auto px-8">
        <div className="flex gap-8">
          {/* Sidebar */}
          <aside className="hidden lg:block w-64 flex-shrink-0">
            <div className="sticky top-24 py-6 max-h-[calc(100vh-6rem)] overflow-y-auto">
              <DocsSidebar />
            </div>
          </aside>

          {/* Main Content */}
          <main className="flex-1 py-6 min-w-0">
            <Breadcrumbs />
            <WipBanner />
            {children}
          </main>
        </div>
      </div>
    </div>
  );
}

