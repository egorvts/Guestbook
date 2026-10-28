import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', 'BACKEND_URL')

  return {
    server: {
      proxy: {
        '/api': {
          target: env.BACKEND_URL || 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
  }
})
