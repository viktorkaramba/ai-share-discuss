package handlers

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// @Summary Add User To Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description add user to group
// @ID add-user-to-group
// @Produce json
// @Param id path int true "group id"
// @Success 200 {object}  models.CreateResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/users [post]
func (h *Handler) createUserGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	platform, err := getUserPlatform(c)
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
	newUserGroup := models.UserGroup{GroupID: groupID, UserID: userID}
	id, err := h.services.UserGroup.Create(platform, newUserGroup)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.CreateResponse{ID: id})
}

// @Summary Get All Users From Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description get all users from group
// @ID get-all-users-from-group
// @Produce json
// @Param id path int true "group id"
// @Success 200 {object} []models.UserGroupData
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/users [get]
func (h *Handler) getAllUserGroups(c *gin.Context) {
	_, err := getUserId(c)
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
	userGroup, err := h.services.UserGroup.GetUsersByGroupId(groupID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userGroup)
}

// @Summary Get User By ID From Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description get user by id from group
// @ID get-user-by-id-from-group
// @Produce json
// @Param id path int true "group id"
// @Param userID path int true "user id"
// @Success 200 {object} models.UserGroupData
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/users/{userID} [get]
func (h *Handler) getUserGroupByUserId(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid groupID param")
		return
	}
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid userID param")
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

	userGroup, err := h.services.UserGroup.GetByGroupIdAndUserIDAllData(groupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if userGroup == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no user with such id in this group")
		return
	}
	c.JSON(http.StatusOK, userGroup)
}

// @Summary Update User By ID From Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description update user by id from group
// @ID update-user-by-id-from-group
// @Accept json
// @Produce json
// @Param id path int true "group id"
// @Param userID path int true "user id"
// @Param input body models.UpdateUserGroupInput true "new user role"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/users/{userID} [put]
func (h *Handler) updateUserGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid userGroupID param")
		return
	}
	updateUserID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid userGroupID param")
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
	isUserExist, err := h.services.UserGroup.GetByGroupIdAndUserID(groupID, updateUserID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isUserExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no user with such id")
		return
	}
	isValidUser, err := h.isValidUserRole(groupID, userID, "SUPER ADMIN")
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser && updateUserID != userID {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		if err := h.validateJSONTags(body, models.UpdateUserGroupInput{}); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		var input models.UpdateUserGroupInput
		if err := c.BindJSON(&input); err != nil {
			models.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
			return
		}
		if err := h.services.UserGroup.Update(groupID, updateUserID, input); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}

// @Summary Delete User By ID From Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description delete user by id from group
// @ID delete-user-by-id-from-group
// @Produce json
// @Param id path int true "group id"
// @Param userID path int true "user id"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/users/{userID} [delete]
func (h *Handler) deleteUserGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid groupID param")
		return
	}
	deleteUserID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid userID param")
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
	isUserExist, err := h.services.UserGroup.GetByGroupIdAndUserIDAllData(groupID, deleteUserID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if isUserExist == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no user with such id in this group")
		return
	}

	isValidUser, err := h.isValidAdmin(groupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "user validation error")
		return
	}
	if isValidUser && userID != deleteUserID && isUserExist.RoleName == "USER" {
		if err := h.services.UserGroup.Delete(isUserExist.Platform, deleteUserID, groupID); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}

// @Summary Leave Group
// @Security ApiKeyAuth
// @Tags user-groups
// @Description leave group
// @ID leave-group
// @Produce json
// @Param id path int true "group id"
// @Success 200 {object} models.StatusResponse
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/groups/{id}/leave [post]
func (h *Handler) leaveGroup(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	platform, err := getUserPlatform(c)
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

	user, err := h.services.UserGroup.GetByGroupIdAndUserIDAllData(groupID, userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if user != nil && user.RoleName != "SUPER ADMIN" {
		if err := h.services.UserGroup.Delete(platform, userID, groupID); err != nil {
			models.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, models.StatusResponse{Status: "ok"})
	} else {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}
}
