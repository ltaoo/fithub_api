import path from "path";
import fs from 'fs';

import { UserConfigExport, defineConfig } from "vite";
import { timeless_jsx } from "./scripts/timeless-jsx";

const pkg = (() => {
  try {
    return JSON.parse(fs.readFileSync(path.resolve(__dirname, "./package.json"), "utf-8"));
  } catch (err) {
    return null;
  }
})();

const config = defineConfig(({ mode }) => {
  return {
    base: "/admin/",
    test: { threads: false, cache: { dir: ".cache/vitest" } },
    plugins: [
      timeless_jsx(),
      {
        name: "assert-timeless-runtime",
        generateBundle() {
          const legacy = [...this.getModuleIds()].filter(id => /node_modules\/(?:\.pnpm\/[^/]+\/node_modules\/)?(?:solid-js|lucide-solid)\//.test(id));
          if (legacy.length) this.error("Solid runtime remains in bundle: " + legacy.join(", "));
          console.log("Timeless runtime verified: no SolidJS or lucide-solid runtime modules.");
        },
      },
    ],
    resolve: {
      alias: {
        "hls.js": "hls.js/dist/hls.min.js",
        "@timeless/timeless": path.resolve(__dirname, "../../timeless/packages/timeless/dist/index.esm.js"),
        "@timeless/timeless-dom": path.resolve(__dirname, "../../timeless/packages/timeless-dom/dist/timeless.dom.esm.js"),
        "@timeless/shadcn/globals.css": path.resolve(__dirname, "../../timeless/packages/shadcn/dist/timeless.shadcn.css"),
        "@timeless/shadcn": path.resolve(__dirname, "../../timeless/packages/shadcn/dist/index.esm.js"),
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      fs: { allow: [path.resolve(__dirname), path.resolve(__dirname, "../../timeless")] },
      proxy: { "/api": { target: process.env.ADMIN_API_URL || "http://127.0.0.1:8080", changeOrigin: true } },
    },
    esbuild: {
      drop: mode === "production" ? ["console", "debugger"] : [],
    },
    define: {
      "process.global.__VERSION__": JSON.stringify(pkg ? pkg.version : "unknown"),
    },
    build: {
      outDir: "../dist/admin",
      emptyOutDir: true,
      target: "esnext",
      rollupOptions: {
        output: {
          manualChunks(filepath) {
            // if (filepath.includes("hls.js")) {
            //   return "hls";
            // }
            if (filepath.includes("node_modules") && !filepath.includes("hls")) {
              return "vendor";
            }
          },
        },
      },
    },
  };
});

export default config;
