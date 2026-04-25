export default defineNuxtConfig({
  devtools: { enabled: true },
  
  runtimeConfig: {
    apiBase: process.env.NUXT_API_BASE || 'http://backend:8080/api',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api'
    }
  },
  
  css: ['~/assets/css/main.css'],
  
  app: {
    head: {
      title: 'Music Playlist Manager',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' }
      ]
    }
  },
  
  compatibilityDate: '2024-01-01'
})
