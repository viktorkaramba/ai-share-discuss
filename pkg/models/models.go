package models

import (
	"io"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username" binding:"required"`
	Email    string `json:"email" db:"email" binding:"required"`
	Platform string `json:"platform" db:"platform" binding:"required"`
}

type Group struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
}

type Role struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" binding:"required"`
}

type UserSpotify struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"userID" db:"userid"`
	SpotifyUri   string `json:"spotifyUri" db:"spotifyuri"`
	AccessToken  string `json:"accessToken" db:"accesstoken"`
	RefreshToken string `json:"refreshToken" db:"refreshtoken"`
}

type UserGroup struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"userID" db:"userid"`
	GroupID    int    `json:"groupID" db:"groupid"`
	RoleID     int    `json:"roleID" db:"roleid"`
	PlayListID string `json:"playListID" db:"playlistid"`
}

type PlayList struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	GroupID     int    `json:"groupID" db:"groupid" binding:"required"`
}

type SpotifyUris struct {
	Uris     []string `json:"uris"`
	Position int      `json:"position"`
}

type Track struct {
	ID             int    `json:"id" db:"id"`
	SpotifyUri     string `json:"spotifyUri" db:"spotifyuri"`
	YouTubeMusicID string `json:"youTubeMusicID" db:"youtubemusicid"`
	Name           string `json:"name" db:"name"`
}

type PlayListTrack struct {
	ID                     int    `json:"id" db:"id"`
	TrackID                int    `json:"trackID" db:"trackid"`
	UserID                 int    `json:"userID" db:"userid"`
	PlayListID             int    `json:"playListID" db:"playlistid"`
	PlayListYouTubeMusicID string `json:"playListYouTubeMusicID" db:"playlistyoutubemusicid"`
	Token                  string `db:"accesstoken"`
}

type CreateTrack struct {
	ID                     int    `json:"id" db:"id"`
	UserID                 int    `json:"userID" db:"userid"`
	SpotifyUri             string `json:"spotifyUri" db:"spotifyuri"`
	YouTubeMusicID         string `json:"youTubeMusicID" db:"youtubemusicid"`
	PlayListYouTubeMusicID string `json:"PlayListYouTubeMusicID" db:"playlistyoutubemusicid"`
	Name                   string `json:"name" db:"name"`
}

type UserGroupToken struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"userID" db:"userid"`
	SpotifyUri string `json:"spotifyUri" db:"spotifyuri"`
	GroupID    int    `json:"groupID" db:"groupid"`
	RoleID     int    `json:"roleID" db:"roleid"`
	Token      string `json:"token" db:"accesstoken"`
	PlayListID string `json:"playListID" db:"playlistid"`
}

type UserClaims struct {
	UserID       int    `json:"userID"`
	UserPlatform string `json:"userPlatform"`
}

type ApiToken struct {
	AccessToken  string `json:"accessToken" db:"accesstoken"`
	RefreshToken string `json:"refreshToken" db:"refreshtoken"`
}

type Token struct {
	ID         int    `json:"id" db:"id"`
	TokenValue string `json:"tokenValue" db:"tokenvalue" binding:"required"`
	Revoked    bool   `json:"revoked" db:"revoked"`
	UserID     int    `json:"userID" db:"userid"`
}

type RefreshToken struct {
	Query string
	Body  io.Reader
}

type SpotifyData struct {
	Token      string
	SpotifyUri string
}

type ApiTrackID struct {
	SpotifyUri   string
	YouTubeMusic string
}
