export interface DocPage {
  title: string;
  path: string;
}

export interface DocSection {
  title: string;
  pages: DocPage[];
}

export const docsNavigation: DocSection[] = [
  {
    title: "Getting Started",
    pages: [
      { title: "Introduction", path: "/docs/getting-started/introduction" },
      { title: "Prerequisites", path: "/docs/getting-started/prerequisites" },
      { title: "Setup", path: "/docs/getting-started/setup" },
      { title: "First Environment", path: "/docs/getting-started/first-environment" },
      { title: "Key Concepts", path: "/docs/getting-started/key-concepts" },
    ],
  },
  {
    title: "Features",
    pages: [
      { title: "Railway Integration", path: "/docs/features/railway" },
      { title: "Environment Management", path: "/docs/features/environments" },
      { title: "Dashboard Overview", path: "/docs/features/dashboard" },
      { title: "Service Management", path: "/docs/features/services" },
    ],
  },
  {
    title: "Railway Integration",
    pages: [
      { title: "Overview", path: "/docs/features/railway/overview" },
      { title: "Connecting Account", path: "/docs/features/railway/connecting-account" },
      { title: "API Tokens", path: "/docs/features/railway/api-tokens" },
      { title: "Browsing Projects", path: "/docs/features/railway/browsing-projects" },
      { title: "Project Details", path: "/docs/features/railway/project-details" },
      { title: "Permissions", path: "/docs/features/railway/permissions" },
      { title: "Troubleshooting", path: "/docs/features/railway/troubleshooting" },
    ],
  },
  {
    title: "Environments",
    pages: [
      { title: "Overview", path: "/docs/features/environments/overview" },
      { title: "Wizard Walkthrough", path: "/docs/features/environments/wizard-walkthrough" },
      { title: "Project Selection", path: "/docs/features/environments/project-selection" },
      { title: "Templates", path: "/docs/features/environments/templates" },
      { title: "Configuration", path: "/docs/features/environments/configuration" },
      { title: "Environment Variables", path: "/docs/features/environments/environment-variables" },
      { title: "Review & Deploy", path: "/docs/features/environments/review-and-deploy" },
      { title: "Progress Tracking", path: "/docs/features/environments/progress-tracking" },
    ],
  },
  {
    title: "Dashboard",
    pages: [
      { title: "Overview", path: "/docs/features/dashboard/overview" },
      { title: "Environment Cards", path: "/docs/features/dashboard/environment-cards" },
      { title: "Quick Actions", path: "/docs/features/dashboard/quick-actions" },
      { title: "Service Management", path: "/docs/features/dashboard/service-management" },
      { title: "Status Indicators", path: "/docs/features/dashboard/status-indicators" },
      { title: "Filtering & Sorting", path: "/docs/features/dashboard/filtering-sorting" },
    ],
  },
  {
    title: "How-To Guides",
    pages: [
      { title: "Create Dev Environment", path: "/docs/how-to/create-dev-environment" },
      { title: "Create Staging Environment", path: "/docs/how-to/create-staging-environment" },
      { title: "Configure Custom Services", path: "/docs/how-to/configure-custom-services" },
      { title: "Manage Environment Variables", path: "/docs/how-to/manage-environment-variables" },
      { title: "Use Templates Effectively", path: "/docs/how-to/use-templates-effectively" },
      { title: "Monitor Environment Health", path: "/docs/how-to/monitor-environment-health" },
      { title: "Clean Up Environments", path: "/docs/how-to/clean-up-environments" },
      { title: "Share Environments", path: "/docs/how-to/share-environments" },
      { title: "Migrate from Railway Direct", path: "/docs/how-to/migrate-from-railway-direct" },
      { title: "Optimize Resource Usage", path: "/docs/how-to/optimize-resource-usage" },
    ],
  },
  {
    title: "Troubleshooting",
    pages: [
      { title: "Overview", path: "/docs/troubleshooting/overview" },
      { title: "Connection Issues", path: "/docs/troubleshooting/connection-issues" },
      { title: "Authentication Errors", path: "/docs/troubleshooting/authentication-errors" },
      { title: "Environment Creation Failures", path: "/docs/troubleshooting/environment-creation-failures" },
      { title: "Service Deployment Issues", path: "/docs/troubleshooting/service-deployment-issues" },
      { title: "Configuration Errors", path: "/docs/troubleshooting/configuration-errors" },
      { title: "Performance Problems", path: "/docs/troubleshooting/performance-problems" },
      { title: "Common Error Messages", path: "/docs/troubleshooting/common-error-messages" },
      { title: "Getting Help", path: "/docs/troubleshooting/getting-help" },
    ],
  },
];

