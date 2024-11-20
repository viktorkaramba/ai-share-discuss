package api_services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"

	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/repositories"
)

type SpotifyApiService struct {
	repoToken repositories.Token
}

func NewSpotifyApiService(repoToken repositories.Token) *SpotifyApiService {
	return &SpotifyApiService{
		repoToken: repoToken,
	}
}

func (s *SpotifyApiService) CreatePlayList(spotifyData models.SpotifyData, playList models.PlayList) (string, error) {
	query := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", spotifyData.SpotifyUri)
	jsonBody := []byte(fmt.Sprintf(`{"name": "%s", "description": "%s", "public": %v}`, playList.Name,
		playList.Description, false))
	bodyReader := bytes.NewReader(jsonBody)
	client := http.Client{}
	request, err := http.NewRequest("POST", query, bodyReader)
	if err != nil {
		return "", err
	}

	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + spotifyData.Token},
	}
	response, err := client.Do(request)
	body, _ := io.ReadAll(response.Body)
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return "", err
	}
	if _, ok := jsonMap["error"]; ok {
		someErrorMap := jsonMap["error"].(map[string]interface{})
		if _, ok := someErrorMap["message"]; ok {
			msg := someErrorMap["message"].(string)
			if strings.Contains(msg, "The access token expired") {
				token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
				if err != nil {
					return "", err
				}
				spotifyData.Token = token
				_, err = s.CreatePlayList(spotifyData, playList)
				if err != nil {
					return "", err
				}
				return "", nil
			} else {
				return "", errors.New(msg)
			}
		}
	}
	playListID := jsonMap["id"].(string)
	return playListID, nil
}

