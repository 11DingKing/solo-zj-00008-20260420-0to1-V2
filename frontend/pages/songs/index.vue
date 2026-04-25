<template>
  <div>
    <div class="section-header">
      <h2 class="section-title">歌曲库</h2>
      <button class="btn btn-primary" @click="showCreateModal = true">
        + 添加歌曲
      </button>
    </div>

    <div class="search-bar">
      <input 
        v-model="searchQuery" 
        type="text" 
        class="search-input"
        placeholder="搜索歌曲名或歌手..."
        @input="debouncedSearch"
      />
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="songs.length === 0" class="empty-state">
      <div class="empty-state-icon">🎵</div>
      <div class="empty-state-text">暂无歌曲</div>
      <div class="empty-state-hint">点击上方按钮添加您的第一首歌曲</div>
    </div>

    <div v-else class="card">
      <div class="song-item song-item-header library-song-item">
        <div class="song-cover-col"></div>
        <div class="song-title-col">歌曲名</div>
        <div class="song-artist-col">歌手</div>
        <div class="song-album-col">专辑</div>
        <div class="song-duration-col">时长</div>
        <div class="song-action-col"></div>
      </div>

      <div class="song-list">
        <div 
          v-for="(song, index) in songs" 
          :key="song.id"
          class="song-item library-song-item"
        >
          <div class="song-cover-col song-cover-small">
            <img v-if="song.cover_url" :src="song.cover_url" alt="" />
            <span v-else>🎵</span>
          </div>
          <div class="song-title-col">
            <div class="song-name">{{ song.name }}</div>
            <div class="song-artist-mobile">{{ song.artist }}</div>
          </div>
          <div class="song-artist-col">{{ song.artist }}</div>
          <div class="song-album-col">{{ song.album || '-' }}</div>
          <div class="song-duration-col">{{ formatDuration(song.duration) }}</div>
          <div class="song-action-col song-actions">
            <button 
              class="btn btn-sm btn-secondary"
              @click="editSong(song)"
            >
              编辑
            </button>
            <button 
              class="btn btn-sm btn-danger"
              @click="deleteSong(song)"
            >
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ showEditModal ? '编辑歌曲' : '添加歌曲' }}</h3>
          <button class="modal-close" @click="closeModals">&times;</button>
        </div>

        <div v-if="formError" class="error-message">{{ formError }}</div>

        <form @submit.prevent="submitSong">
          <div class="form-group">
            <label class="form-label">歌曲名</label>
            <input 
              v-model="songForm.name" 
              type="text" 
              class="form-input"
              placeholder="输入歌曲名"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">歌手</label>
            <input 
              v-model="songForm.artist" 
              type="text" 
              class="form-input"
              placeholder="输入歌手名"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">专辑</label>
            <input 
              v-model="songForm.album" 
              type="text" 
              class="form-input"
              placeholder="输入专辑名"
            />
          </div>

          <div class="form-group">
            <label class="form-label">时长（秒）</label>
            <input 
              v-model.number="songForm.duration" 
              type="number" 
              class="form-input"
              placeholder="输入时长"
              min="1"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">封面图 URL</label>
            <input 
              v-model="songForm.cover_url" 
              type="url" 
              class="form-input"
              placeholder="https://example.com/cover.jpg"
            />
          </div>

          <div class="form-group">
            <label class="form-label">音频文件 URL</label>
            <input 
              v-model="songForm.audio_file_url" 
              type="url" 
              class="form-input"
              placeholder="https://example.com/audio.mp3"
            />
          </div>

          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary"
              @click="closeModals"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="submitting"
            >
              {{ submitting ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="showDeleteConfirm = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">确认删除</h3>
          <button class="modal-close" @click="showDeleteConfirm = false">&times;</button>
        </div>

        <p>确定要删除歌曲 "{{ songToDelete?.name }}" 吗？此操作无法撤销。</p>

        <div class="modal-footer">
          <button 
            type="button" 
            class="btn btn-secondary"
            @click="showDeleteConfirm = false"
          >
            取消
          </button>
          <button 
            type="button" 
            class="btn btn-danger"
            @click="confirmDelete"
            :disabled="deleting"
          >
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Song } from '~/composables/useApi'

const songsApi = useSongsApi()

const songs = ref<Song[]>([])
const loading = ref(true)
const searchQuery = ref('')
const searchTimeout = ref<number | null>(null)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteConfirm = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const formError = ref('')

const songToDelete = ref<Song | null>(null)
const editingSongId = ref<number | null>(null)

const defaultSongForm = {
  name: '',
  artist: '',
  album: '',
  duration: 180,
  cover_url: '',
  audio_file_url: '',
}

const songForm = ref({ ...defaultSongForm })

const loadSongs = async (search?: string) => {
  loading.value = true
  songs.value = []
  try {
    const response = await songsApi.list(search)
    if (response.success && response.data) {
      songs.value = response.data
    }
  } catch (error) {
    console.error('Failed to load songs:', error)
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
  searchTimeout.value = window.setTimeout(() => {
    loadSongs(searchQuery.value)
  }, 300)
}

const editSong = (song: Song) => {
  editingSongId.value = song.id
  songForm.value = {
    name: song.name,
    artist: song.artist,
    album: song.album,
    duration: song.duration,
    cover_url: song.cover_url,
    audio_file_url: song.audio_file_url,
  }
  showEditModal.value = true
  formError.value = ''
}

const deleteSong = (song: Song) => {
  songToDelete.value = song
  showDeleteConfirm.value = true
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingSongId.value = null
  songForm.value = { ...defaultSongForm }
  formError.value = ''
}

const submitSong = async () => {
  if (!songForm.value.name.trim() || !songForm.value.artist.trim()) {
    formError.value = '请填写歌曲名和歌手'
    return
  }

  submitting.value = true
  formError.value = ''

  try {
    let response
    
    if (showEditModal.value && editingSongId.value) {
      response = await songsApi.update(editingSongId.value, songForm.value)
    } else {
      response = await songsApi.create(songForm.value)
    }

    if (response.success) {
      closeModals()
      await loadSongs(searchQuery.value)
    } else {
      formError.value = response.error || '保存失败'
    }
  } catch (error) {
    formError.value = '保存歌曲时出错'
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!songToDelete.value) return

  deleting.value = true

  try {
    const response = await songsApi.delete(songToDelete.value.id)
    if (response.success) {
      showDeleteConfirm.value = false
      songToDelete.value = null
      await loadSongs(searchQuery.value)
    } else {
      alert(response.error || '删除失败')
    }
  } catch (error) {
    alert('删除歌曲时出错')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadSongs()
})
</script>

<style scoped>
img {
  max-width: 100%;
  height: auto;
}

.song-artist-mobile {
  display: none;
}

@media (max-width: 768px) {
  .song-artist-col,
  .song-album-col {
    display: none;
  }
  
  .song-artist-mobile {
    display: block;
    color: var(--text-secondary);
    font-size: 12px;
    margin-top: 2px;
  }
}
</style>
