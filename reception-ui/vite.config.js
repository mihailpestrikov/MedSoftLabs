import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import fs from 'fs'
import path from 'path'

// Only load HTTPS certs if they exist (for local dev)
const certPath = path.resolve('../certs/server.crt')
const keyPath = path.resolve('../certs/server.key')
const httpsConfig = fs.existsSync(certPath) && fs.existsSync(keyPath)
  ? {
      key: fs.readFileSync(keyPath),
      cert: fs.readFileSync(certPath),
    }
  : undefined

export default defineConfig({
  plugins: [svelte()],
  server: {
    https: httpsConfig,
    proxy: {
      '/api': {
        target: 'https://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
