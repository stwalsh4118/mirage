import { MarkdownRenderer } from "@/components/docs/MarkdownRenderer";

const sampleMarkdown = `# Markdown Test Page

This page demonstrates all markdown features with syntax highlighting and custom styling.

## Headings

You can use headings from h1 to h6. This is an h2 heading above.

### This is an h3 heading

#### This is an h4 heading

## Text Formatting

This is a paragraph with **bold text**, *italic text*, and \`inline code\`.

You can also have ~~strikethrough text~~ and [links to other pages](/docs).

External links work too: [Railway Documentation](https://railway.app/docs).

## Lists

### Unordered List

- First item
- Second item
  - Nested item
  - Another nested item
- Third item

### Ordered List

1. First step
2. Second step
3. Third step

## Code Blocks

Here's a TypeScript code block:

\`\`\`typescript
interface Environment {
  id: string;
  name: string;
  status: "active" | "inactive" | "error";
  services: Service[];
}

function createEnvironment(name: string): Environment {
  return {
    id: generateId(),
    name,
    status: "active",
    services: [],
  };
}
\`\`\`

And here's a bash example:

\`\`\`bash
# Install dependencies
pnpm install

# Run the development server
pnpm dev

# Build for production
pnpm build
\`\`\`

## Blockquotes

> This is a blockquote. It's great for highlighting important information or quotes.
>
> You can have multiple paragraphs in a blockquote.

## Tables

| Feature | Status | Description |
| :------ | :----: | :---------- |
| Railway Integration | ✅ | Connect to Railway API |
| Environment Creation | ✅ | Create new environments |
| Service Management | ✅ | Manage services |
| Documentation | 🚧 | In progress |

## Horizontal Rule

---

## Images

Images will be styled with rounded corners and borders:

![Mirage Logo](/mirage_logo.png)

## Task Lists

- [x] Set up routing
- [x] Integrate markdown renderer
- [ ] Write documentation
- [ ] Add search functionality

## Complex Example

Here's a more complex example showing a configuration file:

\`\`\`json
{
  "name": "my-environment",
  "template": "production",
  "services": [
    {
      "name": "api",
      "image": "node:20-alpine",
      "env": {
        "NODE_ENV": "production",
        "PORT": "3000"
      }
    }
  ]
}
\`\`\`

## Inline HTML

<div style="padding: 1rem; background: rgba(var(--primary), 0.1); border-radius: 0.5rem; margin: 1rem 0;">
  This is custom HTML content that can be embedded in markdown.
</div>

That's all for the test!
`;

export default function MarkdownTestPage() {
  return (
    <div className="space-y-6">
      <div className="glass grain rounded-lg p-6 border border-border">
        <h1 className="text-2xl font-bold mb-2">Markdown Renderer Test</h1>
        <p className="text-muted-foreground">
          This page tests all markdown features including syntax highlighting, tables, and more.
        </p>
      </div>

      <div className="glass grain rounded-lg p-8 border border-border">
        <MarkdownRenderer content={sampleMarkdown} />
      </div>
    </div>
  );
}

