import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  esbuild: {
    tsconfigRaw: {
      compilerOptions: {
        skipLibCheck: true, // Optional: Skips type checking of declaration files
        checkJs: false, // Optional: Skips type checking of JS files
      },
    },
  },
});
