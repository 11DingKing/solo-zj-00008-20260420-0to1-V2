export interface Song {
  id: number
  name: string
  artist: string
  album: string
  duration: number
  cover_url: string
  audio_file_url: string
}

export interface Playlist {
  id: number
  name: string
  description: string
  is_public: boolean
  owner_id: number
  song_count?: number
}

export interface PlaylistSongDetail {
  song_id: number
  name: string
  artist: string
  duration: number
  position: number
}

export interface ApiResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
}

const BASE_URL = '/api'

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const response = await $fetch<ApiResponse<T>>(`${BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>),
      },
    })
    return response
  } catch (error: any) {
    return {
      success: false,
      error: error?.data?.error || error?.message || 'Network error',
    }
  }
}

export function useSongsApi() {
  return {
    list: (search?: string) => 
      request<Song[]>(search ? `/songs?search=${encodeURIComponent(search)}` : '/songs', { method: 'GET' }),
    get: (id: number) => request<Song>(`/songs/${id}`, { method: 'GET' }),
    create: (song: Omit<Song, 'id'>) => request<{ id: number }>('/songs', { method: 'POST', body: song }),
    update: (id: number, song: Omit<Song, 'id'>) => request<void>(`/songs/${id}`, { method: 'PUT', body: song }),
    delete: (id: number) => request<void>(`/songs/${id}`, { method: 'DELETE' }),
  }
}

export function usePlaylistsApi() {
  return {
    getMy: () => request<Playlist[]>('/playlists/my', { method: 'GET' }),
    getPopular: () => request<Playlist[]>('/playlists/popular', { method: 'GET' }),
    get: (id: number) => request<Playlist>(`/playlists/${id}`, { method: 'GET' }),
    getSongs: (id: number) => 
      request<{ songs: PlaylistSongDetail[]; total_duration: number }>(`/playlists/${id}/songs`, { method: 'GET' }),
    create: (playlist: { name: string; description: string; is_public: boolean }) =>
      request<{ id: number }>('/playlists', { method: 'POST', body: playlist }),
    update: (id: number, playlist: { name: string; description: string; is_public: boolean }) =>
      request<void>(`/playlists/${id}`, { method: 'PUT', body: playlist }),
    delete: (id: number) => request<void>(`/playlists/${id}`, { method: 'DELETE' }),
    addSong: (playlistId: number, songId: number) =>
      request<void>(`/playlists/${playlistId}/songs`, { method: 'POST', body: { song_id: songId } }),
    removeSong: (playlistId: number, songId: number) =>
      request<void>(`/playlists/${playlistId}/songs/${songId}`, { method: 'DELETE' }),
    updatePositions: (playlistId: number, positions: { song_id: number; position: number }[]) =>
      request<void>(`/playlists/${playlistId}/positions`, { method: 'PUT', body: positions }),
    copy: (id: number) => request<{ id: number }>(`/playlists/${id}/copy`, { method: 'POST', body: {} }),
  }
}

export function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
