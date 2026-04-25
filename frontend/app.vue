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
      </div>
    </header>
    <main class="container">
      <NuxtPage />
    </main>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
</script>
