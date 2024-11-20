package configs

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"
	"os"
)

type Config struct {
	SpotifyLoginConfig    oauth2.Config
	GoogleLoginConfig     oauth2.Config
	AppleMusicLoginConfig oauth2.Config
}

var AppConfig Config

const OauthSpotifyUrlAPI = "https://api.spotify.com/v1/me"
const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

//const OauthAppleMusicUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func LoadAuthConfig() {

	// Oauth configuration for Spotify
	AppConfig.SpotifyLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Endpoint:     spotify.Endpoint,
		RedirectURL:  viper.GetString("domain") + viper.GetString("port") + "/auth/spotify-callback",
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"playlist-modify-public",
			"playlist-modify-private",
		},
	}

	// Oauth configuration for YouTube Music
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  viper.GetString("domain") + viper.GetString("port") + "/auth/youtube-music-callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/youtube.force-ssl",
		},
	}
}
