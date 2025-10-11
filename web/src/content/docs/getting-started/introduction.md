# Introduction to Mirage

Welcome to Mirage! Mirage is a platform for managing Railway environments with ease and efficiency.

## What is Mirage?

Mirage is an environment management tool built specifically for Railway.app users. It provides a streamlined interface for creating and managing Railway environments, with features designed to simplify common workflows.

Think of Mirage as your control panel for Railway environments—it provides an intuitive wizard for environment creation, project visualization, and basic environment management to help you work more efficiently with Railway's platform.

## Why Use Mirage?

### Streamlined Environment Creation

Creating new environments on Railway through the wizard can involve multiple steps. Mirage simplifies this process through:

- **Guided Wizard**: Step-by-step environment creation with clear steps
- **Clone Environments**: Duplicate existing environments with all configurations
- **Monorepo Support**: Automatic service discovery in monorepo repositories
- **Project Management**: Create new Railway projects or use existing ones

### Unified Dashboard

Instead of navigating between multiple Railway projects individually, Mirage provides a dashboard where you can:

- View all your Railway projects in one place
- See environments and services for each project
- Access project details quickly
- Switch between grid and table views

### Environment Management

Mirage helps with basic environment operations:

- **Environment Variables**: Manage global and per-service environment variables
- **Service Discovery**: Automatically detect services in monorepo repositories
- **Deployment Sources**: Configure from GitHub repositories or Docker images
- **Clone Functionality**: Duplicate environments with their configurations

## Who Should Use Mirage?

Mirage is useful for:

- **Development Teams** building applications on Railway who manage multiple projects
- **Developers** who want a centralized view of their Railway infrastructure
- **Teams** using monorepo architectures with multiple deployable services
- **Anyone** who frequently creates similar Railway environments

## Core Features

### Railway Integration

Mirage integrates with Railway's GraphQL API to provide:

- Direct access to your Railway projects
- Service and environment information
- Environment creation and configuration
- Real-time project data

### Environment Creation Wizard

A multi-step wizard guides you through:

1. **Project Selection**: Choose an existing Railway project or create a new one
2. **Source Configuration**: Select a GitHub repository with branch, or Docker image
3. **Service Discovery**: Automatically detect services in monorepo repositories (finds Dockerfiles)
4. **Configuration**: Set environment name and environment variables (global and per-service)
5. **Review**: Confirm settings and provision your environment

### Clone Environments

Duplicate existing environments:

- Copy service configurations
- Clone environment variables (excluding Railway system variables)
- Maintain repository and branch settings
- Quick way to create similar environments

### Dashboard

A clean dashboard interface provides:

- Project cards showing all your Railway projects
- Environment and service counts
- Grid or table view options
- Direct links to project details

## How It Works

Mirage acts as a layer on top of Railway's infrastructure:

```
┌─────────────────┐
│     Mirage      │  ← User Interface & Workflow Management
│   (Frontend)    │
└────────┬────────┘
         │
         │ GraphQL API
         │
┌────────▼────────┐
│   Railway API   │  ← Infrastructure & Deployment
│   (Platform)    │
└─────────────────┘
```

1. **You** interact with Mirage's interface
2. **Mirage** translates your actions into Railway API calls
3. **Railway** provisions and manages the actual infrastructure
4. **Mirage** displays the results and provides access to your resources

## Getting Started

Ready to dive in? Follow these steps:

1. **Check Prerequisites**: Ensure you have a Railway account and workspace token
2. **Complete Setup**: Connect your Railway account to Mirage
3. **Create Your First Environment**: Follow our guided walkthrough
4. **Explore Features**: Learn about the wizard and dashboard

## What's Next?

Continue with the getting started guide to begin using Mirage:

- **Prerequisites**: What you need before starting
- **Setup**: Connect your Railway account
- **First Environment**: Create your first environment
- **Key Concepts**: Understand core terminology

---

> **Need Help?** Visit our [troubleshooting guide](/docs/troubleshooting) or check out the [Railway documentation](https://railway.app/docs) for Railway-specific questions.

