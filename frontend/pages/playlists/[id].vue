<template>
  <div>
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <div v-else-if="errorMessage" class="error-message">{{ errorMessage }}</div>

    <div v-else-if="playlist">
      <div class="playlist-detail-header">
        <div class="playlist-detail-cover">
          🎵
        </div>
        <div class="playlist-detail-info">
          <div class="playlist-detail-type">播放列表</div>
          <h1 class="playlist-detail-name">{{ playlist.name }}</h1>
          <p class="playlist-detail-desc">{{ playlist.description || '暂无描述' }}</p>
          <div class="playlist-detail-meta">
            <span>{{ songCount }} 首歌曲</span>
            <span v-if="playlist.is_public" class="badge badge-public">公开</span>
            <span v-else class="badge badge-private">私有</span>
          </div>
        </div>
      </div>

      <div class="playlist-actions" v-if="isOwner">
        <button class="btn btn-primary" @click="showAddSongModal = true">
          + 添加歌曲
        </button>
        <button class="btn btn-secondary" @click="showEditModal = true">
          编辑
        </button>
        <button class="btn btn-danger" @click="handleDeletePlaylist">
          删除
        </button>
      </div>

      <div class="playlist-actions" v-else-if="playlist.is_public">
        <button class="btn btn-primary" @click="handleCopyPlaylist">
          复制到我的播放列表
        </button>
      </div>

      <div class="card">
        <div class="song-item song-item-header playlist-song-item">
          <div v-if="isOwner" class="drag-header"></div>
          <div class="song-number-col">#</div>
          <div class="song-title-col">歌曲名</div>
          <div class="song-artist-col">歌手</div>
          <div class="song-duration-col">时长</div>
          <div v-if="isOwner" class="song-action-col"></div>
        </div>

        <div ref="songListRef" class="song-list">
          <div 
            v-for="song in songs" 
            :key="song.song_id"
            class="song-item playlist-song-item"
            :data-id="song.song_id"
          >
            <div v-if="isOwner" class="drag-handle"></div>
            <div class="song-number">{{ song.position + 1 }}</div>
            <div class="song-name">{{ song.name }}</div>
            <div class="song-artist">{{ song.artist }}</div>
            <div class="song-duration">{{ formatDuration(song.duration) }}</div>
            <div class="song-actions" v-if="isOwner">
              <button 
                class="btn btn-sm btn-danger"
                @click="handleRemoveSong(song)"
              >
                移除
              </button>
            </div>
          </div>
        </div>

        <div v-if="songs.length === 0" class="empty-state">
          <div class="empty-state-icon">🎵</div>
          <div class="empty-state-text">播放列表为空</div>
          <div v-if="isOwner" class="empty-state-hint">点击上方按钮添加歌曲</div>
        </div>
      </div>

      <div class="total-duration-bar">
        <span class="total-duration-text">
          共 <strong>{{ songCount }}</strong> 首歌曲
        </span>
        <span class="total-duration-text">
          总时长：<strong>{{ formatDuration(totalDuration) }}</strong>
        </span>
      </div>
    </div>

    <div v-if="showAddSongModal" class="modal-overlay" @click.self="showAddSongModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">添加歌曲</h3>
          <button class="modal-close" @click="showAddSongModal = false">&times;</button>
        </div>

        <div class="search-bar">
          <input 
            v-model="addSongSearch" 
            type="text" 
            class="search-input"
            placeholder="搜索歌曲..."
            @input="searchAvailableSongs"
          />
        </div>

        <div v-if="loadingAvailableSongs" class="loading">
          <div class="spinner"></div>
        </div>

        <div v-else-if="availableSongs.length === 0" class="empty-state">
          <div class="empty-state-text">没有可添加的歌曲</div>
        </div>

        <div v-else class="modal-song-list">
          <div 
            v-for="song in availableSongs" 
            :key="song.id"
            class="song-item selectable"
            @click="handleAddSong(song)"
          >
            <div class="song-number-col song-cover-small">
              <img v-if="song.cover_url" :src="song.cover_url" alt="" />
              <span v-else>🎵</span>
            </div>
            <div class="song-title-col">
              <div class="song-name">{{ song.name }}</div>
              <div class="song-artist">{{ song.album || '' }}</div>
            </div>
            <div class="song-artist-col">{{ song.artist }}</div>
            <div class="song-duration-col">{{ formatDuration(song.duration) }}</div>
            <div class="song-action-col">
              <span class="add-indicator">+ 添加</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">编辑播放列表</h3>
          <button class="modal-close" @click="showEditModal = false">&times;</button>
        </div>

        <form @submit.prevent="handleUpdatePlaylist">
          <div class="form-group">
            <label class="form-label">名称</label>
            <input 
              v-model="editForm.name" 
              type="text" 
              class="form-input"
              placeholder="输入播放列表名称"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">描述</label>
            <textarea 
              v-model="editForm.description" 
              class="form-input form-textarea"
              placeholder="描述这个播放列表"
            ></textarea>
          </div>

          <div class="form-group">
            <label class="checkbox">
              <input 
                v-model="editForm.is_public" 
                type="checkbox"
              />
              <span class="checkbox-label">设为公开</span>
            </label>
          </div>

          <div class="modal-footer">
            <button 
              type="button" 
              class="btn btn-secondary"
              @click="showEditModal = false"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="btn btn-primary"
              :disabled="savingEdit"
            >
              {{ savingEdit ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Song, Playlist, PlaylistSongDetail } from '~/composables/useApi'

const route = useRoute()
const router = useRouter()
const songsApi = useSongsApi()
const playlistsApi = usePlaylistsApi()

const songListRef = ref<HTMLElement | null>(null)
const playlist = ref<Playlist | null>(null)
const songs = ref<PlaylistSongDetail[]>([])
const totalDuration = ref(0)
const loading = ref(true)
const errorMessage = ref('')
const isOwner = ref(false)

const songCount = computed(() => songs.value?.length || 0)

const showAddSongModal = ref(false)
const showEditModal = ref(false)
const availableSongs = ref<Song[]>([])
const loadingAvailableSongs = ref(false)
const addSongSearch = ref('')
const savingEdit = ref(false)

const playlistId = computed(() => {
  const id = route.params.id
  return typeof id === 'string' ? parseInt(id, 10) : 0
})

const editForm = ref({
  name: '',
  description: '',
  is_public: false,
})

const currentUserId = 1

const loadPlaylist = async () => {
  loading.value = true
  errorMessage.value = ''
  songs.value = []
  totalDuration.value = 0
  playlist.value = null

  try {
    const [playlistRes, songsRes] = await Promise.all([
      playlistsApi.get(playlistId.value),
      playlistsApi.getSongs(playlistId.value),
    ])

    if (playlistRes.success && playlistRes.data) {
      playlist.value = playlistRes.data
      isOwner.value = playlistRes.data.owner_id === currentUserId
      
      editForm.value = {
        name: playlistRes.data.name,
        description: playlistRes.data.description,
        is_public: playlistRes.data.is_public,
      }
    } else {
      errorMessage.value = playlistRes.error || '加载播放列表失败'
    }

    if (songsRes.success && songsRes.data) {
      songs.value = songsRes.data.songs || []
      totalDuration.value = songsRes.data.total_duration || 0
    }
  } catch (err) {
    errorMessage.value = '加载数据时出错'
  } finally {
    loading.value = false
  }
}

const searchAvailableSongs = async () => {
  loadingAvailableSongs.value = true
  availableSongs.value = []
  try {
    const response = await songsApi.list(addSongSearch.value)
    if (response.success && response.data) {
      const existingSongIds = new Set((songs.value || []).map(s => s.song_id))
      availableSongs.value = response.data.filter(s => !existingSongIds.has(s.id))
    }
  } catch (error) {
    console.error('Failed to search songs:', error)
  } finally {
    loadingAvailableSongs.value = false
  }
}

const handleAddSong = async (song: Song) => {
  try {
    const response = await playlistsApi.addSong(playlistId.value, song.id)
    if (response.success) {
      await loadPlaylist()
    } else {
      alert(response.error || '添加失败')
    }
  } catch (error) {
    alert('添加歌曲时出错')
  }
}

const handleRemoveSong = async (song: PlaylistSongDetail) => {
  if (!confirm(`确定要从播放列表中移除 "${song.name}" 吗？`)) {
    return
  }

  try {
    const response = await playlistsApi.removeSong(playlistId.value, song.song_id)
    if (response.success) {
      await loadPlaylist()
    } else {
      alert(response.error || '移除失败')
    }
  } catch (error) {
    alert('移除歌曲时出错')
  }
}

const handleUpdatePlaylist = async () => {
  if (!editForm.value.name.trim()) {
    alert('请输入播放列表名称')
    return
  }

  savingEdit.value = true
  try {
    const response = await playlistsApi.update(playlistId.value, editForm.value)
    if (response.success) {
      showEditModal.value = false
      await loadPlaylist()
    } else {
      alert(response.error || '更新失败')
    }
  } catch (error) {
    alert('更新播放列表时出错')
  } finally {
    savingEdit.value = false
  }
}

const handleDeletePlaylist = async () => {
  if (!confirm('确定要删除这个播放列表吗？此操作无法撤销！')) {
    return
  }

  try {
    const response = await playlistsApi.delete(playlistId.value)
    if (response.success) {
      router.push('/')
    } else {
      alert(response.error || '删除失败')
    }
  } catch (error) {
    alert('删除播放列表时出错')
  }
}

const handleCopyPlaylist = async () => {
  if (!playlist.value) return

  try {
    const response = await playlistsApi.copy(playlistId.value)
    if (response.success) {
      alert(`已将 "${playlist.value.name}" 复制到您的播放列表中！`)
      router.push('/')
    } else {
      alert(response.error || '复制失败')
    }
  } catch (error) {
    alert('复制播放列表时出错')
  }
}

watch(() => showAddSongModal.value, (val) => {
  if (val) {
    addSongSearch.value = ''
    searchAvailableSongs()
  }
})

watch(playlistId, () => {
  loadPlaylist()
})

onMounted(() => {
  loadPlaylist()
})

if (import.meta.client) {
  nextTick(() => {
    let sortable: any = null

    const initSortable = async () => {
      if (!songListRef.value) return
      
      const { default: Sortable } = await import('sortablejs')
      
      sortable = new Sortable(songListRef.value, {
        animation: 150,
        handle: '.drag-handle',
        ghostClass: 'sortable-ghost',
        chosenClass: 'sortable-chosen',
        onEnd: async (evt: any) => {
          const items = [...(songs.value || [])]
          const oldIndex = evt.oldIndex
          const newIndex = evt.newIndex

          if (oldIndex === newIndex) return

          const [movedItem] = items.splice(oldIndex, 1)
          items.splice(newIndex, 0, movedItem)

          items.forEach((item, index) => {
            item.position = index
          })

          const positions = items.map(item => ({
            song_id: item.song_id,
            position: item.position,
          }))

          songs.value = items

          try {
            await playlistsApi.updatePositions(playlistId.value, positions)
          } catch (error) {
            console.error('Failed to update positions:', error)
            alert('保存排序时出错，请重试')
          }
        },
      })
    }

    watch(isOwner, (val) => {
      if (val && songListRef.value) {
        initSortable()
      }
    }, { immediate: true })
  })
}
</script>

<style scoped>
img {
  max-width: 100%;
  height: auto;
}

.modal-song-list {
  max-height: 400px;
  overflow-y: auto;
}

.selectable {
  cursor: pointer;
}

.selectable:hover {
  background-color: var(--bg-hover) !important;
}

.add-indicator {
  color: var(--primary-color);
  font-size: 12px;
  font-weight: 600;
}

.total-duration-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: var(--bg-card);
  border-radius: 0 0 12px 12px;
  margin-top: -24px;
  border-top: 1px solid var(--border-color);
}

.total-duration-text {
  color: var(--text-secondary);
  font-size: 14px;
}

.total-duration-text strong {
  color: var(--text-primary);
  font-size: 16px;
}
</style>
