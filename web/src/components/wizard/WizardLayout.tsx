"use client";

import { PropsWithChildren } from "react";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

export function WizardLayout(props: PropsWithChildren<{ title: string; breadcrumb?: { href: string; label: string }[] }>) {
  const { children, title, breadcrumb } = props;
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
          {breadcrumb && (
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumb.map((b, i) => (
                  <>
                    <BreadcrumbItem key={`${b.href}-${i}`}>
                      <BreadcrumbLink href={b.href}>{b.label}</BreadcrumbLink>
                    </BreadcrumbItem>
                    {i < breadcrumb.length - 1 && <BreadcrumbSeparator />}
                  </>
                ))}
                <BreadcrumbItem>
                  <BreadcrumbPage>{title}</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          )}
        </div>
      </div>
      <Card className="glass grain border-border/60">
        <CardHeader className="pb-0"></CardHeader>
        <CardContent className="pt-6">{children}</CardContent>
      </Card>
    </div>
  );
}


