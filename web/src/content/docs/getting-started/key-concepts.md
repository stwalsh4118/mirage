
Understanding these core concepts will help you work more effectively with Mirage and Railway. This guide explains the fundamental building blocks and terminology you'll encounter.

## Railway Concepts

### Projects

**A Railway project is the top-level container for your application.**

- Groups related services and environments
- Represents one application or product
- Contains deployments, services, and settings
- Has its own billing and resource allocation

**Example Projects**:
- `my-ecommerce-app`: Contains frontend, API, database
- `company-website`: Contains web server and CMS
- `microservices-platform`: Contains multiple independent services

```
Railway Project: "my-ecommerce-app"
│
├─ Environment: production
│  ├─ Service: frontend
│  ├─ Service: api
│  └─ Service: postgres
│
└─ Environment: staging
   ├─ Service: frontend
   ├─ Service: api
   └─ Service: postgres
```

### Environments

**An environment is an isolated instance of your application.**

Think of environments as separate copies of your application that can run simultaneously without interfering with each other.

**Common Environments**:
- **Development (dev)**: For active development and testing
- **Staging**: Pre-production testing and QA
- **Production (prod)**: Live application serving real users
- **Preview**: Temporary environments for pull requests

**Key Characteristics**:
- Each environment has its own services
- Separate environment variables per environment
- Independent deployments and resources
- Can share a Railway project

### Services

**A service is a single deployable component of your application.**

Services can be:

- **Application Services**: Your code (e.g., API, web frontend)
- **Database Services**: PostgreSQL, MySQL, MongoDB
- **Cache Services**: Redis, Memcached
- **Other Services**: Elasticsearch, RabbitMQ, etc.

**Service Properties**:
- Unique name within an environment
- Source code or Docker image
- Build and start commands
- Port configuration
- Environment variables

```typescript
// Example: API Service
{
  name: "api",
  source: "github.com/user/app",
  buildCommand: "npm run build",
  startCommand: "npm start",
  port: 3000,
  variables: {
    NODE_ENV: "production",
    DATABASE_URL: "\${DATABASE_URL}"
  }
}
```

### Deployments

**A deployment is a specific instance of your service running.**

- Created when you push new code or trigger a manual deploy
- Has a unique deployment ID
- Includes build logs and runtime logs
- Can be rolled back if needed

**Deployment States**:
- **Building**: Code is being built
- **Deploying**: Service is starting
- **Active**: Service is running successfully
- **Failed**: Deployment encountered an error
- **Crashed**: Service started but crashed

## Mirage Concepts

### Environment Types

**Environments can be labeled with a type to indicate their purpose.**

Common environment types:
- **dev**: Development environments for active coding
- **staging**: Pre-production testing and QA
- **prod**: Production environments serving real users
- **ephemeral**: Temporary environments (feature branches, testing)

These types are primarily labels that help you organize and identify your environments. Railway and Mirage use them for categorization and display purposes.

### Clone Functionality

**Mirage allows you to clone existing environments.**

Cloning creates a new environment by copying:
- Source repository and branch settings
- Service configurations
- Environment variables (excluding Railway system variables)
- Basic environment structure

This is useful for:
- Creating staging environments from development
- Duplicating environments for different branches
- Quickly setting up similar configurations

### Environment Cards

**Visual representation of your environments in the Mirage dashboard.**

Each card displays:
- Environment name and type
- Current status (Active, Creating, Error)
- Service count and health
- Quick action buttons
- Last deployment time
- Resource usage indicators

```
┌────────────────────────────┐
│ 🟢 my-app-dev             │
│ Development Environment    │
├────────────────────────────┤
│ Services: 3/3 running      │
│ Status: Active             │
│ Updated: 5 minutes ago     │
├────────────────────────────┤
│ [View] [Edit] [More ▼]    │
└────────────────────────────┘
```

### Configuration Management

**How Mirage handles settings and variables.**

#### Environment Variables

Variables that configure your application:

- **Application Config**: `NODE_ENV`, `LOG_LEVEL\
- **Secrets**: API keys, passwords, tokens
- **Service URLs**: Database connections, API endpoints
- **Feature Flags**: Enable/disable features

**Variable Scopes**:
- **Service-level**: Specific to one service
- **Environment-level**: Shared across services in an environment
- **Project-level**: Shared across all environments (Railway feature)

#### Shared Variables

Variables that reference Railway-provided values:

```bash
# Railway automatically provides these
DATABASE_URL=postgresql://user:pass@host:5432/db
REDIS_URL=redis://host:6379
PORT=3000

# Reference them in your services
API_URL=\${RAILWAY_SERVICE_URL_API}
```

## Workflows

### Environment Creation Workflow

The standard flow for creating a new environment:

```
1. Select Project
   ↓
2. Choose Template
   ↓
3. Configure Services
   ↓
4. Set Variables
   ↓
5. Review
   ↓
6. Create & Deploy
   ↓
7. Monitor
   ↓
8. Verify
```

### Deployment Workflow

How code changes reach your environment:

```
Code Pushed → GitHub
     ↓
Railway Webhook Triggered
     ↓
Build Starts (Railway)
     ↓
