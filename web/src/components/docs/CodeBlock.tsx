"use client";

import { useState } from "react";
import { Check, Copy } from "lucide-react";
import { Button } from "@/components/ui/button";

interface CodeBlockProps {
  children: string;
  className?: string;
  inline?: boolean;
}

export function CodeBlock({ children, className, inline }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  // Extract language from className (format: language-xxx or hljs language-xxx)
  const language = className?.replace(/language-/, "").replace(/hljs\s*/gi, "").trim() || "text";

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(children);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy code:', err);
      // Optionally show a toast notification
    }
  };

  // For inline code, render a simple span
  if (inline) {
    return (
      <code className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
        {children}
      </code>
    );
  }

  // For code blocks
  return (
    <div className="group relative my-4 rounded-lg border border-border bg-muted/50">
      {/* Language label and copy button */}
      <div className="flex items-center justify-between border-b border-border px-4 py-2">
        <span className="text-xs font-medium text-muted-foreground uppercase">
          {language}
        </span>
        <Button
          variant="ghost"
          size="sm"
          className="h-7 w-7 p-0 opacity-70 hover:opacity-100"
          onClick={handleCopy}
        >
          {copied ? (
            <Check className="h-3.5 w-3.5 text-green-500" />
          ) : (
            <Copy className="h-3.5 w-3.5" />
          )}
          <span className="sr-only">{copied ? "Copied" : "Copy code"}</span>
        </Button>
      </div>

      {/* Code content */}
      <div className="overflow-x-auto">
        <pre className="p-4">
          <code className={className}>{children}</code>
        </pre>
      </div>
    </div>
  );
}

