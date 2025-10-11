
This step-by-step guide will walk you through creating your first environment using Mirage. By the end, you'll have a fully configured environment running on Railway.

## What You'll Create

In this guide, we'll create an environment with:

- A Railway project (new or existing)
- Services deployed from a GitHub repository or Docker image
- Environment variables configured
- Automatic deployment to Railway

The process takes about 5-10 minutes to complete.

## Step 1: Start the Environment Wizard

### Access the Creation Wizard

1. From the Mirage dashboard, click the **Create Environment** button
   - Located in the controls bar near the top of the page
2. The environment creation wizard dialog will open

### Wizard Overview

The wizard has 5 steps:

1. **Project**: Choose a Railway project (existing or new) or clone an existing environment
2. **Source**: Select your deployment source (GitHub repository or Docker image)
3. **Discovery**: Automatically detect services in your repository (for monorepos)
4. **Config**: Set environment name and environment variables
5. **Review**: Confirm your configuration and create the environment

## Step 2: Choose Your Project

### Option A: Create From Scratch

If you select "Create from scratch":

#### Select Existing or Create New Project

**Choose an Existing Project:**
```
┌─────────────────────────────────────┐
│ Railway Project                     │
├─────────────────────────────────────┤
│ ○ Use existing project              │
│   [Select: my-backend-api      ▼]   │
│                                      │
│ ○ Create new project                │
│   [Project name: ____________]      │
└─────────────────────────────────────┘
```

1. Select "Use existing project"
2. Choose a project from the dropdown list
3. You'll see the project name displayed
4. Optionally, provide a default environment name

**Create a New Project:**

1. Select "Create new project"
2. Enter a project name (e.g., "my-app", "demo-project")
3. Optionally, provide a default environment name
4. Mirage will create the Railway project for you during provisioning

> **💡 Tip**: For your first environment, using an existing test project is recommended. You can always create new projects later.

### Option B: Clone From Existing Environment

If you select "Clone from existing environment":

1. A dropdown will show all your existing environments across all projects
2. Each environment shows: project name, environment name, and service count
3. Select the environment you want to clone
4. The wizard will pre-populate settings from the cloned environment:
   - Repository URL and branch (or Docker image)
   - Service configurations
   - Environment variables (excluding Railway system variables)

> **⚡ Cloning Tip**: This is the fastest way to create similar environments (e.g., creating a staging environment from dev).

## Step 3: Configure Your Source

Choose how you want to deploy your application:

### Option A: GitHub Repository

```
┌─────────────────────────────────────┐
│ Deployment Source                   │
├─────────────────────────────────────┤
│ ● Repository                         │
│                                      │
│ Repository URL:                      │
│ [https://github.com/user/repo]      │
│                                      │
│ Branch:                              │
│ [main                          ▼]   │
│                                      │
│ GitHub Token (optional):             │
│ [For private repositories]          │
└─────────────────────────────────────┘
```

1. Select "Repository" as your source
2. Enter the full GitHub repository URL
   - Example: `https://github.com/yourusername/your-repo\
3. Specify the branch to deploy (defaults to "main")
4. If it's a private repository, provide a GitHub personal access token

**Repository Requirements:**
- Must be accessible (public or you provide a token)
- Should contain Dockerfiles for your services
- Monorepo structure is supported

### Option B: Docker Image

```
┌─────────────────────────────────────┐
│ Deployment Source                   │
├─────────────────────────────────────┤
│ ○ Docker Image                       │
│                                      │
│ Registry:                            │
│ [docker.io                     ▼]   │
│                                      │
│ Image Name:                          │
│ [nginx                          ]   │
│                                      │
│ Tag or Digest:                       │
│ ● Tag:  [latest              ]      │
│ ○ Digest: [sha256:...]              │
│                                      │
│ Exposed Ports:                       │
│ [80, 443                        ]   │
└─────────────────────────────────────┘
```

1. Select "Docker Image" as your source
2. Choose the registry (Docker Hub, GitHub Container Registry, etc.)
3. Enter the image name
4. Specify either:
   - **Tag**: e.g., `latest`, `v1.0.0`, `production\
   - **Digest**: SHA256 hash for immutable deployments
5. List exposed ports (comma-separated)

## Step 4: Service Discovery (Monorepo)

If you selected a GitHub repository, Mirage will scan for services:

### Automatic Detection

```
Scanning repository for services...

✓ Found 3 services:

┌────────────────────────────────────┐
│ ☑ web-frontend                     │
│   Path: /apps/frontend             │
│   Dockerfile: apps/frontend/Dockerfile
│   Ports: 3000                       │
│                                     │
│ ☑ api-backend                      │
│   Path: /apps/api                  │
│   Dockerfile: apps/api/Dockerfile  │
│   Ports: 8080                       │
│                                     │
│ ☐ worker-service                   │
│   Path: /apps/worker               │
│   Dockerfile: apps/worker/Dockerfile
│   (No exposed ports)                │
└────────────────────────────────────┘
```

**What This Step Does:**
- Scans your repository for Dockerfiles
- Detects build context paths
- Identifies exposed ports from Dockerfiles
- Lists all discovered services

**Your Actions:**
1. Review the discovered services
2. Check or uncheck services to deploy
3. Optionally rename services (click on the name to edit)
4. Services are selected by default if they have exposed ports

### Skip Discovery

If discovery fails or you want to configure manually:

1. Click "Skip Discovery"
2. You'll configure a single generic service in the Config step
3. You can manually specify Dockerfile path and build context later

> **Monorepo Note**: Discovery works best with conventional monorepo structures where each service has its own Dockerfile in its directory.

## Step 5: Configure Environment

Set your environment name and variables:

