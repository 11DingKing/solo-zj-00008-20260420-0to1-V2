<template>
  <div>
    <section class="section">
      <div class="section-header">
        <h2 class="section-title">我的播放列表</h2>
        <button class="btn btn-primary" @click="showCreateModal = true">
          + 创建播放列表
        </button>
      </div>
      
      <div v-if="loadingMy" class="loading">
        <div class="spinner"></div>
      </div>
      
      <div v-else-if="myPlaylists.length === 0" class="empty-state">
        <div class="empty-state-icon">📋</div>
        <div class="empty-state-text">还没有播放列表</div>
        <div class="empty-state-hint">点击上方按钮创建您的第一个播放列表</div>
      </div>
      
      <div v-else class="playlist-grid">
        <div 
          v-for="playlist in myPlaylists" 
          :key="playlist.id"
          class="playlist-card"
          @click="goToPlaylist(playlist.id)"
        >
          <div class="playlist-cover">
            🎵
            <div class="play-btn"></div>
          </div>
          <div class="playlist-name">{{ playlist.name }}</div>
          <div class="playlist-info">
            {{ playlist.song_count }} 首歌曲
            <span v-if="playlist.is_public" class="badge badge-public ml-2">公开</span>
          </div>
        </div>
      </div>
    </section>

    <section class="section mt-8">
      <div class="section-header">
        <h2 class="section-title">热门公开列表</h2>
      </div>
      
      <div v-if="loadingPopular" class="loading">
        <div class="spinner"></div>
      </div>
      
      <div v-else-if="popularPlaylists.length === 0" class="empty-state">
        <div class="empty-state-icon">🔥</div>
        <div class="empty-state-text">暂无热门播放列表</div>
      </div>
      
      <div v-else class="playlist-grid">
        <div 
          v-for="playlist in popularPlaylists" 
          :key="playlist.id"
          class="playlist-card"
          @click="goToPlaylist(playlist.id)"
        >
          <div class="playlist-cover" style="background: linear-gradient(135deg, #ff6b6b, #feca57);">
            🔥
            <div class="play-btn" @click.stop="goToPlaylist(playlist.id)"></div>
          </div>
          <div class="playlist-name">{{ playlist.name }}</div>
          <div class="playlist-info">
            {{ playlist.song_count }} 首歌曲
          </div>
          <div class="flex gap-2 mt-2">
            <button 
              class="btn btn-sm btn-secondary"
              @click.stop="goToPlaylist(playlist.id)"
            >
              查看
            </button>
            <button 
              class="btn btn-sm btn-primary"
              @click.stop="copyPlaylist(playlist)"
            >
              复制
            </button>
          </div>
        </div>
      </div>
    </section>

    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">创建播放列表</h3>
          <button class="modal-close" @click="showCreateModal = false">&times;</button>
        </div>
        
        <div v-if="createError" class="error-message">{{ createError }}</div>
        
        <form @submit.prevent="createPlaylist">
          <div class="form-group">
            <label class="form-label">名称</label>
            <input 
              v-model="newPlaylist.name" 
              type="text" 
              class="form-input"
              placeholder="输入播放列表名称"
              required
            />
          </div>
          
          <div class="form-group">
            <label class="form-label">描述</label>
            <textarea 
              v-model="newPlaylist.description" 
              class="form-input form-textarea"
              placeholder="描述这个播放列表"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label class="checkbox">
              <input 
                v-model="newPlaylist.is_public" 
                type="checkbox"
              />
              <span class="checkbox-label">设为公开</span>
            </label>
          </div>
          
          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary"
              @click="showCreateModal = false"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="creating"
            >
              {{ creating ? '创建中...' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Playlist } from '~/composables/useApi'

const songsApi = useSongsApi()
const playlistsApi = usePlaylistsApi()
const router = useRouter()

const myPlaylists = ref<Playlist[]>([])
const popularPlaylists = ref<Playlist[]>([])
const loadingMy = ref(true)
const loadingPopular = ref(true)

const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')

const newPlaylist = ref({
  name: '',
  description: '',
  is_public: false,
})

const loadMyPlaylists = async () => {
  loadingMy.value = true
  try {
    const response = await playlistsApi.getMy()
    if (response.success && response.data) {
      myPlaylists.value = response.data
    }
  } catch (error) {
    console.error('Failed to load my playlists:', error)
  } finally {
    loadingMy.value = false
  }
}

const loadPopularPlaylists = async () => {
  loadingPopular.value = true
  try {
    const response = await playlistsApi.getPopular()
    if (response.success && response.data) {
      popularPlaylists.value = response.data
    }
  } catch (error) {
    console.error('Failed to load popular playlists:', error)
  } finally {
    loadingPopular.value = false
  }
}

const goToPlaylist = (id: number) => {
  router.push(`/playlists/${id}`)
}

const createPlaylist = async () => {
  if (!newPlaylist.value.name.trim()) {
    createError.value = '请输入播放列表名称'
    return
  }

  creating.value = true
  createError.value = ''

  try {
    const response = await playlistsApi.create(newPlaylist.value)
    if (response.success && response.data) {
      showCreateModal.value = false
      newPlaylist.value = { name: '', description: '', is_public: false }
      await loadMyPlaylists()
      router.push(`/playlists/${response.data.id}`)
    } else {
      createError.value = response.error || '创建失败'
    }
  } catch (error) {
    createError.value = '创建播放列表时出错'
  } finally {
    creating.value = false
  }
}

const copyPlaylist = async (playlist: Playlist) => {
  try {
    const response = await playlistsApi.copy(playlist.id)
    if (response.success) {
      await loadMyPlaylists()
      alert(`已将 "${playlist.name}" 复制到您的播放列表中！`)
    } else {
      alert(response.error || '复制失败')
    }
  } catch (error) {
    alert('复制播放列表时出错')
  }
}

onMounted(() => {
  loadMyPlaylists()
  loadPopularPlaylists()
})
</script>

<style scoped>
.section {
  margin-bottom: 40px;
}

.mt-8 {
  margin-top: 32px;
}

.ml-2 {
  margin-left: 8px;
}

.flex {
  display: flex;
}

.gap-2 {
  gap: 8px;
}

.mt-2 {
  margin-top: 8px;
}
</style>