Build Completes
     ↓
Service Deploys
     ↓
Health Check
     ↓
Traffic Switches
     ↓
Old Instance Stops
```

Mirage monitors this process and displays status updates in real-time.

### Environment Lifecycle

Typical lifecycle of an environment:

1. **Creation**: Environment is set up
2. **Active Development**: Services are deployed and updated frequently
3. **Stable**: Few changes, primarily monitoring
4. **Deprecated**: No longer actively used
5. **Cleanup**: Services stopped, environment removed

## Resource Management

### Resource Allocation

Railway allocates resources to each service:

- **CPU**: Processing power
- **Memory**: RAM allocation
- **Storage**: Disk space
- **Network**: Bandwidth and data transfer

**Resource Tiers** (typical):
- **Small**: 0.5 vCPU, 512 MB RAM
- **Medium**: 1 vCPU, 1 GB RAM
- **Large**: 2 vCPU, 2 GB RAM
- **Custom**: Configure as needed

### Usage and Billing

Understanding Railway costs:

- **Free Tier**: Limited monthly usage credits
- **Pay-as-you-go**: Charged for resources used
- **Metered Services**: Databases, Redis, etc.
- **Execution Time**: Running services are billed by time

**Cost Optimization**:
- Stop unused environments
- Right-size service resources
- Use shared resources where appropriate
- Monitor and adjust based on usage

## Status and Health

### Service Status

**Service Health States**:
- 🟢 **Active**: Running normally
- 🟡 **Creating**: Being provisioned
- 🟡 **Building**: Code is building
- 🟡 **Deploying**: Service is starting
- 🔴 **Error**: Failed to start
- 🔴 **Crashed**: Started but crashed
- ⚪ **Stopped**: Intentionally stopped

### Health Checks

Railway performs health checks to ensure services are responding:

```typescript
// Health check configuration
{
  path: "/health",
  interval: 30,  // seconds
  timeout: 10,   // seconds
  retries: 3
}
```

Mirage displays health check results in real-time.

## Best Practices

### Naming Conventions

Use clear, consistent names:

```bash
# Good naming examples
project: company-product-name
environments: production, staging, dev
services: api, web, worker, postgres

# Be specific when needed
services: 
  - user-api
  - payment-api
  - admin-dashboard
  - customer-web
```

### Environment Separation

Keep environments truly separate:

- **Isolated Data**: Never share databases between prod and dev
- **Different Variables**: Each environment has appropriate config
- **Independent Deployments**: Changes don't affect other environments
- **Resource Boundaries**: Appropriate sizing for each environment

### Configuration Management

- **Version Control**: Store configuration in Git (excluding secrets)
- **Environment Parity**: Keep staging close to production config
- **Secret Management**: Use Railway's secret variables
- **Documentation**: Document non-obvious configuration choices

### Monitoring and Maintenance

- **Regular Checks**: Review environment health daily
- **Log Monitoring**: Watch for errors and warnings
- **Resource Monitoring**: Track usage and optimize
- **Cleanup**: Remove unused environments promptly

## Common Patterns

### Multi-Service Application

Typical setup for a full-stack application:

```
Environment: production
│
├─ Frontend (Next.js)
│  └─ Variables: API_URL, PUBLIC_KEY
│
├─ Backend API (Node.js)
│  └─ Variables: DATABASE_URL, JWT_SECRET
│
├─ Background Worker (Node.js)
│  └─ Variables: DATABASE_URL, QUEUE_URL
│
├─ PostgreSQL (Database)
│  └─ Provides: DATABASE_URL
│
└─ Redis (Cache)
   └─ Provides: REDIS_URL
```

### Development to Production Flow

Standard progression:

1. **Develop** in local environment or dev Railway environment
2. **Push** code to GitHub
3. **Deploy** to dev environment (automatic)
4. **Test** thoroughly in dev
5. **Deploy** to staging (manual or automatic)
6. **QA** in staging environment
7. **Deploy** to production (manual, after approval)
8. **Monitor** production deployment

## Terminology Quick Reference

| Term | Definition |
|------|------------|
| **Project** | Top-level container for your application |
| **Environment** | Isolated instance of your application |
| **Service** | Single deployable component |
| **Deployment** | Specific running instance of a service |
| **Template** | Pre-configured environment blueprint |
| **Variables** | Configuration values for services |
| **Build** | Process of preparing code to run |
| **Health Check** | Automated service health verification |
| **Resource** | CPU, memory, or other allocated capacity |

## Next Steps

Now that you understand the key concepts, you can:

- **Explore Features**: Learn about specific Mirage capabilities
- **Read How-To Guides**: Accomplish specific tasks
- **Review Best Practices**: Optimize your workflow
- **Experiment**: Try creating different environment types

### Recommended Reading

- [Railway Integration Overview](/docs/features/railway/overview)
- [Environment Management](/docs/features/environments/overview)
- [Dashboard Guide](/docs/features/dashboard/overview)
- [Templates Deep Dive](/docs/features/environments/templates)

---

> **💡 Pro Tip**: Understanding these concepts deeply will help you troubleshoot issues faster and design better infrastructure. Take time to experiment with different configurations!
