
Before you start using Mirage, you'll need to have a few things in place. This guide will walk you through everything you need to get started.

## Railway Account

### Required: Active Railway Account

You must have an active Railway.app account to use Mirage. If you don't have one yet:

1. Visit [railway.app](https://railway.app)
2. Sign up using GitHub, email, or another supported authentication method
3. Complete the account verification process
4. Optionally, add a payment method to access additional resources

> **Free Tier Available**: Railway offers a generous free tier that's perfect for getting started with Mirage. You can upgrade later as your needs grow.

### Account Access Level

Ensure you have appropriate access to the Railway projects you want to manage:

- **Owner**: Full access to all project features
- **Admin**: Can create and modify environments
- **Member**: May have limited access depending on team settings

## Railway API Token

### What is an API Token?

An API token is a secure credential that allows Mirage to interact with Railway's API on your behalf. It's like a password specifically for applications.

> **⚠️ Important**: You must create a **Workspace Token**, not a Project or Personal token. Workspace tokens have the necessary permissions to manage environments across all your projects.

### Creating Your Workspace Token

Follow these steps to create a Railway workspace token:

1. Log in to your Railway account at [railway.app](https://railway.app)
2. Click on your profile icon in the top-right corner
3. Select **Account Settings** from the dropdown menu
4. Navigate to the **Tokens** section in the sidebar
5. Click **Create New Token**
6. **Important**: Select **"Workspace"** as the token type (not Project or Personal)
7. Give your token a descriptive name (e.g., "Mirage Integration")
8. Select the appropriate workspace (usually your personal or team workspace)
9. Click **Create** and copy the token immediately

```bash
# Your workspace token will look something like this (example only):
railway_abc123def456ghi789jkl012mno345pqr678
```

> **⚠️ Security Warning**: Copy your workspace token immediately and store it securely. Railway will only show it once. If you lose it, you'll need to create a new one.

### Token Security Best Practices

- **Never** commit API tokens to version control
- **Never** share tokens in screenshots or public forums
- **Do** store tokens in environment variables or secure credential managers
- **Do** rotate tokens periodically (every 90 days recommended)
- **Do** use separate tokens for different applications
- **Do** revoke tokens immediately if compromised

## Browser Requirements

Mirage is a modern web application that works best with up-to-date browsers:

### Recommended Browsers

- **Chrome/Edge**: Version 90 or later
- **Firefox**: Version 88 or later
- **Safari**: Version 14 or later

### Required Browser Features

- JavaScript enabled
- Cookies enabled (for session management)
- LocalStorage support
- Modern ES6+ JavaScript support

## Internet Connection

A stable internet connection is required because Mirage:

- Communicates with Railway's API in real-time
- Monitors service status continuously
- Syncs environment configurations
- Streams deployment logs

Recommended minimum connection speed: 5 Mbps download, 1 Mbps upload

## Knowledge Prerequisites

While Mirage is designed to be user-friendly, some basic knowledge will help you get the most out of it:

### Helpful to Know

- **Railway Basics**: Understanding of Railway projects, services, and environments
- **Environment Variables**: How to configure application settings
- **Service Architecture**: Basic understanding of multi-service applications
- **Command Line**: Basic terminal familiarity (optional, but helpful)

### Not Required, But Useful

- **Docker**: Understanding containerized applications
- **CI/CD**: Continuous integration and deployment concepts
- **GraphQL**: Railway's API uses GraphQL (Mirage handles this for you)

## Optional: Development Tools

If you plan to integrate Mirage into your development workflow:

### Command Line Tools

- **Railway CLI**: For advanced Railway operations
- **Git**: For version control of configurations
- **Docker**: For local development matching Railway environments

### Integrations

Mirage works well alongside:

- **GitHub**: For source code management
- **Vercel**: For frontend deployments
- **Supabase**: For database management
- **Other Railway services**: Postgres, Redis, etc.

## Cost Considerations

### Mirage Costs

**Mirage itself is free to use** (check current pricing at time of use).

### Railway Costs

Be aware of Railway's pricing:

- **Free Tier**: Limited resources, perfect for development
- **Usage-Based**: Pay for what you use beyond free tier
- **Resource Limits**: CPU, memory, and network usage are metered

> **💡 Tip**: Start with the free tier to learn Mirage, then upgrade your Railway plan as needed. Mirage helps you optimize resource usage to keep costs down.

## System Requirements Summary

### Must Have

- ✅ Active Railway account
- ✅ Railway workspace token (not Project or Personal token)
- ✅ Modern web browser
- ✅ Internet connection

### Should Have

- ⚡ Basic Railway knowledge
- ⚡ Understanding of environment variables
- ⚡ Familiarity with service architecture

### Nice to Have

- ⭐ Railway CLI installed
- ⭐ Git for version control
- ⭐ Docker knowledge

## Troubleshooting Prerequisites

### Can't Create Workspace Token?

- Verify your Railway account is fully activated
- Check that you're logged in with the correct account
- Ensure you have access to at least one workspace
- Try a different browser if the token creation page doesn't load

### Workspace Token Not Working?

- **Verify token type**: Ensure you created a **Workspace** token, not a Project or Personal token
- Ensure you copied the entire token without spaces
- Verify the token hasn't been revoked in Railway
- Check the workspace associated with the token is active
- Verify the token has the necessary permissions

### Browser Compatibility Issues?

- Update your browser to the latest version
- Clear browser cache and cookies
- Disable browser extensions that might interfere
- Try a different browser

## Ready to Continue?

Once you have:

- ✅ A Railway account
- ✅ A Railway workspace token ready to use
- ✅ A compatible browser
- ✅ A stable internet connection

You're ready to proceed with the setup!

---

> **Questions?** If you're missing any prerequisites or run into issues, check our [troubleshooting guide](/docs/troubleshooting) or the [Railway documentation](https://railway.app/docs).
