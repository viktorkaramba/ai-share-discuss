package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"PlaylistsSynchronizer.Backend/configs"
	"PlaylistsSynchronizer.Backend/pkg/models"
	"PlaylistsSynchronizer.Backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) spotifyLogin(c *gin.Context) {
	oauthState := utils.GenerateStateOauthCookie(c.Writer)
	u := configs.AppConfig.SpotifyLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
}

func (h *Handler) spotifyCallBack(c *gin.Context) {
	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Writer.Header().Add("content-type", "application/json")
	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		models.NewErrorResponse(c, http.StatusInternalServerError, "invalid oauth spotify state")
		return
	}
	// Exchange Auth Code for Tokens
	token, err := configs.AppConfig.SpotifyLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("falied code exchange: %s", err.Error()))
		return
	}

	// Fetch User Data from spotify server
	client := http.Client{}
	request, err := http.NewRequest("GET", configs.OauthSpotifyUrlAPI, nil)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	request.Header = http.Header{
		"Access-Control-Allow-Origin": {"*"},
		"Content-Type":                {"application/json"},
		"Authorization":               {"Bearer " + token.AccessToken},
	}
	response, err := client.Do(request)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Parse user data JSON Object
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	var oauthResponse map[string]interface{}
	err = json.Unmarshal(contents, &oauthResponse)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := models.User{
		Username: oauthResponse["display_name"].(string),
		Email:    oauthResponse["email"].(string),
		Platform: models.Spotify,
	}
	isUserExist, err := h.services.Authorization.GetUser(user.Email)
	if isUserExist == nil {
		spotifyUri := oauthResponse["id"].(string)
		spotifyTokens := models.ApiToken{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}
		id, err := h.services.Authorization.CreateUserSpotify(spotifyUri, spotifyTokens, user)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		user.ID = id
	} else {
		err = h.services.RevokeAllUserTokens(isUserExist.ID)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		user.ID = isUserExist.ID
	}
	userToken, err := h.services.Authorization.GenerateToken(user.Email)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newToken := models.Token{TokenValue: userToken, Revoked: false, UserID: user.ID}
	_, err = h.services.Token.Create(newToken)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.AccessTokenResponse{AccessToken: newToken.TokenValue})
}

func (h *Handler) youTubeMusicLogin(c *gin.Context) {
	oauthState := utils.GenerateStateOauthCookie(c.Writer)
	u := configs.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
}

func (h *Handler) youTubeMusicCallBack(c *gin.Context) {
	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Writer.Header().Add("content-type", "application/json")
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		models.NewErrorResponse(c, http.StatusInternalServerError, "invalid oauth google state")
		return
	}
	// Exchange Auth Code for Tokens
	token, err := configs.AppConfig.GoogleLoginConfig.Exchange(
		context.Background(), code)

	// ERROR : Auth Code Exchange Failed
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("falied code exchange: %s", err.Error()))
		return
	}
	// Fetch User Data from google server
	response, err := http.Get(configs.OauthGoogleUrlAPI + token.AccessToken)
	// ERROR : Unable to get user data from google
	if err != nil {
		fmt.Fprintf(c.Writer, "failed getting user info: %s", err.Error())
		return
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	var oauthResponse map[string]interface{}
	err = json.Unmarshal(contents, &oauthResponse)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user := models.User{
		Username: oauthResponse["name"].(string),
		Email:    oauthResponse["email"].(string),
		Platform: models.YouTubeMusic,
	}
	isUserExist, err := h.services.Authorization.GetUser(user.Email)
	if isUserExist == nil {
		id, err := h.services.Authorization.CreateUserYouTubeMusic(token.AccessToken, user)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		user.ID = id
	} else {
		err = h.services.Token.RevokeAllUserTokens(isUserExist.ID)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		err = h.services.Token.UpdateYouTubeAccessToken(token.AccessToken, isUserExist.ID)
		if err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		user.ID = isUserExist.ID
	}
	userToken, err := h.services.Authorization.GenerateToken(user.Email)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newToken := models.Token{TokenValue: userToken, Revoked: false, UserID: user.ID}
	_, err = h.services.Token.Create(newToken)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, models.AccessTokenResponse{AccessToken: newToken.TokenValue})
}

func (h *Handler) appleMusicLogin(c *gin.Context) {
	//TODO
}

func (h *Handler) appleMusicCallBack(c *gin.Context) {
	//TODO
}

// @Summary RefreshToken
// @Tags auth
// @Description refresh access token
// @ID refresh-token
// @Accept json
// @Produce json
// @Param input body models.RefreshTokenInput true "user id for token refresh"
// @Success 200 {object} models.AccessTokenResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /refresh-token [post]
func (h *Handler) refreshToken(c *gin.Context) {
	var input models.RefreshTokenInput
	if err := c.BindJSON(&input); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userToken, err := h.services.RefreshToken(input.UserId)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newToken := models.Token{TokenValue: userToken, Revoked: false, UserID: input.UserId}
	_, err = h.services.Token.Create(newToken)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.AccessTokenResponse{AccessToken: newToken.TokenValue})
}

// @Summary Logout
// @Tags auth
// @Description logout from service
// @ID logout
// @Produce json
// @Success 200 {object} models.StatusResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	token, err := checkHeaderToken(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.services.Token.Update(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
}
