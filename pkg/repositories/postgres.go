package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable            = "users"
	groupsTable           = "groups"
	userGroupTable        = "user_group"
	userSpotifyTable      = "users_spotify"
	userYouTubeMusicTable = "user_youtubemusic"
	rolesTable            = "roles"
	playlistsTable        = "playlists"
	tracksTable           = "tracks"
	playlistTrackTable    = "playlist_track"
	tokensTable           = "tokens"
	youtubeMusicTracks    = "youtube_music_tracks"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
