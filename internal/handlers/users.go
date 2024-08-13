package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/sosshik/users-service/internal/dtos"
	"net/http"
)

// HandleCreateUser handles user creation requests
func (h *Handler) HandleCreateUser(c echo.Context) error {
	// Bind the incoming JSON request to CreateUserRequest struct
	var userReq dtos.CreateUserRequest
	if err := c.Bind(&userReq); err != nil {
		log.Warnf("[HandleCreateUser] Unable to decode JSON: %s", err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Invalid request payload"})
	}

	// Validate the request data
	err := userReq.Validate()
	if err != nil {
		log.Warnf("[HandleCreateUser] Invalid request payload: %s", err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Invalid request payload"})
	}

	// Create the user via the service layer
	userResp, err := h.services.CreateUser(userReq)
	if err != nil {
		log.Warnf("[HandleCreateUser] Unable to create user: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to create user: %s", err)})
	}

	// Log success and return the created user response
	log.Infof("[HandleCreateUser] Successfully created user %s with id %s", userResp.Nickname, userResp.ID.String())
	return c.JSON(http.StatusOK, userResp)
}

// HandleUpdateUser handles user update requests
func (h *Handler) HandleUpdateUser(c echo.Context) error {
	// Bind the incoming JSON request to UpdateUserRequest struct
	var userReq dtos.UpdateUserRequest
	if err := c.Bind(&userReq); err != nil {
		log.Warnf("[HandleUpdateUser] Unable to decode JSON: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Update the user by ID via the service layer
	userResp, err := h.services.UpdateUser(c.Param("id"), userReq)
	if err != nil {
		log.Warnf("[HandleUpdateUser] Unable to update user: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to update user: %s", err)})
	}

	// Log success and return the updated user response
	log.Infof("[HandleUpdateUser] Successfully updated user with id %s", userResp.ID.String())
	return c.JSON(http.StatusOK, userResp)
}

// HandleDeleteUser handles user deletion requests
func (h *Handler) HandleDeleteUser(c echo.Context) error {
	// Delete the user by ID via the service layer
	err := h.services.DeleteUser(c.Param("id"))
	if err != nil {
		log.Warnf("[HandleDeleteUser] Unable to delete user: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to delete user: %s", err)})
	}

	// Log success and return a confirmation message
	log.Infof("[HandleDeleteUser] Successfully deleted user with id %s", c.Param("id"))
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Successfully deleted user with id %s", c.Param("id"))})
}

// HandleGetUsers handles requests to retrieve users with optional filtering and pagination
func (h *Handler) HandleGetUsers(c echo.Context) error {
	// Fetch filtered users based on query parameters for pagination and filtering
	response, err := h.services.GetFilteredUsers(c.QueryParam("page"), c.QueryParam("page_size"), c.QueryParam("filter"))
	if err != nil {
		log.Warnf("[HandleGetUsers] Unable to get users: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get users: %s", err)})
	}

	// Return the list of users
	return c.JSON(http.StatusOK, response)
}
