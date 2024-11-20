package services

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	CreateUserSpotify(spotifyUri string, token models.ApiToken, user models.User) (int, error)
	CreateUserYouTubeMusic(token string, user models.User) (int, error)
	GetUser(email string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	GenerateToken(email string) (string, error)
	ParseToken(token string) (*models.UserClaims, error)
	RefreshToken(userId int) (string, error)
}

type Group interface {
	Create(userID int, group models.UserCreateGroup) (int, error)
	GetAll() ([]models.Group, error)
	GetById(id int) (*models.AllGroupData, error)
	Update(id int, group models.UpdateGroupInput) error
	Delete(id int) error
}

type UserGroup interface {
	Create(platform string, userGroup models.UserGroup) (int, error)
	GetAll() ([]models.AllUserGroupData, error)
	GetById(id int) (*models.AllUserGroupData, error)
	GetUsersByGroupId(id int) ([]models.UserGroupData, error)
	GetByGroupIdAndUserIDRole(groupID, userID int, role string) (*models.UserGroup, error)
	GetByGroupIdAndUserIDAllData(groupID, userID int) (*models.UserGroupData, error)
	GetByGroupIdAndUserID(groupID, userID int) (*models.UserGroup, error)
	Update(id, updateUserID int, group models.UpdateUserGroupInput) error
	Delete(platform string, userID, groupID int) error
}

type Role interface {
	Create(role models.Role) (int, error)
	GetAll() ([]models.Role, error)
	GetById(id int) (*models.Role, error)
	Update(id int, role models.UpdateRoleInput) error
	Delete(id int) error
}

type PlayList interface {
	GetAll() ([]models.PlayList, error)
	GetById(id int) (*models.PlayList, error)
	Update(id int, group models.UpdatePlayListInput) error
}

type Track interface {
	Add(input models.AddTrack) (int, error)
	GetByPlayListTrackID(playListID, trackID int) ([]models.PlayListTrack, error)
	DeleteFromPlayList(groupID, playListID, trackID int) error
}

type Token interface {
	Create(token models.Token) (int, error)
	GetByToken(token string) (*models.Token, error)
	Update(token string) error
	UpdateYouTubeAccessToken(token string, userID int) error
	RevokeAllUserTokens(userID int) error
}

type Service struct {
	Authorization
	Group
	UserGroup
	Role
	PlayList
	Track
	Token
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, repos.Token),
		Group: NewGroupService(repos.Authorization, repos.Group, repos.Role, repos.Token, repos.UserGroup,
			repos.Track),
		UserGroup: NewUserGroupService(repos.Authorization, repos.UserGroup, repos.PlayList, repos.Role, repos.Track,
			repos.Token),
		Role:     NewRoleService(repos.Role),
		PlayList: NewPlayListService(repos.PlayList, repos.Token, repos.UserGroup),
		Track:    NewTrackService(repos.Track, repos.UserGroup, repos.Token, repos.PlayList),
		Token:    NewTokenService(repos.Token),
	}
}
