# Mirage Theme System

This document defines the Mirage visual language for UI implementation across the app. It guides component styling (via shadcn), CSS variables, and Tailwind utilities.

## Theme identity
- Mood: warm, dusty, minimal, translucent; refractive shimmer accents evocative of heat haze.
- Surfaces: layered glass panels over blurred dune gradients; subtle grain/noise texture.

## Palette and tokens (light / dark)
- Base neutrals ("Sand" → 50–950): pale ivory, dune, khaki, cumin, burnt umber, char.
- Brand: Primary ochre (sun-warmed); Accent teal (mirage shimmer).
- Semantic: Success agave; Warning desert orange; Destructive rust; Info sky haze.

CSS variables (mapped for shadcn):
- `--background`: sand-50 (dark: umber/char gradient)
- `--foreground`: umber-900 (dark: sand-50)
- `--card`: sand-100 with translucency (dark: lighten over dark bg)
- `--muted`: sand-200
- `--primary`: ochre-500; `--primary-foreground`: umber-950
- `--accent`: teal-400; `--accent-foreground`: umber-950
- `--border`: sand-300 at ~50% opacity
- `--ring`: teal-400 at ~50% opacity
- `--destructive`: rust-500; `--success`: agave-500; `--warning`: desert-500

## Surfaces and glass treatment
- Glass panels: translucent background (`rgb(var(--card)/0.65)`) + `backdrop-blur-md` + warm border.
- Highlight edge: subtle light top edge; soft warm shadow below for depth.
- Noise/grain: faint texture mask at ~4–6% opacity on major surfaces.
- Haze shimmer: low-alpha radial gradients; optional animated displacement on hover.

## Typography
- Typeface: Geist (keep); headings semi-bold with slightly increased letter spacing; generous line-height for body.
- Color: warm foregrounds; ensure minimum 4.5:1 contrast for text.

## Motion and shimmer
- Card hover: scale 1.01, translateY(-1px), and a slow “sheen” sweep (8–12s, staggered).
- Focus: accent ring glow using `--ring`.

## Elevation
- Levels: flat, raised, floating. Use warm translucent shadows (e.g., `0 4px 16px rgb(65 44 20 / 0.15)`).

## Component guidelines (shadcn)
- Button: glassy ochre default; ghost uses dune border with warm hover; destructive uses rust.
- Card: glass base with grain; header underline accent in ochre.
- Badge: Dev = dune with teal text; Prod = solid ochre with umber text.
- Dropdown/Tooltip: glass panels; arrow blurred to match.
- Dialog/AlertDialog: blurred warm gradient backdrop; elevated glass panel.
- Toast: glass chip with accent ring; semantic colored left border.
- Skeleton: warm sand shimmer; avoid cold gray.

## Status colors and badges
- Active: agave; Creating: ochre pulse; Destroying: rust pulse; Error: rust; Unknown: dune/teal.
- Badges: dot + label, with subtle outer glow by state.

## Implementation notes
- Define `:root` and `.dark` CSS variables for tokens above in `web/src/app/globals.css`.
- Add utilities: `.glass` (bg, border, blur), `.grain` (noise), `.sheen` (animated gradient).
- Map variables to shadcn theme tokens (`--background`, `--primary`, etc.) so UI kit inherits theme.
- Respect contrast requirements; raise opacity in dark mode as needed.

