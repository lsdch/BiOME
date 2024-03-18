package accounts

import (
	users "darco/proto/models/people"
	"net/http"

	"github.com/edgedb/edgedb-go"
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
func Current(ctx *gin.Context, db *edgedb.Client, user *users.User) {
	ctx.JSON(http.StatusOK, user)
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
