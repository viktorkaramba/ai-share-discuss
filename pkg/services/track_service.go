package services

import (
	"PlaylistsSynchronizer.Backend/pkg/api_services"
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
	"errors"
)

type TrackService struct {
	repoTrack    repositories.Track
	repoGroup    repositories.UserGroup
	repoToken    repositories.Token
	repoPlaylist repositories.PlayList
	apiService   api_services.ApiService
}

func NewTrackService(repoTrack repositories.Track, repoGroup repositories.UserGroup, repoToken repositories.Token,
	repoPlaylist repositories.PlayList) *TrackService {
	return &TrackService{
		repoTrack:    repoTrack,
		repoGroup:    repoGroup,
		repoPlaylist: repoPlaylist,
		repoToken:    repoToken,
		apiService:   *api_services.NewApiService(repoToken),
	}
}

func (s *TrackService) GetByPlayListTrackID(playListID, trackID int) ([]models.PlayListTrack, error) {
	return s.repoTrack.GetByPlayListTrackID(playListID, trackID)
}

func (s *TrackService) Add(input models.AddTrack) (int, error) {
	spotify := s.apiService.GetSpotifyServiceApi()
	usersSpotify, err := s.repoGroup.GetByGroupIdSpotifyUser(input.GroupID)
	if err != nil {
		return 0, err
	}
	if len(usersSpotify) != 0 {
		err = spotify.GetTrack(models.SpotifyData{Token: usersSpotify[0].Token, SpotifyUri: usersSpotify[0].SpotifyUri},
			models.Track{SpotifyUri: input.SpotifyUri})
		if err != nil {
			return 0, err
		}
	}
	isTrackExist, err := s.repoTrack.GetByTrackApiID(models.ApiTrackID{SpotifyUri: input.SpotifyUri,
		YouTubeMusic: input.YouTubeMusicID})
	if err != nil {
		return 0, err
	}
	var track models.CreateTrack
	track = models.CreateTrack{SpotifyUri: input.SpotifyUri, YouTubeMusicID: input.YouTubeMusicID, Name: input.Name}
	if isTrackExist == nil {
		id, err := s.repoTrack.Create(track)
		if err != nil {
			return 0, err
		}
		track.ID = id
	} else {
		track.ID = isTrackExist.ID
	}
	isTrackPlayListExist, err := s.repoTrack.GetByPlayListIDAndTrackApiID(
		input.PlayListID, models.ApiTrackID{SpotifyUri: input.SpotifyUri,
			YouTubeMusic: input.YouTubeMusicID})
	if err != nil {
		return 0, err
	}
	if isTrackPlayListExist == nil {
		_, err = s.repoTrack.AddSpotifyTrackPlayList(input.PlayListID, track)
		if err != nil {
			return 0, err
		}
		id, err := s.addSpotify(track, usersSpotify)
		if err != nil {
			return 0, err
		}
		_, err = s.addYouTubeMusic(track, input.PlayListID, input.GroupID)
		if err != nil {
			return 0, err
		}
		return id, nil
	} else {
		return 0, errors.New("track already exist in playlist")
	}
}

func (s *TrackService) addSpotify(track models.CreateTrack, usersSpotify []models.UserGroupToken) (int, error) {
	for _, value := range usersSpotify {
		spotify := s.apiService.GetSpotifyServiceApi()
		err := spotify.AddTrack(models.SpotifyData{Token: value.Token, SpotifyUri: value.SpotifyUri}, value.PlayListID,
			[]models.Track{{SpotifyUri: track.SpotifyUri,
				YouTubeMusicID: track.YouTubeMusicID}})
		if err != nil {
			return 0, err
		}
	}
	return track.ID, nil
}

func (s *TrackService) addYouTubeMusic(track models.CreateTrack, playListID, groupID int) (int, error) {
	usersYouTubeMusic, err := s.repoGroup.GetByGroupIdYouTubeMusicUser(groupID)
	if err != nil {
		return 0, err
	}
	var tracksID string
	for _, value := range usersYouTubeMusic {
		youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
		tracksID, err = youTubeMusic.AddTrack(value.Token, value.PlayListID, models.Track{SpotifyUri: track.SpotifyUri,
			YouTubeMusicID: track.YouTubeMusicID})
		if err != nil {
			return 0, err
		}
		newTrack := models.CreateTrack{SpotifyUri: track.SpotifyUri, YouTubeMusicID: track.YouTubeMusicID,
			UserID: value.UserID, PlayListYouTubeMusicID: tracksID, ID: track.ID}
		_, err := s.repoTrack.AddYouTubeMusicTrackPlayList(playListID, newTrack)
		if err != nil {
			return 0, err
		}
	}
	return track.ID, nil
}

func (s *TrackService) DeleteFromPlayList(playListID, groupID, trackID int) error {
	usersSpotify, err := s.repoGroup.GetByGroupIdSpotifyUser(groupID)
	if err != nil {
		return err
	}
	track, err := s.repoTrack.GetByID(trackID)
	if err != nil {
		return err
	}
	err = s.deleteSpotify(*track, usersSpotify)
	if err != nil {
		return err
	}
	playListTracks, err := s.repoTrack.GetByPlayListTrackID(playListID, trackID)
	for _, value := range playListTracks {
		err = s.deleteYouTubeMusic(value.Token, value.PlayListYouTubeMusicID)
		if err != nil {
			return err
		}
		err = s.repoTrack.DeleteFromYouTubeMusicPlayList(value.UserID, playListID, trackID)
		if err != nil {
			return err
		}
	}
	err = s.repoTrack.DeleteFromPlayList(playListID, trackID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TrackService) deleteSpotify(track models.Track, usersSpotify []models.UserGroupToken) error {
	for _, value := range usersSpotify {
		spotify := s.apiService.GetSpotifyServiceApi()
		err := spotify.DeleteTrack(models.SpotifyData{Token: value.Token, SpotifyUri: value.SpotifyUri},
			value.PlayListID, []models.Track{track})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TrackService) deleteYouTubeMusic(token, trackID string) error {
	youTubeMusic := s.apiService.GetYouTubeMusicApiServiceApi()
	err := youTubeMusic.DeleteTrack(token, trackID)
	if err != nil {
		return err
	}
	return nil
}
