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
    // Sourcemaps are ~20MB for this bundle and were being served publicly on
    // demand, which is both slow through a proxy and a full source disclosure.
    // "hidden" still emits them for upload to an error tracker, without the
    // //# sourceMappingURL comment that makes browsers fetch them.
    sourcemap: process.env.VITE_SOURCEMAP === "true" ? true : "hidden",
    cssMinify: "lightningcss",
    rollupOptions: {
      output: {
        // Vite 8 bundles with rolldown rather than rollup/esbuild, and rolldown
        // dropped the `manualChunks` object form in favour of `advancedChunks`.
        // The intent is unchanged: keep React and the router in one long-lived
        // chunk so an app-code change does not invalidate them in client caches.
        advancedChunks: {
          groups: [
            { name: "react", test: /[\\/]node_modules[\\/](react|react-dom|react-router-dom)[\\/]/ },
          ],
        },
        // Was `esbuild: { legalComments: "none" }`. esbuild is no longer the
        // transformer, and rolldown's own `legalComments` is already deprecated
        // in favour of `comments.legal`, so go straight to the latter.
        comments: { legal: false },
      },
    },
  },
})
