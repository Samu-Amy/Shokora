import type { NextConfig } from "next";

const isDev = process.env.NODE_ENV === "development";

const nextConfig: NextConfig = {
  // static export
  output: "export",
};

// reverse proxy (only in dev or prod with node js)
if (isDev) {
  nextConfig.rewrites = async () => {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/api/:path*", // backend reale
      },
    ];
  };
}

export default nextConfig;
