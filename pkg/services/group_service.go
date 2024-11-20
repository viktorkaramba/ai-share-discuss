package services

import (
	"PlaylistsSynchronizer.Backend/pkg/api_services"
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"errors"
)

type GroupService struct {
	reposAuth     repositories.Authorization
	repoGroup     repositories.Group
	repoRole      repositories.Role
	repoToken     repositories.Token
	repoUserGroup repositories.UserGroup
	repoTracks    repositories.Track
	apiService    api_services.ApiService
}

func NewGroupService(reposAuth repositories.Authorization, repoGroup repositories.Group, repoRole repositories.Role,
	repoToken repositories.Token, repoUserGroup repositories.UserGroup, repoTracks repositories.Track) *GroupService {
	return &GroupService{
		reposAuth:     reposAuth,
		repoGroup:     repoGroup,
		repoRole:      repoRole,
		repoToken:     repoToken,
		repoUserGroup: repoUserGroup,
		repoTracks:    repoTracks,
		apiService:    *api_services.NewApiService(repoToken),
	}
}

func (s *GroupService) Create(userID int, group models.UserCreateGroup) (int, error) {
	role, err := s.repoRole.GetByName("SUPER ADMIN")
	if err != nil {
		return 0, err
	}
	var playListID string
	playList := models.PlayList{Name: group.PlayListName, Description: group.PlayListDescription}
	switch group.Platform {
	case models.Spotify:
		spotify := s.apiService.GetSpotifyServiceApi()
		userSpotify, err := s.reposAuth.GetUserSpotifyByID(userID)
		if err != nil {
			return 0, err
		}
		spotifyData := models.SpotifyData{Token: userSpotify.AccessToken, SpotifyUri: userSpotify.SpotifyUri}
		playListID, err = spotify.CreatePlayList(spotifyData, playList)
		if err != nil {
			return 0, err
		}
	case models.YouTubeMusic:
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		token, err := s.repoToken.GetYouTubeMusicToken(userID)
		if err != nil {
			return 0, err
		}
		playListID, err = youTubeMusic.CreatePlayList(token.AccessToken, playList)
		if err != nil {
			return 0, err
		}
	}
	group.PlayListID = playListID
	id, err := s.repoGroup.Create(userID, role.ID, group)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *GroupService) GetAll() ([]models.Group, error) {
	return s.repoGroup.GetAll()
}

func (s *GroupService) GetById(id int) (*models.AllGroupData, error) {
	group, err := s.repoGroup.GetById(id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("there is no group with such id")
	}
	users, err := s.repoUserGroup.GetUsersByGroupId(group.ID)
	if err != nil {
		return nil, err
	}
	var allGroupData models.AllGroupData
	if len(users) != 0 {
		playListID := users[0].PlayListID
		tracks, err := s.repoTracks.GetByPlayListID(playListID)
		if err != nil {
			return nil, err
		}
		allGroupData = models.AllGroupData{Name: group.Name, Description: group.Description, Users: users,
			Tracks: tracks, PlayListID: playListID}
		return &allGroupData, nil
	} else {
		return nil, nil
	}
}

func (s *GroupService) Update(id int, group models.UpdateGroupInput) error {
	if err := group.Validate(); err != nil {
		return err
	}
	return s.repoGroup.Update(id, group)
}

func (s *GroupService) Delete(id int) error {
	err := s.DeleteSpotify(id)
	if err != nil {
		return err
	}
	err = s.DeleteYouTubeMusic(id)
	if err != nil {
		return err
	}
	return s.repoGroup.Delete(id)
}

func (s *GroupService) DeleteSpotify(id int) error {
	usersSpotify, err := s.repoUserGroup.GetByGroupIdSpotifyUser(id)
	if err != nil {
		return err
	}
	for _, value := range usersSpotify {
		spotify := s.apiService.GetSpotifyServiceApi()
		err := spotify.DeletePlayList(models.SpotifyData{Token: value.Token, SpotifyUri: value.SpotifyUri}, value.PlayListID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GroupService) DeleteYouTubeMusic(id int) error {
	usersYouTubeMusic, err := s.repoUserGroup.GetByGroupIdYouTubeMusicUser(id)
	if err != nil {
		return err
	}
	for _, value := range usersYouTubeMusic {
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		err = youTubeMusic.DeletePlayList(value.Token, value.PlayListID)
		if err != nil {
			return err
		}
	}
	return nil
}
