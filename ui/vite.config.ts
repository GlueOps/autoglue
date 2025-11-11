import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": "http://localhost:8080",
      "/swagger": "http://localhost:8080",
      "/db-studio": "http://localhost:8080",
    },
    allowedHosts: [".getexposed.io"],
  },
  build: {
    chunkSizeWarningLimit: 1000,
    outDir: "../internal/web/dist",
    emptyOutDir: true,
    sourcemap: true,
    cssMinify: "lightningcss",
    rollupOptions: {
      output: { manualChunks: { react: ["react", "react-dom", "react-router-dom"] } },
    },
  },
  esbuild: { legalComments: "none" },
})
