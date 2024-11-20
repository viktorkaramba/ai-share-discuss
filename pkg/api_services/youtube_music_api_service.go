package api_services

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type YouTubeMusicApiService struct {
}

func NewYouTubeMusicApiService() *YouTubeMusicApiService {
	return &YouTubeMusicApiService{}
}

func (s *YouTubeMusicApiService) CreatePlayList(token string, playList models.PlayList) (string, error) {
	query := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=snippet%%2Cstatus&key=%s",
		os.Getenv("API_KEY"))
	jsonBody := []byte(fmt.Sprintf(
		`{"snippet": { "title": "%s", "description": "%s", "tags": ["tag"], "defaultLanguage": "uk" }}`,
		playList.Name, playList.Description))
	bodyReader := bytes.NewReader(jsonBody)
	client := http.Client{}
	request, err := http.NewRequest("POST", query, bodyReader)
	if err != nil {
		return "", err
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	var oauthResponse map[string]interface{}
	err = json.Unmarshal(contents, &oauthResponse)
	if err != nil {
		return "", err
	}
	if _, ok := oauthResponse["error"]; ok {
		someErrorMap := oauthResponse["error"].(map[string]interface{})
		if _, ok := someErrorMap["message"]; ok {
			msg := someErrorMap["message"].(string)
			return "", errors.New("api error: " + msg)
		}
	}
	playListID := oauthResponse["id"].(string)
	return playListID, nil
}

func (s *YouTubeMusicApiService) UpdatePlayList(token string, playListID string, updatePlaylist models.UpdatePlayListInput) error {
	query := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=snippet%%2Cstatus&key=%s",
		os.Getenv("API_KEY"))
	description := ""
	if updatePlaylist.Description != nil {
		description = "\"description\": \"" + *updatePlaylist.Description + "\", "
	}
	jsonBody := []byte(fmt.Sprintf(
		`{ "id": "%s", "snippet": { "title": "%s", %s"tags": ["tag"]}}`,
		playListID, *updatePlaylist.Name, description))
	bodyReader := bytes.NewReader(jsonBody)
	client := http.Client{}
	request, err := http.NewRequest("PUT", query, bodyReader)
	if err != nil {
		return err
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
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
			return errors.New("api error: " + msg)
		}
	}
	return nil
}

func (s *YouTubeMusicApiService) DeletePlayList(token string, playListID string) error {
	query := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?id=%s&key=%s", playListID,
		os.Getenv("API_KEY"))
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", query, nil)
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 204 {
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
				return errors.New("api error: " + msg)
			}
		}
	}
	return nil
}

func (s *YouTubeMusicApiService) DeleteTrack(token string, trackID string) error {
	query := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?id=%s&key=%s", trackID,
		os.Getenv("API_KEY"))
	client := http.Client{}
	request, err := http.NewRequest("DELETE", query, nil)
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 204 {
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
				return errors.New("api error: " + msg)
			}
		}
	}
	return nil
}

func (s *YouTubeMusicApiService) AddTrack(token string, playListId string, track models.Track) (string, error) {
	var tracksID string
	query := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&key=%s",
		os.Getenv("API_KEY"))
	jsonBody := []byte(fmt.Sprintf(`{"snippet": { "playlistId": "%s", "resourceId":{ "kind": "youtube#video", "videoId": "%s"} }}`,
		playListId, track.YouTubeMusicID))
	bodyReader := bytes.NewReader(jsonBody)
	client := http.Client{}
	request, err := http.NewRequest("POST", query, bodyReader)
	if err != nil {
		return "", err
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
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
			if strings.Contains(msg, "Video not found") {
				return "", errors.New("invalid youtube music id")
			} else {
				return "", errors.New("api error: " + msg)
			}
		}
	}
	tracksID = jsonMap["id"].(string)
	return tracksID, nil
}
