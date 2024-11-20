package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	CreateUserSpotify(spotifyUri string, token models.ApiToken, user models.User) (int, error)
	CreateUserYouTubeMusic(token string, user models.User) (int, error)
	GetUser(email string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	GetUserSpotifyByID(userID int) (*models.UserSpotify, error)
}

type Group interface {
	Create(userID, roleID int, group models.UserCreateGroup) (int, error)
	GetAll() ([]models.Group, error)
	GetById(id int) (*models.Group, error)
	Update(id int, group models.UpdateGroupInput) error
	Delete(id int) error
}

type UserGroup interface {
	Create(userGroup models.UserGroup) (int, error)
	GetAll() ([]models.AllUserGroupData, error)
	GetById(id int) (*models.AllUserGroupData, error)
	GetUsersByGroupId(id int) ([]models.UserGroupData, error)
	GetByGroupIdAndUserIDRole(groupID, userID int, role string) (*models.UserGroup, error)
	GetByGroupIdAndUserID(groupID, userID int) (*models.UserGroup, error)
	GetByGroupIdAndUserIDAllData(groupID, userID int) (*models.UserGroupData, error)
	GetByGroupIdSpotifyUser(id int) ([]models.UserGroupToken, error)
	GetByGroupIdYouTubeMusicUser(id int) ([]models.UserGroupToken, error)
	Update(id, updateUserID, roleID int) error
	Delete(userID, groupID int) error
}

type Role interface {
	Create(group models.Role) (int, error)
	GetAll() ([]models.Role, error)
	GetById(id int) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Update(id int, group models.UpdateRoleInput) error
	Delete(id int) error
}

type PlayList interface {
	Create(group models.PlayList) (int, error)
	GetAll() ([]models.PlayList, error)
	GetById(id int) (*models.PlayList, error)
	GetByGroupId(id int) (*models.PlayList, error)
	Update(id int, group models.UpdatePlayListInput) error
	Delete(id int) error
}

type Track interface {
	Create(track models.CreateTrack) (int, error)
	AddYouTubeMusicTrackPlayList(playListID int, track models.CreateTrack) (int, error)
	AddSpotifyTrackPlayList(playListID int, track models.CreateTrack) (int, error)
	GetAll() ([]models.Track, error)
	GetByID(ID int) (*models.Track, error)
	GetByPlayListID(playListID int) ([]models.Track, error)
	GetByPlayListTrackID(playListID, trackID int) ([]models.PlayListTrack, error)
	GetByPlayListIDAndTrackApiID(playListID int, apiID models.ApiTrackID) (*models.Track, error)
	GetByTrackApiID(apiID models.ApiTrackID) (*models.Track, error)
	Delete(ID int) error
	DeleteFromPlayList(playListID, trackID int) error
	DeleteFromYouTubeMusicPlayList(userID, playListID, trackID int) error
}

type Token interface {
	GetSpotifyToken(spotifyUri string) (models.ApiToken, error)
	GetYouTubeMusicToken(userID int) (models.ApiToken, error)
	Create(token models.Token) (int, error)
	GetByToken(token string) (*models.Token, error)
	Update(token string) error
	UpdateYouTubeAccessToken(token string, userID int) error
	RevokeAllUserTokens(userID int) error
	UpdateSpotifyTokenBySpotifyUri(accessToken, spotifyUri string) error
}

type Repository struct {
	Authorization
	Group
	UserGroup
	Role
	PlayList
	Track
	Token
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Group:         NewGroupPostgres(db),
		UserGroup:     NewUserGroupPostgres(db),
		Role:          NewRolePostgres(db),
		PlayList:      NewPlayListPostgres(db),
		Track:         NewTrackPostgres(db),
		Token:         NewTokenPostgres(db),
	}
}
