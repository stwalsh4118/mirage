# Dashboard Layout Spec

This document describes the dashboard page layout, components, states, and interactions for Mirage.

## Page structure
- Header: translucent app bar with title "Dashboard" and primary action "New Environment".
- Content: responsive grid of Environment Cards; supports loading, empty, and error states.

## Grid and responsiveness
- Columns: 1 (xs) → 2 (sm) → 3 (md) → 4 (lg+)
- Gap: 24px (Tailwind `gap-6`)
- Card min width: ~320px; height auto with consistent paddings

## Environment Card composition
- Top row: name (primary), type badge (Dev/Prod), status chip with colored dot
- Body: createdAt (relative), endpoint URL, brief description (optional)
- Actions:
  - Primary: Open (link to environment URL)
  - Secondary: View logs (route placeholder)
  - Destructive: Destroy (Alert dialog confirm)
  - Overflow: Dropdown for future actions

## States
- Loading: grid of `Skeleton` cards (shadcn)
- Empty: centered callout with "Create environment" button
- Error: toast + inline retry action

## Data and polling
- Source: GET `/environments`
- Cache key: `['environments']`
- Polling: 5s (`ENV_POLL_INTERVAL_MS` constant)
- Destroy: DELETE `/environments/:id` → invalidate/refetch

## A11y & keyboard
- Buttons and links with ARIA labels
- Focus ring uses theme `--ring` token
- Ensure tab order follows visual order

## shadcn components
- Button, Card, Badge, Dropdown Menu, Tooltip, Alert Dialog, Toast, Skeleton, Dialog

## Page wiring
- `src/app/page.tsx`: renders header + `EnvironmentGrid`
- `src/components/dashboard/EnvironmentGrid.tsx`: data fetching and grid
- `src/components/dashboard/EnvironmentCard.tsx`: card UI and actions
- `src/hooks/useEnvironments.ts`: query and polling
- `src/lib/api/environments.ts`: fetch and destroy helpers
- `src/lib/constants/environments.ts`: tokens and mappings, e.g. `ENV_POLL_INTERVAL_MS = 5000`

## Visual details
- Glass cards using theme utilities (`.glass`, `.grain`, `.sheen`)
- Warm ochre accents; teal ring and shimmer
- Elevation: raised cards on hover (scale 1.01, translateY(-1px))
