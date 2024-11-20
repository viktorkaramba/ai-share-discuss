package services

import (
	"PlaylistsSynchronizer.Backend/pkg/api_services"
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"fmt"
)

type PlayListService struct {
	repoPlayList repositories.PlayList
	repoToken    repositories.Token
	repoGroup    repositories.UserGroup
	apiService   api_services.ApiService
}

func NewPlayListService(repoPlayList repositories.PlayList, repoToken repositories.Token,
	repoGroup repositories.UserGroup) *PlayListService {
	return &PlayListService{
		repoPlayList: repoPlayList,
		repoToken:    repoToken,
		repoGroup:    repoGroup,
		apiService:   *api_services.NewApiService(repoToken),
	}
}

func (s *PlayListService) GetAll() ([]models.PlayList, error) {
	return s.repoPlayList.GetAll()
}

func (s *PlayListService) GetById(id int) (*models.PlayList, error) {
	return s.repoPlayList.GetById(id)
}

func (s *PlayListService) Update(id int, input models.UpdatePlayListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	playList, err := s.repoPlayList.GetById(id)
	if err != nil {
		return err
	}
	err = s.updateSpotify(playList.GroupID, input)
	if err != nil {
		return err
	}
	err = s.updateYouTubeMusic(playList.GroupID, input)
	if err != nil {
		return err
	}
	fmt.Println(input.Name)
	return s.repoPlayList.Update(id, input)
}

func (s *PlayListService) updateSpotify(groupID int, input models.UpdatePlayListInput) error {
	usersSpotify, err := s.repoGroup.GetByGroupIdSpotifyUser(groupID)
	if err != nil {
		return err
	}
	for _, value := range usersSpotify {
		spotify := s.apiService.GetSpotifyServiceApi()
		err := spotify.UpdatePlayList(models.SpotifyData{Token: value.Token, SpotifyUri: value.SpotifyUri},
			value.PlayListID, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PlayListService) updateYouTubeMusic(groupID int, input models.UpdatePlayListInput) error {
	usersYouTubeMusic, err := s.repoGroup.GetByGroupIdYouTubeMusicUser(groupID)
	if err != nil {
		return err
	}
	for _, value := range usersYouTubeMusic {
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		err = youTubeMusic.UpdatePlayList(value.Token, value.PlayListID, input)
		if err != nil {
			return err
		}
	}
	return nil
}
