import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone', // Enables optimized Docker builds
};

export default nextConfig;
