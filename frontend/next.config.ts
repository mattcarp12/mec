import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    return [
      {
        // Whenever the browser requests a URL starting with /api/...
        source: '/api/:path*',
        // Next.js will secretly forward the request to your Go backend
        destination: 'http://127.0.0.1:8080/api/:path*',
      },
    ]
  },
};

export default nextConfig;
