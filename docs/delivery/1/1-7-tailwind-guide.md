# 1-7 Tailwind CSS v4 Guide

Date: 2025-09-06

[Official Docs](https://tailwindcss.com/docs)

## Overview
The project uses Tailwind CSS v4 (via `@tailwindcss/postcss`) configured by create-next-app. Global tokens are defined in `src/app/globals.css` using `@theme inline`.

## Configuration
Key files created by the scaffold:

- `src/app/globals.css`

```css
@import "tailwindcss";

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --font-sans: var(--font-geist-sans);
  --font-mono: var(--font-geist-mono);
}
```

No `tailwind.config.js` is required in v4 for basic usage; the PostCSS plugin auto-detects.

## Usage
Utilities are available globally. Example navbar and container styles live in `src/app/layout.tsx`.

```tsx
<header className="border-b border-black/10 dark:border-white/15">
  <nav className="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
    <div className="text-lg font-semibold">Mirage</div>
  </nav>
 </header>
```

## Recommendations
- Prefer semantic spacing `px-4 py-3`, `max-w-6xl`, and color opacity utilities.
- Keep design tokens in `@theme inline` and reuse via class utilities.
- Avoid custom CSS unless necessary; prefer Tailwind utilities for consistency.


