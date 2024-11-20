package handlers

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	platformCtx         = "platform"
)

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := checkHeaderToken(c)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	isRevoked, err := h.checkIsRevoked(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if isRevoked {
		models.NewErrorResponse(c, http.StatusForbidden, "invalid permission")
		return
	}

	userClaims, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		models.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userClaims.UserID)
	c.Set(platformCtx, userClaims.UserPlatform)
}

// Function to validate JSON tags against the structure
func (h *Handler) validateJSONTags(body []byte, input interface{}) error {

	// Parse the JSON body into a map
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		return err
	}
	structType := reflect.TypeOf(input)

	// Iterate through the fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Get the JSON tag value for the field
		tagValue := field.Tag.Get("json")
		if _, ok := jsonMap["id"]; ok {
			return fmt.Errorf("invalid JSON tags: id")
		}
		// Check if the JSON tag is not empty
		if tagValue != "" {
			// Check if the field exists in the struct
			if _, ok := jsonMap[tagValue]; ok {
				delete(jsonMap, tagValue)
			}
		}
	}
	if len(jsonMap) != 0 {
		var errorTags []string
		for key := range jsonMap {
			errorTags = append(errorTags, key)
		}

		return fmt.Errorf("invalid JSON tags: %s", strings.Join(errorTags[:], ", "))
	}
	return nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}

func (h *Handler) isValidUserRole(groupID, userID int, role string) (bool, error) {
	user, err := h.services.UserGroup.GetByGroupIdAndUserIDRole(groupID, userID, role)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func (h *Handler) isValidAdmin(groupID, userID int) (bool, error) {
	user, err := h.services.UserGroup.GetByGroupIdAndUserIDAllData(groupID, userID)
	if err != nil {
		return false, err
	}
	if user.RoleName == "ADMIN" || user.RoleName == "SUPER ADMIN" {
		return true, nil
	}
	return false, nil
}

func (h *Handler) isValidUser(groupID, userID int) (bool, error) {
	user, err := h.services.UserGroup.GetByGroupIdAndUserID(groupID, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func getUserPlatform(c *gin.Context) (string, error) {
	userPlatform, ok := c.Get(platformCtx)
	if !ok {
		return "", errors.New("user platform not found")
	}
	userPlatformStr, ok := userPlatform.(string)
	if !ok {
		return "", errors.New("user platform not found")
	}

	return userPlatformStr, nil
}

func checkHeaderToken(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}
	return headerParts[1], nil
}

func (h *Handler) checkIsRevoked(tokenValue string) (bool, error) {
	token, err := h.services.Token.GetByToken(tokenValue)
	if err != nil {
		return true, err
	}
	if token.Revoked {
		return true, nil
	} else {
		return false, nil
	}
}
