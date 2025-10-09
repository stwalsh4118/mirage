# Mirage

Platform for managing ephemeral Railway environments.

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker and Docker Compose (for local Vault)
- Railway account and API token

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mirage
   ```

2. **Start Vault (Optional)**
   ```bash
   docker-compose up -d vault
   ```
   See [Vault Setup Guide](docs/VAULT_SETUP.md) for detailed instructions.

3. **Configure API Environment**
   
   Create `api/.env` file:
   ```bash
   # Clerk Authentication
   CLERK_SECRET_KEY=your_clerk_secret_key
   CLERK_WEBHOOK_SECRET=your_webhook_secret
   
   # Railway API (fallback when Vault disabled)
   RAILWAY_API_TOKEN=your_railway_token
   
   # Vault Configuration (if using docker-compose)
   VAULT_ENABLED=true
   VAULT_ADDR=http://localhost:8200
   VAULT_TOKEN=dev-root-token
   VAULT_SKIP_VERIFY=true
   ```
   
   See [Environment Variables Documentation](docs/ENVIRONMENT_VARIABLES.md) for all options.

4. **Start the API**
   ```bash
   cd api
   go run cmd/server/main.go
   ```

5. **Start the Web UI**
   ```bash
   cd web
   npm install
   npm run dev
   ```

6. **Access the Application**
   - Web UI: http://localhost:3000
   - API: http://localhost:8080
   - Vault UI: http://localhost:8200/ui (token: `dev-root-token`)

## Architecture

- **Backend (API):** Go service handling Railway API interactions and secret management
- **Frontend (Web):** Next.js application for user interface
- **Secrets Management:** HashiCorp Vault for secure per-user token storage
- **Authentication:** Clerk for user authentication and management
- **Database:** PostgreSQL for persistence

## Development

### Running Tests

```bash
# API tests
cd api
go test ./...

# Web tests
cd web
npm test
```

### Docker Compose Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f vault

# Stop services
docker-compose down
```

