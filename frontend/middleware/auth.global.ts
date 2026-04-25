export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated, initAuth, fetchCurrentUser, token } = useAuth()

  if (import.meta.client) {
    initAuth()

    if (token.value && !isAuthenticated.value) {
      await fetchCurrentUser()
    }
  }

  const publicPaths = ['/login', '/register']
  const isPublicPath = publicPaths.includes(to.path)

  if (!isPublicPath && !isAuthenticated.value) {
    return navigateTo({
      path: '/login',
      query: { redirect: to.fullPath },
    })
  }

  if (isPublicPath && isAuthenticated.value) {
    return navigateTo('/')
  }
})
