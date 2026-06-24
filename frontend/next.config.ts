import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Next dev blocks cross-origin requests to dev-only assets (HMR/RSC bootstrap)
  // unless the origin is allowlisted. Without this, accessing the dev server via
  // an ngrok tunnel serves SSR HTML but never hydrates — so client effects (e.g.
  // the Telegram widget injection) never run and no login button appears.
  allowedDevOrigins: ["*.ngrok-free.app", "*.ngrok.app", "*.ngrok.io"],
};

export default nextConfig;
