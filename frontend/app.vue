<template>
  <div id="app">
    <header class="header">
      <div class="header-content">
        <h1>
          <span>🎵</span>
          Music Playlist Manager
        </h1>
        <nav class="nav">
          <NuxtLink to="/" :class="{ active: route.path === '/' }">
            首页
          </NuxtLink>
          <NuxtLink to="/songs" :class="{ active: route.path.startsWith('/songs') }">
            歌曲库
          </NuxtLink>
        </nav>
        <div class="user-section" v-if="isAuthenticated">
          <span class="user-greeting">
            <span class="user-avatar">{{ user?.username?.charAt(0)?.toUpperCase() }}</span>
            <span class="user-name">{{ user?.username }}</span>
          </span>
          <button class="btn btn-sm btn-secondary" @click="handleLogout">
            退出
          </button>
        </div>
        <div class="auth-links" v-else>
          <NuxtLink to="/login" class="btn btn-sm btn-secondary">
            登录
          </NuxtLink>
          <NuxtLink to="/register" class="btn btn-sm btn-primary">
            注册
          </NuxtLink>
        </div>
      </div>
    </header>
    <main class="container">
      <NuxtPage />
    </main>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { user, isAuthenticated, logout, initAuth, fetchCurrentUser, token } = useAuth()

onMounted(async () => {
  if (import.meta.client) {
    initAuth()
    if (token.value && !isAuthenticated.value) {
      await fetchCurrentUser()
    }
  }
})

const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    logout()
  }
}
</script>

<style scoped>
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-greeting {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-primary);
  font-size: 14px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
}

.user-name {
  font-weight: 600;
}

.auth-links {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