### Environment Name

```
Environment Name: [my-dev-environment]
```

- Provide a descriptive name for your environment
- Examples: "development", "staging", "feature-auth", "production"
- If you created a new project, this becomes the environment within that project
- If you used an existing project, this is added as a new environment

### Global Environment Variables

Variables that apply to all services:

```
┌─────────────────────────────────────┐
│ Global Environment Variables        │
├─────────────────────────────────────┤
│ Key               Value              │
│ NODE_ENV          development        │
│ LOG_LEVEL         debug              │
│ API_URL           https://...        │
│ [add new...]      [value...]         │
└─────────────────────────────────────┘
```

**Tips:**
- Use the "Import .env" button to paste environment file content
- Variables are automatically parsed from .env format
- Common variables: NODE_ENV, DATABASE_URL, API_KEY, etc.

### Per-Service Variables (Monorepo)

If you discovered multiple services, set service-specific variables:

```
┌─────────────────────────────────────┐
│ Service: web-frontend               │
├─────────────────────────────────────┤
│ Key               Value              │
│ PORT              3000               │
│ NEXT_PUBLIC_API   https://...       │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ Service: api-backend                │
├─────────────────────────────────────┤
│ Key               Value              │
│ PORT              8080               │
│ JWT_SECRET        secret-key-here    │
└─────────────────────────────────────┘
```

**How It Works:**
- Select a service from the dropdown
- Add variables specific to that service
- Variables are scoped to only that service
- Railway automatically provides some variables (DATABASE_URL, REDIS_URL, etc.)

## Step 6: Review and Create

Final confirmation screen:

```
┌─────────────────────────────────────┐
│ Review Your Configuration           │
├─────────────────────────────────────┤
│ Project: my-backend-api (existing)  │
│ Environment: development            │
│                                      │
│ Source: Repository                  │
│ • github.com/user/my-app           │
│ • Branch: main                      │
│                                      │
│ Services (2):                        │
│ • web-frontend                      │
│   - 1 environment variable           │
│ • api-backend                       │
│   - 2 environment variables          │
│                                      │
│ Global Variables: 3                 │
│                                      │
│ [Back] [Create Environment]         │
└─────────────────────────────────────┘
```

### Verification Checklist

Before clicking "Create Environment":

- ✅ Correct project selected or new project name is good
- ✅ Repository URL and branch are correct (or Docker image)
- ✅ Services you want to deploy are selected
- ✅ Environment name is descriptive
- ✅ Required environment variables are set
- ✅ No sensitive data exposed in variable names

### Create and Monitor

1. Click **Create Environment**
2. The wizard switches to a progress view
3. Watch real-time provisioning progress:

```
Creating your environment...

✓ Create Project (if new)
✓ Create Environment
⏳ Create Services
   └─ web-frontend: Building...
   └─ api-backend: Queued...
```

**Stages:**
1. **Create Project**: Creates Railway project (if you chose "new")
2. **Create Environment**: Creates environment within the project
3. **Create Services**: Provisions each service sequentially
   - Uploads Dockerfile path and build context to Railway
   - Sets environment variables
   - Triggers initial deployment

This typically takes 2-5 minutes depending on:
- Number of services
- Build complexity
- Railway platform load

## Step 7: Success and Next Steps

When provisioning completes:

```
✅ Environment Created Successfully!

Your environment "development" is now live!

Services (2/2 running):
• web-frontend
• api-backend

[View in Railway] [Go to Dashboard] [Create Another]
```

### Access Your Environment

1. **View in Railway**: Opens Railway console for this environment
2. **Go to Dashboard**: Returns to Mirage dashboard
3. **Create Another**: Start wizard again

### Verify Deployment

In Railway:
- Check service deployment status
- View build logs
- Access service URLs once deployed
- Monitor resource usage

## Troubleshooting

### Repository Not Found

**Problem**: Can't access the repository

**Solutions**:
- Verify repository URL is correct and accessible
- For private repos, ensure GitHub token has correct permissions
- Token needs `repo` scope for private repositories

### Service Discovery Found Nothing

**Problem**: No services were discovered

**Solutions**:
- Ensure repository contains Dockerfiles
- Check Dockerfile naming (must be named `Dockerfile` or `*.dockerfile`)
- Use "Skip Discovery" and configure manually if needed
- For non-monorepo projects, skip discovery

### Build Failures

**Problem**: Services fail to build on Railway

**Solutions**:
- Check Dockerfile syntax is valid
- Ensure all dependencies are specified
- Verify build context path is correct
- Check Railway service logs for specific errors

### Environment Variables Not Working

**Problem**: Application can't access variables

**Solutions**:
- Verify variable names match what your application expects
- Check for typos in variable names
- Railway auto-injects some variables - don't override them
- Restart services after variable changes

### Clone Didn't Copy Everything

**Problem**: Cloned environment is missing configurations

**Solutions**:
- Railway system variables are intentionally excluded (RAILWAY_*)
- Check that original environment's variables were properly set
- Some settings may need to be reconfigured manually

## What's Next?

### Congratulations! 🎉

You've successfully created your first environment with Mirage!

### Continue Learning

Now explore:

- **Key Concepts**: Understand Railway and Mirage terminology
- **Dashboard Guide**: Learn to navigate and manage your environments
- **How-To Guides**: Accomplish specific tasks
- **Troubleshooting**: Solve common issues

### Try This Next

1. Clone your newly created environment
2. Create a new environment from a different branch
3. Add more services to your environment
4. Modify environment variables and redeploy

---

> **🎓 Pro Tip**: The clone feature is great for creating staging/production environments from your dev setup. Just clone and update the branch and any environment-specific variables!
