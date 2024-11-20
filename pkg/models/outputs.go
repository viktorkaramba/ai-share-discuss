package models

type UserGroupData struct {
	ID         int    `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Email      string `json:"email" db:"email"`
	Platform   string `json:"platform" db:"platform"`
	RoleName   string `json:"roleName" db:"rolename"`
	PlayListID int    `json:"-" db:"playlistid"`
}

type AllUserGroupData struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"userID" db:"userid"`
	GroupID    int    `json:"groupID" db:"groupid"`
	Email      string `json:"email" db:"email"`
	RoleName   string `json:"roleName" db:"rolename"`
	PlayListID string `json:"-" db:"playlistid"`
}

type AllGroupData struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Users       []UserGroupData `json:"users"`
	Tracks      []Track         `json:"tracks"`
	PlayListID  int             `json:"playListID"`
}
