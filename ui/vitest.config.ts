import react from "@vitejs/plugin-react"
import path from "path"
import { defineConfig } from "vitest/config"

// Kept separate from vite.config.ts so the production build config (embedding
// into internal/web/dist, brotli precompression, chunking) is not dragged into
// the test run.
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { "@": path.resolve(__dirname, "./src") },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    // The SDK is generated and the dist output is a build artifact; neither
    // contains tests worth collecting.
    exclude: ["node_modules/**", "dist/**", "src/sdk/**"],
    restoreMocks: true,
  },
})
