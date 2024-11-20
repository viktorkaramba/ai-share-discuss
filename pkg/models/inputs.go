package models

import "errors"

type UserGroupInput struct {
	GroupID int `json:"groupID" binding:"required"`
}

type AddTrackInput struct {
	SpotifyUri     string `json:"spotifyUri" binding:"required,min=1"`
	YouTubeMusicID string `json:"youTubeMusicID" binding:"required,min=1"`
	Name           string `json:"name" binding:"required,min=1"`
}

type AddTrack struct {
	SpotifyUri     string `json:"spotifyUri" binding:"required,min=1"`
	YouTubeMusicID string `json:"youTubeMusicID" binding:"required,min=1"`
	Name           string `json:"name" binding:"required,min=1"`
	PlayListID     int
	GroupID        int
	UserID         int
}

type UserCreateGroupInput struct {
	GroupName           string `json:"groupName" binding:"required,min=1,max=20"`
	PlayListName        string `json:"playListName" binding:"required,min=1,max=20"`
	GroupDescription    string `json:"groupDescription" binding:"required,min=1,max=50"`
	PlayListDescription string `json:"playListDescription" binding:"required,min=1,max=50"`
}

type UserCreateGroup struct {
	ID                  int    `json:"id"`
	GroupName           string `json:"groupName" binding:"required,min=1,max=20"`
	PlayListName        string `json:"playListName" binding:"required,min=1,max=20"`
	GroupDescription    string `json:"groupDescription" binding:"required,min=1,max=50"`
	PlayListDescription string `json:"playListDescription" binding:"required,min=1,max=50"`
	PlayListID          string
	Platform            string
}

type RefreshTokenInput struct {
	UserId int `json:"userId" binding:"required"`
}

type DeleteSpotifyTrack struct {
	URI string `json:"uri"`
}

type DeleteSpotifyTrackInput struct {
	Tracks []DeleteSpotifyTrack `json:"tracks"`
}

type UpdateGroupInput struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=20"`
	Description *string `json:"description" binding:"omitempty,min=1,max=50"`
}

type UpdateUserGroupInput struct {
	Role *string `json:"role" binding:"omitempty,min=1,max=20"`
}

type UpdateRoleInput struct {
	Name *string `json:"name" binding:"omitempty,min=1,max=20"`
}

type UpdatePlayListInput struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=20"`
	Description *string `json:"description" binding:"omitempty,min=1,max=50"`
}

func (i UpdateGroupInput) Validate() error {
	if i.Name == nil && i.Description == nil {
		return errors.New("update structure has no value")
	}
	return nil
}

func (i UpdateUserGroupInput) Validate() error {
	if i.Role == nil {
		return errors.New("update structure has no value")
	}
	return nil
}

func (i UpdateRoleInput) Validate() error {
	if i.Name == nil {
		return errors.New("update structure has no value")
	}
	return nil
}

func (i UpdatePlayListInput) Validate() error {
	if i.Name == nil && i.Description == nil {
		return errors.New("update structure has no value")
	}
	return nil
}
