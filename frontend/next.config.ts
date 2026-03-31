import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // static export
  output: "export",

  // reverse proxy (only in dev or prod with node js)
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/api/:path*", // backend reale
      },
    ];
  },
};

export default nextConfig;
