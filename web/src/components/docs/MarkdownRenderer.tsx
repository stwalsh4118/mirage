"use client";

import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeHighlight from "rehype-highlight";
import rehypeRaw from "rehype-raw";
import Link from "next/link";
import type { Components } from "react-markdown";

// Import highlight.js styles
// We'll import this in globals.css instead to avoid SSR issues

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export function MarkdownRenderer({ content, className = "" }: MarkdownRendererProps) {
  // Custom component mappings for styled markdown elements
  const components: Partial<Components> = {
    // Headings
    h1: ({ children }) => (
      <h1 className="scroll-m-20 text-4xl font-bold tracking-tight mt-8 mb-4 first:mt-0">
        {children}
      </h1>
    ),
    h2: ({ children }) => (
      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight mt-8 mb-4 first:mt-0">
        {children}
      </h2>
    ),
    h3: ({ children }) => (
      <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight mt-6 mb-3">
        {children}
      </h3>
    ),
    h4: ({ children }) => (
      <h4 className="scroll-m-20 text-xl font-semibold tracking-tight mt-6 mb-3">
        {children}
      </h4>
    ),
    h5: ({ children }) => (
      <h5 className="scroll-m-20 text-lg font-semibold tracking-tight mt-4 mb-2">
        {children}
      </h5>
    ),
    h6: ({ children }) => (
      <h6 className="scroll-m-20 text-base font-semibold tracking-tight mt-4 mb-2">
        {children}
      </h6>
    ),

    // Paragraphs
    p: ({ children }) => (
      <p className="leading-7 [&:not(:first-child)]:mt-6">
        {children}
      </p>
    ),

    // Links - use Next.js Link for internal links
    a: ({ href, children }) => {
      const isInternal = href?.startsWith("/") || href?.startsWith("#");
      
      if (isInternal) {
        return (
          <Link
            href={href || "#"}
            className="font-medium text-primary underline underline-offset-4 hover:text-primary/80"
          >
            {children}
          </Link>
        );
      }

      return (
        <a
          href={href || "#"}
          target="_blank"
          rel="noopener noreferrer"
          className="font-medium text-primary underline underline-offset-4 hover:text-primary/80"
        >
          {children}
        </a>
      );
    },

    // Lists
    ul: ({ children }) => (
      <ul className="my-6 ml-6 list-disc [&>li]:mt-2">
        {children}
      </ul>
    ),
    ol: ({ children }) => (
      <ol className="my-6 ml-6 list-decimal [&>li]:mt-2">
        {children}
      </ol>
    ),
    li: ({ children }) => (
      <li className="leading-7">
        {children}
      </li>
    ),

    // Pre blocks (code blocks with language)
    pre: ({ children, ...props }) => {
      // Extract the code element and its props with type guard
      const className = 
        React.isValidElement(children) && typeof children.props === 'object'
          ? ((children.props as { className?: string }).className || '')
          : '';
      const match = /language-(\w+)/.exec(className);
      const language = match ? match[1] : 'text';
      
      return (
        <div className="group relative my-4 rounded-lg border border-border bg-muted/50">
          <div className="flex items-center justify-between border-b border-border px-4 py-2">
            <span className="text-xs font-medium text-muted-foreground uppercase">
              {language}
            </span>
          </div>
          <div className="overflow-x-auto">
            <pre className="p-4 !bg-transparent !m-0" {...props}>
              {children}
            </pre>
          </div>
        </div>
      );
    },

    // Inline code only
    code: (props) => {
      const { children, className } = props;
      // If no className, it's inline code
      if (!className) {
        return (
          <code className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
            {children}
          </code>
        );
      }
      // For code blocks, just return the code element (wrapped by pre above)
      return <code className={className}>{children}</code>;
    },

    // Blockquotes
    blockquote: ({ children }) => (
      <blockquote className="mt-6 border-l-4 border-primary/30 pl-6 italic text-muted-foreground">
        {children}
      </blockquote>
    ),

    // Tables
    table: ({ children }) => (
      <div className="my-6 w-full overflow-y-auto">
        <table className="w-full border-collapse border border-border">
          {children}
        </table>
      </div>
    ),
    thead: ({ children }) => (
      <thead className="bg-muted">
        {children}
      </thead>
    ),
    tbody: ({ children }) => (
      <tbody>
        {children}
      </tbody>
    ),
    tr: ({ children }) => (
      <tr className="border-b border-border">
        {children}
      </tr>
    ),
    th: ({ children }) => (
      <th className="border border-border px-4 py-2 text-left font-semibold [&[align=center]]:text-center [&[align=right]]:text-right">
        {children}
      </th>
    ),
    td: ({ children }) => (
      <td className="border border-border px-4 py-2 text-left [&[align=center]]:text-center [&[align=right]]:text-right">
        {children}
      </td>
    ),

    // Horizontal rule
    hr: () => (
      <hr className="my-8 border-border" />
    ),

    // Images
    img: ({ src, alt }) => (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={src}
        alt={alt || ""}
        className="rounded-lg border border-border my-6 max-w-full h-auto"
      />
    ),
  };

  return (
    <div className={`prose prose-neutral dark:prose-invert max-w-none ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight, rehypeRaw]}
        components={components}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}

