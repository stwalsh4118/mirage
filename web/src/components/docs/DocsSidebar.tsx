"use client";

import { useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronRight, ChevronDown } from "lucide-react";
import { docsNavigation } from "@/config/docsNavigation";
import { useDocsStore } from "@/store/docs";
import { cn } from "@/lib/utils";

export function DocsSidebar() {
  const pathname = usePathname();
  const { openSections, toggleSection, isOpen, setOpenSections } = useDocsStore();

  // Initialize with active section if store is empty (first visit)
  useEffect(() => {
    if (openSections.length === 0) {
      const activeSection = docsNavigation.find((section) =>
        section.pages.some((page) => pathname === page.path)
      );
      if (activeSection) {
        setOpenSections([activeSection.title]);
      }
    }
  }, [pathname, openSections.length, setOpenSections]);

  return (
    <nav className="space-y-2">
      {docsNavigation.map((section) => {
        const isSectionOpen = isOpen(section.title);
        const hasActivePage = section.pages.some((page) => pathname === page.path);

        return (
          <div key={section.title} className="space-y-1">
            <button
              onClick={() => toggleSection(section.title)}
              className={cn(
                "flex items-center justify-between w-full px-3 py-2 text-sm font-medium rounded-md transition-colors",
                "hover:bg-muted",
                hasActivePage && "text-primary"
              )}
            >
              <span>{section.title}</span>
              {isSectionOpen ? (
                <ChevronDown className="h-4 w-4 text-muted-foreground" />
              ) : (
                <ChevronRight className="h-4 w-4 text-muted-foreground" />
              )}
            </button>

            {isSectionOpen && (
              <div className="space-y-0.5 ml-3 pl-3 border-l border-border">
                {section.pages.map((page) => {
                  const isActive = pathname === page.path;

                  return (
                    <Link
                      key={page.path}
                      href={page.path}
                      className={cn(
                        "block px-3 py-2 text-sm rounded-md transition-colors",
                        isActive
                          ? "bg-primary/10 text-primary font-medium"
                          : "text-muted-foreground hover:bg-muted hover:text-foreground"
                      )}
                    >
                      {page.title}
                    </Link>
                  );
                })}
              </div>
            )}
          </div>
        );
      })}
    </nav>
  );
}

