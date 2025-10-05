import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import fs from 'fs'

export default defineConfig({
  plugins: [svelte()],
  server: {
    https: {
      key: fs.readFileSync('../certs/server.key'),
      cert: fs.readFileSync('../certs/server.crt'),
    },
    proxy: {
      '/api': {
        target: 'https://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
