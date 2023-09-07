package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
	"net/http"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func InitAuthRoutes(router *gin.Engine) {

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", login)
		authRoutes.GET("/authenticate", Auth(dtos.READER), authenticate)
	}

}

// @Tags Auth
// @Produce json
// @Success 200 {object} dtos.PunqToken
// @Router /auth/login [post]
// @Param body body LoginInput true "LoginInput"
func login(c *gin.Context) {
	input := LoginInput{}

	err := c.MustBindWith(&input, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := services.GetUserByEmail(input.Email)

	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	_, err = user.PasswordCheck(input.Password)
	if err != nil {
		utils.Unauthorized(c, "username or password is incorrect")
		return
	}

	token, err := services.GenerateToken(user)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

// @Tags Auth
// @Produce json
// @Success 200 {object} dtos.PunqToken
// @Router /auth/authenticate [get]
// @Security Bearer
func authenticate(c *gin.Context) {
	if temp, exists := c.Get("user"); exists {
		user, ok := temp.(dtos.PunqUser)
		if !ok {
			utils.MalformedMessage(c, "Type Assertion failed")
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}
	utils.Unauthorized(c, "Unauthorized")
}
