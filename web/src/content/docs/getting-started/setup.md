
This guide will walk you through connecting your Railway account to Mirage and getting everything ready to create your first environment.

## Step 1: Access Mirage

### Navigate to Mirage

1. Open your web browser
2. Navigate to the Mirage application (URL depends on deployment)
3. You should see the Mirage landing page

### First-Time Access

On your first visit, you'll see:

- Welcome message
- Brief overview of features
- Sign-in or get started options

## Step 2: Sign In to Mirage

### Authentication Options

Mirage may support multiple authentication methods:

- **Email/Password**: Create a Mirage account
- **OAuth**: Sign in with GitHub or other providers
- **SSO**: Enterprise single sign-on (if configured)

### Create Your Account

If you're new to Mirage:

1. Click **Sign Up** or **Get Started**
2. Choose your preferred authentication method
3. Complete the registration process
4. Verify your email if required
5. Log in to Mirage

## Step 3: Connect Your Railway Account

Once logged in to Mirage, you need to connect it to your Railway account.

### Navigate to Settings

1. Click on your profile icon or menu
2. Select **Settings** or **Account Settings**
3. Find the **Integrations** or **Railway Connection** section

### Add Railway Workspace Token

```bash
# Your Railway workspace token (from Prerequisites step)
railway_abc123def456ghi789jkl012mno345pqr678
```

1. Click **Connect Railway Account** or **Add Workspace Token**
2. Paste your Railway workspace token into the input field
3. Optionally, give this connection a name (e.g., "Personal Account")
4. Click **Connect** or **Save**

> **Important**: Make sure you're using a **Workspace token**, not a Project or Personal token, as only workspace tokens have the necessary permissions to manage environments across all your projects.

### Verify Connection

Mirage will test the connection by:

1. Validating the token format
2. Making a test API call to Railway
3. Fetching your basic account information
4. Displaying a success or error message

✅ **Success**: You should see "Connected to Railway" with your account details

❌ **Error**: Check the troubleshooting section below if connection fails

## Step 4: Verify Project Access

After connecting your Railway account, verify access to your projects:

1. Navigate to **Dashboard** in the main navigation
2. You should see a list of your Railway projects
3. Each project card shows:
   - Project name
   - Number of environments
   - Number of services
4. Verify the projects you expect to manage are visible

If no projects appear:
- Check your workspace token is correctly configured
- Ensure the workspace has projects
- Try refreshing the page
- Check browser console for any errors

## Step 5: Explore the Dashboard

Let's familiarize yourself with Mirage's interface.

### Dashboard Overview

The main dashboard shows:

- **Environment Cards**: All your managed environments
- **Quick Actions**: Buttons for common tasks
- **Status Indicators**: Health of services
- **Filter/Sort**: Options to organize your view

### Navigation Bar

The top navigation provides access to:

- **Dashboard**: Main view of all projects and environments
- **Settings**: Railway credentials management
- **Documentation**: This documentation
- **User Menu**: Sign out and account options

### Dashboard Views

The dashboard offers two view modes:

- **Grid View**: Project cards with environment details
- **Table View**: Compact table format with sortable columns

You can switch between views using the controls bar.

## Verification Checklist

Before proceeding, verify you've completed:

- ✅ Signed in to Mirage successfully
- ✅ Connected your Railway workspace token
- ✅ Connection test passed
- ✅ Can see your Railway projects
- ✅ Dashboard loads without errors
- ✅ Comfortable with basic navigation

## Troubleshooting Setup Issues

### Can't Sign In to Mirage

**Problem**: Login page shows errors or doesn't respond

**Solutions**:
- Clear browser cache and cookies
- Try a different browser
- Disable browser extensions temporarily
- Check if Mirage service is operational

### Workspace Token Rejected

**Problem**: "Invalid token" or "Connection failed" error

**Solutions**:
- **Verify token type**: Ensure you created a **Workspace** token, not a Project or Personal token
- Verify you copied the entire token without spaces
- Check the token wasn't accidentally truncated
- Ensure the token hasn't been revoked in Railway
- Verify the workspace associated with the token is active
- Create a new workspace token if needed
- Confirm Railway account is active

### No Projects Showing

**Problem**: Project list is empty after connecting

**Solutions**:
- Refresh the page
- Verify you have projects in your Railway account
- Check token permissions (should have read access)
- Wait a moment (initial sync can take a few seconds)
- Check Railway API status

### Connection Timeout

**Problem**: Connection attempt times out

**Solutions**:
- Check your internet connection
- Verify Railway API is accessible
- Try again in a moment (may be temporary)
- Check firewall or proxy settings
- Try from a different network

### Dashboard Not Loading

**Problem**: Dashboard shows loading spinner indefinitely

**Solutions**:
- Refresh the page
- Clear browser cache
- Check browser console for errors (F12)
- Verify Railway connection is still active
- Try signing out and back in

## Advanced Setup Options

### Multiple Railway Accounts

If you manage multiple Railway accounts:

1. Add additional workspace tokens in Settings
2. Give each connection a descriptive name
3. Switch between accounts using the account selector
4. Each environment is linked to its source account

### Team Setup

For team usage:

1. Ensure all team members have Railway access
2. Each member connects their own Railway workspace token
3. Share environment configurations via export/import
4. Use consistent naming conventions

### Automation Setup (Optional)

For advanced users:

- **Webhooks**: Configure webhooks for status updates
- **API Access**: Generate Mirage API keys for automation
- **CI/CD**: Integrate Mirage with deployment pipelines

## Security Best Practices

### Token Management

- **Rotate Regularly**: Change your Railway workspace token every 90 days
- **Monitor Usage**: Check for unexpected API calls
- **Use Workspace Tokens**: Always use workspace tokens, not project or personal tokens
- **Revoke Unused**: Remove old or unused tokens

### Account Security

- **Strong Password**: Use a unique, strong password for Mirage
- **2FA**: Enable two-factor authentication if available
- **Session Management**: Log out on shared computers
- **Review Activity**: Check access logs regularly

## What's Next?

Now that setup is complete, you're ready to create your first environment!

The next guide will walk you through:

- Starting the environment creation wizard
- Selecting a Railway project
- Choosing a template
- Configuring services
- Deploying your first environment

---

> **Setup Complete!** 🎉 Your Railway account is now connected to Mirage. Continue to the next guide to create your first environment.
