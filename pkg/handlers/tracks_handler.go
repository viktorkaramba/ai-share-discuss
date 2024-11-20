package handlers

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// @Summary Add Track
// @Security ApiKeyAuth
// @Tags tracks
// @Description add new track
// @ID add-track
// @Accept json
// @Produce json
// @Param id path int true "playlist id"
// @Param input body models.AddTrackInput true "track info"
// @Success 200 {object} models.CreateResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/playlists/{id}/tracks [post]
func (h *Handler) addTrack(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	playListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid playListID param")
		return
	}
	var input models.AddTrackInput

	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := h.validateJSONTags(body, models.AddTrackInput{}); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.BindJSON(&input); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	playList, err := h.services.PlayList.GetById(playListID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if playList == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no playlist with such id")
		return
	}
	isValidUser, err := h.isValidUser(playList.GroupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser {
		newTrack := models.AddTrack{SpotifyUri: input.SpotifyUri, YouTubeMusicID: input.YouTubeMusicID,
			Name: input.Name, PlayListID: playList.ID, GroupID: playList.GroupID, UserID: userID}
		id, err := h.services.Track.Add(newTrack)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, models.CreateResponse{ID: id})
	} else {
		models.NewErrorResponse(c, http.StatusBadRequest, "you are not member of the group")
		return
	}
}

// @Summary Delete Track
// @Security ApiKeyAuth
// @Tags tracks
// @Description delete track
// @ID delete-track
// @Produce json
// @Param id path int true "playlist id"
// @Param trackID path int true "track id"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/playlists/{id}/tracks/{trackID} [delete]
func (h *Handler) deleteTrack(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	playListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid playListID param")
		return
	}
	trackID, err := strconv.Atoi(c.Param("trackID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid trackID param")
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
	isTrackExist, err := h.services.Track.GetByPlayListTrackID(playListID, trackID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isTrackExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no track with such id")
		return
	}
	isValidUser, err := h.isValidAdmin(isPlayListExist.GroupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser {
		err := h.services.Track.DeleteFromPlayList(isPlayListExist.GroupID, playListID, trackID)
		if err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}
