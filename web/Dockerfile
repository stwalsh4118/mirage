# Build stage
FROM node:20-alpine AS builder

ARG NODE_ENV=production
ARG NEXT_PUBLIC_API_BASE_URL

# Clerk environment variables needed during build
ARG NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY
ARG CLERK_SECRET_KEY

WORKDIR /app

# Copy package files from web directory (build context is repo root)
COPY web/package.json web/pnpm-lock.yaml ./

# Install pnpm
RUN npm install -g pnpm

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy web source code
COPY web/ ./

# Build the application
RUN pnpm build

# Runtime stage
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Copy standalone build output
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Expose web port
EXPOSE 3000

# Run the standalone server
CMD ["node", "server.js"]

