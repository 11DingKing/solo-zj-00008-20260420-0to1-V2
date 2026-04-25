package models

type Song struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int    `json:"duration"`
	CoverURL    string `json:"cover_url"`
	AudioFileURL string `json:"audio_file_url"`
}

type CreateSongRequest struct {
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int    `json:"duration"`
	CoverURL    string `json:"cover_url"`
	AudioFileURL string `json:"audio_file_url"`
}

type Playlist struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsPublic      bool   `json:"is_public"`
	OwnerID       int    `json:"owner_id"`
	OwnerUsername string `json:"owner_username,omitempty"`
	SongCount     int    `json:"song_count,omitempty"`
}

type CreatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type PlaylistSong struct {
	PlaylistID int `json:"playlist_id"`
	SongID     int `json:"song_id"`
	Position   int `json:"position"`
}

type PlaylistSongDetail struct {
	SongID     int    `json:"song_id"`
	Name       string `json:"name"`
	Artist     string `json:"artist"`
	Duration   int    `json:"duration"`
	Position   int    `json:"position"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type UserWithPassword struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
