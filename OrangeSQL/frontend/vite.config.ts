import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import monaco from "vite-plugin-monaco-editor-esm";
import { fileURLToPath, URL } from "node:url";

export default defineConfig({
  plugins: [
    vue(),
    monaco(),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
