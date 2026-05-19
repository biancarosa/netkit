import type { NextConfig } from "next";

const normalizeBasePath = (value?: string): string => {
  if (!value || value === "/") {
    return "";
  }

  const withLeadingSlash = value.startsWith("/") ? value : `/${value}`;
  return withLeadingSlash.replace(/\/+$/, "");
};

const basePath = normalizeBasePath(
  process.env.NEXT_PUBLIC_NETKIT_BASE_PATH ||
  process.env.NETKIT_DASHBOARD_BASE_PATH
);

const nextConfig: NextConfig = {
  // Enable static export for embedding
  output: 'export',
  trailingSlash: true,
  distDir: 'out',
  ...(basePath ? { basePath, assetPrefix: basePath } : {}),
  env: {
    NEXT_PUBLIC_NETKIT_BASE_PATH: basePath,
  },
  
  experimental: {
    staleTimes: {
      dynamic: 0,
      static: 0,
    },
  },
  
  // Configure static export
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
