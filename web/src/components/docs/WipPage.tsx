import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft, FileText } from "lucide-react";

interface WipPageProps {
  title: string;
  backLink?: string;
  backLabel?: string;
}

export function WipPage({ title, backLink = "/docs", backLabel = "Back to Documentation" }: WipPageProps) {
  return (
    <div className="space-y-6">
      <div className="glass grain rounded-lg p-12 border border-border text-center">
        <FileText className="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
        <h1 className="text-3xl font-bold mb-4">{title}</h1>
        <p className="text-lg text-muted-foreground mb-6">
          This page is coming soon. We're working on this documentation section.
        </p>
        <p className="text-sm text-muted-foreground mb-8">
          In the meantime, check out the{" "}
          <Link href="/docs/getting-started/introduction" className="text-primary hover:underline">
            Getting Started Guide
          </Link>
          {" "}to learn the basics.
        </p>
        <Button variant="outline" asChild>
          <Link href={backLink}>
            <ArrowLeft className="mr-2 h-4 w-4" /> {backLabel}
          </Link>
        </Button>
      </div>
    </div>
  );
}

