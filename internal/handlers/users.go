package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/sosshik/users-service/pkg/dtos"
	"net/http"
)

// HandleCreateUser handles user creation requests
// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dtos.CreateUserRequest true "User data"
// @Success 200 {object} dtos.CreateUserResponse
// @Failure 422 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Unable to create user"
// @Router /users [post]
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
// @Summary Update an existing user
// @Description Update the user with the given ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body dtos.UpdateUserRequest true "Updated user data"
// @Success 200 {object} dtos.UpdateUserResponse
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Unable to update user"
// @Router /users/{id} [put]
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
// @Summary Delete a user
// @Description Delete the user with the given ID
// @Tags users
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "Successfully deleted user"
// @Failure 500 {object} map[string]string "Unable to delete user"
// @Router /users/{id} [delete]
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
// @Summary Get a list of users
// @Description Retrieve a list of users with optional filtering and pagination. Filter must look like this and be URL encoded: field=value
// @Tags users
// @Produce  json
// @Param page query string false "Page number"
// @Param page_size query string false "Page size"
// @Param filter query string false "Filter query"
// @Success 200 {object} dtos.GetUserResponse
// @Failure 500 {object} map[string]string "Unable to get users"
// @Router /users [get]
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
