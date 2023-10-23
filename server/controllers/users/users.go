package accounts

import (
	"darco/proto/models/users"

	"github.com/gin-gonic/gin"
)

// @Summary Authenticated user details
// @Description Get details of currently authenticated user
// @id CurrentUser
// @tags People
// @Accept json
// @Produce json
// @Success 200 {object} users.User "Authenticated user details"
// @Failure 400 "User is not authenticated"
// @Router /account [get]
func Current(ctx *gin.Context) {

}

type PasswordUpdateRequest struct {
	Credentials users.UserCredentials `json:"credentials" validate:"required"`
	NewPassword users.PasswordInput   `json:"password" validate:"required"`
}

// @Summary Update user details
// @Description
// @tags People
// @Accept json
// @Produce json
// @Success 200 {object} users.User
// @Failure 300
// @Param data body PasswordUpdateRequest true "Current credentials and new password with confirmation"
func UpdatePassword(ctx *gin.Context) {

}
