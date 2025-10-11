import { MarkdownRenderer } from "@/components/docs/MarkdownRenderer";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import { getDocsContent } from "@/lib/docs";

export default async function IntroductionPage() {
  const content = await getDocsContent("getting-started/introduction.md");

  return (
    <div className="space-y-6">
      <div className="glass grain rounded-lg p-8 border border-border">
        <MarkdownRenderer content={content} />
      </div>

      <div className="flex justify-between items-center glass grain rounded-lg p-4 border border-border">
        <span className="text-sm text-muted-foreground">Next: Prerequisites</span>
        <Button asChild>
          <Link href="/docs/getting-started/prerequisites">
            Continue <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
