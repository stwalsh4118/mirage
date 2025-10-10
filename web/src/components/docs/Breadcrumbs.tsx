"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronRight } from "lucide-react";
import { docsNavigation } from "@/config/docsNavigation";

export function Breadcrumbs() {
  const pathname = usePathname();

  // Build breadcrumb trail
  const breadcrumbs: Array<{ title: string; path: string }> = [
    { title: "Docs", path: "/docs" },
  ];

  // Find the current page in navigation
  if (pathname !== "/docs") {
    for (const section of docsNavigation) {
      const page = section.pages.find((p) => p.path === pathname);
      if (page) {
        breadcrumbs.push({
          title: section.title,
          path: "#", // Section doesn't have its own page
        });
        breadcrumbs.push({
          title: page.title,
          path: page.path,
        });
        break;
      }
    }

    // If not found in navigation, build from URL path
    if (breadcrumbs.length === 1) {
      const pathParts = pathname.split("/").filter(Boolean);
      let currentPath = "";
      
      pathParts.forEach((part, index) => {
        if (part !== "docs") {
          currentPath += `/${part}`;
          const title = part
            .split("-")
            .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(" ");
          
          breadcrumbs.push({
            title,
            path: index === pathParts.length - 1 ? pathname : currentPath,
          });
        }
      });
    }
  }

  return (
    <nav className="flex items-center space-x-1 text-sm text-muted-foreground mb-6">
      {breadcrumbs.map((crumb, index) => {
        const isLast = index === breadcrumbs.length - 1;
        const isSection = crumb.path === "#";

        return (
          <div key={crumb.path + index} className="flex items-center">
            {index > 0 && (
              <ChevronRight className="h-4 w-4 mx-1 flex-shrink-0" />
            )}
            {isLast || isSection ? (
              <span className={isLast ? "text-foreground font-medium" : ""}>
                {crumb.title}
              </span>
            ) : (
              <Link
                href={crumb.path}
                className="hover:text-foreground transition-colors"
              >
                {crumb.title}
              </Link>
            )}
          </div>
        );
      })}
    </nav>
  );
}

