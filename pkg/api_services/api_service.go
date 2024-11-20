package api_services

import "PlaylistsSynchronizer.Backend/pkg/repositories"

type ApiService struct {
	spotify      SpotifyApiService
	youTubeMusic YouTubeMusicApiService
}

func NewApiService(repoToken repositories.Token) *ApiService {
	return &ApiService{
		spotify:      *NewSpotifyApiService(repoToken),
		youTubeMusic: *NewYouTubeMusicApiService(),
	}
}

func (a *ApiService) GetSpotifyServiceApi() SpotifyApiService {
	return a.spotify
}

func (a *ApiService) GetYouTubeMusicApiServiceApi() YouTubeMusicApiService {
	return a.youTubeMusic
}