func (s *SpotifyApiService) UpdatePlayList(spotifyData models.SpotifyData, playListID string,
	updatePlaylist models.UpdatePlayListInput) error {
	query := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playListID)
	description := ""
	if updatePlaylist.Description != nil {
		description = ", \"description\": \"" + *updatePlaylist.Description + "\""
	}
	jsonBody := []byte(fmt.Sprintf(
		`{ "name": "%s"%s }`, *updatePlaylist.Name, description))
	bodyReader := bytes.NewReader(jsonBody)
	client := http.Client{}
	request, err := http.NewRequest("PUT", query, bodyReader)
	if err != nil {
		return err
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + spotifyData.Token},
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			return err
		}
		if _, ok := jsonMap["error"]; ok {
			someErrorMap := jsonMap["error"].(map[string]interface{})
			if _, ok := someErrorMap["message"]; ok {
				msg := someErrorMap["message"].(string)
				if strings.Contains(msg, "The access token expired") {
					token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
					if err != nil {
						return err
					}
					spotifyData.Token = token
					err = s.UpdatePlayList(spotifyData, playListID, updatePlaylist)
					if err != nil {
						return err
					}
				} else {
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

func (s *SpotifyApiService) DeletePlayList(spotifyData models.SpotifyData, playListID string) error {
	query := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", playListID)
	client := http.Client{}
	request, err := http.NewRequest("DELETE", query, nil)
	if err != nil {
		return err
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + spotifyData.Token},
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			return err
		}
		if _, ok := jsonMap["error"]; ok {
			someErrorMap := jsonMap["error"].(map[string]interface{})
			if _, ok := someErrorMap["message"]; ok {
				msg := someErrorMap["message"].(string)
				if strings.Contains(msg, "The access token expired") {
					token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
					if err != nil {
						return err
					}
					spotifyData.Token = token
					err = s.DeletePlayList(spotifyData, playListID)
					if err != nil {
						return err
					}
					return nil
				} else {
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

func (s *SpotifyApiService) AddTrack(spotifyData models.SpotifyData, playListId string, tracks []models.Track) error {
	query := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playListId)
	var spotifyUris []string

	for _, track := range tracks {
		spotifyUris = append(spotifyUris, "spotify:track:"+track.SpotifyUri)
	}
	spotifyUrisObj := models.SpotifyUris{
		Uris:     spotifyUris,
		Position: 0,
	}
	countOfAdding := len(tracks) / 50
	countOfAdding = int(math.Ceil(float64(countOfAdding*100)) / 100)
	if countOfAdding == 0 {
		countOfAdding = 1
	}
	for i := 0; i < countOfAdding; i++ {
		jsonData, err := json.MarshalIndent(spotifyUrisObj, "", "    ")
		if err != nil {
			return err
		}
		bodyReader := bytes.NewReader(jsonData)
		client := http.Client{}
		request, err := http.NewRequest("POST", query, bodyReader)
		if err != nil {
			return err
		}
		request.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + spotifyData.Token},
		}
		response, err := client.Do(request)
		body, _ := io.ReadAll(response.Body)
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			return err
		}
		if _, ok := jsonMap["error"]; ok {
			someErrorMap := jsonMap["error"].(map[string]interface{})
			if _, ok := someErrorMap["message"]; ok {
				msg := someErrorMap["message"].(string)
				if strings.Contains(msg, "The access token expired") {
					token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
					if err != nil {
						return err
					}
					spotifyData.Token = token
					err = s.AddTrack(spotifyData, playListId, tracks)
					if err != nil {
						return err
					}
					return nil
				} else if strings.Contains(msg, "Invalid track uri:") {
					return errors.New("invalid spotify track uri")
				} else {
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

func (s *SpotifyApiService) DeleteTrack(spotifyData models.SpotifyData, playListId string, tracks []models.Track) error {
	query := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playListId)
	var spotifyUris []models.DeleteSpotifyTrack

	for _, track := range tracks {
		spotifyUris = append(spotifyUris, models.DeleteSpotifyTrack{URI: "spotify:track:" + track.SpotifyUri})
	}
	spotifyUrisObj := models.DeleteSpotifyTrackInput{
		Tracks: spotifyUris,
	}
	countOfAdding := len(tracks) / 100
	countOfAdding = int(math.Ceil(float64(countOfAdding*100)) / 100)
	if countOfAdding == 0 {
		countOfAdding = 1
	}
	for i := 0; i < countOfAdding; i++ {
		jsonData, err := json.MarshalIndent(spotifyUrisObj, "", "    ")
		if err != nil {
			return err
		}
		bodyReader := bytes.NewReader(jsonData)
		client := http.Client{}
		request, err := http.NewRequest("DELETE", query, bodyReader)
		if err != nil {
			return err
		}
		request.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + spotifyData.Token},
		}
		response, err := client.Do(request)
		body, _ := io.ReadAll(response.Body)
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			return err
		}
		if _, ok := jsonMap["error"]; ok {
			someErrorMap := jsonMap["error"].(map[string]interface{})
			if _, ok := someErrorMap["message"]; ok {
				msg := someErrorMap["message"].(string)
				if strings.Contains(msg, "The access token expired") {
					token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
					if err != nil {
						return err
					}
					spotifyData.Token = token
					err = s.DeleteTrack(spotifyData, playListId, tracks)
					if err != nil {
						return err
					}
					return nil
				} else if strings.Contains(msg, "Invalid track uri:") {
					return errors.New("invalid spotify track uri")
				} else {
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

func (s *SpotifyApiService) GetTrack(spotifyData models.SpotifyData, track models.Track) error {
	query := "https://api.spotify.com/v1/tracks/" + track.SpotifyUri
	client := http.Client{}
	request, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return err
	}
	request.Header = http.Header{
		"Access-Control-Allow-Origin": {"*"},
		"Content-Type":                {"application/json"},
		"Authorization":               {"Bearer " + spotifyData.Token},
	}
	response, err := client.Do(request)
	body, _ := io.ReadAll(response.Body)
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return err
	}
	if _, ok := jsonMap["error"]; ok {
		someErrorMap := jsonMap["error"].(map[string]interface{})
		if _, ok := someErrorMap["message"]; ok {
			msg := someErrorMap["message"].(string)
			if strings.Contains(msg, "The access token expired") {
				token, err := s.RefreshSpotifyToken(spotifyData.SpotifyUri)
				if err != nil {
					return err
				}
				spotifyData.Token = token
				err = s.GetTrack(spotifyData, track)
				if err != nil {
					return err
				}
				return nil
			} else if strings.Contains(msg, "invalid id") {
				return errors.New("invalid spotify track uri")
			} else {
				return errors.New(msg)
			}
		}
	}
	return nil
}

func (s *SpotifyApiService) RefreshSpotifyToken(spotifyUri string) (string, error) {
	token, err := s.repoToken.GetSpotifyToken(spotifyUri)
	if err != nil {
		return "", err
	}
	query := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", token.RefreshToken)
	encodedData := data.Encode()
	refreshToken := models.RefreshToken{Query: query, Body: strings.NewReader(encodedData)}
	client := http.Client{}
	request, err := http.NewRequest("POST", refreshToken.Query, refreshToken.Body)
	if err != nil {
		return "", err
	}
	base64Token := base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" +
		os.Getenv("SPOTIFY_CLIENT_SECRET")))
	header := "Basic " + base64Token
	request.Header = http.Header{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {header},
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(response.Body)
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return "", err
	}
	err = s.repoToken.UpdateSpotifyTokenBySpotifyUri(jsonMap["access_token"].(string), spotifyUri)
	if err != nil {
		return "", err
	}
	return jsonMap["access_token"].(string), err
}
