package handlers

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get User By ID
// @Security ApiKeyAuth
// @Tags users
// @Description get user by id
// @ID get-user-by-id
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} models.User
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/users/{id} [get]
func (h *Handler) getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}
	user, err := h.services.Authorization.GetUserByID(id)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if user == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no user with such id")
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Get Current User
// @Security ApiKeyAuth
// @Tags users
// @Description get current user
// @ID get-current-user
// @Produce json
// @Success 200 {object} models.User
// @Failure 400,401,403,404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure default {object} models.ErrorResponse
// @Router /api/users/me [get]
func (h *Handler) getMe(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	user, err := h.services.Authorization.GetUserByID(userID)
	if err != nil {
		models.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if user == nil {
		models.NewErrorResponse(c, http.StatusBadRequest, "there is no user with such id")
		return
	}
	c.JSON(http.StatusOK, user)
}