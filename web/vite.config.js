import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
export default defineConfig({
    plugins: [react()],
    server: {
        host: "127.0.0.1",
        proxy: {
            "/api": "http://127.0.0.1:8080",
            "/ws": {
                target: "ws://127.0.0.1:8080",
                ws: true,
            },
        },
    },
    test: {
        environment: "jsdom",
        globals: true,
        setupFiles: "./src/vitest.setup.ts",
    },
});
