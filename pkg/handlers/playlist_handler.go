package handlers

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// @Summary Get All Playlist
// @Security ApiKeyAuth
// @Tags playlists
// @Description get all playlists
// @ID get-all-playlists
// @Produce json
// @Success 200 {object} []models.PlayList
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/playlists [get]
func (h *Handler) getAllPlayList(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	playLists, err := h.services.PlayList.GetAll()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, playLists)
}

// @Summary Get Playlist By ID
// @Security ApiKeyAuth
// @Tags playlists
// @Description get playlist by id
// @ID get-playlist-by-id
// @Produce json
// @Param id path int true "playlist id"
// @Success 200 {object} models.PlayList
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/playlists/{id} [get]
func (h *Handler) getPlayListById(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	isPlayListExist, err := h.services.PlayList.GetById(id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isPlayListExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no playlist with such id")
		return
	}
	playlist, err := h.services.PlayList.GetById(id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// @Summary Update Playlist
// @Security ApiKeyAuth
// @Tags playlists
// @Description update playlist
// @ID update-playlist
// @Accept json
// @Produce json
// @Param id path int true "playlist id"
// @Param input body models.UpdatePlayListInput true "playlist info"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/playlists/{id} [put]
func (h *Handler) updatePlayList(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	playListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid playListID param")
		return
	}
	isPlayListExist, err := h.services.PlayList.GetById(playListID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isPlayListExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no playlist with such id")
		return
	}
	isValidUser, err := h.isValidAdmin(isPlayListExist.GroupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		// Check if there are any additional fields in the JSON body
		if err := h.validateJSONTags(body, models.UpdatePlayListInput{}); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		var input models.UpdatePlayListInput

		if err := c.BindJSON(&input); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
			return
		}
		if input.Name == nil {
			input.Name = &isPlayListExist.Name
		}
		if err := h.services.PlayList.Update(playListID, input); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}
