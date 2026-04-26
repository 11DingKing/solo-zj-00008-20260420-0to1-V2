export interface User {
  id: number
  username: string
  email?: string
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterCredentials {
  username: string
  email: string
  password: string
}

export interface LoginResult {
  token: string
  user: User
}

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
  owner_username?: string
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

const getToken = (): string | null => {
  if (import.meta.client) {
    return localStorage.getItem('auth_token')
  }
  return null
}

const clearAuth = () => {
  if (import.meta.client) {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
  }
}

let onAuthError: (() => void) | null = null

export const setAuthErrorHandler = (handler: () => void) => {
  onAuthError = handler
}

function extractErrorMessage(error: any): string {
  if (error?.data?.error) {
    return error.data.error
  }
  if (error?.response?._data?.error) {
    return error.response._data.error
  }
  if (error?.response?.data?.error) {
    return error.response.data.error
  }
  if (error?.message) {
    return error.message
  }
  return 'Network error'
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = getToken()

  try {
    const response = await $fetch<ApiResponse<T>>(`${BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(options.headers as Record<string, string>),
      },
    })
    return response
  } catch (error: any) {
    const status = error?.status || error?.response?.status
    const errorMessage = extractErrorMessage(error)

    if (status === 401) {
      clearAuth()
      if (onAuthError) {
        onAuthError()
      }
    }

    return {
      success: false,
      error: errorMessage,
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

export function useAuthApi() {
  return {
    login: (credentials: LoginCredentials) =>
      request<LoginResult>('/auth/login', { method: 'POST', body: credentials }),
    register: (credentials: RegisterCredentials) =>
      request<LoginResult>('/auth/register', { method: 'POST', body: credentials }),
    getCurrentUser: () =>
      request<User>('/auth/me', { method: 'GET' }),
  }
}

export function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
