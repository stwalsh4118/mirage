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
      { title: "Railway Integration", path: "/docs/wip" },
      { title: "Environment Management", path: "/docs/wip" },
      { title: "Dashboard Overview", path: "/docs/wip" },
      { title: "Service Management", path: "/docs/wip" },
    ],
  },
  {
    title: "Railway Integration",
    pages: [
      { title: "Overview", path: "/docs/wip" },
      { title: "Connecting Account", path: "/docs/wip" },
      { title: "API Tokens", path: "/docs/wip" },
      { title: "Browsing Projects", path: "/docs/wip" },
      { title: "Project Details", path: "/docs/wip" },
      { title: "Permissions", path: "/docs/wip" },
      { title: "Troubleshooting", path: "/docs/wip" },
    ],
  },
  {
    title: "Environments",
    pages: [
      { title: "Overview", path: "/docs/wip" },
      { title: "Wizard Walkthrough", path: "/docs/wip" },
      { title: "Project Selection", path: "/docs/wip" },
      { title: "Templates", path: "/docs/wip" },
      { title: "Configuration", path: "/docs/wip" },
      { title: "Environment Variables", path: "/docs/wip" },
      { title: "Review & Deploy", path: "/docs/wip" },
      { title: "Progress Tracking", path: "/docs/wip" },
    ],
  },
  {
    title: "Dashboard",
    pages: [
      { title: "Overview", path: "/docs/wip" },
      { title: "Environment Cards", path: "/docs/wip" },
      { title: "Quick Actions", path: "/docs/wip" },
      { title: "Service Management", path: "/docs/wip" },
      { title: "Status Indicators", path: "/docs/wip" },
      { title: "Filtering & Sorting", path: "/docs/wip" },
    ],
  },
  {
    title: "How-To Guides",
    pages: [
      { title: "Create Dev Environment", path: "/docs/wip" },
      { title: "Create Staging Environment", path: "/docs/wip" },
      { title: "Configure Custom Services", path: "/docs/wip" },
      { title: "Manage Environment Variables", path: "/docs/wip" },
      { title: "Use Templates Effectively", path: "/docs/wip" },
      { title: "Monitor Environment Health", path: "/docs/wip" },
      { title: "Clean Up Environments", path: "/docs/wip" },
      { title: "Share Environments", path: "/docs/wip" },
      { title: "Migrate from Railway Direct", path: "/docs/wip" },
      { title: "Optimize Resource Usage", path: "/docs/wip" },
    ],
  },
  {
    title: "Troubleshooting",
    pages: [
      { title: "Overview", path: "/docs/wip" },
      { title: "Connection Issues", path: "/docs/wip" },
      { title: "Authentication Errors", path: "/docs/wip" },
      { title: "Environment Creation Failures", path: "/docs/wip" },
      { title: "Service Deployment Issues", path: "/docs/wip" },
      { title: "Configuration Errors", path: "/docs/wip" },
      { title: "Performance Problems", path: "/docs/wip" },
      { title: "Common Error Messages", path: "/docs/wip" },
      { title: "Getting Help", path: "/docs/wip" },
    ],
  },
];

