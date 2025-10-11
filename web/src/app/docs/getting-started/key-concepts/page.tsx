import { MarkdownRenderer } from "@/components/docs/MarkdownRenderer";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { getDocsContent } from "@/lib/docs";

export default async function KeyConceptsPage() {
  const content = await getDocsContent("getting-started/key-concepts.md");

  return (
    <div className="space-y-6">
      <div className="glass grain rounded-lg p-8 border border-border">
        <MarkdownRenderer content={content} />
      </div>

      <div className="flex justify-between items-center glass grain rounded-lg p-4 border border-border">
        <Button variant="ghost" asChild>
          <Link href="/docs/getting-started/first-environment">
            <ArrowLeft className="mr-2 h-4 w-4" /> Previous: First Environment
          </Link>
        </Button>
        <Button asChild>
          <Link href="/docs/features">
            Explore Features
          </Link>
        </Button>
      </div>
    </div>
  );
}
