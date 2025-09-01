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
      "/debug/pprof": "http://localhost:8080",
    },
  },
  build: {
    outDir: "../internal/ui/dist",
    emptyOutDir: true,
  },
})
