import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { visualizer } from "rollup-plugin-visualizer"
import { defineConfig } from "vite"

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    visualizer({
      filename: "dist/stats.html",
      template: "treemap",
      gzipSize: true,
      brotliSize: true,
    }),
  ],
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
      allowedHosts: ['.getexposed.io']
  },
  build: {
    chunkSizeWarningLimit: 1000,
    outDir: "../internal/ui/dist",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes("node_modules")) return

          if (id.includes("react-router")) return "router"
          if (id.includes("@radix-ui")) return "radix"
          if (id.includes("lucide-react") || id.includes("react-icons")) return "icons"
          if (id.includes("recharts") || id.includes("d3")) return "charts"
          if (id.includes("date-fns") || id.includes("dayjs")) return "dates"

          return "vendor"
        },
      },
    },
  },
  optimizeDeps: {
    include: ["react", "react-dom", "react-router-dom"],
  },
})
