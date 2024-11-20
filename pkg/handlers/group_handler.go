package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"PlaylistsSynchronizer.Backend/pkg/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create Group
// @Security ApiKeyAuth
// @Tags groups
// @Description create new group
// @ID create-group
// @Accept json
// @Produce json
// @Param input body models.UserCreateGroupInput true "group and playlist info"
// @Success 200 {object} models.CreateResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups [post]
func (h *Handler) createGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	platform, err := getUserPlatform(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	// Check if there are any additional fields in the JSON body
	if err := h.validateJSONTags(body, models.UserCreateGroupInput{}); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input models.UserCreateGroupInput
	if err := c.BindJSON(&input); err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	newInput := models.UserCreateGroup{GroupName: input.GroupName, PlayListName: input.PlayListName,
		GroupDescription: input.GroupDescription, PlayListDescription: input.PlayListDescription, Platform: platform}
	id, err := h.services.Group.Create(userID, newInput)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.CreateResponse{ID: id})
}

// @Summary Get All Groups
// @Security ApiKeyAuth
// @Tags groups
// @Description get all groups
// @ID get-all-groups
// @Produce json
// @Success 200 {object} []models.Group
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups [get]
func (h *Handler) getAllGroups(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}
	groups, err := h.services.Group.GetAll()
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(c.Request.Header.Values("Access-Control-Allow-Origin"))
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	//c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, groups)
}

// @Summary Get Group By ID
// @Security ApiKeyAuth
// @Tags groups
// @Description get group by id
// @ID get-group-by-id
// @Produce json
// @Param id path int true "group id"
// @Success 200 {object} models.AllGroupData
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id} [get]
func (h *Handler) getGroupById(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	group, err := h.services.Group.GetById(id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, group)
}

// @Summary Update Group
// @Security ApiKeyAuth
// @Tags groups
// @Description update group
// @ID update-group
// @Accept json
// @Produce json
// @Param id path int true "group id"
// @Param input body models.UpdateGroupInput true "group info"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id} [put]
func (h *Handler) updateGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid groupID param")
		return
	}
	isGroupExist, err := h.services.Group.GetById(groupID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isGroupExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no group with such id")
		return
	}
	isValidUser, err := h.isValidAdmin(groupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		if err := h.validateJSONTags(body, models.UpdateGroupInput{}); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		var input models.UpdateGroupInput
		if err := c.BindJSON(&input); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
			return
		}

		if err := h.services.Group.Update(groupID, input); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}

// @Summary Delete Group
// @Security ApiKeyAuth
// @Tags groups
// @Description delete group
// @ID delete-group
// @Produce json
// @Param id path int true "group id"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id} [delete]
func (h *Handler) deleteGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid groupID param")
		return
	}
	isGroupExist, err := h.services.Group.GetById(groupID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isGroupExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no group with such id")
		return
	}
	isValidUser, err := h.isValidUserRole(groupID, userID, "SUPER ADMIN")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser {
		if err := h.services.Group.Delete(groupID); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}
