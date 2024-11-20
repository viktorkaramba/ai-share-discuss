package services

import (
	"PlaylistsSynchronizer.Backend/pkg/api_services"
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"errors"
)

type UserGroupService struct {
	repoAuth      repositories.Authorization
	repoUserGroup repositories.UserGroup
	repoPlayList  repositories.PlayList
	repoRole      repositories.Role
	repoTrack     repositories.Track
	repoToken     repositories.Token
	apiService    api_services.ApiService
}

func NewUserGroupService(repoAuth repositories.Authorization, repoUserGroup repositories.UserGroup,
	repoPlayList repositories.PlayList, repoRole repositories.Role, repoTrack repositories.Track,
	repoToken repositories.Token) *UserGroupService {
	return &UserGroupService{
		repoAuth:      repoAuth,
		repoUserGroup: repoUserGroup,
		repoPlayList:  repoPlayList,
		repoRole:      repoRole,
		repoTrack:     repoTrack,
		repoToken:     repoToken,
		apiService:    *api_services.NewApiService(repoToken),
	}
}

func (s *UserGroupService) Create(platform string, group models.UserGroup) (int, error) {
	isAlreadyExist, err := s.isUserAlreadyExist(group)
	if err != nil {
		return 0, err
	}
	if isAlreadyExist {
		return 0, errors.New("you are already member of this group")
	}
	role, err := s.repoRole.GetByName("USER")
	group.RoleID = role.ID
	if err != nil {
		return 0, err
	}
	playList, err := s.repoPlayList.GetByGroupId(group.GroupID)
	if err != nil {
		return 0, err
	}
	tracks, err := s.repoTrack.GetByPlayListID(playList.ID)
	if err != nil {
		return 0, err
	}
	var playListID string
	switch platform {
	case models.Spotify:
		spotify := s.apiService.GetSpotifyServiceApi()
		userSpotify, err := s.repoAuth.GetUserSpotifyByID(group.UserID)
		if err != nil {
			return 0, err
		}
		spotifyData := models.SpotifyData{Token: userSpotify.AccessToken, SpotifyUri: userSpotify.SpotifyUri}
		playListID, err = spotify.CreatePlayList(spotifyData, *playList)
		if err != nil {
			return 0, err
		}
		if len(tracks) != 0 {
			err = spotify.AddTrack(spotifyData, playListID, tracks)
			if err != nil {
				return 0, err
			}
		}
	case models.YouTubeMusic:
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		token, err := s.repoToken.GetYouTubeMusicToken(group.UserID)
		playListID, err = youTubeMusic.CreatePlayList(token.AccessToken, *playList)
		if err != nil {
			return 0, err
		}
		if len(tracks) != 0 {
			for _, value := range tracks {
				playListAPIID, err := youTubeMusic.AddTrack(token.AccessToken, playListID, value)
				if err != nil {
					return 0, err
				}
				newTrack := models.CreateTrack{SpotifyUri: value.SpotifyUri, YouTubeMusicID: value.YouTubeMusicID,
					UserID: group.UserID, PlayListYouTubeMusicID: playListAPIID, ID: value.ID}
				_, err = s.repoTrack.AddYouTubeMusicTrackPlayList(playList.ID, newTrack)
				if err != nil {
					return 0, err
				}
			}
		}
	}
	group.PlayListID = playListID
	return s.repoUserGroup.Create(group)
}

func (s *UserGroupService) GetAll() ([]models.AllUserGroupData, error) {
	return s.repoUserGroup.GetAll()
}

func (s *UserGroupService) GetById(id int) (*models.AllUserGroupData, error) {
	return s.repoUserGroup.GetById(id)
}

func (s *UserGroupService) GetUsersByGroupId(id int) ([]models.UserGroupData, error) {
	return s.repoUserGroup.GetUsersByGroupId(id)
}

func (s *UserGroupService) GetByGroupIdAndUserIDRole(groupID, userID int, role string) (*models.UserGroup, error) {
	return s.repoUserGroup.GetByGroupIdAndUserIDRole(groupID, userID, role)
}

func (s *UserGroupService) GetByGroupIdAndUserIDAllData(groupID, userID int) (*models.UserGroupData, error) {
	return s.repoUserGroup.GetByGroupIdAndUserIDAllData(groupID, userID)
}

func (s *UserGroupService) GetByGroupIdAndUserID(groupID, userID int) (*models.UserGroup, error) {
	return s.repoUserGroup.GetByGroupIdAndUserID(groupID, userID)
}

func (s *UserGroupService) Update(id, updateUserID int, group models.UpdateUserGroupInput) error {
	if err := group.Validate(); err != nil {
		return err
	}
	role, err := s.repoRole.GetByName(*group.Role)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("such a role does not exist")
	}
	return s.repoUserGroup.Update(id, updateUserID, role.ID)
}

func (s *UserGroupService) Delete(platform string, userID, groupID int) error {
	userGroupData, err := s.repoUserGroup.GetByGroupIdAndUserID(groupID, userID)
	if err != nil {
		return err
	}
	switch platform {
	case models.Spotify:
		spotify := s.apiService.GetSpotifyServiceApi()
		userSpotify, err := s.repoAuth.GetUserSpotifyByID(userID)
		if err != nil {
			return err
		}
		spotifyData := models.SpotifyData{Token: userSpotify.AccessToken, SpotifyUri: userSpotify.SpotifyUri}
		err = spotify.DeletePlayList(spotifyData, userGroupData.PlayListID)
		if err != nil {
			return err
		}
	case models.YouTubeMusic:
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		token, err := s.repoToken.GetYouTubeMusicToken(userID)
		err = youTubeMusic.DeletePlayList(token.AccessToken, userGroupData.PlayListID)
		if err != nil {
			return err
		}
	}
	return s.repoUserGroup.Delete(userID, groupID)
}

func (s *UserGroupService) isUserAlreadyExist(group models.UserGroup) (bool, error) {
	isUserMember, err := s.repoUserGroup.GetByGroupIdAndUserID(group.GroupID, group.UserID)
	if err != nil {
		return false, err
	}
	if isUserMember != nil {
		return true, nil
	}
	return false, nil
}
