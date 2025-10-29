import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import https from 'https'

export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      '/api': {
        target: 'https://localhost:8081',
        changeOrigin: true,
        secure: false,
        agent: new https.Agent({
          rejectUnauthorized: false
        })
      },
      '/ws': {
        target: 'wss://localhost:8081',
        ws: true,
        changeOrigin: true,
        secure: false,
        agent: new https.Agent({
          rejectUnauthorized: false
        })
      },
    },
  },
})
