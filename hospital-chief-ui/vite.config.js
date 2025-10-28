import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import https from 'https'

export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      '/api': {
        target: 'https://localhost:9090',
        changeOrigin: true,
        secure: false,
        agent: new https.Agent({
          rejectUnauthorized: false
        })
      },
      '/fhir': {
        target: 'https://localhost:9090',
        changeOrigin: true,
        secure: false,
        agent: new https.Agent({
          rejectUnauthorized: false
        })
      },
    },
  },
})
