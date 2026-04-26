export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const method = getMethod(event)
  const body = method !== 'GET' && method !== 'HEAD' ? await readBody(event) : undefined
  
  const backendUrl = config.apiBase || 'http://backend:8080/api'
  const path = getRequestPath(event).replace(/^\/api/, '')
  const fullUrl = `${backendUrl}${path}`
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  
  const authorization = getHeader(event, 'authorization')
  if (authorization) {
    headers['Authorization'] = authorization
  }
  
  try {
    const response = await $fetch(fullUrl, {
      method,
      body,
      headers,
      onResponse({ response }) {
        setResponseStatus(event, response.status)
      },
    })
    return response
  } catch (error: any) {
    const status = error?.response?.status || 500
    const message = error?.data?.error || error?.message || 'Proxy request failed'
    
    setResponseStatus(event, status)
    return {
      success: false,
      error: message,
    }
  }
})
